package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type UI struct {
	win     *gtk.Window
	display *gdk.Display
	monitor gdk.Monitor
}

func NewUI() *UI {
	return &UI{}
}

func (ui *UI) Init() error {
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

	ui.win, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	ui.win.SetTitle("notification-panel")
	ui.win.SetDefaultSize(WINDOW_WIDTH, height-32)
	ui.win.Move(width-WINDOW_WIDTH, 32)
	ui.win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)
	ui.win.SetResizable(false)
	visual, _ := screen.GetRGBAVisual()
	ui.win.SetVisual(visual)
	ui.win.SetDecorated(false)

	if err != nil {
		return err
	}

	ui.win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Create a new CSS provider and load the CSS file
	provider, err := gtk.CssProviderNew()
	if err != nil {
		return err
	}
	if err := provider.LoadFromPath("style.css"); err != nil {
		return err
	}

	// Apply the CSS provider to the window
	style, err := ui.win.GetStyleContext()
	if err != nil {
		return err
	}
	style.AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	return nil
}

func (ui *UI) Run() {
	// call SetPosition in a loop to keep resetting the position
	ui.win.SetPosition(gtk.WIN_POS_MOUSE)

	ui.win.ShowAll()
	gtk.Main()
}
