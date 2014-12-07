package iptables_test

import (
	"errors"
	"net"
	"os/exec"

	. "github.com/cloudfoundry-incubator/garden-linux/network/iptables"
	"github.com/cloudfoundry/gunk/command_runner/fake_command_runner"
	. "github.com/cloudfoundry/gunk/command_runner/fake_command_runner/matchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Iptables", func() {
	var fakeRunner *fake_command_runner.FakeCommandRunner
	var subject *Chain

	BeforeEach(func() {
		fakeRunner = fake_command_runner.New()
		subject = &Chain{"foo-bar-baz", fakeRunner}
	})

	Describe("NATRule", func() {
		Context("creating a rule", func() {
			It("runs iptables to create the rule with the correct parameters", func() {
				_, source, _ := net.ParseCIDR("1.3.5.0/28")
				subject.Create(&Rule{
					Source: source,
					Jump:   Return,
					To:     net.ParseIP("1.2.3.4"),
				})

				Ω(fakeRunner).Should(HaveExecutedSerially(fake_command_runner.CommandSpec{
					Path: "/sbin/iptables",
					Args: []string{"-w", "-t", "nat", "-A", "foo-bar-baz", "--source", "1.3.5.0/28", "--jump", "RETURN", "--to", "1.2.3.4"},
				}))
			})

			Context("when the command returns an error", func() {
				It("returns an error", func() {
					someError := errors.New("badly laid iptable")
					fakeRunner.WhenRunning(
						fake_command_runner.CommandSpec{Path: "/sbin/iptables"},
						func(cmd *exec.Cmd) error {
							return someError
						},
					)

					_, source, _ := net.ParseCIDR("1.3.5.0/28")
					Ω(subject.Create(&Rule{Source: source})).ShouldNot(Succeed())
				})
			})
		})

		Context("deleting a rule", func() {
			It("runs iptables to delete the rule with the correct parameters", func() {
				_, source, _ := net.ParseCIDR("1.3.5.0/28")
				subject.Destroy(&Rule{
					Source: source,
					Jump:   Return,
					To:     net.ParseIP("1.2.3.4"),
				})

				Ω(fakeRunner).Should(HaveExecutedSerially(fake_command_runner.CommandSpec{
					Path: "/sbin/iptables",
					Args: []string{"-w", "-t", "nat", "-D", "foo-bar-baz", "--source", "1.3.5.0/28", "--jump", "RETURN", "--to", "1.2.3.4"},
				}))
			})

			Context("when the command returns an error", func() {
				It("returns an error", func() {
					someError := errors.New("badly laid iptable")
					fakeRunner.WhenRunning(
						fake_command_runner.CommandSpec{Path: "/sbin/iptables"},
						func(cmd *exec.Cmd) error {
							return someError
						},
					)

					_, source, _ := net.ParseCIDR("1.3.5.0/28")
					Ω(subject.Destroy(&Rule{Source: source})).ShouldNot(Succeed())
				})
			})
		})
	})
})
