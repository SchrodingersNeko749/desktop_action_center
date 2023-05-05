package main

import (
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type ActionCenterHandler interface {
	SetVisible(bool)
}
type ActionCenter struct {
	win *gtk.Window
}

func NewActionCenter() *ActionCenter {
	return &ActionCenter{}
}

func (app *ActionCenter) Init() error {
	gtk.Init(nil)

	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		return err
	}
	display, _ := screen.GetDisplay()
	monitor, _ := display.GetPrimaryMonitor()
	geometry := monitor.GetGeometry()
	width := geometry.GetWidth()
	height := geometry.GetHeight()

	app.win, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	app.win.SetTitle("notification-panel")
	app.win.SetDefaultSize(WINDOW_WIDTH, height-32)
	app.win.Move(width-WINDOW_WIDTH, 32)
	app.win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)
	app.win.SetResizable(false)
	visual, _ := screen.GetRGBAVisual()
	app.win.SetVisual(visual)
	app.win.SetDecorated(false)

	if err != nil {
		return err
	}

	app.win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Create a new CSS provider and load the CSS file
	provider, err := gtk.CssProviderNew()
	if err != nil {
		return err
	}
	if err := provider.LoadFromPath("window.css"); err != nil {
		return err
	}

	// Apply the CSS provider to the window
	style, err := app.win.GetStyleContext()
	if err != nil {
		return err
	}
	style.AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	return nil
}

func (app *ActionCenter) Run() {
	// call SetPosition in a loop to keep resetting the position
	app.win.SetPosition(gtk.WIN_POS_NONE)

	app.win.ShowAll()
	gtk.Main()
}

func (a *ActionCenter) SetVisible(visible bool) {
	if a.win == nil {
		log.Println("Action center window is nil")
		return
	}

	if visible {
		a.win.ShowAll()
	} else {
		a.win.Hide()
	}
}
