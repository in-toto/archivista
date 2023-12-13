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
	@docker build -t archivista:dev .
	@docker compose -f compose-dev.yml up --remove-orphans


.PHONY: stop
stop:  ## Stop the dev server
	@docker-compose down -v


.PHONY: clean
clean: ## Clean up the dev server
	$(MAKE) stop
	@docker compose rm --force
	@docker rmi archivista-archivista --force


.PHONY: test
test: ## Run tests
	@bash ./test/test.sh

help:  ## Show this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
