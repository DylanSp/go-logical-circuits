[private]
default:
    @just --list

build:
    go build ./...

run:
    go run main.go

test:
    go test ./...

fuzz: fuzz-not fuzz-and fuzz-full-adder

fuzz-not:
    go test -fuzz=FuzzNot32 -fuzztime 10s "github.com/DylanSp/go-logical-circuits/circuit"

fuzz-and:
    go test -fuzz=FuzzAnd32 -fuzztime 10s "github.com/DylanSp/go-logical-circuits/circuit"

fuzz-full-adder:
    go test -fuzz=FuzzFullAdder32 -fuzztime 10s "github.com/DylanSp/go-logical-circuits/circuit"