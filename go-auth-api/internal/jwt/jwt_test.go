package jwt_test

import (
	. "go-auth-api/internal/jwt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("jwt", func() {
	It("should successfully get a token and parse it", func() {
		token := GetToken("admin", "someorg", []byte(`{"someval": true}`))
		Expect(token).NotTo(HaveLen(0))

		user, org, i, err := IsTokenValid(token)
		Expect(err).To(BeNil())
		Expect(user).To(Equal("admin"))
		Expect(org).To(Equal("someorg"))
		Expect(string(i)).To(Equal(`{"someval": true}`))
	})
})
