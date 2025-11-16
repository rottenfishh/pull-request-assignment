all: tidy build run
tidy:
	go mod tidy
build:
	go build -o main ./cmd
run:
	go run ./cmd