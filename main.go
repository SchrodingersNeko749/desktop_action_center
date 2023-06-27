package main

import (
	"github.com/actionCenter/Data"
	"github.com/actionCenter/Service"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)
	Data.LoadConfig()
	Service.NewActionCenter().Init()
}
