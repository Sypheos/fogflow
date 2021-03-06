RELEASE_DIR ?= release#
GO_ENV = CGO_ENABLED=0
GO_FLAGS =
LOG_NAME ?= $(shell basename $$(pwd))
log = echo -e "$(LOG_NAME) "

# Infer GOOS and GOARCH
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# default vendor folder
VENDOR_DIR ?= $(PWD)/vendor
VENDOR_FILE ?= Gopkg.toml

LAZY_GOOS = `echo $@ | sed 's:$(RELEASE_DIR)/.*-\(.*\)-\(.*\):\1:'`
LAZY_GOARCH = `echo $@ | sed 's:$(RELEASE_DIR)/.*-\(.*\)-\(.*\):\2:'`
LAZY_GOEXE = $$(GOOS=$(LAZY_GOOS) go env GOEXE)

# the first entry of the go path
GO_PATH ?= $(shell echo $(GOPATH) | awk -F':' '{ print $$1 }')
LD_FLAGS =

# init initializes go
init:
	@$(log) "Initializing go"
	@make dev-deps
	@make deps

# get tools required for development
dev-deps:
	@$(log) "Installing go dev dependencies"
	@$(log) "Getting dep" && go get -u github.com/golang/dep/cmd/dep

# install dependencies
deps:
	@$(log) "Installing go dependencies"
	@dep ensure -v

## initialize go dep
$(VENDOR_FILE):
	@$(log) "Initializing go deps"
	@mkdir -p $(VENDOR_DIR) && cd $(VENDOR_DIR)/.. && dep init

test:
	@$(log) "Testing `$(TEST_PACKAGES) | $(count)` go packages"
	@$(GO) test $(GO_TEST_FLAGS) `$(TEST_PACKAGES)`

broker:MAIN=./broker/cmd/main.go
broker:$(RELEASE_DIR)/broker-$(GOOS)-$(GOARCH)
broker.docker: broker
	docker build -f ./broker/Dockerfile -t "fogflow/broker" .

discovery:MAIN=./discovery/cmd/main.go
discovery:$(RELEASE_DIR)/discovery-$(GOOS)-$(GOARCH)
discovery.docker: discovery
	docker build -f ./discovery/Dockerfile -t "fogflow/discovery" .

master:MAIN=./master/cmd/main.go
master:$(RELEASE_DIR)/master-$(GOOS)-$(GOARCH)
master.docker: master
	docker build -f ./master/Dockerfile -t "fogflow/master" .

worker:MAIN=./worker/cmd/main.go
worker:$(RELEASE_DIR)/worker-$(GOOS)-$(GOARCH)
worker.docker: worker
	docker build -f ./worker/Dockerfile -t "fogflow/worker" .

# Build the executable
$(RELEASE_DIR)/%:
	@$(log) "Building" [$(GO_ENV) GOOS=$(LAZY_GOOS) GOARCH=$(LAZY_GOARCH) $(GO) build $(GO_FLAGS) ...] to "$@$(LAZY_GOEXE)"
	@$(GO_ENV) go build -o "$@$(LAZY_GOEXE)" $(GO_FLAGS) $(LD_FLAGS) $(MAIN)
	@chmod +x $@

all: broker discovery

all.docker: broker.docker discovery.docker
clean:
	rm -rf ./release/*

re: clean all