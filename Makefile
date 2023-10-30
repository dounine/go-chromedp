APP="autodump"
.PHONY: clean build run
PID=/tmp/.$(APP)-server.pid

build: clean
	CGO_ENABLED=1 go build -o $(APP) -ldflags '-linkmode "external" -extldflags "-static"' main.go
run: clean
	go run main.go
clean:
	@go clean
