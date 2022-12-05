package main

import (
	"fmt"

	"github.com/pjsoftware/update-wallpaper/pkg/vc"
)

func main() {
	s := vc.Detect(".")
	fmt.Printf("Detected: %s\n", s.Detected)
}
