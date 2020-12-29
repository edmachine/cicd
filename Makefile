# This is the makefile for the cicd project

run: mod ; $(info Running cicd)
	go run cmd/cicd/main.go

build: mod ; $(info Running cicd)
	go build -o cicd cmd/cicd/main.go

test: mod ; $(info Unit Testing Project)
	go test -cover -v -race ./...

testcoverage: mod ; $(info Unit Test Coverage HTML)
	go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

clean: ; $(info Cleaning build artifacts)
	rm coverage.out cicd

mod: setup ; $(info Running setup commands)

setup: ; $(info Getting development tools)

.PHONY: server

