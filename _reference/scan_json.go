package main

import (
	"./spotlight"
)

func main() {
	metadata := new(spotlight.MetaData)
	metadata.ImportAll()
}
