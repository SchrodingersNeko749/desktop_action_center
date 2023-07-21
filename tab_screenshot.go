package main

import (
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/kbinani/screenshot"
)

type ScreenTab struct {
	container    *gtk.Box
	listbox      *gtk.ListBox
	ActionCenter *ActionCenter
	x            float64
	y            float64
}

func (app *ScreenTab) Create(actionCenter *ActionCenter) (*gtk.Box, error) {
	app.ActionCenter = actionCenter
	container, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	label, _ := gtk.LabelNew("Wifi")
	fullscreenScreenshotButton, _ := gtk.ButtonNewWithLabel("Screenshot")
	fullscreenScreenshotButton.Connect("clicked", func() {
		go func() {
			app.ActionCenter.ToggleVisiblity()
			time.Sleep(time.Duration(100) * time.Millisecond)
			path := "/tmp/screenshot.png"
			img := ImgGetScreenshot(path)
			img = ImgResize(img, Conf.WINDOW_WIDTH-70, 0)
			app.ActionCenter.ToggleVisiblity()
			actionCenter.AddNotification(NewNotification("", 0, "", "", "", nil, nil, 0, img))
		}()
	})

	regionScreenshotButton, _ := gtk.ButtonNewWithLabel("Region Screenshot")
	regionScreenshotButton.Connect("clicked", func() {
		bounds := screenshot.GetDisplayBounds(0)

		selectionWindow, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

		visual, _ := selectionWindow.GetScreen().GetRGBAVisual()
		selectionWindow.SetSizeRequest(bounds.Size().X, bounds.Size().Y)
		selectionWindow.SetDecorated(false)
		selectionWindow.SetPosition(gtk.WIN_POS_CENTER)
		selectionWindow.SetKeepAbove(true)
		selectionWindow.SetAppPaintable(true)
		selectionWindow.SetVisual(visual)
		selectionWindow.Connect("draw", func(window *gtk.Window, ctx *cairo.Context) {
			ctx.SetSourceRGBA(0.0, 0.0, 0.0, 0.25)
			ctx.SetOperator(cairo.OPERATOR_SOURCE)
			ctx.Paint()
		})
		da, _ := gtk.DrawingAreaNew()
		selectionWindow.AddEvents(4) // accept only integer ...
		unitSize := 20.0
		selectionWindow.Connect("event", app.daEvent)
		da.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
			cr.SetSourceRGB(111, 0, 0)
			cr.Rectangle(app.x*unitSize, app.y*unitSize, unitSize, unitSize)
			cr.Fill()
		})
		// Setting parameter for drawing area
		da.SetHAlign(gtk.ALIGN_FILL)
		da.SetVAlign(gtk.ALIGN_FILL)
		da.SetHExpand(true)
		da.SetVExpand(true)
		da.SetSizeRequest(bounds.Size().X, bounds.Size().Y)
		selectionWindow.Add(da)
		selectionWindow.ShowAll()
	})

	container.Add(fullscreenScreenshotButton)
	container.Add(regionScreenshotButton)
	container.Add(label)
	return container, nil
}

func (app ScreenTab) daEvent(da *gtk.Window, event *gdk.EventMotion) bool {
	app.x, app.y = event.MotionVal()
	return false
}
