# RedisScanner

A command-line tool for scanning and processing Redis keys by pattern. Supports both standalone and cluster modes, with multiple built-in processors for counting, TTL analysis, memory size estimation, and safe deletion.

## Features

- **Pattern-based key scanning** — Scan Redis keys matching a glob pattern
- **Multiple processors** — Count keys, analyze TTL, estimate memory size, or delete matched keys
- **Cluster support** — Automatically iterates over all master nodes in a Redis Cluster
- **Configurable batch size** — Control scan throughput with `--batch` flag
- **YAML configuration** — Store connection details and defaults in a config file
- **Colorized output** — Colored terminal output (toggle with `--no-color`)

## Installation

### From source

```bash
git clone https://github.com/dellnoantechnp/RedisScanner.git
cd RedisScanner
go build -o redisScan
```

### Pre-built binaries

Download the latest release from the [GitHub Releases page](https://github.com/dellnoantechnp/RedisScaner/releases).

## Configuration

Generate a default configuration file:

```bash
./redisScan config create
```

This creates `~/.config/redisScan.yaml` (or `./redisScan.yaml` in the current directory) with the following defaults:

```yaml
address: "127.0.0.1"
port: 6379
password: ""
pattern: "*"
prefer_master: false
dial_timeout: 10
```

You can also override any config value via environment variables with the prefix `REDISSCAN_`:

```bash
export REDISSCAN_ADDRESS=redis.example.com
export REDISSCAN_PORT=6380
```

Or specify a custom config file path via the `--config` flag (if wired).

## Usage

### Global flags

| Flag | Default | Description |
|------|---------|-------------|
| `--batch, -b` | `300` | Scan batch size per loop (1–1000) |
| `--no-color` | `false` | Disable colored output |

### Commands

#### `count`

Count the number of keys matching the configured pattern.

```bash
./redisScan count
```

#### `ttl`

Count keys that have a TTL set (non-expiring keys are excluded).

```bash
./redisScan ttl
```

#### `memsize`

Estimate the total memory size of keys matching the pattern.

```bash
./redisScan memsize
```

#### `delete`

Delete all keys matching the pattern. **This action cannot be undone.**

```bash
# Interactive mode (prompts for confirmation)
./redisScan delete

# Force mode (skips confirmation)
./redisScan delete --force
```

#### `config`

Display or manage configuration.

```bash
# Show current configuration values
./redisScan config

# Create a default config file
./redisScan config create
```

#### `version`

Display the build version and timestamp.

```bash
./redisScan version
```

## Architecture

```
redisScan/
├── cmd/commands/        # CLI commands (Cobra)
│   ├── root.go          # Root command & global flags
│   ├── count.go         # Count processor command
│   ├── ttl.go           # TTL processor command
│   ├── memsize.go       # Memory size processor command
│   ├── delete.go        # Delete processor command
│   ├── config.go        # Config management
│   ├── configStructs.go # Config schema & defaults
│   └── version.go       # Version info
├── Processor/           # Key processors
│   ├── CountProcessor.go
│   ├── TTLProcessor.go
│   ├── SizeProcessor.go
│   └── DeleteProcessor.go
├── pkg/                 # Core logic
│   ├── KeyProcessorInterface.go  # Processor interface
│   └── RunScanner.go           # Scan engine & cluster iteration
└── utils/               # Utilities
    └── branchPrefix.go
```

### Processor Interface

New processors can be added by implementing the `KeyProcessor` interface:

```go
type KeyProcessor interface {
    Name() string
    Process(ctx context.Context, client redis.Cmdable, keys []string) error
    PrintSummary()
}
```

## Build

```bash
# Basic build
go build -o redisScan

# Build with version info (via ldflags)
bash build.sh
```

## License

Copyright 2025-2026 dellnoantechnp
