package main

import (
	"fmt"

	"github.com/pjsoftware/win-spotlight/vc"
)

func main() {
	s := vc.Detect(".")
	fmt.Printf("Detected: %s\n", s.Detected)
}
