# 변수 정의
BINARY_NAME=renamatic
BUILD_DIR=build
CMD_DIR=cmd/renamatic

.PHONY: all build clean install test

all: build

build:
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "✅ Build $(BUILD_DIR)/$(BINARY_NAME)"

install:
	@echo "📦 Install renamatic"
	@go install ./$(CMD_DIR)
	@echo "✅ Install renamatic"

clean:
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean build files"

test:
	@go test -v ./...
	@echo "✅ Run tests"
