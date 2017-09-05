all: clean build run

build: main.go
	go build -o not-fancy-telnet

run: not-fancy-telnet
	./not-fancy-telnet

clean: 
	rm -rf ./not-fancy-telnet
	rm -rf ./server.log
