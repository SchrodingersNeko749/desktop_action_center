package main

import (
	"image/png"
	"os"

	"github.com/gotk3/gotk3/cairo"
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
	fullscreenScreenshotButton, _ := gtk.ButtonNewWithLabel("Screenshot")
	fullscreenScreenshotButton.Connect("clicked", func() {
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

	regionScreenshotButton, _ := gtk.ButtonNewWithLabel("Region Screenshot")
	regionScreenshotButton.Connect("clicked", func() {
		bounds := screenshot.GetDisplayBounds(0)

		selectionWindow, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

		visual, _ := selectionWindow.GetScreen().GetRGBAVisual()
		var alphaSupported bool
		selectionWindow.SetSizeRequest(bounds.Size().X, bounds.Size().Y)
		selectionWindow.SetDecorated(false)
		selectionWindow.SetPosition(gtk.WIN_POS_CENTER)
		selectionWindow.SetKeepAbove(true)
		selectionWindow.SetAppPaintable(true)
		if visual != nil {
			alphaSupported = true
		} else {
			println("Alpha not supported")
			alphaSupported = false
		}
		selectionWindow.SetVisual(visual)
		sCtx, _ := selectionWindow.GetStyleContext()
		sCtx.AddClass("selection-window")
		selectionWindow.Connect("draw", func(window *gtk.Window, ctx *cairo.Context) {
			if alphaSupported {
				ctx.SetSourceRGBA(0.0, 0.0, 0.0, 0.25)
			} else {
				ctx.SetSourceRGB(0.0, 0.0, 0.0)
			}

			ctx.SetOperator(cairo.OPERATOR_SOURCE)
			ctx.Paint()
		})
		selectionWindow.ShowAll()

		// img, err := screenshot.CaptureRect(bounds)
		// if err != nil {
		// 	panic(err)
		// }
		// file, err := os.Create("test.png")
		// if err != nil {
		// 	panic(err)
		// }
		// defer file.Close()
		// png.Encode(file, img)
	})

	container.Add(fullscreenScreenshotButton)
	container.Add(regionScreenshotButton)
	container.Add(label)
	return container, nil
}
