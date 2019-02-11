VERSION=2.2.0
COMMIT=$(shell git rev-parse --verify HEAD)

PACKAGES=$(shell go list ./... | grep -v /vendor/)
BUILD_FLAGS=-ldflags "-X main.VERSION=$(VERSION) -X main.COMMIT=$(COMMIT)"

.PHONY: all
all: build

.PHONY: build
build: vendor
	go build $(BUILD_FLAGS) cmd/h2spec/h2spec.go

.PHONY: test
test:
	go test -v $(PACKAGES)
	go vet $(PACKAGES)

.PHONY: clean
clean:
	rm -rf h2spec
	rm -rf release

.PHONY: container
container:
	GOARCH=amd64 GOOS=linux go build $(BUILD_FLAGS) cmd/h2spec/h2spec.go
	docker build -t summerwind/h2spec:latest -t summerwind/h2spec:$(VERSION) .
	rm -rf h2spec

release:
	mkdir -p release
	
	GOARCH=amd64 GOOS=darwin go build $(BUILD_FLAGS) cmd/h2spec/h2spec.go
	tar -czf release/h2spec_darwin_amd64.tar.gz h2spec
	rm -rf h2spec
	
	GOARCH=amd64 GOOS=windows go build $(BUILD_FLAGS) cmd/h2spec/h2spec.go
	zip release/h2spec_windows_amd64.zip -r h2spec.exe
	rm -rf h2spec.exe
	
	GOARCH=amd64 GOOS=linux go build $(BUILD_FLAGS) cmd/h2spec/h2spec.go
	tar -czf release/h2spec_linux_amd64.tar.gz h2spec
	rm -rf h2spec

vendor:
	dep ensure -v
