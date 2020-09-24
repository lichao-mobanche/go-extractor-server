PROJECT     := github.com/lichao-mobanche/go-extractor-server
BINNAME     ?= extractor-server
BINDIR      := $(CURDIR)/bin
BUILDTIME   := $(shell date +'%Y-%m-%d %H:%M:%S')

GOPATH      := $(shell go env GOPATH)
GOVERSION   := $(shell go version)
GOIMPORTS   := $(GOPATH)/bin/goimports
GOLINT      := $(GOPATH)/bin/golangci-lint
INSTALLPATH := $(GOPATH)/bin

PKG         := ./...
TESTS       := .
LDFLAGS     :=
TESTFLAGS   :=
SRC         := $(shell find . -type f -name '*.go' -print)

GIT_COMMIT  := $(shell git rev-parse HEAD)
GIT_SHA     := $(shell git rev-parse --short HEAD)
GIT_TAG     := $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_STATUS  := $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

VERSION     ?= $(GIT_TAG)
VERSIONMOD  := $(PROJECT)/pkg/version

ifneq ($(VERSION),)
	LDFLAGS +=  -X "$(VERSIONMOD).version=$(VERSION)"
endif

LDFLAGS +=  -X "$(VERSIONMOD).goVersion=$(GOVERSION)"
LDFLAGS +=  -X "$(VERSIONMOD).buildTime=$(BUILDTIME)"
LDFLAGS +=  -X "$(VERSIONMOD).gitCommit=$(GIT_COMMIT)"
LDFLAGS +=  -X "$(VERSIONMOD).gitTag=$(GIT_TAG)"
LDFLAGS +=  -X "$(VERSIONMOD).gitStatus=$(GIT_STATUS)"

.PHONY: all
all: build

.PHONY: build
build: generate
build: $(BINDIR)/$(BINNAME)
$(BINDIR)/$(BINNAME): $(SRC)
		@echo
		@echo  "==> Building ./cmd/extractor-server $(BINDIR)/$(BINNAME) <=="
		GO111MODULE=on go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BINNAME) ./cmd/extractor-server
.PHONY: format
format: $(GOIMPORTS)
	@echo
	@echo  "==> Formatting <=="
	GO111MODULE=on go list -f '{{.Dir}}' ./... | xargs $(GOIMPORTS) -w

$(GOIMPORTS):
	@echo
	@echo  "==> Installing goimports <=="
	(cd /; GO111MODULE=on go get -u golang.org/x/tools/cmd/goimports)

.PHONY: install
install:
	GO111MODULE=on go build -i $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(INSTALLPATH)/$(BINNAME) ./cmd/extractor-server

.PHONY: clean
clean:
	@rm -rf $(BINDIR)

.PHONY: info
info:
	 @echo "Version:    ${VERSION}"
	 @echo "Go Version: ${GOVERSION}"
	 @echo "Git Tag:    ${GIT_TAG}"
	 @echo "Git Commit: ${GIT_COMMIT}"
	 @echo "Git Status: ${GIT_STATUS}"
