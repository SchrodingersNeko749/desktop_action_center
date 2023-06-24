package main

import (
	"log"

	"github.com/actionCenter/Service"
)

func main() {

	actionCenter := Service.NewActionCenter()

	if err := actionCenter.Init(); err != nil {
		log.Fatal("Unable to initialize action center:", err)
	}
	actionCenter.Run()
}
