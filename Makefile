PACKAGES=$(shell go list ./... | grep -v /vendor/)

.PHONY: all
all: build

.PHONY: build
build:
	go build cmd/h2spec/h2spec.go

.PHONY: test
test:
	go test -v $(PACKAGES)
	go vet $(PACKAGES)

.PHONY: clean
clean:
	rm -rf h2spec
	rm -rf release

.PHONY: release
release:
	mkdir -p release
	
	GOARCH=amd64 GOOS=darwin go build cmd/h2spec/h2spec.go
	zip release/h2spec_darwin_amd64.zip -r h2spec
	rm -rf h2spec
	
	GOARCH=amd64 GOOS=windows go build cmd/h2spec/h2spec.go
	zip release/h2spec_windows_amd64.zip -r h2spec.exe
	rm -rf h2spec.exe
	
	GOARCH=amd64 GOOS=linux go build cmd/h2spec/h2spec.go
	tar -czf release/h2spec_linux_amd64.tar.gz h2spec
	rm -rf h2spec

