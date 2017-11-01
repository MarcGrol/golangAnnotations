GO := $(shell which go)
GO_VERSION := $(shell $(GO) version)

GOLANG_ANNOT_ROOT := $(shell echo "${GOPATH}/src/github.com/MarcGrol/golangAnnotations")

all: gen test install

help:
	@echo "\tdeps: installs all dependencies"
	@echo "\tgen: generates boilerplate code"
	@echo "\ttest: Run all tests"

deps:
	@echo "---------------------------"
	@echo "Performing dependency check"
	@echo "---------------------------"
	go get -u golang.org/x/tools/cmd/goimports
	go get -u -t ./...                                  # get the application with all its deps

verify:
	@echo "----------------------------------"
	@echo "Run static analysis on source-code"
	@echo "----------------------------------"
	$(GO) vet ./...
	golint ./...

generate:
	@echo "----------------------"
	@echo "Generating source-code"
	@echo "----------------------"
	$(GO) generate ./...

imports:
	@echo "------------------"
	@echo "Optimizing imports"
	@echo "------------------"
	find . -name '*.go' -exec goimports -l -w -local github.com/ {} \;

format: imports
	@echo "----------------------"
	@echo "Formatting source-code"
	@echo "----------------------"
	find . -name '*.go' -exec gofmt -l -s -w {} \;

gen: generate imports format

test:
	@echo "---------------------"
	@echo "Running backend tests"
	@echo "---------------------"
	$(GO) test ./...                        # run unit tests
	make format

citest:
	@echo "---------------------"
	@echo "Running backend tests"
	@echo "---------------------"
	$(GO) get -u golang.org/x/tools/cmd/goimports
	$(GO) generate -tags ci  ./...
	make imports
	$(GO) test -tags ci ./...                        # run unit tests
	make format

coverage:
	@echo "----------------"
	@echo "Running coverage"
	@echo "----------------"
	$(GOLANG_ANNOT_ROOT)/scripts/coverage.sh --html

install: clean
	@echo "----------------------------"
	@echo "Installing for $(GO_VERSION)"
	@echo "----------------------------"
	$(GO) install ./...

clean:
	find . -name '\$$*.go' -exec rm -rfv {} +
	rm -rf ./examples/rest/restTestLog/ ./generator/rest/testData/
	$(GO) clean ./...

.PHONY:
	help deps gen test citest coverage install clean all
