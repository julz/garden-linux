package netfence_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNetFence(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "net-fence Suite")
}
