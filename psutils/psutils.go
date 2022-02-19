package psutils

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Process interface {
	Pid() int
	PPid() int
	Executable() string
	Owner() string
	Arch() string
}

type WindowsProcess struct {
	pid   int
	ppid  int
	exe   string
	owner string
	arch  string
}

func (p *WindowsProcess) Pid() int {
	return p.pid
}

func (p *WindowsProcess) PPid() int {
	return p.ppid
}

func (p *WindowsProcess) Executable() string {
	return p.exe
}

func (p *WindowsProcess) Owner() string {
	return p.owner
}

func (p *WindowsProcess) Arch() string {
	return p.arch
}

func newWindowsProcess(e *syscall.ProcessEntry32) *WindowsProcess {
	// Find when the string ends for decoding
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}
	account, _ := getProcessOwner(e.ProcessID)

	pHandle, _ := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, e.ProcessID)
	defer syscall.CloseHandle(pHandle)
	isWow64Process, err := IsWow64Process(pHandle)

	arch := "x86"
	if !isWow64Process {
		arch = "x64"
	}
	if err != nil {
		arch = "err"
	}

	return &WindowsProcess{
		pid:   int(e.ProcessID),
		ppid:  int(e.ParentProcessID),
		exe:   syscall.UTF16ToString(e.ExeFile[:end]),
		owner: account,
		arch:  arch,
	}
}

// getInfo retrieves a specified type of information about an access token.
func getInfo(t syscall.Token, class uint32, initSize int) (unsafe.Pointer, error) {
	n := uint32(initSize)
	for {
		b := make([]byte, n)
		e := syscall.GetTokenInformation(t, class, &b[0], uint32(len(b)), &n)
		if e == nil {
			return unsafe.Pointer(&b[0]), nil
		}
		if e != syscall.ERROR_INSUFFICIENT_BUFFER {
			return nil, e
		}
		if n <= uint32(len(b)) {
			return nil, e
		}
	}
}

// getTokenOwner retrieves access token t owner account information.
func getTokenOwner(t syscall.Token) (*syscall.Tokenuser, error) {
	i, e := getInfo(t, syscall.TokenOwner, 50)
	if e != nil {
		return nil, e
	}
	return (*syscall.Tokenuser)(i), nil
}

func getProcessOwner(pid uint32) (owner string, err error) {
	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, pid)
	if err != nil {
		return
	}
	defer syscall.CloseHandle(handle)
	var token syscall.Token
	if err = syscall.OpenProcessToken(handle, syscall.TOKEN_QUERY, &token); err != nil {
		return
	}
	tokenUser, err := getTokenOwner(token)
	if err != nil {
		return
	}
	owner, domain, _, err := tokenUser.User.Sid.LookupAccount("")
	owner = fmt.Sprintf("%s\\%s", domain, owner)
	return
}

// IsWow64Process determines the process architecture
// https://github.com/shenwei356/rush/blob/master/process/process_windows.go
func IsWow64Process(processHandle syscall.Handle) (bool, error) {
	var wow64Process bool
	kernel32 := windows.NewLazySystemDLL("kernel32")
	procIsWow64Process := kernel32.NewProc("IsWow64Process")

	r1, _, e1 := procIsWow64Process.Call(
		uintptr(processHandle),
		uintptr(unsafe.Pointer(&wow64Process)))
	if int(r1) == 0 {
		return false, e1
	}
	return wow64Process, nil
}

func Processes() ([]Process, error) {
	handle, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(handle)

	var entry syscall.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	if err = syscall.Process32First(handle, &entry); err != nil {
		return nil, err
	}

	procs := make([]Process, 0, 50)
	for {
		procs = append(procs, newWindowsProcess(&entry))

		err = syscall.Process32Next(handle, &entry)
		if err != nil {
			break
		}
	}

	return procs, nil
}

func PS() string {
	var result strings.Builder
	processList, err := Processes()
	if err != nil {
		return err.Error()
	}

	result.WriteString("\nPID\tPPID\tARCH\tOWNER\tNAME\n")
	for p := range processList {
		process := processList[p]
		result.WriteString(fmt.Sprintf("%d\t%d\t%s\t%s\t%s\n", process.Pid(), process.PPid(), process.Arch(), process.Owner(), process.Executable()))
	}
	return result.String()
}

func Kill(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}
