GOVERSION=$(shell go version)
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
VERSION=$(patsubst "%",%,$(lastword $(shell grep 'const Version' internal/config/version.go)))
ARTIFACTS_DIR=$(CURDIR)/artifacts/$(VERSION)
RELEASE_DIR=$(CURDIR)/release/$(VERSION)
SRC_FILES=$(shell find . -type f -name '*.go')
GITHUB_USERNAME=shogo82148

help: ## Show this text.
	# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: build-windows-amd64 build-linux-amd64 build-darwin-amd64 ## Build binaries.

.PHONY: all test clean help

test: ## Run test.
	go test -v -race ./...
	go vet ./...

clean: ## Remove built files.


##### build settings

.PHONY: build build-windows-amd64 build-linux-amd64 build-darwin-amd64

$(ARTIFACTS_DIR)/rdsmysql_$(GOOS)_$(GOARCH):
	@mkdir -p $@

$(ARTIFACTS_DIR)/rdsmysql_$(GOOS)_$(GOARCH)/rdsmysql$(SUFFIX): $(ARTIFACTS_DIR)/rdsmysql_$(GOOS)_$(GOARCH) $(SRC_FILES)
	@echo " * Building binary for $(GOOS)/$(GOARCH)..."
	@CGO_ENABLED=0 go build -o $@ cmd/rdsmysql/*.go

$(ARTIFACTS_DIR)/rdsmysql_$(GOOS)_$(GOARCH)/rdsmysqldump$(SUFFIX): $(ARTIFACTS_DIR)/rdsmysql_$(GOOS)_$(GOARCH) $(SRC_FILES)
	@echo " * Building binary for $(GOOS)/$(GOARCH)..."
	@CGO_ENABLED=0 go build -o $@ cmd/rdsmysqldump/*.go

build: $(ARTIFACTS_DIR)/rdsmysql_$(GOOS)_$(GOARCH)/rdsmysql$(SUFFIX) $(ARTIFACTS_DIR)/rdsmysql_$(GOOS)_$(GOARCH)/rdsmysqldump$(SUFFIX)

build-windows-amd64:
	@$(MAKE) build GOOS=windows GOARCH=amd64 SUFFIX=.exe

build-linux-amd64:
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-darwin-amd64:
	@$(MAKE) build GOOS=darwin GOARCH=amd64


##### release settings

.PHONY: release-windows-amd64 release-linux-amd64 release-darwin-amd64
.PHONY: release-targz release-zip release-files release-upload

$(RELEASE_DIR)/rdsmysql_$(GOOS)_$(GOARCH):
	@mkdir -p $@

release-windows-amd64:
	@$(MAKE) release-zip GOOS=windows GOARCH=amd64 SUFFIX=.exe

release-linux-amd64:
	@$(MAKE) release-targz GOOS=linux GOARCH=amd64

release-darwin-amd64:
	@$(MAKE) release-targz GOOS=darwin GOARCH=amd64

release-targz: build $(RELEASE_DIR)/rdsmysql_$(GOOS)_$(GOARCH)
	@echo " * Creating tar.gz for $(GOOS)/$(GOARCH)"
	tar -czf $(RELEASE_DIR)/rdsmysql_$(GOOS)_$(GOARCH).tar.gz -C $(ARTIFACTS_DIR) rdsmysql_$(GOOS)_$(GOARCH)

release-zip: build $(RELEASE_DIR)/rdsmysql_$(GOOS)_$(GOARCH)
	@echo " * Creating zip for $(GOOS)/$(GOARCH)"
	cd $(ARTIFACTS_DIR) && zip -9 $(RELEASE_DIR)/rdsmysql_$(GOOS)_$(GOARCH).zip rdsmysql_$(GOOS)_$(GOARCH)/*

release-files: release-windows-amd64 release-linux-amd64 release-darwin-amd64

release-upload: release-files
	ghr -u $(GITHUB_USERNAME) --draft --replace v$(VERSION) $(RELEASE_DIR)