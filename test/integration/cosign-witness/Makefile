.PHONY: help validate test clean

help: ## Show this help message
	@echo "Cosign-Witness-Archivista Integration Test"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

validate: ## Validate prerequisites are met
	./validate-setup.sh

test: ## Run the full integration test
	./run-test.sh

clean: ## Clean up test artifacts
	./clean.sh

test-clean: clean test ## Clean and run test

.DEFAULT_GOAL := help
