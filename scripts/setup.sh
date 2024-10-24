#!/bin/bash
go mod tidy
go build -o translators cmd/translators/main.go