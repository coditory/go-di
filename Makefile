SHELL := /usr/bin/env bash -o errexit -o pipefail -o nounset
MAKEFLAGS := --warn-undefined-variables --no-builtin-rules

BINARY_NAME = $(shell sed -nE "s| *module *([^ ]*).*|\1|p" go.mod)
OUT_DIR = out
BIN_DIR = $(OUT_DIR)/bin
REPORT_DIR = $(OUT_DIR)/report
VERSION = $(shell git tag --list --sort=-version:refname "v*" | head -n 1 | grep "." || echo "v0.0.0")
SERVICE_PORT ?= 8080# used in watch mode
DOCKER_REGISTRY ?= # finish with /
EXPORT_RESULT ?= false # on CI, set EXPORT_RESULT = true

COLORS ?= true
RED    := $(if $(findstring $(COLORS),true),$(shell tput -Txterm setaf 1))
GREEN  := $(if $(findstring $(COLORS),true),$(shell tput -Txterm setaf 2))
YELLOW := $(if $(findstring $(COLORS),true),$(shell tput -Txterm setaf 3))
WHITE  := $(if $(findstring $(COLORS),true),$(shell tput -Txterm setaf 7))
CYAN   := $(if $(findstring $(COLORS),true),$(shell tput -Txterm setaf 6))
RESET  := $(if $(findstring $(COLORS),true),$(shell tput -Txterm sgr0))

# Tools
GOFUMPT_CMD = go run mvdan.cc/gofumpt@latest
GOJUNITREP_CMD = go run github.com/jstemmer/go-junit-report/v2@v2.0.0
GOCOVXML_CMD = go run github.com/AlekSi/gocov-xml@v1.1.0
GOCOV_CMD = go run github.com/axw/gocov/gocov@v1.1.0
YAMLLINT_CHECKSTYLE_CMD = go run github.com/thomaspoignant/yamllint-checkstyle@v1.0.2
# Dockerized tools
COSMTREK_AIR_VERSION = latest
GOLANGCI_LINT_VERSION = v1.51.1

define task
 @echo "${CYAN}>>> $(1)${RESET}"
endef

.PHONY: all
all: clean lint test build

## Build:
.PHONY: build
build: ## Build project (deafult: ./out/bin/$binary)
	$(call task,build)
	@rm -rf $(BIN_DIR)
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) .

.PHONY: ci
ci: clean lint coverage build ## Build project for CI (create lint and test reports)

.PHONY: clean
clean: ## Remove build related files (default ./out)
	$(call task,clean)
	@rm -fr $(OUT_DIR)

.PHONY: run
run: ## Run the project
	$(call task,run)
	@go run .

.PHONY: watch
watch: ## Run the project in watch mode
	$(call task,watch)
	$(eval PACKAGE_NAME=$(shell head -n 1 go.mod | cut -d ' ' -f2))
	@mkdir -p $(OUT_DIR)/watch
	@docker run -it --rm -w /go/src/$(PACKAGE_NAME) -v $(shell pwd):/go/src/$(PACKAGE_NAME) \
		-p $(SERVICE_PORT):$(SERVICE_PORT) cosmtrek/air:$(COSMTREK_AIR_VERSION) \
		--tmp_dir=$(OUT_DIR)/watch \
		--build.cmd="go build -o $(OUT_DIR)/watch/main ." \
		--build.bin="$(OUT_DIR)/watch/main"


.PHONY: format
format: ## Format go source files
	$(call task,format)
	@$(GOFUMPT_CMD) -l -w .
	@go mod tidy -e

## Test:
.PHONY: test
test: ## Run tests
	$(call task,test)
	@rm -rf $(REPORT_DIR)/tests
	@mkdir -p $(REPORT_DIR)/tests
	@go test -v -race ./... \
		| tee >($(GOJUNITREP_CMD) -set-exit-code > $(REPORT_DIR)/tests/junit-report.xml)

