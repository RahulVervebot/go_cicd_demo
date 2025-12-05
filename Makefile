.PHONY: tidy fmt vet test ci run

tidy:
	go mod tidy

fmt:
	gofmt -w .

vet:
	go vet ./...

test:
	go test ./... -race

ci: tidy vet test

run:
	go run ./cmd/go_cicd_demo 2 3
