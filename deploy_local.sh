#!/bin/bash

go build ./cmd/protoc-gen-connect-go-mpc
mv ./protoc-gen-connect-go-mpc ~/.local/bin/protoc-gen-connect-go-mpc
