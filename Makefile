VERSION=2.6.0
COMMIT=$(shell git rev-parse --verify HEAD)
BUILD_FLAGS=-ldflags "-X main.VERSION=$(VERSION) -X main.COMMIT=$(COMMIT)"

all: build

build:
	go build $(BUILD_FLAGS) cmd/h2spec/h2spec.go

test:
	go vet ./...
	go test -v ./...

clean:
	rm -rf h2spec release

build-container:
	docker build --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) -t summerwind/h2spec:latest -t summerwind/h2spec:$(VERSION) .

push-container:
	docker push summerwind/h2spec:latest

push-release-container:
	docker push summerwind/h2spec:$(VERSION)

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

github-release: release
	ghr v$(VERSION) release/
