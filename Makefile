all: build

build:
	go build .

tidy:
	go mod tidy

clean:
	go clean

run: build
	@echo Running with awscli, run \'./superkey_cleaner -access xxx -secret yyy\' to clean a specific account
	./superkey_cleaner -cli

lint:
	go vet ./...
	golangci-lint run -E gofmt,gci,bodyclose,forcetypeassert,misspell

gci:
	golangci-lint run -E gci --fix

.PHONY: build run clean debug lint
