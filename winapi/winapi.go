package winapi

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	GMEM_MOVEABLE                 = 0x0002
	SRCCOPY                       = 0x00CC0020
	DIB_RGB_COLORS                = 0
	CCHDEVICENAME                 = 32
	ENUM_CURRENT_SETTINGS         = 0xFFFFFFFF
	LOGON32_LOGON_NEW_CREDENTIALS = 9
	LOGON32_PROVIDER_DEFAULT      = 0
)

type RECT struct {
	Left, Top, Right, Bottom int32
}

type POINT struct {
	X, Y int32
}
type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}
type RGBQUAD struct {
	RgbBlue     byte
	RgbGreen    byte
	RgbRed      byte
	RgbReserved byte
}

type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors *RGBQUAD
}

type MONITORINFO struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
}
type MONITORINFOEX struct {
	MONITORINFO
	DeviceName [CCHDEVICENAME]uint16
}

type DEVMODE struct {
	_            [68]byte
	DmSize       uint16
	_            [6]byte
	DmPosition   POINT
	_            [86]byte
	DmPelsWidth  uint32
	DmPelsHeight uint32
	_            [40]byte
}

var (
	libuser32                   = windows.NewLazySystemDLL("user32.dll")
	libgdi32                    = windows.NewLazySystemDLL("gdi32.dll")
	libkernel32                 = windows.NewLazySystemDLL("kernel32.dll")
	libadvapi32                 = windows.NewLazySystemDLL("advapi32.dll")
	openEventLog                = libadvapi32.NewProc("OpenEventLogW")
	closeEventLog               = libadvapi32.NewProc("CloseEventLog")
	clearEventLog               = libadvapi32.NewProc("ClearEventLogW")
	getNumberOfEventLogRecords  = libadvapi32.NewProc("GetNumberOfEventLogRecords")
	getDesktopWindow            = libuser32.NewProc("GetDesktopWindow")
	enumDisplayMonitors         = libuser32.NewProc("EnumDisplayMonitors")
	getMonitorInfoW             = libuser32.NewProc("GetMonitorInfoW")
	enumDisplaySettingsW        = libuser32.NewProc("EnumDisplaySettingsW")
	getDC                       = libuser32.NewProc("GetDC")
	releaseDC                   = libuser32.NewProc("ReleaseDC")
	createCompatibleDC          = libgdi32.NewProc("CreateCompatibleDC")
	deleteDC                    = libgdi32.NewProc("DeleteDC")
	deleteObject                = libgdi32.NewProc("DeleteObject")
	createCompatibleBitmap      = libgdi32.NewProc("CreateCompatibleBitmap")
	globalAlloc                 = libkernel32.NewProc("GlobalAlloc")
	globalFree                  = libkernel32.NewProc("GlobalFree")
	globalLock                  = libkernel32.NewProc("GlobalLock")
	globalUnlock                = libkernel32.NewProc("GlobalUnlock")
	selectObject                = libgdi32.NewProc("SelectObject")
	bitBlt                      = libgdi32.NewProc("BitBlt")
	getDIBits                   = libgdi32.NewProc("GetDIBits")
	createMutex                 = libkernel32.NewProc("CreateMutexW")
	lookupPrivilegeDisplayNameW = libadvapi32.NewProc("LookupPrivilegeDisplayNameW")
	lookupPrivilegeNameW        = libadvapi32.NewProc("LookupPrivilegeNameW")
	impersonateLoggedOnUser     = libadvapi32.NewProc("ImpersonateLoggedOnUser")
	logonUserW                  = libadvapi32.NewProc("LogonUserW")
	getLogicalDriveStrings      = libkernel32.NewProc("GetLogicalDriveStringsA")
	getDriveType                = libkernel32.NewProc("GetDriveTypeA")
)

type (
	HANDLE   uintptr
	HGDIOBJ  HANDLE
	HDC      HANDLE
	HWND     HANDLE
	HBITMAP  HGDIOBJ
	HGLOBAL  HANDLE
	HMONITOR HANDLE
)

func CreateMutex(mstring *uint16) error {
	_, _, err := syscall.Syscall(createMutex.Addr(), 3,
		0,
		0,
		uintptr(unsafe.Pointer(mstring)))
	return err
}

