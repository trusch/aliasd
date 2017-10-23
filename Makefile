SRC=$(shell find ./generator ./manager ./storage -type f -name "*.go") main.go

all: docker

aliasd: $(SRC) vendor
	docker run \
		-v $(shell pwd):/go/src/github.com/trusch/aliasd \
		-w /go/src/github.com/trusch/aliasd \
		-u $(shell stat -c '%u:%g' .) \
		-e CGO_ENABLED=0 \
		-e GOOS=linux \
		golang:1.9 \
			go build -v -a -ldflags '-extldflags "-static"' .

vendor: glide.yaml
	docker run \
		-v $(shell pwd):/go/src/github.com/trusch/aliasd \
		-w /go/src/github.com/trusch/aliasd \
		-u $(shell stat -c '%u:%g' .) \
		-e CGO_ENABLED=0 \
		-e GOOS=linux \
		golang:1.9 bash -c \
			"(curl https://glide.sh/get | sh) && glide --home /tmp update"

docker: aliasd Dockerfile
	docker build -t trusch/aliasd .

clean:
	rm -rf aliasd vendor
