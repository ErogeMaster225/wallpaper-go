package wallpaper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/schollz/progressbar/v3"
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
func DownloadWallpaper(fileurl string) (string, error) {
	res, err := http.Get(fileurl)
	if err != nil {
		return "", err
	}
	contentType := res.Header.Get("Content-Type")
	if (contentType != "image/jpeg") &&
		(contentType != "image/png") &&
		(contentType != "image/bmp") {
		return "", errors.New("this file type is not supported")
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return "", errors.New(fmt.Sprint(res.StatusCode, " ", http.StatusText(res.StatusCode)))
	}
	userDir, _ := os.UserHomeDir()
	fileLocation := path.Join(userDir, "Pictures", filepath.Base(fileurl))
	file, _ := os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	bar := progressbar.NewOptions64(res.ContentLength-1,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetDescription("[cyan]Downloading[reset]"),
	)
	io.Copy(io.MultiWriter(file, bar), res.Body)
	bar.Describe("[cyan]Finished[reset]")
	fmt.Println("\n\033[32mImage saved to", fileLocation, string("\033[0m"))
	return fileLocation, nil
}
func SetLocalWallpaper(filename string) error {
	filenameUTF16, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
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
func SetWallpaperUrl(url string) error {
	filename, err := DownloadWallpaper(url)
	if err != nil {
		return err
	}
	return SetLocalWallpaper(filename)
}
