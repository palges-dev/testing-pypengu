# PyPengu

![GitHub all releases](https://img.shields.io/github/downloads/pengu-apm/PyPengu/total)

PyPengu is the `default attack-path collector` for BloodPengu.

PyPengu is the attack-path creator to collects privilege escalation data from Linux systems and outputs JSON for BloodPengu to digest and became visualization.

## Get PyPengu

The latest Project of PyPengu is in this repo ⚙️. 

The Build is personal. 

## Features

- Static binary with no dependencies
- Cross-platform (x86_64, ARM64, x86)
- 8 core collectors (sudo, SUID, docker, services, cron, kernel, groups, users)
- JSON output compatible with BloodPengu

## Quick Start

```
┌──(kali㉿kali)-[~]
└─$ go build -o pypengu ./cmd/pypenguo build -o pypengu cmd/pypengu/main.go
```

## Build Cross Multiple Platforms

```bash
GOOS=linux GOARCH=amd64 go build -o pypengu-linux-amd64 cmd/pypengu/main.go
```

```bash
GOOS=linux GOARCH=arm64 go build -o pypengu-linux-arm64 cmd/pypengu/main.go
```

```bash
GOOS=linux GOARCH=386 go build -o pypengu-linux-386 cmd/pypengu/main.go
```

## CLI Arguments

```
┌──(kali㉿kali)-[~]
└─$ ./pypengu -h                     
Usage of ./pypengu:
  -o string   Output JSON file (default "pypengu-output.json")
  -v Verbose  output
  -version    Show version
```

## License

MIT License
