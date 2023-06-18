package View

import (
	"log"
	"time"

	"github.com/actionCenter/Command"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const WINDOW_WIDTH = 700

type ActionCenterUI struct {
	win                    *gtk.Window
	containerStyleProvider *gtk.CssProvider
	container              *gtk.Box
	actionCenter           Command.ActionCenterInterface
	notifications          NotificationList
}

func (app *ActionCenterUI) ToggleVisiblity() {
	if app.win.IsVisible() {
		app.win.Hide()
	} else {
		app.win.ShowAll()
	}
}

func (app *ActionCenterUI) CreateUI(ac Command.ActionCenterInterface) error {
	// Initialize the window
	if err := app.initWindow(); err != nil {
		return err
	}
	// make the actioncenter handler
	app.actionCenter = ac

	c, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}
	app.container = c
	// Add Containers
	if err := app.createHeaderContainer(); err != nil {
		return err
	}

	if err := app.createTabViewerContainer(); err != nil {
		return err
	}

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	label, _ := gtk.LabelNew("test")
	box.Add(label)
	box.AddEvents(int(gdk.POINTER_MOTION_MASK))
	box.AddEvents(int(gdk.POINTER_MOTION_HINT_MASK))

	app.win.Add(app.container)
	return nil
}
func (app *ActionCenterUI) createHeaderContainer() error {

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}
	l, err := gtk.LabelNew(time.Now().Format("Mon 3:04 PM"))
	l.SetName("clock")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	// Apply the CSS provider to the label widget's style context
	lStyle, err := l.GetStyleContext()
	if err != nil {
		return err
	}
	lStyle.AddProvider(app.containerStyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	// Apply the CSS provider to the box container's style context
	vboxStyle, err := vbox.GetStyleContext()
	if err != nil {
		return err
	}
	vboxStyle.AddProvider(app.containerStyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
	vbox.SetHAlign(gtk.ALIGN_START)
	// Add the box container to the window

	vbox.Add(l)
	app.container.Add(vbox)

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

	// Initialize all Css containers
	app.initCssStyles()

	// Apply the CSS provider to the window style context
	style, err := app.win.GetStyleContext()
	if err != nil {
		return err
	}
	style.AddProvider(app.containerStyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	return nil
}
func (app *ActionCenterUI) initCssStyles() error {
	// Create a new CSS provider and load the CSS file
	provider, err := gtk.CssProviderNew()
	if err != nil {
		return err
	}

	if err := provider.LoadFromPath("assets/window.css"); err != nil {
		return err
	}
	app.containerStyleProvider = provider
	return nil
}
func (app *ActionCenterUI) Run() {
	app.win.SetPosition(gtk.WIN_POS_NONE)
	app.win.ShowAll()
}
