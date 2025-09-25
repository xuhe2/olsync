EXE_NAME="olsync"
MAIN_FILE="./cmd/main.go"

build:
	go build -o ./bin/$(EXE_NAME)-$(shell go env GOOS)-$(shell go env GOARCH) $(MAIN_FILE)

test:
	go test -v ./...

fmt:
	gofmt -w $(shell find . -name "*.go" -type f) 

# 检查代码中可能存在的错误和可疑构造
vet:
	go vet ./...

# golangci-lint 是一个功能更强大的第三方代码检查工具
lint:
	golangci-lint run