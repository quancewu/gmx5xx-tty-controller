buildarm:
	@env GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ CGO_ENABLED=1 go build -a -ldflags '-extldflags "-static"'  -o bin/gmx5xx-tty-controller

build:
	@go build -a -ldflags '-extldflags "-static"'  -o bin/gmx5xx-tty-controller

run: build
	@./bin/gmx5xx-tty-controller

test:
	@go test -v ./...