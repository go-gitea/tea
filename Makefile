DIST := dist
IMPORT := code.gitea.io/tea
export GO111MODULE=off

GO ?= go
SED_INPLACE := sed -i
SHASUM ?= shasum -a 256

export PATH := $($(GO) env GOPATH)/bin:$(PATH)

ifeq ($(OS), Windows_NT)
	EXECUTABLE := tea.exe
else
	EXECUTABLE := tea
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		SED_INPLACE := sed -i ''
	endif
endif

GOFILES := $(shell find . -name "*.go" -type f ! -path "./vendor/*" ! -path "*/bindata.go")
GOFMT ?= gofmt -s

GOFLAGS := -i -v
EXTRA_GOFLAGS ?=

MAKE_VERSION := $(shell make -v | head -n 1)

ifneq ($(DRONE_TAG),)
	VERSION ?= $(subst v,,$(DRONE_TAG))
	TEA_VERSION ?= $(VERSION)
else
	ifneq ($(DRONE_BRANCH),)
		VERSION ?= $(subst release/v,,$(DRONE_BRANCH))
	else
		VERSION ?= master
	endif
	TEA_VERSION ?= $(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')
endif

LDFLAGS := -X "main.Version=$(TEA_VERSION)" -X "main.Tags=$(TAGS)"

PACKAGES ?= $(shell $(GO) list ./... | grep -v /vendor/)
SOURCES ?= $(shell find . -name "*.go" -type f)

TAGS ?=

ifeq ($(OS), Windows_NT)
	EXECUTABLE := tea.exe
else
	EXECUTABLE := tea
endif

# $(call strip-suffix,filename)
strip-suffix = $(firstword $(subst ., ,$(1)))

.PHONY: all
all: build

.PHONY: clean
clean:
	$(GO) clean -i ./...
	rm -rf $(EXECUTABLE) $(DIST)

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: vet
vet:
	$(GO) vet $(PACKAGES)

.PHONY: lint
lint:
	@hash revive > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/mgechev/revive; \
	fi
	revive -config .revive.toml -exclude=./vendor/... ./... || exit 1

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error -i unknwon,destory $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w -i unknwon $(GOFILES)

.PHONY: fmt-check
fmt-check:
	# get all go files and run go fmt on them
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: test
test:
	GO111MODULE=on $(GO) test -mod=vendor -tags='sqlite sqlite_unlock_notify' $(PACKAGES)

.PHONY: unit-test-coverage
unit-test-coverage:
	$(GO) test -tags='sqlite sqlite_unlock_notify' -cover -coverprofile coverage.out $(PACKAGES) && echo "\n==>\033[32m Ok\033[m\n" || exit 1

.PHONY: vendor
vendor:
	GO111MODULE=on $(GO) mod tidy && GO111MODULE=on $(GO) mod vendor

.PHONY: test-vendor
test-vendor: vendor
	@diff=$$(git diff vendor/); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make vendor' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: check
check: test

.PHONY: install
install: $(wildcard *.go)
	$(GO) install -v -tags '$(TAGS)' -ldflags '-s -w $(LDFLAGS)'

.PHONY: build
build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	GO111MODULE=on $(GO) build -mod=vendor $(GOFLAGS) $(EXTRA_GOFLAGS) -tags '$(TAGS)' -ldflags '-s -w $(LDFLAGS)' -o $@

.PHONY: release
release: release-dirs release-windows release-linux release-darwin release-copy release-compress release-check

.PHONY: release-dirs
release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

.PHONY: release-windows
release-windows:
	@hash xgo > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u src.techknowlogick.com/xgo; \
	fi
	xgo -dest $(DIST)/binaries -tags 'netgo $(TAGS)' -ldflags '-linkmode external -extldflags "-static" $(LDFLAGS)' -targets 'windows/*' -out tea-$(VERSION) .
ifeq ($(CI),drone)
	cp /build/* $(DIST)/binaries
endif

.PHONY: release-linux
release-linux:
	@hash xgo > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u src.techknowlogick.com/xgo; \
	fi
	xgo -dest $(DIST)/binaries -tags 'netgo $(TAGS)' -ldflags '-linkmode external -extldflags "-static" $(LDFLAGS)' -targets 'linux/*' -out tea-$(VERSION) .
ifeq ($(CI),drone)
	cp /build/* $(DIST)/binaries
endif

.PHONY: release-darwin
release-darwin:
	@hash xgo > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u src.techknowlogick.com/xgo; \
	fi
	xgo -dest $(DIST)/binaries -tags 'netgo $(TAGS)' -ldflags '$(LDFLAGS)' -targets 'darwin/*' -out tea-$(VERSION) .
ifeq ($(CI),drone)
	cp /build/* $(DIST)/binaries
endif

.PHONY: release-copy
release-copy:
	cd $(DIST); for file in `find /build -type f -name "*"`; do cp $${file} ./release/; done;

.PHONY: release-compress
release-compress:
	@hash gxz > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/ulikunitz/xz/cmd/gxz; \
	fi
	cd $(DIST)/release/; for file in `find . -type f -name "*"`; do echo "compressing $${file}" && gxz -k -9 $${file}; done;

.PHONY: release-check
release-check:
	cd $(DIST)/release/; for file in `find . -type f -name "*"`; do echo "checksumming $${file}" && $(SHASUM) `echo $${file} | sed 's/^..//'` > $${file}.sha256; done;
