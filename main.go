package main

import (
	"log"
)

const WINDOW_WIDTH = 700

func main() {

	actionCenter := NewActionCenter()
	notificationCenter := NewNotificationCenter(actionCenter)
	if err := actionCenter.Init(); err != nil {
		log.Fatal("Unable to initialize UI:", err)
	}

	go notificationCenter.Run()
	actionCenter.Run()
}
