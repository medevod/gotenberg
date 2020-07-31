# Build.
GOTENBERG_VERSION=snapshot
TINI_VERSION=0.19.0
GOTENBERG_USER_GID=1001
GOTENBERG_USER_UID=1001
DOCKER_REGISTRY=thecodingmachine

#DOCKER_CLI_EXPERIMENTAL=enabled \
#docker buildx build \
#--platform linux/arm64,linux/amd64 \

image:
	docker build \
	--build-arg GOTENBERG_VERSION=$(GOTENBERG_VERSION) \
	--build-arg GOTENBERG_USER_GID=$(GOTENBERG_USER_GID) \
	--build-arg GOTENBERG_USER_UID=$(GOTENBERG_USER_UID) \
	--build-arg TINI_VERSION=$(TINI_VERSION) \
	-t $(DOCKER_REGISTRY)/gotenberg:$(GOTENBERG_VERSION) \
	-f build/package/Dockerfile .

fmt:
	go fmt ./...
	go mod tidy

test:
	go test -race -cover ./...

lint:
	 golangci-lint run --enable-all --disable gomnd --disable goerr113 --disable gochecknoglobals --disable gochecknoinits