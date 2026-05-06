#!/usr/bin/env bash

flags="-X RedisScanner/cmd/commands.githash=$(git describe --tags --long --dirty)
-X RedisScanner/cmd/commands.buildstamp=$(date '+%s')"

CGO_ENABLE=0 go build -ldflags="${flags}" -o RedisScanner main.go