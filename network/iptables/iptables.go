package iptables

import (
	"net"
	"os/exec"

	"github.com/cloudfoundry/gunk/command_runner"
)

type Chain struct {
	Name   string
	Runner command_runner.CommandRunner
}

type Rule struct {
	Source *net.IPNet
	To     net.IP
	Jump   Action
}

func (n *Rule) create(chain string, runner command_runner.CommandRunner) error {
	return runner.Run(exec.Command("/sbin/iptables", "-w", "-t", "nat", "-A", chain, "--source", n.Source.String(), "--jump", string(n.Jump), "--to", n.To.String()))
}

func (n *Rule) destroy(chain string, runner command_runner.CommandRunner) error {
	return runner.Run(exec.Command("/sbin/iptables", "-w", "-t", "nat", "-D", chain, "--source", n.Source.String(), "--jump", string(n.Jump), "--to", n.To.String()))
}

type Destroyable interface {
	Destroy() error
}

type Action string

const (
	Return    Action = "RETURN"
	SourceNAT        = "SNAT"
)

type creater interface {
	create(chain string, runner command_runner.CommandRunner) error
}

type destroyer interface {
	destroy(chain string, runner command_runner.CommandRunner) error
}

func (c *Chain) Create(rule creater) error {
	return rule.create(c.Name, c.Runner)
}

func (c *Chain) Destroy(rule destroyer) error {
	return rule.destroy(c.Name, c.Runner)
}
