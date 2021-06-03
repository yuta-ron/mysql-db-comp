build-linux-amd64: 
	GOARCH="amd64" GOOS="linux" go build -v -o ./dbcmp main.go
build-darwin-amd64: 
	GOARCH="amd64" GOOS="darwin" go build -v -o ./dbcmp main.go
.PHONY: install
install:
	chmod +x dbcmp
	sudo mv dbcmp /usr/local/bin