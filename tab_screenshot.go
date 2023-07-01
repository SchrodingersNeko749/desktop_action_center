package main

import (
	"image/png"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/kbinani/screenshot"
)

type ScreenTab struct {
	container *gtk.Box
	listbox   *gtk.ListBox
}

func (app *ScreenTab) Create() (*gtk.Box, error) {
	container, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	label, _ := gtk.LabelNew("Wifi")
	button, _ := gtk.ButtonNewWithLabel("Screenshot")
	button.Connect("clicked", func() {
		bounds := screenshot.GetDisplayBounds(0)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		file, err := os.Create("test.png")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		png.Encode(file, img)
	})
	container.Add(button)
	container.Add(label)
	return container, nil
}
