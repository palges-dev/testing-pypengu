# PyPengu Usage

```
./pypengu [options]

Options:
  -o, --output    Output JSON file (default: pypengu-output.json)
  -v, --verbose   Verbose output
  -h, --help      Show help
```

Real life attack paths-collector

```
┌──(kali㉿kali)-[~/Documents/Project/pengu-project]
└─$ sudo ./pypengu -v
PyPengu: Linux PrivEsc Data Collector
Version: 1.0.0
Collecting data...

[*] Collecting users...
    Found 59 users
[*] Collecting groups...
    Found 93 groups
[*] Collecting sudo configuration...
    Found 9 sudo entries
[*] Scanning for SUID binaries...
    Found 34 SUID binaries
[*] Checking Docker access...
    Docker socket is writable
[*] Scanning services...
    Found 655 writable services
[*] Scanning cron jobs...
    Found 9 cron jobs
[*] Checking kernel version...
    Kernel: 6.12.25-amd64
[*] Building graph...

[+] Data collected successfully
[+] Output: pypengu-output.json
[+] Nodes: 152 | Edges: 23 | Paths to root: 0
```

## CONTACT

For more, come to the documentation for use cases and write-ups [here](https://pengu-apm.github.io), if there's any security concern, please contact me at <byt3n33dl3@pm.me>