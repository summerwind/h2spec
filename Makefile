all:
	mkdir h2spec_darwin_amd64
	GOARCH=amd64 GOOS=darwin go build cmd/h2spec/h2spec.go
	mv h2spec h2spec_darwin_amd64/h2spec
	zip h2spec_darwin_amd64.zip -r h2spec_darwin_amd64
	rm -rf h2spec_darwin_amd64
	
	mkdir h2spec_linux_amd64
	GOARCH=amd64 GOOS=linux go build cmd/h2spec/h2spec.go
	mv h2spec h2spec_linux_amd64/h2spec
	zip h2spec_linux_amd64.zip -r h2spec_linux_amd64
	rm -rf h2spec_linux_amd64
	
	mkdir h2spec_windows_amd64
	GOARCH=amd64 GOOS=windows go build cmd/h2spec/h2spec.go
	mv h2spec.exe h2spec_windows_amd64/h2spec.exe
	zip h2spec_windows_amd64.zip -r h2spec_windows_amd64
	rm -rf h2spec_windows_amd64

