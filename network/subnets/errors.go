package subnets

import "errors"

var (
	// ErrInsufficientSubnets is returned by AllocateDynamically if no more subnets can be allocated.
	ErrInsufficientSubnets = errors.New("insufficient subnets remaining in the pool")

	// ErrInsufficientIPs is returned by AllocateDynamically if no more IPs can be allocated.
	ErrInsufficientIPs = errors.New("insufficient IPs remaining in the pool")

	// ErrReleasedUnallocatedNetwork is returned by Release if the subnet is not allocated.
	ErrReleasedUnallocatedSubnet = errors.New("subnet is not allocated")

	// ErrOverlapsExistingSubnet is returned if a requested subnet overlaps an existing, allocated subnet
	ErrOverlapsExistingSubnet = errors.New("subnet overlaps an existing subnet")

	// ErrInvalidRange is returned by AllocateStatically and by Recover if the subnet range is invalid.
	ErrInvalidRange = errors.New("subnet has invalid range")

	// ErrNotAllowed is returned by AllocateStatically if the subnet range overlaps the dynamic allocation range
	// and by Recover if the subnet range contains the dynamic allocation range.
	ErrNotAllowed = errors.New("the requested range cannot be allocated statically")

	// ErrInvalidIP is returned if a static IP is requested inside a subnet
	// which does not contain that IP
	ErrInvalidIP = errors.New("the requested IP is not within the subnet")

	// ErrIPAlreadyAllocated is returned if a static IP is requested which has already been allocated
	ErrIPAlreadyAllocated = errors.New("the requested IP is already allocated")

	// ErrIpCannotBeNil is returned by Release(..) and Recover(..) if a nil
	// IP address is passed.
	ErrIpCannotBeNil = errors.New("the IP field cannot be empty")

	ErrIPEqualsGateway   = errors.New("a container IP must not equal the gateway IP")
	ErrIPEqualsBroadcast = errors.New("a container IP must not equal the broadcast IP")
)
