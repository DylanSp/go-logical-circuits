[private]
default:
    @just --list

build:
    go build ./...

run:
    go run main.go
