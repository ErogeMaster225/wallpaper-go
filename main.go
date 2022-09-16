package main

import (
	"fmt"
	"os"

	wallpaper "github.com/ErogeMaster225/wallpaper-go/src"
)

func main() {
	var err error
	fmt.Println("Current wallpaper:", wallpaper.GetWallpaper())
	args := os.Args
	if len(args) > 1 {
		err = wallpaper.SetWallpaper(args[1])
	}
	if err != nil {
		fmt.Printf("\nCould not set wallpaper: %s\n", err.Error())
	}
}
