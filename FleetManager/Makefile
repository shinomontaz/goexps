SHELL:=/bin/bash
# current project in last GOPATH folder
LAST_GOPATH:=$(lastword $(subst :, ,$(GOPATH)))
# glide and deps in first GOPATH folder
FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GOBIN:=$(LAST_GOPATH)/bin
TEST_ARGS?=$$($(GLIDE_BIN) nv)

BIN?=./bin/FleetManager
PROJECT_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
GIT_REV:=$(shell git describe --all --long HEAD)
GO_VERSION:=$(shell go version)
GO_VERSION_SHORT:=$(shell echo $(GO_VERSION)|sed -E 's/.* go(.*) .*/\1/g')
BUILD_TS:=$(shell date +"%s")
LDFLAGS:=-X 'main.appVersion=$(GIT_REV)' -X 'main.goVersion=$(GO_VERSION_SHORT)' -X 'main.buildTimestamp=$(BUILD_TS)'

GLIDE_PATH:=$(FIRST_GOPATH)/src/github.com/Masterminds/glide
GLIDE_VERSION:=v0.13.1
GLIDE_BIN:=$(GOBIN)/glide-$(GLIDE_VERSION)
VENDOR_DIR=vendor
GLIDE_CHANGE=$(VENDOR_DIR)/.glide_change

# install glide
$(GLIDE_BIN):
ifeq ($(wildcard $(GLIDE_BIN)),)
	$(info #Installing glide version $(GLIDE_VERSION)...)
ifeq ($(wildcard $(GLIDE_PATH)),)
	mkdir -p $(GLIDE_PATH) && cd $(GLIDE_PATH) ;\
	git clone https://github.com/Masterminds/glide.git .
endif
	cd $(GLIDE_PATH) && git fetch --tags && git checkout $(GLIDE_VERSION) ;\
	git reset --hard && git clean -fd ;\
	make clean && make build && mv ./glide $(GLIDE_BIN)
else
	$(info #Found glide $(GLIDE_VERSION) in $(GLIDE_BIN)...)
endif

glide.lock:
ifeq ($(wildcard glide.lock),)
	@$(MAKE) -B deps-update
endif

# check glide.lock, glide.yaml, vendor dir changes and reinstall dependencies if needed
$(GLIDE_CHANGE): glide.lock glide.yaml $(GLIDE_BIN) | $(VENDOR_DIR)
	@cd $(PROJECT_PATH)
	$(GLIDE_BIN) install --force
	touch $(GLIDE_CHANGE)

$(VENDOR_DIR):
	@mkdir -p $(VENDOR_DIR)


# install project dependencies (and install glide if you have no glide with right version)
deps: $(GLIDE_BIN)
	$(info #Install dependencies...)
	@cd $(PROJECT_PATH)
	@$(MAKE) -B $(GLIDE_CHANGE)

# update project dependencies and glide.lock (and install glide if you have no glide)
deps-update: $(GLIDE_BIN)
	$(info #Update dependencies...)
	GOPATH=$(LAST_GOPATH) $(GLIDE_BIN) update --force
	touch $(GLIDE_CHANGE)

fast-build-noglide:
	$(info #Building goblin binary to $(BIN)...)
	GOGC=off go build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/*.go

# just build the binary, work with dependencies on changes
fast-build: $(GLIDE_CHANGE)
	$(info #Building goblin binary to $(BIN)...)
	GOGC=off go build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/*.go

# build the binary after install dependencies and running tests (you don't need to run "make deps" before "make build")
build: $(GLIDE_CHANGE) test fast-build

# run unit tests
test: $(GLIDE_CHANGE)
	$(info #Running tests...)
	GOGC=off go test $(TEST_ARGS)

protoc-gen-go-http-server:
	@go get -u gitlab.ozon.ru/platform/protobuf/protoc-gen-go-http-server

generate: # install protoc-gen-go-http-server
	protoc -I. -I/usr/local/include -I$(FIRST_GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gofast_out=plugins=grpc:. --go-http-server_out=verbose=true,impl=true:. ./FleetManager/FleetManager.proto

generate-impl: # install protoc-gen-go-http-server
	protoc -I. -I/usr/local/include -I$(FIRST_GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gofast_out=plugins=grpc:. --go-http-server_out=verbose=true,impl=true:. ./FleetManager/FleetManager.proto


.PHONY: deps deps-update fast-build build test gen
