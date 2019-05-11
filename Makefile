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
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_AMD64) -v
build-go-arm:
	make deps
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o $(BINARY_ARM) -v

build-docker-linux:
	docker build . -t selexin/hue-im-home:amd64-latest -f docker/Dockerfile
build-docker-arm:
	# Download Qemu binary if they have not been donwloaded yet
	[ -f qemu-arm-static ] || (curl -L https://github.com/balena-io/qemu/releases/download/v3.0.0%2Bresin/qemu-3.0.0+resin-arm.tar.gz | tar zxvf - -C . && mv qemu-3.0.0+resin-arm/qemu-arm-static .)

	[ ! -d qemu-3.0.0+resin-arm ] || rm -fr qemu-3.0.0+resin-arm

	docker build . -t selexin/hue-im-home:arm32v7-latest -f docker/Dockerfile.arm32v7