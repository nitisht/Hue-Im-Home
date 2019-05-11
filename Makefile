 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=hue-im-home
BINARY_ARM=$(BINARY_NAME)_arm
BINARY_AMD64=$(BINARY_NAME)_amd64

all: deps clean build-go-linux build-go-arm
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_ARM)
	rm -f $(BINARY_AMD64)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)
deps:
	$(GOGET) github.com/heatxsink/go-hue/configuration
	$(GOGET) github.com/heatxsink/go-hue/lights
	$(GOGET) github.com/heatxsink/go-hue/portal

# Cross compilation
build-go-linux:
	make deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_AMD64)
build-go-arm:
	make deps
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o $(BINARY_ARM)

build-docker-linux:
	docker build . -t selexin/hue-im-home:amd64-latest -f docker/Dockerfile
build-docker-arm:
	docker build . -t selexin/hue-im-home:arm32v7-latest -f docker/Dockerfile.arm32v7