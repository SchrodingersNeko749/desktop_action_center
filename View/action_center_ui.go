package View

import (
	"fmt"
	"time"

	"github.com/actionCenter/Command"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const WINDOW_WIDTH = 680

type ActionCenterUI struct {
	win                    *gtk.Window
	componentStyleProvider *gtk.CssProvider
	container              *gtk.Box
	actionCenterHandler    Command.ActionCenterInterface
	notifications          NotificationList
}

func (app *ActionCenterUI) ToggleVisiblity() {
	app.win.SetVisible(!app.win.GetVisible())
}

func (app *ActionCenterUI) CreateUI(ac Command.ActionCenterInterface, filename string) error {
	// Initialize the window
	if err := app.initWindow(); err != nil {
		return err
	}
	// make the actioncenter handler
	app.actionCenterHandler = ac

	ws, err := GetWidgetsFromConfig("View/test.json")
	if err != nil {
		return err
	}
	c, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}
	app.container = c
	// Add Containers
	for _, widget := range ws {
		widgetContainer, err := app.createComponent(widget)
		if err != nil {
			return err
		}
		app.container.Add(widgetContainer)
	}
	app.win.Add(app.container)
	return nil
}
func (app *ActionCenterUI) createComponent(widget Widget) (*gtk.Box, error) {
	var component *gtk.Box
	var notebook *gtk.Notebook // for tabviewer
	var err error
	switch widget.Type {
	case "header":
		if component, err = app.createHeaderComponent(); err != nil {
			return nil, err
		}
	case "tab-viewer":
		if component, notebook, err = app.createTabViewerContainer(widget); err != nil {
			return nil, err
		}
	case "wifi":
		if component, err = app.createWifiComponent(); err != nil {
			return nil, err
		}
	case "radio":
		if component, err = app.createRadioComponent(); err != nil {
			return nil, err
		}
	case "ai":
		if component, err = app.createAiComponent(); err != nil {
			return nil, err
		}
	case "notification":
		if component, err = app.createNotificationComponent(); err != nil {
			return nil, err
		}
	case "capture":
		if component, err = app.createScreenCaptureComponent(); err != nil {
			return nil, err
		}
	default:
		// Handle unrecognized widget types
		return nil, fmt.Errorf("unrecognized widget type: %s", widget.Type)
	}

	// Recursively call the method for the children of the widget
	for _, child := range widget.Children {
		if widget.Type == "tab-viewer" {
			childComponent, err := app.createComponent(*child)
			if err != nil {
				return nil, err
			}
			app.addTab(notebook, child.Properties.Label, childComponent)
		} else {
			childContainer, err := app.createComponent(*child)
			if err != nil {
				return nil, err
			}
			component.Add(childContainer)
		}
	}

	return component, nil
}
func (app *ActionCenterUI) createHeaderComponent() (*gtk.Box, error) {

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	l, err := gtk.LabelNew(time.Now().Format("Mon 3:04 PM"))
	l.SetName("clock")
	if err != nil {
		return nil, err
	}

	// Apply the CSS provider to the label widget's style context
	lStyle, err := l.GetStyleContext()
	if err != nil {
		return nil, err
	}
	lStyle.AddProvider(app.componentStyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	// Apply the CSS provider to the box container's style context
	vboxStyle, err := vbox.GetStyleContext()
	if err != nil {
		return nil, err
	}
	vboxStyle.AddProvider(app.componentStyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
	vbox.SetHAlign(gtk.ALIGN_START)
	// Add the box container to the window

	vbox.Add(l)
	return vbox, nil
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
	style.AddProvider(app.componentStyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

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
	app.componentStyleProvider = provider
	return nil
}
func (app *ActionCenterUI) Run() {
	app.win.SetPosition(gtk.WIN_POS_NONE)
	app.win.ShowAll()
}
func (app *ActionCenterUI) ShowAll() {
	app.win.ShowAll()
}
