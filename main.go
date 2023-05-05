package main

import (
	"log"
)

const WINDOW_WIDTH = 700

func main() {

	ui := NewUI()

	if err := ui.Init(); err != nil {
		log.Fatal("Unable to initialize UI:", err)
	}

	ui.Run()
}
