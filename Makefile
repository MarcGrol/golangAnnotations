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
	@echo "----------------------------"
	@echo "Run static analysis on source-code"
	@echo "----------------------------"
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
	for i in `find . -name "*.go"`; do goimports -w -local github.com/ $${i}; done

format:
	@echo "----------------------"
	@echo "Formatting source-code"
	@echo "----------------------"
	for i in `find . -name "*.go"`; do gofmt -s -w $${i}; done

gen: generate imports format

test:
	@echo "-------------"
	@echo "Running backend tests"
	@echo "-------------"
	$(GO) test ./...                        # run unit tests
	make format

citest:
	@echo "-------------"
	@echo "Running backend tests"
	@echo "-------------"
	$(GO) generate -tags ci  ./...
	$(GO) test -tags ci ./...                        # run unit tests
	make format

coverage:
	@echo "-------------"
	@echo "Running coverage"
	@echo "-------------"
	$(GOLANG_ANNOT_ROOT)/scripts/coverage.sh --html

install:
	@echo "----------------"
	@echo "Installing for $(GO_VERSION)"
	@echo "----------------"
	$(GO) install ./...

clean:
	rm -rf ./examples/structExample/aggregates.go ./examples/structExample/wrappers.go \
		./examples/structExample/wrappers_test.go \
		./examples/structExample/structExample_json.go ./examples/structExample/enumExample_json.go \ 
	    ./examples/rest/httpTourService.go \
	    ./examples/rest/httpTourServiceHelpers_test.go ./examples/rest/httpClientForTourService.go \
		./examples/rest/restTestLog/ ./examples/store/structExampleEventStore.go	
	$(GO) clean ./...

.PHONY:
	help deps gen test coverage install clean all
