package screen

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"syscall"
	"unsafe"

	"github.com/kensh1ro/willie/winapi"
)

const (
	BI_RGB       = 0
	BI_RLE8      = 1
	BI_RLE4      = 2
	BI_BITFIELDS = 3
	BI_JPEG      = 4
	BI_PNG       = 5
)

func CreateImage(rect image.Rectangle) (img *image.RGBA, e error) {
	img = nil
	defer func() {
		err := recover()
		if err == nil {
			e = nil
		}
	}()
	// image.NewRGBA may panic if rect is too large.
	img = image.NewRGBA(rect)

	return img, e
}

func Capture(x, y, width, height int) (*image.RGBA, error) {
	rect := image.Rect(0, 0, width, height)
	img, err := CreateImage(rect)
	if err != nil {
		return nil, err
	}

	hwnd := winapi.GetDesktopWindow()
	hdc := winapi.GetDC(hwnd)
	if hdc == 0 {
		return nil, errors.New("GetDC failed")
	}
	defer winapi.ReleaseDC(hwnd, hdc)

	memory_device := winapi.CreateCompatibleDC(hdc)
	if memory_device == 0 {
		return nil, errors.New("CreateCompatibleDC failed")
	}
	defer winapi.DeleteDC(memory_device)

	bitmap := winapi.CreateCompatibleBitmap(hdc, int32(width), int32(height))
	if bitmap == 0 {
		return nil, errors.New("CreateCompatibleBitmap failed")
	}
	defer winapi.DeleteObject(winapi.HGDIOBJ(bitmap))

	var header winapi.BITMAPINFOHEADER
	header.BiSize = uint32(unsafe.Sizeof(header))
	header.BiPlanes = 1
	header.BiBitCount = 32
	header.BiWidth = int32(width)
	header.BiHeight = int32(-height)
	header.BiCompression = BI_RGB
	header.BiSizeImage = 0

	// GetDIBits balks at using Go memory on some systems. The MSDN example uses
	// GlobalAlloc, so we'll do that too. See:
	// https://docs.microsoft.com/en-gb/winapi.ows/desktop/gdi/capturing-an-image
	bitmapDataSize := uintptr(((int64(width)*int64(header.BiBitCount) + 31) / 32) * 4 * int64(height))
	hmem := winapi.GlobalAlloc(winapi.GMEM_MOVEABLE, bitmapDataSize)
	defer winapi.GlobalFree(hmem)
	memptr := winapi.GlobalLock(hmem)
	defer winapi.GlobalUnlock(hmem)

	old := winapi.SelectObject(memory_device, winapi.HGDIOBJ(bitmap))
	if old == 0 {
		return nil, errors.New("SelectObject failed")
	}
	defer winapi.SelectObject(memory_device, old)

	if !winapi.BitBlt(memory_device, 0, 0, int32(width), int32(height), hdc, int32(x), int32(y), winapi.SRCCOPY) {
		return nil, errors.New("BitBlt failed")
	}

	if winapi.GetDIBits(hdc, bitmap, 0, uint32(height), (*uint8)(memptr), (*winapi.BITMAPINFO)(unsafe.Pointer(&header)), winapi.DIB_RGB_COLORS) == 0 {
		return nil, errors.New("GetDIBits failed")
	}

	i := 0
	src := uintptr(memptr)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			v0 := *(*uint8)(unsafe.Pointer(src))
			v1 := *(*uint8)(unsafe.Pointer(src + 1))
			v2 := *(*uint8)(unsafe.Pointer(src + 2))

			// BGRA => RGBA, and set A to 255
			img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = v2, v1, v0, 255

			i += 4
			src += 4
		}
	}

	return img, nil
}

func NumActiveDisplays() int {
	var count int = 0
	winapi.EnumDisplayMonitors(winapi.HDC(0), nil, syscall.NewCallback(countupMonitorCallback), uintptr(unsafe.Pointer(&count)))
	return count
}

func GetDisplayBounds(displayIndex int) image.Rectangle {
	var ctx getMonitorBoundsContext
	ctx.Index = displayIndex
	ctx.Count = 0
	winapi.EnumDisplayMonitors(winapi.HDC(0), nil, syscall.NewCallback(getMonitorBoundsCallback), uintptr(unsafe.Pointer(&ctx)))
	return image.Rect(
		int(ctx.Rect.Left), int(ctx.Rect.Top),
		int(ctx.Rect.Right), int(ctx.Rect.Bottom))
}

func countupMonitorCallback(hMonitor winapi.HMONITOR, hdcMonitor winapi.HDC, lprcMonitor *winapi.RECT, dwData uintptr) uintptr {
	var count *int
	count = (*int)(unsafe.Pointer(dwData))
	*count = *count + 1
	return uintptr(1)
}

type getMonitorBoundsContext struct {
	Index int
	Rect  winapi.RECT
	Count int
}

func getMonitorBoundsCallback(hMonitor winapi.HMONITOR, hdcMonitor winapi.HDC, lprcMonitor *winapi.RECT, dwData uintptr) uintptr {
	var ctx *getMonitorBoundsContext
	ctx = (*getMonitorBoundsContext)(unsafe.Pointer(dwData))
	if ctx.Count != ctx.Index {
		ctx.Count = ctx.Count + 1
		return uintptr(1)
	}

	if realSize := getMonitorRealSize(hMonitor); realSize != nil {
		ctx.Rect = *realSize
	} else {
		ctx.Rect = *lprcMonitor
	}

	return uintptr(0)
}

func getMonitorRealSize(hMonitor winapi.HMONITOR) *winapi.RECT {
	info := winapi.MONITORINFOEX{}
	info.CbSize = uint32(unsafe.Sizeof(info))

	winapi.GetMonitorInfo(hMonitor, &info)

	devMode := winapi.DEVMODE{}
	devMode.DmSize = uint16(unsafe.Sizeof(devMode))

	winapi.EnumDisplaySettings(&info.DeviceName[0], winapi.ENUM_CURRENT_SETTINGS, &devMode)

	return &winapi.RECT{
		Left:   devMode.DmPosition.X,
		Right:  devMode.DmPosition.X + int32(devMode.DmPelsWidth),
		Top:    devMode.DmPosition.Y,
		Bottom: devMode.DmPosition.Y + int32(devMode.DmPelsHeight),
	}
}

func ScreenShot() []byte {
	nDisplays := NumActiveDisplays()

	var height, width int = 0, 0
	for i := 0; i < nDisplays; i++ {
		rect := GetDisplayBounds(i)
		if rect.Dy() > height {
			height = rect.Dy()
		}
		width += rect.Dx()
	}
	img, err := Capture(0, 0, width, height)

	var buf bytes.Buffer
	if err != nil {
		return buf.Bytes()
	}
	//var opt jpeg.Options
	//opt.Quality = 100
	png.Encode(&buf, img)
	//jpeg.Encode(&buf, img, &opt)
	return buf.Bytes()
}
