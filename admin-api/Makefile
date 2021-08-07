BRANCHNAME:=$(shell git rev-parse --abbrev-ref HEAD)

# Testing
test-unit:
	$(eval PACKAGES=$(shell go list ./...))
	@go test ${PACKAGES} -v -ginkgo.failFast -test.parallel 10

test-integration:
	$(eval PACKAGES=$(shell go list ./...))
	@go test ${PACKAGES} -v -tags=integration

test-deep:
	$(eval PACKAGES=$(shell go list ./...))
	@go test ${PACKAGES} -v -tags=deep -ginkgo.focus="Deep Model Tests"

# USAGE - make name=AddSomeDBTable db.migrate.create
db.migrate.create:
	goose -dir db/migrations create $(name) sql

db.migrate.up:
	goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:5431/admin-api?connect_timeout=180&sslmode=disable" up
	goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:5431/admin-api-test?connect_timeout=180&sslmode=disable" up

db.migrate.down:
	goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:5431/admin-api?connect_timeout=180&sslmode=disable" down
	goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:5431/admin-api-test?connect_timeout=180&sslmode=disable" down

db.migrate.reset:
	goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:5431/admin-api?connect_timeout=180&sslmode=disable" reset
	goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:5431/admin-api-test?connect_timeout=180&sslmode=disable" reset
