package main

import (
	"log"

	"github.com/actionCenter/Data"
	"github.com/actionCenter/Service"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)
	Data.LoadConfig()
	actionCenter := Service.NewActionCenter()

	if err := actionCenter.Init(); err != nil {
		log.Fatal("Unable to initialize action center:", err)
	}
	actionCenter.Run()
}
