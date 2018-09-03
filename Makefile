GO_IMAGE_TAG=115100/poc-go-client-side-grpc-lb:latest
PHP_IMAGE_TAG=115100/poc-php-client-side-grpc-lb:latest

default: all

docker: all
	docker build --pull -t $(GO_IMAGE_TAG) -f Dockerfile.golang .
	docker push $(GO_IMAGE_TAG)

	docker build --pull -t $(PHP_IMAGE_TAG) -f Dockerfile.php php
	docker push $(PHP_IMAGE_TAG)

all: proto go/vendor/ bin/greeter-client bin/greeter-server

proto: go/greeterpb/greeter.pb.go php/Greeterpb/GreeterClient.php php/Greeterpb/GreetReply.php php/Greeterpb/GreetRequest.php

clean:
	rm -rf bin/

.PHONY: default docker all proto clean

go/vendor/:
	cd go/ && dep ensure -v -vendor-only

bin/%: go/cmd/%/main.go
	@mkdir -p $(dir $@)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s" -a -installsuffix cgo -o $@ $<

GO_GREETERPB_DIR=go/greeterpb
$(GO_GREETERPB_DIR)/%.pb.go: greeter.proto
	@mkdir -p $(dir $@)
	protoc --go_out=plugins=grpc:go/greeterpb $<

PHP_GREETERPB_DIR=php/Greeterpb
$(PHP_GREETERPB_DIR)/%.php: greeter.proto
	@mkdir -p $(dir $@)
	protoc --php_out=php/ --grpc_out=php/ --plugin=protoc-gen-grpc=$(shell which grpc_php_plugin) $<
