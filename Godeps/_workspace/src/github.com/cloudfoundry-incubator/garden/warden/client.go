package warden

import (
	"time"
)

type Client interface {
	Ping() error

	Capacity() (Capacity, error)

	Create(ContainerSpec) (Container, error)
	Destroy(handle string) error
	Containers(Properties) ([]Container, error)
	Lookup(handle string) (Container, error)
}

type ContainerSpec struct {
	Handle     string
	GraceTime  time.Duration
	RootFSPath string
	BindMounts []BindMount

	// The container's IP address. This allocates a network with a corresponding /30 CIDR block.
	// The remaining 2 bits denote:
	// 00 network
	// 01 host IP
	// 10 container IP
	// 11 subnet mask
	Network    string

	Properties Properties
	Env        []string
}

const ContainerNetworkCIDRPrefixSize = 30

type BindMount struct {
	SrcPath string
	DstPath string
	Mode    BindMountMode
	Origin  BindMountOrigin
}

type Capacity struct {
	MemoryInBytes uint64
	DiskInBytes   uint64
	MaxContainers uint64
}

type Properties map[string]string

type BindMountMode uint8

const BindMountModeRO BindMountMode = 0
const BindMountModeRW BindMountMode = 1

type BindMountOrigin uint8

const BindMountOriginHost BindMountOrigin = 0
const BindMountOriginContainer BindMountOrigin = 1