func OpenEventLog(name *uint16) HANDLE {
	ret, _, _ := syscall.Syscall(openEventLog.Addr(), 2, uintptr(0), uintptr(unsafe.Pointer(name)),
		0)
	return HANDLE(ret)
}

func ClearEventLog(hEvent HANDLE) {
	syscall.Syscall(clearEventLog.Addr(), 2, uintptr(hEvent), uintptr(0), 0)
}

func GetNumberOfEventLogRecords(hEvent HANDLE, NumberOfRecords *uint32) {
	syscall.Syscall(getNumberOfEventLogRecords.Addr(), 2, uintptr(hEvent), uintptr(unsafe.Pointer(NumberOfRecords)), 0)
}

func CloseEventLog(hEvent HANDLE) {
	syscall.Syscall(closeEventLog.Addr(), 1, uintptr(hEvent), 0, 0)
}

func GetDC(hWnd HWND) HDC {
	ret, _, _ := syscall.Syscall(getDC.Addr(), 1,
		uintptr(hWnd),
		0,
		0)

	return HDC(ret)
}

func ReleaseDC(hWnd HWND, hDC HDC) bool {
	ret, _, _ := syscall.Syscall(releaseDC.Addr(), 2,
		uintptr(hWnd),
		uintptr(hDC),
		0)

	return ret != 0
}

func CreateCompatibleDC(hdc HDC) HDC {
	ret, _, _ := syscall.Syscall(createCompatibleDC.Addr(), 1,
		uintptr(hdc),
		0,
		0)

	return HDC(ret)
}

func DeleteDC(hdc HDC) bool {
	ret, _, _ := syscall.Syscall(deleteDC.Addr(), 1,
		uintptr(hdc),
		0,
		0)

	return ret != 0
}
func CreateCompatibleBitmap(hdc HDC, nWidth, nHeight int32) HBITMAP {
	ret, _, _ := syscall.Syscall(createCompatibleBitmap.Addr(), 3,
		uintptr(hdc),
		uintptr(nWidth),
		uintptr(nHeight))

	return HBITMAP(ret)
}

func DeleteObject(hObject HGDIOBJ) bool {
	ret, _, _ := syscall.Syscall(deleteObject.Addr(), 1,
		uintptr(hObject),
		0,
		0)
	return ret != 0
}
func GlobalAlloc(uFlags uint32, dwBytes uintptr) HGLOBAL {
	ret, _, _ := syscall.Syscall(globalAlloc.Addr(), 2,
		uintptr(uFlags),
		dwBytes,
		0)

	return HGLOBAL(ret)
}

func GlobalFree(hMem HGLOBAL) HGLOBAL {
	ret, _, _ := syscall.Syscall(globalFree.Addr(), 1,
		uintptr(hMem),
		0,
		0)

	return HGLOBAL(ret)
}
func GlobalLock(hMem HGLOBAL) unsafe.Pointer {
	ret, _, _ := syscall.Syscall(globalLock.Addr(), 1,
		uintptr(hMem),
		0,
		0)

	return unsafe.Pointer(ret)
}

func GlobalUnlock(hMem HGLOBAL) bool {
	ret, _, _ := syscall.Syscall(globalUnlock.Addr(), 1,
		uintptr(hMem),
		0,
		0)

	return ret != 0
}

func SelectObject(hdc HDC, hgdiobj HGDIOBJ) HGDIOBJ {
	ret, _, _ := syscall.Syscall(selectObject.Addr(), 2,
		uintptr(hdc),
		uintptr(hgdiobj),
		0)

	return HGDIOBJ(ret)
}
func BitBlt(hdcDest HDC, nXDest, nYDest, nWidth, nHeight int32, hdcSrc HDC, nXSrc, nYSrc int32, dwRop uint32) bool {
	ret, _, _ := syscall.Syscall9(bitBlt.Addr(), 9,
		uintptr(hdcDest),
		uintptr(nXDest),
		uintptr(nYDest),
		uintptr(nWidth),
		uintptr(nHeight),
		uintptr(hdcSrc),
		uintptr(nXSrc),
		uintptr(nYSrc),
		uintptr(dwRop))

	return ret != 0
}

