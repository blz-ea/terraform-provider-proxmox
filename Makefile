OS = $(shell go env GOOS)
ARCH = $(shell go env GOARCH)
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
ACC_TEST?=$$(go list ./proxmoxtf/acceptancetests |grep -v 'vendor')
NAME=$$(grep TerraformProviderName proxmoxtf/version.go | grep -o -e 'terraform-provider-[a-z]*')
TARGETS=darwin linux windows
TERRAFORM_PLUGIN_EXTENSION=
VERSION=$$(grep TerraformProviderVersion proxmoxtf/version.go | grep -o -e '[0-9]\.[0-9]\.[0-9]')
DEV_VERSION = 99.0.0
PROVIDER_PATH = registry.terraform.io/blz-ea/proxmox/$(DEV_VERSION)/$(OS)_$(ARCH)/
PROVIDER_PATH_WINDOWS = registry.terraform.io\blz-ea\proxmox\$(DEV_VERSION)\$(OS)_$(ARCH)\

ifeq ($(OS),Windows_NT)
	TERRAFORM_CACHE_DIRECTORY=$$(cygpath -u "$(APPDATA)")/terraform.d/plugins
	TERRAFORM_PLATFORM=windows_amd64
	TERRAFORM_PLUGIN_EXTENSION=.exe
else
	TERRAFORM_CACHE_DIRECTORY=$(HOME)/terraform.d/plugins
	UNAME_S=$$(shell uname -s)

	ifeq ($(UNAME_S),Darwin)
		TERRAFORM_PLATFORM=darwin_amd64
	else
		TERRAFORM_PLATFORM=linux_amd64
	endif
endif

default: build

build:
	fmtcheck
	go install

build-install-dev:
	go build -o terraform-provider-proxmox_$(DEV_VERSION)
ifeq ($(OS), darwin)
	mkdir -p ~/.terraform.d/plugins/$(PROVIDER_PATH)
	mv terraform-provider-proxmox_$(DEV_VERSION) ~/.terraform.d/plugins/$(PROVIDER_PATH)
endif
ifeq ($(OS), linux)
	mkdir -p ~/.terraform.d/plugins/$(PROVIDER_PATH)
	mv terraform-provider-proxmox_$(DEV_VERSION) ~/.terraform.d/plugins/$(PROVIDER_PATH)
endif
ifeq ($(OS), windows)
	mkdir %APPDATA%\terraform.d\plugins\$(PROVIDER_PATH_WINDOWS)
	mv terraform-provider-proxmox_$(DEV_VERSION) %APPDATA%\terraform.d\plugins\$(PROVIDER_PATH_WINDOWS)
endif

fmt:
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

init:
	go get ./...

targets: $(TARGETS)

test:
	go test -v ./...

testacc:
	@echo "==> Sourcing .env file if available"
	if [ -f .env ]; then set -o allexport; . ./.env; set +o allexport; fi; \
	TF_ACC=1 go test -timeout 120m -run ^TestAcc* -tags "${*:-all}" -v $(ACC_TEST) || echo "Build finished in error due to failed tests"

$(TARGETS):
	GOOS=$@ GOARCH=amd64 CGO_ENABLED=0 go build \
		-o "dist/$@/$(NAME)_v$(VERSION)-custom_x4" \
		-a -ldflags '-extldflags "-static"'
	zip \
		-j "dist/$(NAME)_v$(VERSION)-custom_$@_amd64.zip" \
		"dist/$@/$(NAME)_v$(VERSION)-custom_x4"

.PHONY: build build-and-install-dev-version example example-apply example-destroy example-init example-plan fmt init targets test testacc $(TARGETS)
