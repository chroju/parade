BINARY_NAME=parade

.PHONY: install test lint crossbuild clean

install:
	go install

lint:
	go mod tidy
	gofmt -s -l .
	golint ./...
	go vet ./...

test: lint
	go test -v ./...

crossbuild: test
	gox -os="linux darwin windows" -arch="386 amd64" -output "bin/remo_{{.OS}}_{{.Arch}}/{{.Dir}}"

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f bin/