.PHONY: coverage
coverage: ## Run tests and create coverage report
	$(call task,coverage)
	@rm -rf $(REPORT_DIR)/tests
	@mkdir -p $(REPORT_DIR)/tests
	@go test -cover -covermode=count -coverprofile=$(REPORT_DIR)/tests/profile.cov -v ./... \
		| tee >($(GOJUNITREP_CMD) -set-exit-code > $(REPORT_DIR)/tests/junit-report.xml)
	@go tool cover -func $(REPORT_DIR)/tests/profile.cov
	@go tool cover -html=$(REPORT_DIR)/tests/profile.cov -o $(REPORT_DIR)/tests/coverage.html
	@$(GOCOV_CMD) convert $(REPORT_DIR)/tests/profile.cov \
		| $(GOCOVXML_CMD) > $(REPORT_DIR)/tests/coverage.xml

## Lint:
.PHONY: lint
lint: ## Lint go source files
	$(call task,lint-go)
	@rm -f $(REPORT_DIR)/checktyle/format-go-*
	@rm -f $(REPORT_DIR)/checktyle/checkstyle-go.*
	@mkdir -p $(REPORT_DIR)/checkstyle
ifneq ($(shell $(GOFUMPT_CMD) -l . | wc -l),0)
	@echo "${YELLOW}Detected unformatted code${RESET} (fix: make format)"
	@$(GOFUMPT_CMD) -l . | tee $(REPORT_DIR)/checkstyle/format-go-files.txt
	@$(GOFUMPT_CMD) -d . | tee $(REPORT_DIR)/checkstyle/format-go-diff.txt
	@exit 1
endif
	@docker run --rm -t -v $(shell pwd):/app -v ~/.cache/golangci-lint/v1.51.1:/root/.cache -w /app \
		golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) \
		golangci-lint run \
		--deadline=65s \
		--out-format checkstyle:$(REPORT_DIR)/checkstyle/checkstyle-go.xml,colored-line-number \
		./...

## Docker:
.PHONY: docker-build
docker-build: ## Build with Dockerfile
	$(call task,docker-build)
	docker build . \
		--tag $(BINARY_NAME):latest \
		--tag $(BINARY_NAME):latest-default \
		--tag $(BINARY_NAME):$(VERSION)

.PHONY: docker-build-labeled
docker-build-labeled: ## Build with Dockerfile.labeled
	$(call task,docker-build-labeled)
	docker build . -f Dockerfile.labeled \
		--build-arg VERSION="$(VERSION)" \
		--build-arg VCS_REF="$(shell git rev-parse --short HEAD)" \
		--build-arg BUILD_DATE="$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')" \
		--tag $(BINARY_NAME):latest \
		--tag $(BINARY_NAME):latest-min \
		--tag $(BINARY_NAME):$(VERSION)

.PHONY: docker-build-ci
docker-build-ci: ## Build with Dockerfile.ci
	$(call task,docker-build-ci)
	docker build . -f Dockerfile.ci \
		--tag $(BINARY_NAME):latest \
		--tag $(BINARY_NAME):latest-ci \
		--tag $(BINARY_NAME):$(VERSION)


.PHONY: docker-build-base
docker-build-base: ## Build with Dockerfile.base
	$(call task,docker-build-base)
	docker build . -f Dockerfile.base \
		--tag $(BINARY_NAME):latest \
		--tag $(BINARY_NAME):latest-base \
		--tag $(BINARY_NAME):$(VERSION)

.PHONY: docker-run
docker-run: ## Run docker container of latest image
	$(call task,docker-run)
	docker run -p 8080:8080 $(BINARY_NAME):latest

## Help:
.PHONY: help
help: ## Show this help
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} - run whole build process (clean, lint, test, build)'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET} - run single target'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

.PHONY: version
version: ## Print project version
	@echo '${VERSION}'
