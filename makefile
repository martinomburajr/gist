PROJECT_NAME := "gist"
PKG := "github.com/martinomburajr/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/...)
GO_FILES := $(shell find . -name '.*go')

.PHONY = all dep build clean test coverage coverhtml lint

all: build

build: ## Runs the build command that creates a OS specific binary named gist
	@go build -i -v $(PKG)

lint: ## Runs a linter on the code to ensure it meets standards
	@golint -set_exit_status ${PKG_LIST}

test: ## Runs unit tests
	@go test -short ${PKG_LIST}

race: dep ## Run the race detector
	@go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage roprt
	./tools/coverage.sh

coverhtml: ## Generate global code coverage report HTML
	./tools/coverage.sh html

dep: ## Get dependencies
	@go get -v -d ./...

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'