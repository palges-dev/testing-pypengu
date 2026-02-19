# PyPengu

**Linux Privilege Escalation Data Collector**

PyPengu collects privilege escalation data from Linux systems and outputs JSON for BloodPengu visualization.

## Features

- Static binary (no dependencies)
- Cross-platform (x86_64, ARM64, x86)
- 8 core collectors (sudo, SUID, docker, services, cron, kernel, groups, users)
- JSON output compatible with BloodPengu

## Quick Start

```bash
go build -o pypengu cmd/pypengu/main.go
./pypengu -o output.json
```

## Usage

```bash
./pypengu [options]

Options:
  -o, --output    Output JSON file (default: pypengu-output.json)
  -v, --verbose   Verbose output
  -h, --help      Show help
```

## Build for Multiple Platforms

```bash
GOOS=linux GOARCH=amd64 go build -o pypengu-linux-amd64 cmd/pypengu/main.go
```

```bash
GOOS=linux GOARCH=arm64 go build -o pypengu-linux-arm64 cmd/pypengu/main.go
```

```bash
GOOS=linux GOARCH=386   go build -o pypengu-linux-386 cmd/pypengu/main.go
```

## Project Structure

```
PyPengu/
├── cmd/pypengu/          # Main application entry point
├── internal/
│   ├── collector/        # Data collectors (sudo, SUID, etc.)
│   ├── graph/            # Graph builder (nodes + edges)
│   ├── models/           # Data structures
│   └── utils/            # Helper functions
└── pkg/output/           # JSON output formatter
```

## Collectors

- **sudo** - Parse /etc/sudoers and sudo -l
- **suid** - Find SUID/SGID binaries
- **docker** - Check docker group membership and socket access
- **service** - Find writable systemd services
- **cron** - Find writable cron jobs
- **kernel** - Detect kernel version and CVEs
- **group** - Enumerate group memberships
- **user** - Enumerate users and permissions

## License

MIT License
