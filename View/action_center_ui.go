package View

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/actionCenter/Command"
	"github.com/actionCenter/Data"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var WINDOW_WIDTH = 550
var ICON_SIZE = 64
var HORIZONTAL_SPACING = 24
var VERTICAL_SPACING = 32
var CSS_THEME_FILE = "neko.css"

type ActionCenterUI struct {
	win                 *gtk.Window
	container           *gtk.Box
	actionCenterHandler Command.ActionCenterInterface
	notifications       NotificationList
	aimessages          NotificationList
}

func (app *ActionCenterUI) initWindow() error {
	WINDOW_WIDTH = Data.Conf.WINDOW_WIDTH
	ICON_SIZE = Data.Conf.ICON_SIZE
	HORIZONTAL_SPACING = Data.Conf.HORIZONTAL_SPACING
	VERTICAL_SPACING = Data.Conf.VERTICAL_SPACING
	CSS_THEME_FILE = Data.Conf.CSS_THEME_FILE

	screen, _ := gdk.ScreenGetDefault()
	visual, _ := screen.GetRGBAVisual()
	display, _ := screen.GetDisplay()
	monitor, _ := display.GetPrimaryMonitor()
	width := monitor.GetGeometry().GetWidth()
	height := monitor.GetGeometry().GetHeight()

	app.win, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	app.win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)
	app.win.SetTitle("action-center-panel")
	app.win.SetDefaultSize(WINDOW_WIDTH, height-32)
	app.win.Move(width-WINDOW_WIDTH, 32)
	app.win.SetResizable(false)
	app.win.SetVisual(visual)
	app.win.SetDecorated(false)

	app.win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	style, err := app.win.GetStyleContext()
	style.AddProvider(Data.StyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	return err
}

func (app *ActionCenterUI) ToggleVisiblity() {
	app.win.SetVisible(!app.win.GetVisible())
	if app.win.GetVisible() {
		app.win.ShowAll()
	}
}

func (app *ActionCenterUI) CreateUI(ac Command.ActionCenterInterface) error {
	if err := app.initWindow(); err != nil {
		return err
	}
	app.actionCenterHandler = ac

	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}
	app.container = container
	for _, widget := range Data.WidgetConfs {
		widgetContainer, err := app.createComponent(&widget)
		if err != nil {
			return err
		}
		app.container.Add(widgetContainer)
	}
	app.win.Add(app.container)
	return nil
}
func (app *ActionCenterUI) createComponent(widget *Data.WidgetConfig) (*gtk.Box, error) {
	var component *gtk.Box
	var notebook *gtk.Notebook // for tabviewer
	var err error
	switch widget.Type {
	case "header":
		if component, err = app.createHeaderComponent(); err != nil {
			return nil, err
		}
	case "brightness":
		fmt.Println("test")
		if component, err = app.createBrightnessComponent(widget); err != nil {
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
		if component, err = app.Create(); err != nil {
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
			childComponent, err := app.createComponent(child)
			childComponent.SetCanFocus(true)
			if err != nil {
				return nil, err
			}
			app.addTab(notebook, child.Properties.Label, childComponent)
		} else {
			childContainer, err := app.createComponent(child)
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
	clockLabel, err := gtk.LabelNew("")
	go func() {
		glib.TimeoutAdd(uint(1000), func() bool {
			clockLabel.SetText(time.Now().Format("Mon 3:04 PM"))
			return true
		})
	}()
	clockLabel.SetName("clock")
	if err != nil {
		return nil, err
	}

	lStyle, err := clockLabel.GetStyleContext()
	if err != nil {
		return nil, err
	}
	lStyle.AddProvider(Data.StyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	vboxStyle, err := vbox.GetStyleContext()
	if err != nil {
		return nil, err
	}
	vboxStyle.AddProvider(Data.StyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
	vbox.SetHAlign(gtk.ALIGN_START)
	vbox.PackStart(clockLabel, true, true, 0)
	return vbox, nil
}
func (app *ActionCenterUI) createBrightnessComponent(configWidget *Data.WidgetConfig) (*gtk.Box, error) {
	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	hbox.SetHAlign(gtk.ALIGN_CENTER)
	style, _ := hbox.GetStyleContext()
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style.AddClass("notification-container-header")

	brightnessBar, err := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	brightnessBar.SetHExpand(true)
	brightnessBar.SetSizeRequest(500, -1)
	cmd := exec.Command("./getbrightness.sh")
	output, _ := cmd.Output()
	brightnessBar.SetValue(float64(output[0]))
	style, _ = brightnessBar.GetStyleContext()
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	brightnessBar.Connect("value-changed", func() {
		v := brightnessBar.GetValue()
		cmd := exec.Command("./setbrightness.sh", fmt.Sprintf("%d", int(v)))
		output, err := cmd.Output()
		fmt.Println(string(output), err)
	})

	label, _ := gtk.LabelNew(configWidget.Properties.Label)

	hbox.PackStart(brightnessBar, true, true, 0)
	hbox.PackEnd(label, true, true, 0)

	return hbox, err
}
func (app *ActionCenterUI) Run() {
	app.ShowAll()
}

func (app *ActionCenterUI) ShowAll() {
	app.win.ShowAll()
}
