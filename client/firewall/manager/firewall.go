package manager

import (
	"fmt"
	"net"
	"net/netip"

	log "github.com/sirupsen/logrus"
)

const (
	ForwardingFormatPrefix = "netbird-fwd-"
	ForwardingFormat       = "netbird-fwd-%s-%t"
	NatFormat              = "netbird-nat-%s-%t"
)

// Rule abstraction should be implemented by each firewall manager
//
// Each firewall type for different OS can use different type
// of the properties to hold data of the created rule
type Rule interface {
	// GetRuleID returns the rule id
	GetRuleID() string
}

// RuleDirection is the traffic direction which a rule is applied
type RuleDirection int

const (
	// RuleDirectionIN applies to filters that handlers incoming traffic
	RuleDirectionIN RuleDirection = iota
	// RuleDirectionOUT applies to filters that handlers outgoing traffic
	RuleDirectionOUT
)

// Action is the action to be taken on a rule
type Action int

const (
	// ActionAccept is the action to accept a packet
	ActionAccept Action = iota
	// ActionDrop is the action to drop a packet
	ActionDrop
)

// Manager is the high level abstraction of a firewall manager
//
// It declares methods which handle actions required by the
// Netbird client for ACL and routing functionality
type Manager interface {
	// AllowNetbird allows netbird interface traffic
	AllowNetbird() error

	// AddPeerFiltering adds a rule to the firewall
	//
	// If comment argument is empty firewall manager should set
	// rule ID as comment for the rule
	AddPeerFiltering(
		ip net.IP,
		proto Protocol,
		sPort *Port,
		dPort *Port,
		direction RuleDirection,
		action Action,
		ipsetName string,
		comment string,
	) ([]Rule, error)

	// DeletePeerRule from the firewall by rule definition
	DeletePeerRule(rule Rule) error

	// IsServerRouteSupported returns true if the firewall supports server side routing operations
	IsServerRouteSupported() bool

	AddRouteFiltering(
		source netip.Prefix,
		destination netip.Prefix,
		proto Protocol,
		sPort *Port,
		dPort *Port,
		action Action,
	) (Rule, error)

	// DeleteRouteRule deletes a routing rule
	DeleteRouteRule(rule Rule) error

	// AddNatRule inserts a routing NAT rule
	AddNatRule(pair RouterPair) error

	// RemoveNatRule removes a routing NAT rule
	RemoveNatRule(pair RouterPair) error

	// SetLegacyManagement sets the legacy management mode
	SetLegacyManagement(legacy bool) error

	// Reset firewall to the default state
	Reset() error

	// Flush the changes to firewall controller
	Flush() error
}

func GenKey(format string, pair RouterPair) string {
	return fmt.Sprintf(format, pair.ID, pair.Inverse)
}

// LegacyManager defines the interface for legacy management operations
type LegacyManager interface {
	RemoveAllLegacyRouteRules() error
	GetLegacyManagement() bool
	SetLegacyManagement(bool)
}

// SetLegacyManagement sets the route manager to use legacy management
func SetLegacyManagement(router LegacyManager, isLegacy bool) error {
	oldLegacy := router.GetLegacyManagement()

	if oldLegacy != isLegacy {
		router.SetLegacyManagement(isLegacy)
		log.Debugf("Set legacy management to %v", isLegacy)
	}

	// client reconnected to a newer mgmt, we need to clean up the legacy rules
	if !isLegacy && oldLegacy {
		if err := router.RemoveAllLegacyRouteRules(); err != nil {
			return fmt.Errorf("remove legacy routing rules: %v", err)
		}

		log.Debugf("Legacy routing rules removed")
	}

	return nil
}
