APP="chromedp"
.PHONY: clean build run
PID=/tmp/.$(APP)-server.pid

build: clean
	CGO_ENABLED=1 go build -o $(APP) -ldflags '-linkmode "external" -extldflags "-static"' main.go
build-win: clean
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="x86_64-w64-mingw32-gcc" go build -o $(APP).exe -ldflags '-linkmode "external" -extldflags "-static"' main.go
run: clean
	go run main.go
clean:
	@go clean
