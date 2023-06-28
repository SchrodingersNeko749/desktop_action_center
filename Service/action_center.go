package Service

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/actionCenter/Data"
	"github.com/actionCenter/Model"
	"github.com/actionCenter/View"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type ActionCenter struct {
	win                *gtk.Window
	container          *gtk.Box
	notificationServer *NotificationServer

	HeaderUI        *View.HeaderUI
	NotificationTab *View.NotificationTab
	WifiTab         *View.WifiTab
	ScreenTab       *View.ScreenTab
	RadioTab        *View.RadioTab
	AITab           *View.AITab
}

func NewActionCenter() *ActionCenter {
	return &ActionCenter{
		notificationServer: &NotificationServer{},

		HeaderUI:        &View.HeaderUI{},
		NotificationTab: &View.NotificationTab{},
		WifiTab:         &View.WifiTab{},
		ScreenTab:       &View.ScreenTab{},
		RadioTab:        &View.RadioTab{},
		AITab:           &View.AITab{},
	}
}

func (app *ActionCenter) Init() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGTERM)
	go func() {
		fmt.Println("Monitoring signals")
		for {
			sig := <-sigs
			fmt.Println(sig)
			switch sig {
			case syscall.SIGUSR1:
				app.ToggleVisiblity()
			case syscall.SIGTERM:
				fmt.Println("Closing dbus conn")
				app.notificationServer.conn.Close()
				os.Exit(0)
			}
		}
	}()

	go app.notificationServer.Init(app)

	app.initWindow()
	app.win.ShowAll()
	gtk.Main()
}

func (app *ActionCenter) initWindow() {
	screen, _ := gdk.ScreenGetDefault()
	visual, _ := screen.GetRGBAVisual()
	display, _ := screen.GetDisplay()
	monitor, _ := display.GetPrimaryMonitor()
	width := monitor.GetGeometry().GetWidth()
	height := monitor.GetGeometry().GetHeight()

	app.win, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	// app.win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)
	app.win.SetTitle("action-center-panel")
	app.win.SetDefaultSize(Data.Conf.WINDOW_WIDTH, height-32)
	app.win.Move(width-Data.Conf.WINDOW_WIDTH, 32)
	app.win.SetResizable(false)
	app.win.SetVisual(visual)
	app.win.SetDecorated(false)
	app.container, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	for _, widget := range Data.WidgetConfs {
		widgetContainer, _ := app.createComponent(&widget)
		app.container.Add(widgetContainer)
	}
	app.win.Add(app.container)

	app.win.Connect("configure-event", func(win *gtk.Window, event *gdk.Event) {
		win.Move(width-Data.Conf.WINDOW_WIDTH, 32)
	})
	app.win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	style, _ := app.win.GetStyleContext()
	style.AddProvider(Data.StyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
}

func (app *ActionCenter) createComponent(widget *Data.WidgetConfig) (*gtk.Box, error) {
	var component *gtk.Box
	var notebook *gtk.Notebook // for tabviewer
	var err error

	switch widget.Type {
	case "header":
		component, err = app.HeaderUI.Create()
	case "brightness":
		component, err = app.createBrightnessComponent(widget)
	case "tab-viewer":
		component, notebook, err = app.createTabViewerContainer(widget)
	case "wifi":
		component, err = app.WifiTab.Create()
	case "radio":
		component, err = app.RadioTab.Create()
	case "ai":
		component, err = app.AITab.Create()
	case "notification":
		component, err = app.NotificationTab.Create(*app.win)
	case "capture":
		component, err = app.ScreenTab.Create()
	default:
		return nil, fmt.Errorf("unrecognized widget type: %s", widget.Type)
	}

	if err != nil {
		return nil, err
	}

	// Recursively call the method for the children of the widget
	for _, child := range widget.Children {
		if widget.Type == "tab-viewer" {
			childComponent, err := app.createComponent(child)
			childComponent.SetCanFocus(true)
			if err != nil {
				return nil, err
			}
			tabLabel, _ := gtk.LabelNew(child.Properties.Label)
			tabLabel.SetSizeRequest(50, 50)
			notebook.AppendPage(childComponent, tabLabel)
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

func (app *ActionCenter) createBrightnessComponent(configWidget *Data.WidgetConfig) (*gtk.Box, error) {
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

func (app *ActionCenter) createTabViewerContainer(configWidget *Data.WidgetConfig) (*gtk.Box, *gtk.Notebook, error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	box.SetSizeRequest(Data.Conf.WINDOW_WIDTH, -1)
	if err != nil {
		return nil, nil, err
	}

	notebook, err := gtk.NotebookNew()
	if err != nil {
		return nil, nil, err
	}

	notebook.SetHExpand(true)
	notebook.SetHAlign(gtk.ALIGN_CENTER)

	stylectx, _ := notebook.GetStyleContext()
	stylectx.AddClass("tab-viewer")
	stylectx.AddProvider(Data.StyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
	notebook.SetCurrentPage(0)

	box.Add(notebook)
	return box, notebook, nil
}

func (app *ActionCenter) GetNotifications() ([]Model.Notification, error) {
	return app.notificationServer.GetNotifications()
}

func (app *ActionCenter) AddNotification(n Model.Notification) {
	notifictation, _ := Model.CreateNotificationComponent(n)
	app.NotificationTab.AddNotification(notifictation)
	if app.win.GetVisible() {
		app.win.ShowAll()
	}
}

func (app *ActionCenter) ToggleVisiblity() {
	app.win.SetVisible(!app.win.GetVisible())
	if app.win.GetVisible() {
		app.win.ShowAll()
	}
}
