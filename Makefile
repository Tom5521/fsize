
release:
	make clean
	# Linux
	GOOS=linux GOARCH=amd64 go build -v -o builds/fsize-linux-amd64
	GOOS=linux GOARCH=arm64 go build -v -o builds/fsize-linux-arm64
	# Windows
	GOOS=windows GOARCH=amd64 go build -v -o builds/fsize-windows-x64.exe
	GOOS=windows GOARCH=arm64 go build -v -o builds/fsize-windows-arm64.exe
	# Darwin
	GOOS=darwin GOARCH=amd64 go build -v -o builds/fsize-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -v -o builds/fsize-darwin-arm64
build:
	go build -v
clean:
	rm -rf builds
