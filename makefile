
default: help

.PHONY: help
help:
	@echo 'hippo'
	@echo 'usage: make [target] ...'

.PHONY: install-tool
install-tool:
	go get -u github.com/golang/mock/gomock
	go get -u github.com/golang/mock/mockgen

.PHONY: install-dependency
install-dependency:
	go mod tidy
	go mod verify
	go mod vendor

.PHONY: clean-dependency
clean-dependency:
	rm -f go.sum
	rm -rf vendor
	go clean -modcache

.PHONY: install
install:
	go install -v ./...

.PHONY: test
test:
	go test ./.../ -p 1 -race -coverprofile coverage.out
	go tool cover -func coverage.out | grep ^total:

.PHONY: test-coverage
test-coverage:
	ginkgo -r -v -p -race --progress --randomize-all --randomize-suites -cover -coverprofile="coverage.out"

.PHONY: test-watch
test-watch:
	ginkgo watch -r -v -p -race --trace

.PHONY: test-unit
test-unit:
	ginkgo -r -v -p -race --label-filter="unit" -cover -coverprofile="coverage.out"

.PHONY: test-integration
test-integration:
	ginkgo -r -v -p -race --label-filter="integration" -cover -coverprofile="coverage.out"

.PHONY: test-watch-unit
test-watch-unit:
	ginkgo watch -r -v -p -race --trace --label-filter="unit"

.PHONY: test-watch-integration
test-watch-integration:
	ginkgo watch -r -v -p -race --trace --label-filter="integration"

.PHONY: generate-mock
generate-mock:
	mockgen -package=mock_grpcapp -source api/grpcapp/file_grpc.pb.go -destination=api/grpcapp/mock/file_grpc_mock.go
	mockgen -package=mock_auth -source internal/auth/basic.go -destination=internal/auth/mock/basic_mock.go
	mockgen -package=mock_file -source internal/file/file.go -destination=internal/file/mock/file_mock.go
	mockgen -package=mock_file -source internal/file/location.go -destination=internal/file/mock/location_mock.go
	mockgen -package=mock_filesystem -source internal/filesystem/file.go -destination=internal/filesystem/mock/file_mock.go
	mockgen -package=mock_filesystem -source internal/filesystem/directory.go -destination=internal/filesystem/mock/directory_mock.go
	mockgen -package=mock_grpcapp -source internal/grpcapp/server.go -destination=internal/grpcapp/mock/server_mock.go
	mockgen -package=mock_healthcheck -source internal/healthcheck/health.go -destination=internal/healthcheck/mock/health_mock.go
	mockgen -package=mock_repository -source internal/repository/repository.go -destination=internal/repository/mock/repository_mock.go
	mockgen -package=mock_repository -source internal/repository/file.go -destination=internal/repository/mock/file_mock.go
	mockgen -package=mock_repository -source internal/repository/auth.go -destination=internal/repository/mock/auth_mock.go
	mockgen -package=mock_restapp -source internal/restapp/server.go -destination=internal/restapp/mock/server_mock.go

.PHONY: generate-proto
generate-proto:
	buf generate api/

.PHONY: verify-swagger
verify-swagger:
	swagger-cli bundle api/restapp/main.yml --type json > api/restapp/main.all.json
	swagger-cli validate api/restapp/main.all.json 

.PHONY: generate-swagger
generate-swagger:
	swagger-cli bundle api/restapp/main.yml --type yaml > api/restapp/main.all.yml

.PHONY: generate-oapi
generate-oapi:
	make generate-swagger
	make generate-oapi-type

.PHONY: generate-oapi-type
generate-oapi-type:
	oapi-codegen -old-config-style -config api/restapp/type.gen.yaml api/restapp/main.all.yml

.PHONY: run-grpcapp
run-grpcapp:
	go run cmd/grpcapp/main.go

.PHONY: run-restapp
run-restapp:
	go run cmd/restapp/main.go

.PHONY: run-hybridapp
run-hybridapp:
	go run cmd/hybridapp/main.go

.PHONY: build-grpcapp
build-grpcapp:
	go build -o ./build/grpcapp/ ./cmd/grpcapp/main.go

.PHONY: build-restapp
build-restapp:
	go build -o ./build/restapp/ ./cmd/restapp/main.go

.PHONY: build-hybridapp
build-hybridapp:
	go build -o ./build/hybridapp/ ./cmd/hybridapp/main.go

ifeq (migrate-mysql,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "migrate-mysql"
  MIGRATE_MYSQL_RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(MIGRATE_MYSQL_RUN_ARGS):dummy;@:)
endif

ifeq (migrate-mysql-create,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "migrate-mysql-create"
  MIGRATE_MYSQL_RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(MIGRATE_MYSQL_RUN_ARGS):dummy;@:)
endif

ifeq (migrate-mongo,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "migrate-mongo"
  MIGRATE_MONGO_RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(MIGRATE_MONGO_RUN_ARGS):dummy;@:)
endif

ifeq (migrate-mongo-create,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "migrate-mongo-create"
  MIGRATE_MONGO_RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(MIGRATE_MONGO_RUN_ARGS):dummy;@:)
endif

dummy: ## used by migrate script as do-nothing targets
	@:


MYSQL_DB_URI=mysql://admin:123456@tcp(localhost:3308)/hippo?x-tls-insecure-skip-verify=true
MONGO_DB_URI=mongodb://admin:123456@localhost:27030/hippo

.PHONY: migrate-mysql
migrate-mysql:
	migrate -database "$(MYSQL_DB_URI)" -path ./migration/mysql $(MIGRATE_MYSQL_RUN_ARGS)

.PHONY: migrate-mysql-create
migrate-mysql-create:
	migrate create -dir migration/mysql -ext .sql $(MIGRATE_MYSQL_RUN_ARGS)

.PHONY: migrate-mongo
migrate-mongo:
	migrate -database "$(MONGO_DB_URI)" -path ./migration/mongo $(MIGRATE_MONGO_RUN_ARGS)

.PHONY: migrate-mongo-create
migrate-mongo-create:
	migrate create -dir migration/mongo -ext .json $(MIGRATE_MONGO_RUN_ARGS)
