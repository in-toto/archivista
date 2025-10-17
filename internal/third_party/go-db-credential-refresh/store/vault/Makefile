MODULE=vault

include ./../../tools/tools.mk

# Special handling so we're not generating (non) coverage for our helpers or examples
test:
	go test github.com/davepgreene/go-db-credential-refresh/store/vault \
	github.com/davepgreene/go-db-credential-refresh/store/vault/auth \
	github.com/davepgreene/go-db-credential-refresh/store/vault/credentials \
	-count=1 -coverprofile=cover.out
