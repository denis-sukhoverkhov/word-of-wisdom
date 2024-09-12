# Variables for Docker image names and tags
DOCKER_SERVER_IMAGE ?= word-of-wisdom-server
DOCKER_CLIENT_IMAGE ?= word-of-wisdom-client
DOCKER_TAG ?= latest

.PHONY: install-lint
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint: install-lint
	golangci-lint run ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

.PHONY: build-server
build-server:
	go build -o bin/server cmd/server/main.go

.PHONY: build-client
build-client:
	go build -o bin/client cmd/client/main.go

.PHONY: run-server
run-server:
	source .env && go run cmd/server/main.go

.PHONY: run-client
run-client:
	source .env && go run cmd/client/main.go

# Build the Docker image for the server
.PHONY: docker-build-server
docker-build-server:
	docker build -t $(DOCKER_SERVER_IMAGE):$(DOCKER_TAG) -f Dockerfile.server .

# Build the Docker image for the client
.PHONY: docker-build-client
docker-build-client:
	docker build -t $(DOCKER_CLIENT_IMAGE):$(DOCKER_TAG) -f Dockerfile.client .

# Build both server and client Docker images
.PHONY: docker-build-all
docker-build-all: 
	docker-build-server docker-build-client

.PHONY: up
up: 
	docker-compose up --build
