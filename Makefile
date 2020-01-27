GIT_COMMIT  = $(shell git rev-list -1 HEAD)
GIT_VERSION = $(shell git describe --always --abbrev=7 --dirty)
# By default, disable CGO_ENABLED. See the details on https://golang.org/cmd/cgo
CGO         ?= 0

DOCKER        := docker
DOCKERFILE_DIR?= ./docker
DOCKERFILE    := Dockerfile

RELEASE_NAME  ?= rawpotion-go

# Binaries
BINARIES ?= rawpotion-go

# Output dir
OUT_DIR := ./dist

# Architecture
LOCAL_ARCH := $(shell uname -m)
ifeq ($(LOCAL_ARCH),x86_64)
	TARGET_ARCH_LOCAL=amd64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),armv8)
	TARGET_ARCH_LOCAL=arm64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 4),armv)
	TARGET_ARCH_LOCAL=arm
else
	TARGET_ARCH_LOCAL=amd64
endif
export GOARCH ?= $(TARGET_ARCH_LOCAL)

ifeq ($(GOARCH),amd64)
	LATEST_TAG=latest
else
	LATEST_TAG=latest-$(GOARCH)
endif

# OS
LOCAL_OS := $(shell uname)
ifeq ($(LOCAL_OS),Linux)
   TARGET_OS_LOCAL = linux
else ifeq ($(LOCAL_OS),Darwin)
   TARGET_OS_LOCAL = darwin
else
   TARGET_OS_LOCAL ?= windows
endif
export GOOS ?= $(TARGET_OS_LOCAL)

ifeq ($(GOOS),windows)
BINARY_EXT_LOCAL:=.exe
GOLANGCI_LINT:=golangci-lint.exe
export ARCHIVE_EXT = .zip
else
BINARY_EXT_LOCAL:=
GOLANGCI_LINT:=golangci-lint
export ARCHIVE_EXT = .tar.gz
endif
export BINARY_EXT ?= $(BINARY_EXT_LOCAL)

DEFAULT_LDFLAGS:=-X $(BASE_PACKAGE_NAME)/pkg/version.commit=$(GIT_VERSION) -X $(BASE_PACKAGE_NAME)/pkg/version.version=$(DAPR_VERSION)

# Build Type
ifeq ($(origin DEBUG), undefined)
  BUILDTYPE_DIR:=release
  LDFLAGS:="$(DEFAULT_LDFLAGS) -s -w"
else ifeq ($(DEBUG),0)
  BUILDTYPE_DIR:=release
  LDFLAGS:="$(DEFAULT_LDFLAGS) -s -w"
else
  DOCKERFILE:=debug.Dockerfile
  BUILDTYPE_DIR:=debug
  GCFLAGS:=-gcflags="all=-N -l"
  LDFLAGS:="$(DEFAULT_LDFLAGS)"
  $(info Build with debugger information)
endif

# Go Build Details 
# Inspiration taken from dapr/dapr

PACKAGE := github.com/kjuulh/rawpotion-go

RP_OUT_DIR := $(OUT_DIR)/$(GOOS)_$(GOARCH)/$(BUILDTYPE_DIR)
RP_LINUX_OUT_DIR := $(OUT_DIR)/linux_$(GOARCH)/$(BUILDTYPE_DIR)

all: test build

.PHONY: build
RP_BINS := $(foreach ITEM,$(BINARIES),$(RP_OUT_DIR)/$(ITEM)$(BINARY_EXT))
build: $(RP_BINS)

define genBinariesForTarget
.PHONY: $(5)/$(1)
$(5)/$(1):
	CGO_ENABLED=$(CGO) GOOS=$(3) GOARCH=$(4) go build $(GCFLAGS) -ldflags=$(LDFLAGS) \
	-o $(5)/$(1) -mod=vendor \
	$(2)/main.go;
endef

# Generate binary targets
$(foreach ITEM,$(BINARIES),$(eval $(call genBinariesForTarget,$(ITEM)$(BINARY_EXT),./cmd/$(ITEM),$(GOOS),$(GOARCH),$(RP_OUT_DIR))))

BUILD_LINUX_BINS:=$(foreach ITEM,$(BINARIES),$(RP_LINUX_OUT_DIR)/$(ITEM))
build-linux: $(BUILD_LINUX_BINS)

# Generate linux binaries targets to build linux docker image
ifneq ($(GOOS), linux)
$(foreach ITEM,$(BINARIES),$(eval $(call genBinariesForTarget,$(ITEM),./cmd/$(ITEM),linux,$(GOARCH),$(RP_LINUX_OUT_DIR))))
endif

.PHONY: test
test:
		go test ./pkg/... -mod=vendor
		#go test ./tests/... -mod=vendor

# ARCHIVE
ARCHIVE_OUT_DIR   ?= $(RP_OUT_DIR)
ARCHIVE_FILE_EXTS := $(foreach ITEM,$(BINARIES),archive-$(ITEM)$(ARCHIVE_EXT))

archive: $(ARCHIVE_FILE_EXTS)
define genArchiveBinary
ifeq ($(GOOS),windows)
archive-$(1).zip:
	7z.exe a -tzip "$(2)\\$(1)_$(GOOS)_$(GOARCH)$(ARCHIVE_EXT)" "$(RP_OUT_DIR)\\$(1)$(BINARY_EXT)"
else
archive-$(1).tar.gz:
	tar czf "$(2)/$(1)_$(GOOS)_$(GOARCH)$(ARCHIVE_EXT)" -C "$(RP_OUT_DIR)" "$(1)$(BINARY_EXT)"
endif
endef

# Generate archive-*.[zip|tar.gz] targets
$(foreach ITEM,$(BINARIES),$(eval $(call genArchiveBinary,$(ITEM),$(ARCHIVE_OUT_DIR))))

# DOCKER

LINUX_BINS_OUT_DIR=$(OUT_DIR)/linux_$(GOARCH)
DOCKER_IMAGE_TAG=$(RP_REGISTRY)/$(RELEASE_NAME):$(RP_TAG)

ifeq ($(LATEST_RELEASE),true)
DOCKER_IMAGE_LATEST_TAG=$(RP_REGISTRY)/$(RELEASE_NAME):$(LATEST_TAG)
endif

check-docker-env:
ifeq ($(RP_REGISTRY),)
	$(error RP_REGISTRY environment variable must be set)
endif
ifeq ($(RP_TAG),)
	$(error RP_TAG environment variable must be set)
endif

docker-build: check-docker-env
	$(info Building $(DOCKER_IMAGE_TAG) docker image ...)
	$(DOCKER) build -f $(DOCKERFILE_DIR)/$(DOCKERFILE) . -t $(DOCKER_IMAGE_TAG)
ifeq ($(LATEST_RELEASE),true)
	$(info Building $(DOCKER_IMAGE_LATEST_TAG) docker image ...)
	$(DOCKER) tag $(DOCKER_IMAGE_TAG) $(DOCKER_IMAGE_LATEST_TAG)
endif

# push docker image to the registry
docker-push: docker-build
	$(info Pushing $(DOCKER_IMAGE_TAG) docker image ...)
	$(DOCKER) push $(DOCKER_IMAGE_TAG)
ifeq ($(LATEST_RELEASE),true)
	$(info Pushing $(DOCKER_IMAGE_LATEST_TAG) docker image ...)
	$(DOCKER) push $(DOCKER_IMAGE_LATEST_TAG)
endif

docker-test:
	$(DOCKER) build -f $(DOCKERFILE_DIR)/test.$(DOCKERFILE) .