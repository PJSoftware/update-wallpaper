package ui

import (
	"bufio"
	"fmt"
	"os"
)

func Pause() {
	fmt.Printf("Press [ENTER] to continue:\n")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}