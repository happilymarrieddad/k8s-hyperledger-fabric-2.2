package k8client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("k8sclient namespace", func() {

	var (
		nmspName string
	)

	BeforeEach(func() {
		nmspName = "some-value"
	})

	Context("Create", func() {
		It("should successfully create a namespace", func() {
			nmsp, err := client.CreateNamespace("some-value")
			Expect(err).To(BeNil())
			Expect(nmsp.GetName()).To(Equal("some-value"))
		})
	})

	Context("Delete", func() {
		It("should successfully create a namespace", func() {
			Expect(client.DeleteNamespace(nmspName)).To(Succeed())
		})
	})

})
