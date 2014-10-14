package subnets_test

import (
	"net"
	"runtime"

	"github.com/cloudfoundry-incubator/garden-linux/net_fence/subnets"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Subnets", func() {
	var manager subnets.Manager
	var defaultSubnetPool *net.IPNet

	Describe("Creating with a /32 IP range", func() {
		It("returns an error", func() {
			var err error
			_, defaultSubnetPool, err = net.ParseCIDR("10.2.3.0/32")
			Ω(err).ShouldNot(HaveOccurred())

			manager, err = subnets.New(defaultSubnetPool)
			Ω(err).Should(HaveOccurred())
		})
	})

	Describe("Allocating and Releasing", func() {
		JustBeforeEach(func() {
			var err error
			manager, err = subnets.New(defaultSubnetPool)
			Ω(err).ShouldNot(HaveOccurred())
		})

		Context("when the pool does not have sufficient IPs to allocate a subnet", func() {
			BeforeEach(func() {
				var err error
				_, defaultSubnetPool, err = net.ParseCIDR("10.2.3.0/31")
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("the first request returns an error", func() {
				_, err := manager.AllocateDynamically()
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("when the pool has sufficient IPs to allocate a single subnet", func() {
			BeforeEach(func() {
				var err error
				_, defaultSubnetPool, err = net.ParseCIDR("10.2.3.0/30")
				Ω(err).ShouldNot(HaveOccurred())
			})

			Context("the first request", func() {
				It("succeeds, and returns a /30 network within the subnet", func() {
					network, err := manager.AllocateDynamically()
					Ω(err).ShouldNot(HaveOccurred())

					Ω(network).ShouldNot(BeNil())
					Ω(network.String()).Should(Equal("10.2.3.0/30"))
				})
			})

			Context("subsequent requests", func() {
				It("fail, and return an err", func() {
					_, err := manager.AllocateDynamically()
					Ω(err).ShouldNot(HaveOccurred())

					_, err = manager.AllocateDynamically()
					Ω(err).Should(HaveOccurred())
				})
			})

			Context("when an allocated network is released", func() {
				It("a subsequent allocation succeeds, and returns the first network again", func() {
					// first
					allocated, err := manager.AllocateDynamically()
					Ω(err).ShouldNot(HaveOccurred())

					// second - will fail (sanity check)
					_, err = manager.AllocateDynamically()
					Ω(err).Should(HaveOccurred())

					// release
					err = manager.Release(allocated)
					Ω(err).ShouldNot(HaveOccurred())

					// third - should work now because of release
					network, err := manager.AllocateDynamically()
					Ω(err).ShouldNot(HaveOccurred())

					Ω(network).ShouldNot(BeNil())
					Ω(network.String()).Should(Equal(allocated.String()))
				})
			})

			Context("when a network is released twice", func() {
				It("returns an error", func() {
					// first
					allocated, err := manager.AllocateDynamically()
					Ω(err).ShouldNot(HaveOccurred())

					// release
					err = manager.Release(allocated)
					Ω(err).ShouldNot(HaveOccurred())

					// release again
					err = manager.Release(allocated)
					Ω(err).Should(HaveOccurred())
					Ω(err).Should(Equal(subnets.ErrReleasedUnallocatedNetwork))
				})
			})

			It("allocates distinct networks concurrently", func() {
				prev := runtime.GOMAXPROCS(2)
				defer runtime.GOMAXPROCS(prev)

				Consistently(func() bool {
					_, network, err := net.ParseCIDR("10.0.0.0/29")
					Ω(err).ShouldNot(HaveOccurred())

					pool, err := subnets.New(network)
					Ω(err).ShouldNot(HaveOccurred())

					out := make(chan *net.IPNet)
					go func(out chan *net.IPNet) {
						defer GinkgoRecover()
						n1, err := pool.AllocateDynamically()
						Ω(err).ShouldNot(HaveOccurred())
						out <- n1
					}(out)

					go func(out chan *net.IPNet) {
						defer GinkgoRecover()
						n1, err := pool.AllocateDynamically()
						Ω(err).ShouldNot(HaveOccurred())
						out <- n1
					}(out)

					a := <-out
					b := <-out
					return a.IP.Equal(b.IP)
				}, "100ms", "2ms").ShouldNot(BeTrue())
			})

		})

		Context("when the pool has sufficient IPs to allocate two subnets", func() {
			BeforeEach(func() {
				var err error
				_, defaultSubnetPool, err = net.ParseCIDR("10.2.3.0/29")
				Ω(err).ShouldNot(HaveOccurred())

			})

			Context("the second request", func() {
				It("succeeds", func() {
					_, err := manager.AllocateDynamically()
					Ω(err).ShouldNot(HaveOccurred())

					_, err = manager.AllocateDynamically()
					Ω(err).ShouldNot(HaveOccurred())
				})

				It("returns the second /30 network within the subnet", func() {
					_, err := manager.AllocateDynamically()
					Ω(err).ShouldNot(HaveOccurred())

					network, err := manager.AllocateDynamically()
					Ω(err).ShouldNot(HaveOccurred())

					Ω(network).ShouldNot(BeNil())
					Ω(network.String()).Should(Equal("10.2.3.4/30"))
				})
			})
		})
	})

})
