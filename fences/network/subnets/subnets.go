// The subnets package provides a subnet pool from which networks may be dynamically acquired or
// statically reserved.
package subnets

import (
	"errors"
	"net"
	"sync"
)

// A Subnets provides a means of allocating /30 subnets.
type Subnets interface {
	// Dynamically allocates a /30 subnet, or returns an error if no more subnets can be allocated.
	AllocateDynamically() (*net.IPNet, error)

	// Statically allocates the given subnet, and returns an error if the subnet cannot be newly allocated
	AllocateStatically(subnet *net.IPNet) error

	// Releases an allocated network.
	Release(*net.IPNet) error

	// Recovers an unallocated network so it appears to be allocated.
	Recover(*net.IPNet) error

	// Returns the number of /30 subnets which can be Allocate(d)Dynamically.
	Capacity() int
}

type subnetpool struct {
	mutex sync.Mutex

	dynamicAllocationNet *net.IPNet
	pool                 []*net.IPNet // Unallocated /30 subnets in dynamicAllocationNet
	capacity             int          // Number of /30 subnets in dynamicAllocationNet

	static []*net.IPNet // Statically allocated subnets, disjoint from dynamicAllocationNet
}

var (
	// ErrInsufficientSubnets is returned by AllocateDynamically if no more subnets can be allocated.
	ErrInsufficientSubnets = errors.New("insufficient subnets remaining in the pool")

	// ErrReleasedUnallocatedNetwork is returned by Release if the subnet is not allocated.
	ErrReleasedUnallocatedSubnet = errors.New("subnet is not allocated")

	// ErrAlreadyAllocated is returned by AllocateStatically and by Recover if the subnet is already allocated.
	ErrAlreadyAllocated = errors.New("subnet is already allocated")

	// ErrInvalidRange is returned by AllocateStatically and by Recover if the subnet range is invalid.
	ErrInvalidRange = errors.New("subnet has invalid range")

	// ErrNotAllowed is returned by AllocateStatically if the subnet range overlaps the dynamic allocation range
	// and by Recover if the subnet range contains the dynamic allocation range.
	ErrNotAllowed = errors.New("the requested range cannot be allocated statically")
)

var slash30mask net.IPMask

func init() {
	_, maskedNetwork, err := net.ParseCIDR("1.1.1.1/30")
	if err != nil {
		panic("Does not compute")
	}

	slash30mask = maskedNetwork.Mask
}

// New creates a Subnets implementation from a dynamic allocation range.
//
// All dynamic allocations come from the range, static allocations are prohibited from the range.
func New(ipNet *net.IPNet) (Subnets, error) {
	pool := poolOfSubnets(ipNet)
	return &subnetpool{dynamicAllocationNet: ipNet, pool: pool, capacity: len(pool)}, nil
}

func poolOfSubnets(ipNet *net.IPNet) []*net.IPNet {
	min := ipNet.IP
	pool := make([]*net.IPNet, 0)
	for ip := min; ipNet.Contains(ip); ip = next(ip) {
		subnet := &net.IPNet{ip, slash30mask}
		ip = next(next(next(ip)))
		if ipNet.Contains(ip) {
			pool = append(pool, subnet)
		}
	}

	return pool
}

func (m *subnetpool) Capacity() int {
	return m.capacity
}

func (m *subnetpool) AllocateStatically(ipNet *net.IPNet) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if overlaps(ipNet, m.dynamicAllocationNet) {
		return ErrNotAllowed
	}

	for _, s := range m.static {
		if overlaps(s, ipNet) {
			return ErrAlreadyAllocated
		}
	}

	m.static = append(m.static, ipNet)

	return nil
}

func overlaps(net1, net2 *net.IPNet) bool {
	return net1.Contains(net2.IP) || net2.Contains(net1.IP)
}

func (m *subnetpool) Recover(ipNet *net.IPNet) error {
	if !m.dynamicAllocationNet.Contains(ipNet.IP) {
		return m.AllocateStatically(ipNet)
	}

	found := -1
	for i, s := range m.pool {
		if s.IP.Equal(ipNet.IP) {
			found = i
		}
	}

	if found > -1 {
		m.pool = append(m.pool[:found], m.pool[found+1:]...)
		return nil
	}

	return ErrAlreadyAllocated
}

func (m *subnetpool) AllocateDynamically() (*net.IPNet, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.pool) == 0 {
		return nil, ErrInsufficientSubnets
	}

	acquired := m.pool[0]
	m.pool = m.pool[1:]

	return acquired, nil
}

func (m *subnetpool) Release(ipNet *net.IPNet) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, n := range m.pool {
		if n.IP.Equal(ipNet.IP) {
			return ErrReleasedUnallocatedSubnet
		}
	}

	found := -1
	for i, s := range m.static {
		if s.IP.Equal(ipNet.IP) {
			found = i
		}
	}

	if found > -1 {
		m.static = append(m.static[:found], m.static[found+1:]...)
	}

	m.pool = append(m.pool, ipNet)
	return nil
}

func next(ip net.IP) net.IP {
	next := clone(ip)
	for i := len(next) - 1; i >= 0; i-- {
		next[i]++
		if next[i] != 0 {
			return next
		}
	}

	panic("overflowed maximum IP")
}

func clone(ip net.IP) net.IP {
	clone := make([]byte, len(ip))
	copy(clone, ip)
	return clone
}