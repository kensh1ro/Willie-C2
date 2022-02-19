package hgsyscall

import (
	"encoding/binary"
	"errors"

	"github.com/Binject/debug/pe"
)

//rvaToOffset converts an RVA value from a PE file into the file offset. When using binject/debug, this should work fine even with in-memory files.
func rvaToOffset(pefile *pe.File, rva uint32) uint32 {
	for _, hdr := range pefile.Sections {
		baseoffset := uint64(rva)
		if baseoffset > uint64(hdr.VirtualAddress) &&
			baseoffset < uint64(hdr.VirtualAddress+hdr.VirtualSize) {
			return rva - hdr.VirtualAddress + hdr.Offset
		}
	}
	return rva
}

//getSysIDFromMemory takes values to resolve, and resolves from disk.
func getSysIDFromDisk(funcname string, ord uint32, useOrd bool) (uint16, error) {
	l := "C:\\Windows\\System32\\ntdll.dll"
	p, e := pe.Open(l)

	if e != nil {
		return 0, e
	}

	ex, _ := p.Exports()
	for _, exp := range ex {
		if (useOrd && exp.Ordinal == ord) || // many bothans died for this feature
			exp.Name == funcname {
			offset := rvaToOffset(p, exp.VirtualAddress)
			b, e := p.Bytes()
			if e != nil {
				return 0, e
			}
			buff := b[offset : offset+10]

			return sysIDFromRawBytes(buff), nil
		}
	}
	return 0, errors.New("could not find syscall ID")
}

func sysIDFromRawBytes(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b[4:8])
}

/*func WriteMemory(inbuf []byte, destination uintptr) {
	for index := uint32(0); index < uint32(len(inbuf)); index++ {
		writePtr := unsafe.Pointer(destination + uintptr(index))
		v := (*byte)(writePtr)
		*v = inbuf[index]
	}
}*/
