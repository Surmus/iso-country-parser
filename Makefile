GO ?= $(shell which go)
GOFMT := $(shell which gofmt) "-s"
GOLINT ?= ${GOPATH}/bin/golint
PACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go" -type f)
TESTFOLDER := $(shell $(GO) list ./... | grep -v test)

all: test build_linux build_windows

.PHONY: install
install: deps
	$(GO) install ./cmd/parser

.PHONY: build_linux
build_linux: deps
	env GOOS=linux GOARCH=amd64 $(GO) build -o ./build/linux64/parser --tags "linux" -a -tags netgo -ldflags '-w -extldflags "-static"' ./cmd/parser

.PHONY: build_windows
build_windows: deps
	env GOOS=windows GOARCH=amd64 $(GO) build -o ./build/win64/parser.exe ./cmd/parser

.PHONY: test
test:
	echo "mode: atomic" > coverage.out
	for d in $(TESTFOLDER); do \
		$(GO) test -v -coverpkg=./... -covermode=atomic -coverprofile=profile.out $$d > tmp.out; \
		cat tmp.out; \
		if grep -q "FAIL" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		elif grep -q "build failed" tmp.out; then \
			rm tmp.out; \
			exit; \
		fi; \
		if [ -f profile.out ]; then \
			cat profile.out | grep -v "mode:" >> coverage.out; \
			rm profile.out; \
		fi; \
	done

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: vet
vet:
	$(GO) vet $(PACKAGES)

.PHONY: lint
lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u golang.org/x/lint/golint; \
	fi
	for PKG in $(PACKAGES); do $(GOLINT) -set_exit_status $$PKG || exit 1; done;

.PHONY: deps
deps:
	@hash go > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "Install Go language before running this!"; \
	fi
