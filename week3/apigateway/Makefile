#project variables
APP_NAME := gateway
SRC := ./cmd/$(APP_NAME)

#go commands
GO := go
GOBUILD := $(GO) build
GORUN := $(GO) run


#default target
.PHONY: all
all: build

#build the binary
.PHONY: build
build:
	$(GOBUILD) -o bin/$(APP_NAME) $(SRC)

#run the app
.PHONY: run
run:
	$(GORUN) $(SRC)

.PHONY: service1proto
service1proto:
	protoc \
  --go_out=pkg/proto/service1 \
  --go-grpc_out=pkg/proto/service1 \
  -I pkg/proto/service1 \
  pkg/proto/service1/service1.proto


#clean binaries
.PHONY: clean
clean:
	rm -rf bin/$(APP_NAME)

#rebuild (clean + build)
.PHONY: rebuild
rebuild: clean build
