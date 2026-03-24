VERSION ?= dev
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build test clean install fmt vet build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-windows-arm64 build-all

build:
	go build -ldflags="$(LDFLAGS)" -o zai-quota ./cmd/zai-quota

test:
	go test -v -cover ./...

clean:
	rm -f zai-quota
	rm -rf dist/

install: build
	cp zai-quota ~/.local/bin/

fmt:
	gofmt -s -w .

vet:
	go vet ./...

build-linux-amd64:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/zai-quota-linux-amd64 ./cmd/zai-quota

build-linux-arm64:
	mkdir -p dist
	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/zai-quota-linux-arm64 ./cmd/zai-quota

build-darwin-amd64:
	mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/zai-quota-darwin-amd64 ./cmd/zai-quota

build-darwin-arm64:
	mkdir -p dist
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/zai-quota-darwin-arm64 ./cmd/zai-quota

build-windows-amd64:
	mkdir -p dist
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/zai-quota-windows-amd64.exe ./cmd/zai-quota

build-windows-arm64:
	mkdir -p dist
	GOOS=windows GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/zai-quota-windows-arm64.exe ./cmd/zai-quota

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-windows-arm64

.PHONY: mock mock-server

mock:
	cd tests/uat/mock && go build -o mock-server .

mock-server:
	./tests/uat/mock/mock-server

.PHONY: uat uat-run

uat: mock
	./tests/uat/run_uat.sh

