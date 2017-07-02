GITHUB_REPO       = github.com/JonPulsifer/go-check-dns
DOCKER_REGISTRY   = gcr.io/kubesec
DOCKER_IMAGE_NAME = go-check-dns
DOCKER_IMAGE_TAG  = latest
IMAGE = $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

all: build_native

clean:
	/bin/rm -v go-check-dns

test:
	go test -v

container:
	docker build -t builder:tmp -f Dockerfile.build .
	docker run  --rm -v $(shell pwd):/go/bin -e "CGO_ENABLED=0" builder go get -v $(GITHUB_REPO)
	docker build -t $(IMAGE) -f Dockerfile .
	docker rmi -f builder:tmp

push:
	gcloud docker -- push $(IMAGE)

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-check-dns

build_native:
	CGO_ENABLED=0 go build -v

run: build_native
	DATADOG_URL=127.0.0.1:8125 ./go-check-dns -project $(PROJECT)

.PHONY: all build_linux build_native clean container run test     
