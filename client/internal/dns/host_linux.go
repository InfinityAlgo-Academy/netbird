package dns

import (
	"bufio"
	"context"
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/miekg/dns"
	"github.com/netbirdio/netbird/iface"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"net"
	"net/netip"
	"os"
	"strings"
	"time"
)

const (
	defaultResolvConfPath = "/etc/resolv.conf"
)

const (
	fileManager osManagerType = iota
	networkManager
	systemdManager
	resolvConfManager
)

const (
	systemdDbusManagerInterface            = "org.freedesktop.resolve1.Manager"
	systemdResolvedDest                    = "org.freedesktop.resolve1"
	systemdDbusObjectNode                  = "/org/freedesktop/resolve1"
	systemdDbusGetLinkMethod               = systemdDbusManagerInterface + ".GetLink"
	systemdDbusFlushCachesMethod           = systemdDbusManagerInterface + ".FlushCaches"
	systemdDbusLinkInterface               = "org.freedesktop.resolve1.Link"
	systemdDbusRevertMethodSuffix          = systemdDbusLinkInterface + ".Revert"
	systemdDbusSetDNSMethodSuffix          = systemdDbusLinkInterface + ".SetDNS"
	systemdDbusSetDefaultRouteMethodSuffix = systemdDbusLinkInterface + ".SetDefaultRoute"
	systemdDbusSetDomainsMethodSuffix      = systemdDbusLinkInterface + ".SetDomains"
	systemdDbusDefaultFlag                 = 0
)

type systemdDbusConfigurator struct {
	dbusLinkInterface    dbus.ObjectPath
	createdLinkedDomains map[string]systemdDbusLinkDomainsInput
}

// https://dbus.freedesktop.org/doc/dbus-specification.html
// https://www.freedesktop.org/software/systemd/man/org.freedesktop.resolve1.html
type systemdDbusDNSInput struct {
	Family  int32
	Address []byte
}

type systemdDbusLinkDomainsInput struct {
	Domain    string
	MatchOnly bool
}

type osManagerType int

func newHostManager(wgInterface *iface.WGIface) hostManager {
	switch getOSDNSManagerType() {
	default:
		log.Debugf("discovered mode is: %d", getOSDNSManagerType())
		return newSystemdDbusConfigurator(wgInterface)
	}
}

func newSystemdDbusConfigurator(wgInterface *iface.WGIface) hostManager {
	iface, err := net.InterfaceByName(wgInterface.GetName())
	if err != nil {
		// todo add proper error handling
		panic(err)
	}

	obj, closeConn, err := getDbusObject(systemdResolvedDest, systemdDbusObjectNode)
	if err != nil {
		panic(err)
	}
	defer closeConn()
	var s string
	err = obj.Call(systemdDbusGetLinkMethod, systemdDbusDefaultFlag, iface.Index).Store(&s)
	if err != nil {
		// todo add proper error handling
		panic(err)
	}

	log.Debugf("got dbus Link interface: %s from net interface %s and index %d", s, iface.Name, iface.Index)

	return &systemdDbusConfigurator{
		dbusLinkInterface:    dbus.ObjectPath(s),
		createdLinkedDomains: make(map[string]systemdDbusLinkDomainsInput),
	}
}

func getOSDNSManagerType() osManagerType {
	file, err := os.Open(defaultResolvConfPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text[0] != '#' {
			return fileManager
		}
		if strings.Contains(text, "NetworkManager") {
			return networkManager
		}
		if strings.Contains(text, "systemd-resolved") && isDbusListenerRunning(systemdResolvedDest, systemdDbusObjectNode) {
			return systemdManager
		}
		if strings.Contains(text, "resolvconf") {
			return resolvConfManager
		}
	}
	return fileManager
}

func (s *systemdDbusConfigurator) applyDNSSettings(domains []string, ip string, port int) error {
	parsedIP := netip.MustParseAddr(ip).As4()
	defaultLinkInput := systemdDbusDNSInput{
		Family:  unix.AF_INET,
		Address: parsedIP[:],
	}
	err := s.callLinkMethod(systemdDbusSetDNSMethodSuffix, defaultLinkInput)
	if err != nil {
		return fmt.Errorf("setting the interface DNS server %s:%d failed with error: %s", ip, port, err)
	}

	var domainsInput []systemdDbusLinkDomainsInput

	for _, domain := range domains {
		if isRootZoneDomain(domain) {
			err = s.callLinkMethod(systemdDbusSetDefaultRouteMethodSuffix, true)
			if err != nil {
				log.Errorf("setting link as default dns router, failed with error: %s", err)
			}
		}
		domainsInput = append(domainsInput, systemdDbusLinkDomainsInput{
			Domain:    dns.Fqdn(domain),
			MatchOnly: true,
		})
	}
	err = s.addDNSStateForDomain(domainsInput)
	if err != nil {
		log.Error(err)
	}
	return nil
}

