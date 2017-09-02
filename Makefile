all: build run

build: main.go
	go build -o not-fancy-telnet

run: not-fancy-telnet
	./not-fancy-telnet
