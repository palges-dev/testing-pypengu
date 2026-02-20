# PyPengu

![GitHub all releases](https://img.shields.io/github/downloads/pengu-apm/PyPengu/total)

PyPengu is the `default attack-path collector` for BloodPengu.

PyPengu is the attack-path creator to collects privilege escalation data from Linux systems and outputs JSON for BloodPengu to digest and became visualization.

## Get PyPengu

The latest Project of PyPengu is in this repo ⚙️. 

The Build is personal. 

## Maintainer

PyPengu is considered to be Properties under [AdverXarial](https://byt3n33dl3.github.io/) and well-maintained by [@byt3n33dl3.](https://github.com/byt3n33dl3/)

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

## Custom Build

```bash
GOOS=linux GOARCH=amd64 go build -o pypengu-linux-amd64 cmd/pypengu/main.go
```

```bash
GOOS=linux GOARCH=arm64 go build -o pypengu-linux-arm64 cmd/pypengu/main.go
```

```bash
GOOS=linux GOARCH=386 go build -o pypengu-linux-386 cmd/pypengu/main.go
```

## Documentation

- [PyPengu Documentation](https://pengu-apm.github.io/)

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

```
MIT License

Copyright (c) 2026 byt3n33dl3 && AdverXarial

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

And other LICENSES.

## CONTACT

For more, come to the documentation for use cases and write-ups [here](https://pengu-apm.github.io), if there's any security concern, please contact me at <byt3n33dl3@pm.me>