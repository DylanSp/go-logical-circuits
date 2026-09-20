[private]
default:
    @just --list

build:
    go build ./...

run:
    go run main.go

test:
    go test ./...

fuzz-circuit:
    go test -fuzz=Fuzz "github.com/DylanSp/go-logical-circuits/circuit"