package k8client_test

import (
	"os"
	"strings"
	"testing"

	. "admin-api/internal/k8client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	client Client
)

var _ = BeforeSuite(func() {
	var isInCluster bool
	if val := os.Getenv("ADMIN_API_IN_CLUSTER"); len(val) > 0 && strings.ToLower(val) == "true" {
		isInCluster = true
	}

	var err error
	client, err = NewClient(&Config{
		IsIncluster: isInCluster,
		DebugLogs:   true,
	})
	Expect(err).To(BeNil())
	Expect(client).NotTo(BeNil())
})

func TestK8client(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8client Suite")
}
