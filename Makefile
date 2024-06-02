# Copyright 2023 The Archivista Contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Nothing to be run by default as the CodeQL autobuild tries to run make
# See https://docs.github.com/en/code-security/code-scanning/creating-an-advanced-setup-for-code-scanning/codeql-code-scanning-for-compiled-languages#autobuild-for-go
all: help

# Run the dev server
run-dev:
	@echo "Running dev server. It will refresh automatically when you change code."
	@docker compose -f compose-dev.yml up --remove-orphans

.PHONY: stop
# Stop the dev server
stop:
	@docker-compose -f compose-dev.yml down -v

.PHONY: clean
# Clean up the dev server
clean:
	$(MAKE) stop
	@docker compose -f compose-dev.yml rm --force
	@docker rmi archivista-archivista --force

.PHONY: test
# Run tests
test:
	@go test ./... -covermode atomic -coverprofile=cover.out -v

.PHONY: coverage
# Show html coverage
coverage:
	@go tool cover -html=cover.out

.PHONY: lint
# Run linter
lint:
	@golangci-lint run
	@go fmt ./...
	@go vet ./...

.PHONY: docs
# Generate swagger docs
docs: check_docs

.PHONY: check_docs
check_docs:
	@server_mod_time=$$(stat -c %Y internal/server/server.go); \
	swagger_mod_time=$$(stat -c %Y docs/swagger.json); \
	if [ $$server_mod_time -gt $$swagger_mod_time ]; then \
		echo "Swagger documentation needs to be updated"; \
		make update_docs; \
	else \
		echo "Swagger documentation is up to date"; \
	fi

.PHONY: update_docs
update_docs:
	@go install github.com/swaggo/swag/cmd/swag@v1.16.2
	@swag init -o docs -d internal/server -g server.go -pd

.PHONY: db-migrations
# Run the migrations for the database
db-migrations:
	@go generate ./...
	@atlas migrate diff mysql --dir "file://ent/migrate/migrations/mysql" --to "ent://ent/schema" --dev-url "docker://mysql/8/dev"
	@atlas migrate diff pgsql --dir "file://ent/migrate/migrations/pgsql" --to "ent://ent/schema" --dev-url "docker://postgres/16/dev?search_path=public"

# Show this help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
