IMAGE_TAG=115100/poc-client-side-grpc-lb:latest

default: all

docker: all
	docker build --pull -t $(IMAGE_TAG) .
	docker push $(IMAGE_TAG)

all: proto vendor/ bin/greeter-client bin/greeter-server

proto: greeterpb/greeter.pb.go

clean:
	rm -rf bin/

.PHONY: default docker test all clean

vendor/:
	dep ensure -v

bin/%: cmd/%/main.go
	@mkdir -p $(dir $@)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s" -a -installsuffix cgo -o $@ $<

greeterpb/%.pb.go: greeterpb/%.proto
	protoc --go_out=plugins=grpc:. greeterpb/greeter.proto
