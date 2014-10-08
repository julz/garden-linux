package ip_pool_test

import (
	"net"
	"runtime"

	"github.com/cloudfoundry-incubator/garden-linux/net_fence/ip_pool"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IpPool", func() {
	It("returns an error if maximum IP is less than minimum IP", func() {
		min := net.ParseIP("1.1.1.1")
		max := net.ParseIP("1.1.1.0")
		_, err := ip_pool.New(min, max)
		Ω(err).Should(Equal(ip_pool.ErrInvalidRange))
	})

	Describe(".Allocate", func() {
		Context("when there is more than one IP address in the pool", func() {
			var (
				min  net.IP
				max  net.IP
				pool ip_pool.IPPool
				err  error
			)

			BeforeEach(func() {
				min = net.ParseIP("1:1:1:1:1:1:1:1")
				max = net.ParseIP("1:1:1:1:1:1:1:2")
				pool, err = ip_pool.New(min, max)
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("can allocate two IP addresses", func() {
				ip, err := pool.Allocate()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(ip).ShouldNot(BeNil())

				ip, err = pool.Allocate()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(ip).ShouldNot(BeNil())
			})

			It("allocates two distinct IP addresses", func() {
				ip1, err := pool.Allocate()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(ip1).ShouldNot(BeNil())

				ip2, err := pool.Allocate()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(ip2).ShouldNot(Equal(ip1))
			})

			It("allocates distinct IP addresses concurrently", func() {
				prev := runtime.GOMAXPROCS(2)
				defer runtime.GOMAXPROCS(prev)

				Consistently(func() bool {
					pool, err = ip_pool.New(min, max)
					out := make(chan net.IP)
					go func(out chan net.IP) {
						defer GinkgoRecover()
						ip1, err := pool.Allocate()
						Ω(err).ShouldNot(HaveOccurred())
						out <- ip1
					}(out)

					go func(out chan net.IP) {
						defer GinkgoRecover()
						ip2, err := pool.Allocate()
						Ω(err).ShouldNot(HaveOccurred())
						out <- ip2
					}(out)

					a := <-out
					b := <-out
					return a.Equal(b)
				}, "100ms", "2ms").ShouldNot(BeTrue())
			})

			It("allocates all the IP addresses in the pool", func() {
				ip1, err := pool.Allocate()
				Ω(err).ShouldNot(HaveOccurred())

				ip2, err := pool.Allocate()
				Ω(err).ShouldNot(HaveOccurred())

				Ω([]net.IP{ip1, ip2}).Should(ContainElement(min))
				Ω([]net.IP{ip1, ip2}).Should(ContainElement(max))
			})

			It("returns an error when attempting to allocate three IP addresses", func() {
				pool.Allocate()
				pool.Allocate()
				_, err = pool.Allocate()
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("when the pool spans multiple bytes", func() {
			var (
				min  net.IP
				max  net.IP
				pool ip_pool.IPPool
				err  error
			)

			BeforeEach(func() {
				min = net.ParseIP("1.1.1.1")
				max = net.ParseIP("1.1.2.5")
				pool, err = ip_pool.New(min, max)
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("allocates as many IP address as there are in the range before returning an error", func() {
				var count int = -1
				for ; err == nil; _, err = pool.Allocate() {
					count++
				}
				Ω(count).Should(BeNumerically("==", 255+6))
			})

			It("never allocates the same IP address twice", func() {
				var ip net.IP
				seen := make(map[string]bool, 255+6)
				for ; err == nil; ip, err = pool.Allocate() {
					Ω(seen).NotTo(HaveKey(ip.String()))
					seen[ip.String()] = true
				}
			})
		})

		Context("when the pool ends on the maximal IP address", func() {
			var (
				max  net.IP
				pool ip_pool.IPPool
				err  error
			)

			BeforeEach(func() {
				max = net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")
				pool, err = ip_pool.New(max, max)
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("allocates up to the maximum", func() {
				Ω(func() { pool.Allocate() }).ShouldNot(Panic())
			})
		})

		Describe(".Release", func() {
			Context("when there is one IP address in the pool", func() {
				var (
					max  net.IP
					pool ip_pool.IPPool
					err  error
				)

				BeforeEach(func() {
					max = net.ParseIP("1:1:1:1:1:1:1:1")
					pool, err = ip_pool.New(max, max)
					Ω(err).ShouldNot(HaveOccurred())
				})

				It("allows the IP address to be allocated again", func() {
					ip, err := pool.Allocate()
					Ω(err).ShouldNot(HaveOccurred())

					err = pool.Release(ip)
					Ω(err).ShouldNot(HaveOccurred())

					ip2, err := pool.Allocate()
					Ω(err).ShouldNot(HaveOccurred())
					Ω(ip2).Should(Equal(ip))
				})
			})

			Context("when there are multiple IPs in the pool", func() {
				var (
					min  net.IP
					max  net.IP
					pool ip_pool.IPPool
					err  error
				)

				BeforeEach(func() {
					min = net.ParseIP("1:1:1:1:1:1:1:1")
					max = net.ParseIP("1:1:1:1:1:1:1:2")
					pool, err = ip_pool.New(min, max)
					Ω(err).ShouldNot(HaveOccurred())
				})

				It("releases the IP back to the pool", func() {
					ip, err := pool.Allocate()
					Ω(err).ShouldNot(HaveOccurred())

					err = pool.Release(ip)
					Ω(err).ShouldNot(HaveOccurred())

					_, err = pool.Allocate()
					Ω(err).ShouldNot(HaveOccurred())
				})

				It("allows the second IP address to be allocated again", func() {
					_, err := pool.Allocate()
					Ω(err).ShouldNot(HaveOccurred())

					ip, err := pool.Allocate()
					Ω(err).ShouldNot(HaveOccurred())

					err = pool.Release(ip)
					Ω(err).ShouldNot(HaveOccurred())

					ip2, err := pool.Allocate()
					Ω(err).ShouldNot(HaveOccurred())
					Ω(ip2).Should(Equal(ip))
				})

				It("returns an error if release is called on an unallocated IP", func() {
					err := pool.Release(min)
					Ω(err).Should(HaveOccurred())
				})
			})
		})
	})
})