func GetDIBits(hdc HDC, hbmp HBITMAP, uStartScan uint32, cScanLines uint32, lpvBits *byte, lpbi *BITMAPINFO, uUsage uint32) int32 {
	ret, _, _ := syscall.Syscall9(getDIBits.Addr(), 7,
		uintptr(hdc),
		uintptr(hbmp),
		uintptr(uStartScan),
		uintptr(cScanLines),
		uintptr(unsafe.Pointer(lpvBits)),
		uintptr(unsafe.Pointer(lpbi)),
		uintptr(uUsage),
		0,
		0)
	return int32(ret)
}

func GetDesktopWindow() HWND {
	ret, _, _ := syscall.Syscall(getDesktopWindow.Addr(), 0,
		0,
		0,
		0)

	return HWND(ret)
}

func EnumDisplayMonitors(hdc HDC, lprcClip *RECT, lpfnEnum uintptr, dwData uintptr) bool {
	ret, _, _ := syscall.Syscall6(enumDisplayMonitors.Addr(), 4,
		uintptr(hdc),
		uintptr(unsafe.Pointer(lprcClip)),
		lpfnEnum,
		dwData,
		0,
		0)
	return int(ret) != 0
}

func GetMonitorInfo(hMonitor HMONITOR, lpmi *MONITORINFOEX) bool {
	ret, _, _ := syscall.Syscall(getMonitorInfoW.Addr(), 2,
		uintptr(hMonitor),
		uintptr(unsafe.Pointer(lpmi)),
		0)

	return ret != 0
}

func EnumDisplaySettings(DeviceName *uint16, iModeNum int, DevMode *DEVMODE) bool {
	ret, _, _ := syscall.Syscall(enumDisplaySettingsW.Addr(), 3, uintptr(unsafe.Pointer(DeviceName)), uintptr(iModeNum), uintptr(unsafe.Pointer(DevMode)))
	return ret != 0
}

func ImpersonateLoggedOnUser(hToken windows.Token) bool {
	ret, _, _ := syscall.Syscall(impersonateLoggedOnUser.Addr(), 1, uintptr(hToken), 0, 0)
	return ret != 0
}

func LogonUser(lpszUsername *uint16, lpszDomain *uint16, lpszPassword *uint16, dwLogonType uint32, dwLogonProvider uint32, phToken *windows.Token) bool {
	ret, _, _ := syscall.Syscall6(logonUserW.Addr(), 6, uintptr(unsafe.Pointer(lpszUsername)), uintptr(unsafe.Pointer(lpszDomain)), uintptr(unsafe.Pointer(lpszPassword)), uintptr(dwLogonType), uintptr(dwLogonProvider), uintptr(unsafe.Pointer(phToken)))
	return ret != 0
}

func LookupPrivilegeDisplayNameW(systemName *uint16, privilegeName *uint16, buffer *uint16, size *uint32, languageId *uint32) bool {
	ret, _, _ := syscall.Syscall6(lookupPrivilegeDisplayNameW.Addr(), 5, uintptr(unsafe.Pointer(systemName)), uintptr(unsafe.Pointer(privilegeName)), uintptr(unsafe.Pointer(buffer)), uintptr(unsafe.Pointer(size)), uintptr(unsafe.Pointer(languageId)), 0)
	return ret != 0
}

func LookupPrivilegeNameW(systemName *uint16, luid *uint64, buffer *uint16, size *uint32) bool {
	ret, _, _ := syscall.Syscall6(lookupPrivilegeNameW.Addr(), 4, uintptr(unsafe.Pointer(systemName)), uintptr(unsafe.Pointer(luid)), uintptr(unsafe.Pointer(buffer)), uintptr(unsafe.Pointer(size)), 0, 0)
	return ret != 0
}

func GetLogicalDriveStrings(nBufferLegnth uint32, lpBuffer *byte) uint32 {
	ret, _, _ := syscall.Syscall(getLogicalDriveStrings.Addr(), 2, uintptr(nBufferLegnth), uintptr(unsafe.Pointer(lpBuffer)), 0)
	return uint32(ret)
}

func GetDriveType(lpRootPathName *byte) uint32 {
	ret, _, _ := syscall.Syscall(getDriveType.Addr(), 1, uintptr(unsafe.Pointer(lpRootPathName)), 0, 0)
	return uint32(ret)
}
