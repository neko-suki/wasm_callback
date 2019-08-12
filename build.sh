#!/bin/bash
GOOS=js GOARCH=wasm go build -o callback.wasm callback.go
curl -sO https://raw.githubusercontent.com/golang/go/go1.12.7/misc/wasm/wasm_exec.js
