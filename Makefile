all: build run
build:
	go build -o main.exe .
run:
	main.exe