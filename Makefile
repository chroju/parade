BINARY_NAME=parade

.PHONY: install test lint build crossbuild clean mod test-coverage

install:
	go install

lint:
	go mod tidy
	gofmt -s -l .
	golint ./...
	go vet ./...

test: lint
	go test -v ./...

build:
	go build -o bin/$(BINARY_NAME)

crossbuild: test
	gox -os="linux darwin windows" -arch="386 amd64" -output "bin/remo_{{.OS}}_{{.Arch}}/{{.Dir}}"

mod:
	go mod download

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f bin/

test-coverage: mod
	go test -race -covermode atomic -coverprofile=covprofile ./...
