package main

import (
	"fmt"
	"net/url"
	"os"

	wallpaper "github.com/ErogeMaster225/wallpaper-go/src"
)

func printHelp() {
	fmt.Printf(`
Wallpaper-go is a utility for setting the desktop wallpaper on Windows.
You can give the local filepath or any image url.
Usage:
	wallpaper-go -f C:\Users\your_user_name\Pictures\made-a-botw-vector-wallpaper-4k-2560×1600.jpg
	wallpaper-go -i https://i.redd.it/l1764nd9h3721.jpg
	-f [filepath] local filepath.
	-i [url] download image from internet and set it as wallpaper
	--help, -h for help.
	`)
}
func main() {
	var err error
	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "-i":
			_, err = url.ParseRequestURI(args[2])
			if err == nil {
				err = wallpaper.SetWallpaperUrl(args[2])
			}
		case "-f", "--file":
			err = wallpaper.SetLocalWallpaper(args[2])
		case "--help", "-h":
			printHelp()
		default:
			fmt.Println("Invalid argument!")
			printHelp()
		}
		if err != nil {
			fmt.Printf("\n\033[31mCould not set wallpaper: %s\n\033[0m", err.Error())
		}
	} else {
		printHelp()
	}
}
