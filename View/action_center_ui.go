package View

import (
	"log"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const WINDOW_WIDTH = 700

type ActionCenterUI struct {
	win *gtk.Window
}

func (app *ActionCenterUI) ToggleVisiblity() {
	vis := app.win.IsVisible()
	app.setVisible(!vis)
}
func (a *ActionCenterUI) setVisible(visible bool) {
	if visible {
		a.win.ShowAll()
	} else {
		a.win.Hide()
	}
}
func (app *ActionCenterUI) CreateUI() error {
	// Initialize the window
	if err := app.initWindow(); err != nil {
		return err
	}

	// Create and add the label widget
	if err := app.addLabel(); err != nil {
		return err
	}

	return nil
}
func (app *ActionCenterUI) addLabel() error {
	// Create the label widget
	l, err := gtk.LabelNew(time.Now().Format("Mon 3:04 PM"))
	l.SetName("header")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	l2, err := gtk.LabelNew("Action Center")
	l2.SetName("header")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	// Create the vertical box container and add the label widget to it
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}
	vbox.Add(l)
	vbox.Add(l2)

	// Create a new CSS provider and load the CSS file
	provider, err := gtk.CssProviderNew()
	if err != nil {
		return err
	}

	if err := provider.LoadFromPath("/home/neko/Projects/programming/go/desktop_action_center/window.css"); err != nil {
		return err
	}

	// Apply the CSS provider to the label widget's style context
	lStyle, err := l.GetStyleContext()
	if err != nil {
		return err
	}
	lStyle.AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	// Apply the CSS provider to the box container's style context
	vboxStyle, err := vbox.GetStyleContext()
	if err != nil {
		return err
	}
	vboxStyle.AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
	vbox.SetHAlign(gtk.ALIGN_START)
	// Add the box container to the window
	app.win.Add(vbox)

	return nil
}
func (app *ActionCenterUI) initWindow() error {
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

	app.win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return err
	}

	app.win.SetTitle("action-center-panel")
	app.win.SetDefaultSize(WINDOW_WIDTH, height-32)
	app.win.Move(width-WINDOW_WIDTH, 32)
	app.win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)
	app.win.SetResizable(false)
	visual, _ := screen.GetRGBAVisual()
	app.win.SetVisual(visual)
	app.win.SetDecorated(false)

	// Create a new CSS provider and load the CSS file
	provider, err := gtk.CssProviderNew()
	if err != nil {
		return err
	}

	if err := provider.LoadFromPath("/home/neko/Projects/programming/go/desktop_action_center/window.css"); err != nil {
		return err
	}

	// Apply the CSS provider to the window style context
	style, err := app.win.GetStyleContext()
	if err != nil {
		return err
	}
	style.AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	return nil
}
func (app *ActionCenterUI) Run() {
	app.win.SetPosition(gtk.WIN_POS_NONE)
	app.win.ShowAll()
}