func (s *systemdDbusConfigurator) addDNSSetupForAll() error {
	err := s.callLinkMethod(systemdDbusSetDefaultRouteMethodSuffix, true)
	if err != nil {
		return fmt.Errorf("setting link as default dns router, failed with error: %s", err)
	}
	return nil
}

func (s *systemdDbusConfigurator) addDNSStateForDomain(domainsInput []systemdDbusLinkDomainsInput) error {
	err := s.callLinkMethod(systemdDbusSetDomainsMethodSuffix, domainsInput)
	if err != nil {
		return fmt.Errorf("setting domains configuration failed with error: %s", err)
	}
	for _, input := range domainsInput {
		s.createdLinkedDomains[input.Domain] = input
	}
	return nil
}

func (s *systemdDbusConfigurator) addSearchDomain(domain string, ip string, port int) error {
	var newDomainsInput []systemdDbusLinkDomainsInput

	fqdnDomain := dns.Fqdn(domain)

	existingDomain, found := s.createdLinkedDomains[fqdnDomain]
	if found && !existingDomain.MatchOnly {
		return nil
	}

	delete(s.createdLinkedDomains, fqdnDomain)
	for _, existingInput := range s.createdLinkedDomains {
		newDomainsInput = append(newDomainsInput, existingInput)
	}

	newDomainsInput = append(newDomainsInput, systemdDbusLinkDomainsInput{
		Domain:    fqdnDomain,
		MatchOnly: false,
	})

	err := s.addDNSStateForDomain(newDomainsInput)
	if err != nil {
		return fmt.Errorf("setting domains configuration with search domain %s failed with error: %s", domain, err)
	}

	return s.flushCaches()
}
func (s *systemdDbusConfigurator) removeDomainSettings(domains []string) error {
	var err error
	for _, domain := range domains {
		if isRootZoneDomain(domain) {
			err = s.callLinkMethod(systemdDbusSetDefaultRouteMethodSuffix, false)
			if err != nil {
				log.Errorf("setting link as non default dns router, failed with error: %s", err)
			}
			break
		}
	}

	// cleaning the configuration as it gets rebuild
	emptyList := make([]systemdDbusLinkDomainsInput, 0)

	err = s.callLinkMethod(systemdDbusSetDomainsMethodSuffix, emptyList)
	if err != nil {
		log.Error(err)
	}

	s.createdLinkedDomains = make(map[string]systemdDbusLinkDomainsInput)

	return s.flushCaches()
}
func (s *systemdDbusConfigurator) removeDNSSettings() error {
	err := s.callLinkMethod(systemdDbusRevertMethodSuffix, nil)
	if err != nil {
		return fmt.Errorf("unable to revert link configuration, got error: %s", err)
	}
	return s.flushCaches()
}

func (s *systemdDbusConfigurator) flushCaches() error {
	obj, closeConn, err := getDbusObject(systemdResolvedDest, systemdDbusObjectNode)
	if err != nil {
		return fmt.Errorf("got error while attempting to retrieve the object %s, err: %s", systemdDbusObjectNode, err)
	}
	defer closeConn()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	err = obj.CallWithContext(ctx, systemdDbusFlushCachesMethod, systemdDbusDefaultFlag).Store()
	if err != nil {
		return fmt.Errorf("got error while calling the FlushCaches method with context, err: %s", err)
	}

	return nil
}

func (s *systemdDbusConfigurator) callLinkMethod(method string, value any) error {
	obj, closeConn, err := getDbusObject(systemdResolvedDest, s.dbusLinkInterface)
	if err != nil {
		return fmt.Errorf("got error while attempting to retrieve the object, err: %s", err)
	}
	defer closeConn()

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	err = obj.CallWithContext(ctx, method, systemdDbusDefaultFlag, value).Store()
	if err != nil {
		return fmt.Errorf("got error while calling command with context, err: %s", err)
	}

	return nil
}

func isDbusListenerRunning(dest string, path dbus.ObjectPath) bool {
	obj, closeConn, err := getDbusObject(dest, path)
	if err != nil {
		return false
	}
	defer closeConn()

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	err = obj.CallWithContext(ctx, "org.freedesktop.DBus.Peer.Ping", 0).Store()
	if err != nil {
		return false
	}
	return true
}

func getDbusObject(dest string, path dbus.ObjectPath) (dbus.BusObject, func() error, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, nil, err
	}
	obj := conn.Object(dest, path)

	return obj, conn.Close, nil
}
