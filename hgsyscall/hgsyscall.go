package hgsyscall

func Syscall(callid uint16, argh ...uintptr) uint32 {
	errcode := hgSyscall(callid, argh...)
	return errcode
}

func hgSyscall(callid uint16, argh ...uintptr) (errcode uint32)

func GetSysIDFromDisk(funcname string) (uint16, error) {
	return getSysIDFromDisk(funcname, 0, false)
}
