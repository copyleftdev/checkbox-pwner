# Define the binary name
BINARY_NAME = checkbox-pwner
PACKAGE_NAME = github.com/copyleftdev/checkbox-pwner

# Go build parameters
BUILD_FLAGS = -ldflags="-s -w"

# Default target
all: build

# Target to install dependencies
deps:
	@echo "==> Installing dependencies..."
	@go get -v $(PACKAGE_NAME)

# Target to build for Linux/OSX
build: deps
	@echo "==> Building for $(GOOS)/$(GOARCH)..."
	@go build $(BUILD_FLAGS) -o $(BINARY_NAME)

# Target to build for Windows
build-windows: deps
	@echo "==> Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME).exe

# Target to run the application
run:
	@echo "==> Running application..."
	@./$(BINARY_NAME)

# Target to clean the build
clean:
	@echo "==> Cleaning build..."
	@rm -f $(BINARY_NAME) $(BINARY_NAME).exe

# Target to build and run for Linux/OSX
build-run: build run

# Target to build and run for Windows
build-run-windows: build-windows
	@./$(BINARY_NAME).exe

.PHONY: all deps build build-windows run clean build-run build-run-windows
