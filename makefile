
default: help

.PHONY: help
help:
	@echo 'local-storage'
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
	go test ./... -coverprofile coverage.out
	go tool cover -func coverage.out | grep ^total:

.PHONY: test-coverage
test-coverage:
	ginkgo -r -v -p -race --progress --randomize-all --randomize-suites -cover -coverprofile="coverage.out"

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
	mockgen -package=mock_grpcv1 -source generated/proto/api/grpc/v1/file_grpc.pb.go -destination=generated/proto/api/grpc/v1/mock/file_grpc_mock.go
	mockgen -package=mock_auth -source internal/auth/basic.go -destination=internal/auth/mock/basic_mock.go
	mockgen -package=mock_config -source internal/config/config.go -destination=internal/config/mock/config_mock.go
	mockgen -package=mock_context -source internal/context/context.go -destination=internal/context/mock/context_mock.go
	mockgen -package=mock_datetime -source internal/datetime/clock.go -destination=internal/datetime/mock/clock_mock.go
	mockgen -package=mock_dbmongo -source internal/db-mongo/client.go -destination=internal/db-mongo/mock/client_mock.go
	mockgen -package=mock_dbmysql -source internal/db-mysql/client.go -destination=internal/db-mysql/mock/client_mock.go
	mockgen -package=mock_encoding -source internal/encoding/encoder.go -destination=internal/encoding/mock/encoder_mock.go
	mockgen -package=mock_file -source internal/file/file.go -destination=internal/file/mock/file_mock.go
	mockgen -package=mock_file -source internal/file/location.go -destination=internal/file/mock/location_mock.go
	mockgen -package=mock_filesystem -source internal/filesystem/file.go -destination=internal/filesystem/mock/file_mock.go
	mockgen -package=mock_filesystem -source internal/filesystem/directory.go -destination=internal/filesystem/mock/directory_mock.go
	mockgen -package=mock_grpcapp -source internal/grpc-app/server.go -destination=internal/grpc-app/mock/server_mock.go
	mockgen -package=mock_grpclog -source internal/grpc-log/stream.go -destination=internal/grpc-log/mock/stream_mock.go
	mockgen -package=mock_hashing -source internal/hashing/hasher.go -destination=internal/hashing/mock/hasher_mock.go
	mockgen -package=mock_healthcheck -source internal/healthcheck/checker.go -destination=internal/healthcheck/mock/checker_mock.go
	mockgen -package=mock_healthcheck -source internal/healthcheck/go_health.go -destination=internal/healthcheck/mock/go_health_mock.go
	mockgen -package=mock_healthcheck -source internal/healthcheck/health.go -destination=internal/healthcheck/mock/health_mock.go
	mockgen -package=mock_io -source internal/io/io.go -destination=internal/io/mock/io_mock.go
	mockgen -package=mock_logging -source internal/logging/log.go -destination=internal/logging/mock/log_mock.go
	mockgen -package=mock_repository -source internal/repository/provider.go -destination=internal/repository/mock/provider_mock.go
	mockgen -package=mock_repository -source internal/repository/file.go -destination=internal/repository/mock/file_mock.go
	mockgen -package=mock_repository -source internal/repository/auth.go -destination=internal/repository/mock/auth_mock.go
	mockgen -package=mock_restapp -source internal/rest-app/server.go -destination=internal/rest-app/mock/server_mock.go
	mockgen -package=mock_retrieving -source internal/retrieving/retriever.go -destination=internal/retrieving/mock/retriever_mock.go
	mockgen -package=mock_serialization -source internal/serialization/serializer.go -destination=internal/serialization/mock/serializer_mock.go
	mockgen -package=mock_text -source internal/text/id.go -destination=internal/text/mock/id_mock.go

.PHONY: run-grpc-app
run-grpc-app:
	go run cmd/grpc-app/main.go

.PHONY: run-rest-app
run-rest-app:
	go run cmd/rest-app/main.go

.PHONY: run-hybrid-app
run-hybrid-app:
	go run cmd/hybrid-app/main.go

.PHONY: build-grpc-app
build-grpc-app:
	go build -o ./build/grpc-app/ ./cmd/grpc-app/main.go

.PHONY: build-rest-app
build-rest-app:
	go build -o ./build/rest-app/ ./cmd/rest-app/main.go

.PHONY: build-hybrid-app
build-hybrid-app:
	go build -o ./build/hybrid-app/ ./cmd/hybrid-app/main.go

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


MYSQL_DB_URI=mysql://admin:123456@tcp(localhost:3308)/goseidon_local?x-tls-insecure-skip-verify=true
MONGO_DB_URI=mongodb://admin:123456@localhost:27030/goseidon_local

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
