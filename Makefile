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

run-dev:  ## Run the dev server
	@echo "Running dev server. It will refresh automatically when you change code."
	@docker compose -f compose-dev.yml up --remove-orphans


.PHONY: stop
stop:  ## Stop the dev server
	@docker-compose -f compose-dev.yml down -v


.PHONY: clean
clean: ## Clean up the dev server
	$(MAKE) stop
	@docker compose -f compose-dev.yml rm --force
	@docker rmi archivista-archivista --force


.PHONY: test
test: ## Run tests
	@go test ./... -covermode atomic -coverprofile=cover.out -v

.PHONY: coverage
coverage:  ## Show html coverage
	@go tool cover -html=cover.out


.PHONY: lint
lint:  ## Run linter
	@golangci-lint run
	@go fmt ./...
	@go vet ./...


.PHONY: docs
docs:  ## Generate swagger docs
	@go install github.com/swaggo/swag/cmd/swag@v1.16.2
	@swag init -o docs -d internal/server -g server.go -pd

.PHONY: db-migrations
db-migrations:  ## Run the migrations for the database
	@atlas migrate diff mysql --dir "file://ent/migrate/migrations/mysql" --to "ent://ent/schema" --dev-url "docker://mysql/8/dev"
	@atlas migrate diff pgsql --dir "file://ent/migrate/migrations/pgsql" --to "ent://ent/schema" --dev-url "docker://postgres/16/dev?search_path=public"


help:  ## Show this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
