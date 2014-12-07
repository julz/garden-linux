package network

import (
	"net"

	"github.com/cloudfoundry-incubator/garden-linux/network/subnets"
)

//
type BridgingSubnet struct {
	pool subnets.Subnets
}

func (p *BridgingSubnet) Allocate(sn subnets.SubnetSelector, i subnets.IPSelector) (subnet *net.IPNet, ip net.IP, err error) {
	return nil, nil, nil
}

func (p *BridgingSubnet) Release(subnet *net.IPNet, containerIP net.IP) (bool, error) {
	return false, nil
}

func (p *BridgingSubnet) Recover(subnet *net.IPNet, containerIP net.IP) error {
	return nil
}
