.PHONY: run build test clean dev

# Application variables
APP_NAME=image-identify
BUILD_DIR=build
TMP_DIR=tmp

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean

# Air variables
AIR_VERSION=v1.49.0
AIR_BINARY=$(TMP_DIR)/air

# Build flags
LDFLAGS=-ldflags "-w -s"

all: clean build

build:
	@echo "Building application..."
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./main.go

run:
	@echo "Running application..."
	$(GOCMD) run main.go

dev: ensure-air
	@echo "Running in development mode..."
	$(AIR_BINARY)

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(TMP_DIR)

deps:
	@echo "Installing dependencies..."
	$(GOCMD) mod tidy
	$(GOCMD) mod verify

ensure-air: $(AIR_BINARY)

$(AIR_BINARY):
	@echo "Downloading air binary..."
	@mkdir -p $(TMP_DIR)
	@curl -fsSL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | BINDIR=$(TMP_DIR) sh -s $(AIR_VERSION)