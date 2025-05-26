#!/usr/bin/env bash


lint() {
    golangci-lint run
}


lint-fix() {
    golangci-lint run --fix
}

test() {
    go test -v ./...
}

CI() {
    lint
    test
}

CI
