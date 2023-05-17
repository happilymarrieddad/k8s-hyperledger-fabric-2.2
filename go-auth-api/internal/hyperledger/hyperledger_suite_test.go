package hyperledger_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHyperledger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hyperledger Suite")
}
