BINARY=coxmos
BUILD_DIR=.build

.PHONY: build install clean

build:
	go build -o $(BUILD_DIR)/$(BINARY) .

install:
	go build -o $(shell go env GOPATH)/bin/$(BINARY) .

clean:
	rm -rf $(BUILD_DIR)

# cross-compile for common platforms
release:
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-linux-amd64 .
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY)-linux-arm64 .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY)-darwin-arm64 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe .
