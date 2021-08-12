package models_test

import (
	"admin-api/internal/k8client"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"xorm.io/xorm"
)

var db *xorm.Engine
var k8c k8client.Client

func ClearDatabase(db *xorm.Engine) (err error) {
	results, err := db.Query(`
		SELECT table_name FROM information_schema.tables 
		WHERE table_schema='public'
		AND table_name NOT IN ('geography_columns','geometry_columns','spatial_ref_sys','goose_db_version','us_lex','us_gaz','us_rules')
	`)
	Expect(err).To(BeNil())
	for _, r := range results {
		tbn, ok := r["table_name"]
		if ok {
			if _, err = db.Query(fmt.Sprintf(`TRUNCATE %s CASCADE`, string(tbn))); err != nil {
				return err
			}
		}
	}

	return nil
}

func DestroyKubernetesNamespaces(k8c k8client.Client) error {
	list, err := k8c.GetNamespaces()
	if err != nil {
		return err
	}

	for _, n := range list {
		if err = k8c.DeleteNamespace(n.Name); err != nil {
			return err
		}
	}

	return nil
}

var _ = BeforeSuite(func() {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("ADMIN_API_DB_USER"),
		os.Getenv("ADMIN_API_DB_PASSWORD"),
		os.Getenv("ADMIN_API_DB_HOST"),
		os.Getenv("ADMIN_API_DB_PORT"),
		os.Getenv("ADMIN_API_DB_NAME"),
	)
	fmt.Printf("Using DB Conn: '%s' for testing\n", conn)
	var err error
	db, err = xorm.NewEngine("pgx", conn)
	Expect(err).To(BeNil())

	var isInCluster bool
	if val := os.Getenv("ADMIN_API_IN_CLUSTER"); len(val) > 0 && strings.ToLower(val) == "true" {
		isInCluster = true
	}

	k8c, err = k8client.NewClient(&k8client.Config{
		IsIncluster: isInCluster,
		DebugLogs:   true,
	})
	Expect(err).To(BeNil())
	Expect(k8c).NotTo(BeNil())
})

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}
