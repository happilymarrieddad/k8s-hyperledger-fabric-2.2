package models_test

import (
	. "admin-api/internal/models"
	"admin-api/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Organizations", func() {

	var (
		globalModel GlobalModels
		orgsModel   Organizations

		org *types.Organization
	)

	BeforeEach(func() {
		Expect(ClearDatabase(db)).To(Succeed())
		Expect(DestroyKubernetesResource(k8c)).To(Succeed())

		globalModel = NewGlobalModels(db, k8c)
		Expect(globalModel).NotTo(BeNil())

		orgsModel = globalModel.Organizations()
		Expect(orgsModel).NotTo(BeNil())

		org = &types.Organization{
			Name:        "Some Org",
			NetworkName: "someorg",
		}
	})

	AfterEach(func() {
		Expect(ClearDatabase(db)).To(Succeed())
		Expect(DestroyKubernetesResource(k8c)).To(Succeed())
	})

	Context("Create", func() {
		It("should successfully create an org and a namespace", func() {
			Expect(orgsModel.Create(org)).To(Succeed())

			nmsp, err := k8c.GetNamespace(org.NetworkName)
			Expect(err).To(BeNil())
			Expect(nmsp.Name).To(Equal(org.NetworkName))
		})
	})

})
