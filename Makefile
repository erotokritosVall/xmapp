lint:
	golangci-lint run -c ./golangci.yml ./...

generate:
	go generate ./...

api: generate
	go build -o ./cmd/api/server ./cmd/api

clean:
	del ./cmd/api/server