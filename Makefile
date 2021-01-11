.PHONY: build
build:
	go build -o ./build/jwt ./cmd/jwt/main.go

run:
	go run ./cmd/jwt/main.go

start:
	./build/jwt

test:
	go test -cover -coverprofile=coverage.html -timeout 30s ./...

.PHONY: coverage
coverage:
	go tool cover -html=coverage.html
