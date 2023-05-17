package testing_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s-hyperledger-fabric-2.2/go-api/testing"
)

var _ = Describe("testing client", func() {

	It("should successfully enroll", func() {
		Expect(testing.NewClient(
			"/home/nick/Projects/k8s-hyperledger-fabric-2.2/go-api/testing/config.yaml",
			"admin", "adminpw", "ibm", "mainchannel",
		)).To(Succeed())
	})

})
