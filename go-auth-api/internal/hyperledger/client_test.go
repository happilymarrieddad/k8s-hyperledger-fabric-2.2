package hyperledger_test

import (
	. "go-auth-api/internal/hyperledger"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("hyperledger", func() {
	var client Client

	BeforeEach(func() {
		var err error
		client, err = NewClient("/home/nick/Projects/k8s-hyperledger-fabric-2.2/go-auth-api/internal/hyperledger/local/config.yaml")
		Expect(err).To(BeNil())
		Expect(client).NotTo(BeNil())
	})

	Context("FindOrCreateAffilation", func() {
		It("should ensure an affiliation exists", func() {
			res, err := client.FindOrCreateAffilation("ibm", "garbage")
			Expect(err).To(BeNil())
			Expect(res).NotTo(BeNil())
		})
	})
})
