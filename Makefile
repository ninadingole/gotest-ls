NO_COLOR=\033[0m
OK_COLOR=\033[32;01m

BINARY = gotest-ls
BINDIR = bin

export GO111MODULE=on

GO ?= go
GOFLAGS ?= -v
GOLANGCI_LINT ?= golangci-lint-$(GOLANGCI_LINT_VERSION)

.PHONY: all
all: deps test build

.PHONY: deps
deps:
	$(GO) mod vendor -v

.PHONE: build
build:
	@echo "$(OK_COLOR)==> Building$(NO_COLOR)"
	@$(GO) build $(GOFLAGS) -o $(BINDIR)/$(BINARY) main.go

## run unit tests
test:
	@echo "$(OK_COLOR)==> Running tests$(NO_COLOR)"
	$(GO) test -covermode=atomic -coverprofile=coverage.out -race -shuffle=on ./...

.PHONY: lint
lint: bin/$(GOLANGCI_LINT)
	@echo "$(OK_COLOR)==> Running lint$(NO_COLOR)"
	@bin/$(GOLANGCI_LINT) run -c .golangci.yml

# Download and install golangci-lint
bin/$(GOLANGCI_LINT):
	@echo "$(OK_COLOR)==> Installing golangci-lint $(GOLANGCI_LINT_VERSION)$(NO_COLOR)"; \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin "$(GOLANGCI_LINT_VERSION)"
	@mv ./bin/golangci-lint bin/$(GOLANGCI_LINT)
