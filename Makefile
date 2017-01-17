# This is a Makefile which maintains files automatically generated but to be
# shipped together with other files.
CFLAGS=-g

PKG=github.com/grpc-ecosystem/grpc-gateway
GOOGLEAPIS_DIR=third_party/googleapis
GO_SRC_PATH=$(GOPATH)/src
GOOGLEAPIS_DIR=third_party/googleapis
OPTIONS_PROTO=$(GOOGLEAPIS_DIR)/google/api/annotations.proto $(GOOGLEAPIS_DIR)/google/api/http.proto
PKGMAP=Mgoogle/api/annotations.proto=$(PKG)/$(GOOGLEAPIS_DIR)/google/api

# SVCSRCS_GEN is the paratemers of generating service code
SVCSRCS_GEN=protoc -I . \
	-I $(GO_SRC_PATH) \
	-I $(GO_SRC_PATH)/$(PKG)/$(GOOGLEAPIS_DIR) \
	--go_out=$(PKGMAP),plugins=grpc:. \
	$(PROTOC_SRC)

# GWSRCS_GEN is the paratemers of generating gateway code
GWSRCS_GEN=protoc -I . \
	-I $(GO_SRC_PATH) \
	-I $(GO_SRC_PATH)/$(PKG)/$(GOOGLEAPIS_DIR) \
	--grpc-gateway_out=logtostderr=true:. \
	$(PROTOC_SRC)

# PROTOC_SRC is the list of proto buf files
PROTOC_SRC=pb/fingerprint.proto \
	pb/result.proto

protoc:
	$(SVCSRCS_GEN)
	$(GWSRCS_GEN)

pb: protoc
	go build -o build/findpro main.go

build:
	go build -o build/findpro main.go

install:
	go build -o $(GOPATH)/bin/findpro main.go

server:
	./build/findpro server

gateway:
	./build/findpro gateway

.PHONY: server gateway protoc build install pb
