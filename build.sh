#!/usr/bin/env bash

flags="-X RedisScanTask/cmd/commands.githash=$(git describe --tags --long --dirty)
-X RedisScanTask/cmd/commands.builtstamp=$(date '+%s')"

CGO_ENABLE=0 go build -ldflags="${flags}" -o RedisScanTask main.go