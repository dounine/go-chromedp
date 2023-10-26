APP="autodump"
.PHONY: clean build run
PID=/tmp/.$(APP)-server.pid

build: clean
	go build -o $(APP) main.go
run: clean
	go run main.go
clean:
	@go clean
