lint:
	golangci-lint run -c ./golangci.yml ./...

generate:
	go generate ./...

api: generate
	go build -o server ./cmd/api

clean:
	del server