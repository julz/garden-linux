// Package ip_pool provides IP address pools.
package ip_pool

import (
	"bytes"
	"errors"
	"net"
	"sync"
)

// An IPPool is a collection of IP addresses.
type IPPool interface {
	// Allocates an IP address from the pool unless the pool is empty in
	// which case returns an error. In the successful case, the IP address
	// is removed from the pool. Otherwise the pool is not modified.
	Allocate() (net.IP, error)

	// Returns a previously allocated IP address back to the pool. If the
	// IP address was not allocated from the pool, returns an error and
	// does not modify the pool.
	Release(net.IP) error
}

// Errors.
var (
	ErrInvalidRange         = errors.New("invalid IP pool range")
	ErrUnallocatedIPAddress = errors.New("cannot release an unallocated IP address")
	ErrPoolEmpty            = errors.New("no more IP addresses are available")
)

type pool struct {
	mutex sync.Mutex

	min, max net.IP
	prev     net.IP // Invariant: min <= prev <= max

	alloc map[string]bool
}

func NewFromIPNet(ip *net.IPNet) (IPPool, error) {
	return nil, nil
}

// New creates a new IP pool containing all the IP addresses between the given minimum and maximum, inclusive.
// Returns an error if the maximum IP address is less than the minimum.
func New(min, max net.IP) (IPPool, error) {
	if bytes.Compare(min, max) > 0 {
		return nil, ErrInvalidRange
	}

	return &pool{min: min, max: max, prev: max, alloc: make(map[string]bool)}, nil
}

func (p *pool) Allocate() (net.IP, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if next := p.nextAvailable(); next != nil {
		p.allocate(next)
		p.prev = next
		return next, nil
	}

	return nil, ErrPoolEmpty
}

func (p *pool) nextAvailable() net.IP {
	return p.find(p.prev, func(i net.IP) bool {
		if !p.allocated(i) {
			return true
		}

		return false
	})
}

func (p *pool) find(end net.IP, match func(net.IP) bool) net.IP {
	for i := p.next(end); ; i = p.next(i) {
		if match(i) {
			return i
		}
		if i.Equal(end) {
			break
		}
	}

	return nil
}

func (p *pool) Release(ip net.IP) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.allocated(ip) {
		return ErrUnallocatedIPAddress
	}

	p.release(ip)
	return nil
}

func (p *pool) allocated(ip net.IP) bool {
	return p.alloc[ip.String()]
}

func (p *pool) allocate(ip net.IP) {
	p.alloc[ip.String()] = true
}

func (p *pool) release(ip net.IP) {
	delete(p.alloc, ip.String())
}

// Returns the next IP in the pool after the given IP, wrapping from p.max to p.min if necessary.
func (p *pool) next(ip net.IP) net.IP {
	if ip.Equal(p.max) {
		return p.min
	}

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
