#!/usr/bin/env bash

go test -v -race ./... -coverprofile=coverage.txt -covermode=atomic -timeout=10m

go tool cover -html=coverage.txt -o coverage.html
