BINARY_NAME=koda

build:
	go build -o dist/$(BINARY_NAME) ./cmd

clean:
	go clean
	rm -f dist/$(BINARY_NAME)

test:
	go test -v ./... --count=1

prep:
	betteralign --apply --test_files ./...
