#!/usr/bin/env bash

go install golang.org/x/vuln/cmd/govulncheck@latest

govulncheck ./...