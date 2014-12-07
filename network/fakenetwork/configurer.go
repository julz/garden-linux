package fakenetwork

import "net"

type FakeConfigurer struct {
	ConfiguredSubnets            []ConfiguredSubnet
	ConfigureSubnetsShouldReturn error
}

type ConfiguredSubnet struct {
	ExternalIP   string
	BridgeName   string
	BridgeIP     string
	BridgeSubnet string
}

func (c *FakeConfigurer) ConfigureSubnet(bridgeName string, externalIP, bridgeIP net.IP, subnet *net.IPNet) error {
	c.ConfiguredSubnets = append(c.ConfiguredSubnets, ConfiguredSubnet{
		ExternalIP:   externalIP.String(),
		BridgeName:   bridgeName,
		BridgeIP:     bridgeIP.String(),
		BridgeSubnet: subnet.String(),
	})

	return c.ConfigureSubnetsShouldReturn
}
