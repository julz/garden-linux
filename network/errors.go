package network

import (
	"fmt"
	"net"

	"github.com/cloudfoundry-incubator/garden-linux/network/iptables"
)

// VethPairCreationError is returned if creating a virtual ethernet pair fails
type VethPairCreationError struct {
	Cause                         error
	HostIfcName, ContainerIfcName string
}

func (err VethPairCreationError) Error() string {
	return fmtErr("failed to create veth pair with host interface name '%s', container interface name '%s': %v", err.HostIfcName, err.ContainerIfcName, err.Cause)
}

// MTUError is returned if setting the Mtu on an interface fails
type MTUError struct {
	Cause error
	Intf  *net.Interface
	MTU   int
}

func (err MTUError) Error() string {
	return fmtErr("failed to set interface '%v' mtu to %d", err.Intf, err.MTU, err.Cause)
}

type SetNsFailedError struct {
	Cause error
	Intf  *net.Interface
	Pid   int
}

func (err SetNsFailedError) Error() string {
	return fmtErr("failed to move interface %v in to pid namespace %d: %v", err.Intf, err.Pid, err.Cause)
}

// BridgeCreationError is returned if an error occurs while creating a bridge
type BridgeCreationError struct {
	Cause  error
	Name   string
	IP     net.IP
	Subnet *net.IPNet
}

func (err BridgeCreationError) Error() string {
	return fmtErr("failed to create bridge with name '%s', IP '%s', subnet '%s': %v", err.Name, err.IP, err.Subnet, err.Cause)
}

// AddToBridgeError is returned if an error occurs while adding an interface to a bridge
type AddToBridgeError struct {
	Cause         error
	Bridge, Slave *net.Interface
}

func (err AddToBridgeError) Error() string {
	return fmtErr("failed to add slave %s to bridge %s: %v", err.Slave.Name, err.Bridge.Name, err.Cause)
}

// LinkUpError is returned if brinding an interface up fails
type LinkUpError struct {
	Cause error
	Link  *net.Interface
	Role  string
}

func (err LinkUpError) Error() string {
	return fmtErr("failed to bring %s link %s up: %v", err.Role, err.Link.Name, err.Cause)
}

// FindLinkError is returned if an expected interface cannot be found inside the container
type FindLinkError struct {
	Cause error // may be nil if no error occurred other than the link not existing
	Role  string
	Name  string
}

func (err FindLinkError) Error() string {
	return fmtErr("failed to find %s link with name %s", err.Role, err.Name)
}

// ConfigureLinkError is returned if configuring a link fails
type ConfigureLinkError struct {
	Cause          error
	Role           string
	Interface      *net.Interface
	IntendedIP     net.IP
	IntendedSubnet *net.IPNet
}

func (err ConfigureLinkError) Error() string {
	return fmtErr("failed to configure %s link (%v) to IP %v, subnet %v", err.Role, err.Interface, err.IntendedIP, err.IntendedSubnet)
}

// ConfigureDefaultGWError is returned if the default gateway cannot be updated
type ConfigureDefaultGWError struct {
	Cause     error
	Interface *net.Interface
	IP        net.IP
}

func (err ConfigureDefaultGWError) Error() string {
	return fmtErr("failed to set default gateway to IP %v via device %v", err.IP, err.Interface, err.Cause)
}

// IPTablesError is returned if an iptable command fails
type IPTablesError struct {
	Cause  error
	Action string
	Rule   iptables.Rule
}

func (err IPTablesError) Error() string {
	return fmtErr("failed to %s iptable rule %v: %v", err.Action, err.Rule, err.Cause)
}

func fmtErr(msg string, args ...interface{}) string {
	return fmt.Sprintf("network: "+msg, args...)
}
