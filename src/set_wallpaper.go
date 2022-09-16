package wallpaper

import (
	"errors"
	"os"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// Set and get wallpaper using SystemParametersInfoW function from user32.dll
var (
	user32                = syscall.NewLazyDLL("user32.dll")
	systemParametersInfoW = user32.NewProc("SystemParametersInfoW")
)

// SystemParametersInfoW's parameters. More info: https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-systemparametersinfow
const (
	spiGetDeskWallpaper = 0x0073
	spiSetDeskWallpaper = 0x0014
	uiParam             = 0x0000
	spifUpdateIniFile   = 0x01
	spifSendChange      = 0x02
)

func GetWallpaper() string {
	var filename [256]uint16
	systemParametersInfoW.Call(
		uintptr(spiGetDeskWallpaper),
		uintptr(cap(filename)),
		uintptr(unsafe.Pointer(&filename[0])),
		uintptr(0),
	)
	return string(utf16.Decode(filename[:]))
}
func SetWallpaper(filename string) error {
	filenameUTF16, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}
	_, err1 := os.Stat(filename)
	if os.IsNotExist(err1) {
		return errors.New("file does not exists")
	}
	systemParametersInfoW.Call(
		uintptr(spiSetDeskWallpaper),
		uintptr(uiParam),
		uintptr(unsafe.Pointer(filenameUTF16)),
		uintptr(spifUpdateIniFile|spifSendChange),
	)
	return nil
}
