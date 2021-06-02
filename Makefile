build-linux: 
	GOARCH="amd64" GOOS="linux" go build -v -o ./dbcmp main.go
.PHONY: install
install:
	chmod +x dbcmp
	sudo mv dbcmp /usr/local/bin