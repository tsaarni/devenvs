IMAGE ?= envoy-ext-authz-server:latest

build:
	CGO_ENABLED=0 go build -o envoy-ext-authz-server .

container:
	docker build -t $(IMAGE) .

kind-load: container
	kind load docker-image $(IMAGE) --name contour

update-deps:
	go get -u ./...
	go mod tidy

clean:
	rm -f envoy-ext-authz-server
