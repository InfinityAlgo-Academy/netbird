//go:build (linux || darwin) && !android

package iface

import (
	"net"
	"os"

	"github.com/pion/transport/v2"
	"golang.zx2c4.com/wireguard/ipc"

	"github.com/netbirdio/netbird/iface/bind"

	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
)

type tunDevice struct {
	name         string
	address      WGAddress
	mtu          int
	netInterface NetInterface
	iceBind      *bind.ICEBind
	uapi         net.Listener
	tunDevice    *device.Device
	close        chan struct{}
}

func newTunDevice(name string, address WGAddress, mtu int, transportNet transport.Net) *tunDevice {
	return &tunDevice{
		name:    name,
		address: address,
		mtu:     mtu,
		iceBind: bind.NewICEBind(transportNet),
		close:   make(chan struct{}),
	}
}

func (c *tunDevice) UpdateAddr(address WGAddress) error {
	c.address = address
	return c.assignAddr()
}

func (c *tunDevice) WgAddress() WGAddress {
	return c.address
}

func (c *tunDevice) DeviceName() string {
	return c.name
}

func (c *tunDevice) Close() error {

	select {
	case c.close <- struct{}{}:
	default:
	}

	var err1, err2, err3 error
	if c.netInterface != nil {
		err1 = c.netInterface.Close()
	}

	if c.uapi != nil {
		err2 = c.uapi.Close()
	}

	sockPath := "/var/run/wireguard/" + c.name + ".sock"
	if _, err3 = os.Stat(sockPath); err3 == nil {
		err3 = os.Remove(sockPath)
	}

	if err1 != nil {
		return err1
	}

	if err2 != nil {
		return err2
	}

	return err3
}

// createWithUserspace Creates a new Wireguard interface, using wireguard-go userspace implementation
func (c *tunDevice) createWithUserspace() (NetInterface, error) {
	tunIface, err := tun.CreateTUN(c.name, c.mtu)
	if err != nil {
		return nil, err
	}

	// We need to create a wireguard-go device and listen to configuration requests
	tunDevice := device.NewDevice(tunIface, c.iceBind, device.NewLogger(device.LogLevelSilent, "[wiretrustee] "))
	err = tunDevice.Up()
	if err != nil {
		return nil, c.Close()
	}

	c.uapi, err = c.getUAPI(c.name)
	if err != nil {
		return tunIface, c.Close()
	}

	go func() {
		for {
			select {
			case <-c.close:
				log.Debugf("exit uapi.Accept()")
				return
			default:
			}
			uapiConn, uapiErr := c.uapi.Accept()
			if uapiErr != nil {
				log.Traceln("uapi Accept failed with error: ", uapiErr)
				continue
			}
			go func() {
				tunDevice.IpcHandle(uapiConn)
				log.Debugf("exit tunDevice.IpcHandle")
			}()
		}
	}()

	log.Debugln("UAPI listener started")
	return tunIface, nil
}

// getUAPI returns a Listener
func (c *tunDevice) getUAPI(iface string) (net.Listener, error) {
	tunSock, err := ipc.UAPIOpen(iface)
	if err != nil {
		return nil, err
	}
	return ipc.UAPIListen(iface, tunSock)
}
