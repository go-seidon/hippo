
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
	mockgen -package=mock_app -source internal/app/server.go -destination=internal/app/mock/server_mock.go
	mockgen -package=mock_auth -source internal/auth/basic.go -destination=internal/auth/mock/basic_mock.go
	mockgen -package=mock_config -source internal/config/config.go -destination=internal/config/mock/config_mock.go
	mockgen -package=mock_datetime -source internal/datetime/clock.go -destination=internal/datetime/mock/clock_mock.go
	mockgen -package=mock_dbmongo -source internal/db-mongo/client.go -destination=internal/db-mongo/mock/client_mock.go
	mockgen -package=mock -source internal/text/id.go -destination=internal/mock/text_id_mock.go
	mockgen -package=mock -source internal/logging/log.go -destination=internal/mock/logging_log_mock.go
	mockgen -package=mock -source internal/serialization/serializer.go -destination=internal/mock/serialization_serializer_mock.go
	mockgen -package=mock -source internal/encoding/encoder.go -destination=internal/mock/encoding_encoder_mock.go
	mockgen -package=mock -source internal/hashing/hasher.go -destination=internal/mock/hashing_hasher_mock.go
	mockgen -package=mock -source internal/filesystem/file.go -destination=internal/mock/filesystem_file_mock.go
	mockgen -package=mock -source internal/filesystem/directory.go -destination=internal/mock/filesystem_directory_mock.go
	mockgen -package=mock -source internal/repository/file.go -destination=internal/mock/repository_file_mock.go
	mockgen -package=mock -source internal/repository/auth.go -destination=internal/mock/repository_auth_mock.go
	mockgen -package=mock -source internal/healthcheck/health.go -destination=internal/mock/healthcheck_health_mock.go
	mockgen -package=mock -source internal/healthcheck/go_health.go -destination=internal/mock/healthcheck_go_health_mock.go
	mockgen -package=mock -source internal/deleting/deleter.go -destination=internal/mock/deleting_deleter_mock.go
	mockgen -package=mock -source internal/retrieving/retriever.go -destination=internal/mock/retrieving_retriever_mock.go
	mockgen -package=mock -source internal/uploading/uploader.go -destination=internal/mock/uploading_uploader_mock.go
	mockgen -package=mock -source internal/uploading/location.go -destination=internal/mock/uploading_location_mock.go
	mockgen -package=mock -source internal/repository/provider.go -destination=internal/repository/mock/provider_mock.go

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
