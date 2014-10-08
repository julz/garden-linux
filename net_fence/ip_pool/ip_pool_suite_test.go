package ip_pool_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIpPool(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IpPool Suite")
}
