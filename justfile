release:
  just clean
  # Linux
  GOOS=linux GOARCH=amd64 go build -v -o builds/fsize-linux-amd64
  GOOS=linux GOARCH=arm64 go build -v -o builds/fsize-linux-arm64
  GOOS=linux GOARCH=386 go build -v -o builds/fsize-linux-386
  # Windows
  GOOS=windows GOARCH=amd64 go build -v -o builds/fsize-windows-x64.exe
  GOOS=windows GOARCH=arm64 go build -v -o builds/fsize-windows-arm64.exe
  GOOS=windows GOARCH=386 go build -v -o builds/fsize-windows-x86.exe
  # Darwin
  GOOS=darwin GOARCH=amd64 go build -v -o builds/fsize-darwin-amd64
  GOOS=darwin GOARCH=arm64 go build -v -o builds/fsize-darwin-arm64
build:
  go build -v
clean:
  rm -rf builds
install:
  go install -v github.com/Tom5521/fsize@latest
