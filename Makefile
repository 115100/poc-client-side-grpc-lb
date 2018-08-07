GO_IMAGE_TAG=115100/poc-client-side-grpc-lb:latest

default: all

docker: all
	docker build --pull -t $(GO_IMAGE_TAG) -f Dockerfile.go .
	docker push $(GO_IMAGE_TAG)

all: proto go/vendor/ bin/greeter-client bin/greeter-server

proto: go/greeterpb/greeter.pb.go

clean:
	rm -rf bin/

.PHONY: default docker test all clean

go/vendor/:
	cd go/ && dep ensure -v

bin/%: go/cmd/%/main.go
	@mkdir -p $(dir $@)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s" -a -installsuffix cgo -o $@ $<

GREETERPB_DIR=go/greeterpb
$(GREETERPB_DIR)/%.pb.go: %.proto
	protoc --go_out=plugins=grpc:go $<
