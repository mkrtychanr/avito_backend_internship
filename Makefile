build: clean
	go build -v ./cmd/server

clean:
	rm -rf server