// Package net_fence provides Garden's networking function.
package net_fence

import (
	"flag"
	"fmt"
	"github.com/cloudfoundry-incubator/garden-linux/net_fence/ip_pool"
	"net"
)

var config = struct {
	network string
}{}

const (
	DefaultNetworkPool = "10.254.0.0/22"
)

func InitializeFlags(flagset *flag.FlagSet) {
	flagset.StringVar(&config.network, "networkPool",
		DefaultNetworkPool,
		"Pool of IP addresses for container networks")
}

func Initialize() error {
	_, network, err := net.ParseCIDR(config.network)
	if err != nil {
		return fmt.Errorf("Invalid networkPool flag: %s", err)
	}

	NewIpPoolFromIPNet(network)
	return nil
}

var NewIpPoolFromIPNet = ip_pool.NewFromIPNet
