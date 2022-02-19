# Willie C2
### Go implant using Discord as a C2 team server
### Still in development
<p align="center">
<img src="willie.png" width="250" >
</p>

## Features

- Executing shell commands (cmd, powershell, wmic)
- Files and Directories (listing, removing...)
- Uploading and Downloading files
- Taking screenshots
- List attatched Drives
- Simple network scanner
- Processes listing
- Registry operations
- Clearing event logs
- Shellcode injection using HellsGate technique
- CLR binary injection
  
## Installation

- Create a discord Bot, and invite it to your Discord server.
- Change "TOKEN" and "CHANNEL_ID" variables in "setup_config.py" to your Bot's settings accordingly.
- run "setup_config.py".
- Compile "willie.go" by running `go build -ldflags="-s -w -H windowsgui" -trimpath`.

Everytime you run "setup_config.py" a new "config.go" file gets created with a randomized RC4 key responsible for decrypting strings to avoid triggering AV/EDR during static checks.
The program uses Mutex lock to prevent it from running more than once simultaneously.

## Demo

![Demo](demo.gif)


## TODO

- [ ] Process Migration
- [ ] Process Token Manipulation
- [ ] Priviliage Escilation

### PRs are welcome