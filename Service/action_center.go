package Service

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/actionCenter/Data"
	"github.com/actionCenter/Model"
	"github.com/actionCenter/View"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type ActionCenter struct {
	win                *gtk.Window
	container          *gtk.Box
	notificationServer *NotificationServer
	NotificationTab    *View.NotificationTab
	WifiTab            *View.WifiTab
	ScreenTab          *View.ScreenTab
	RadioTab           *View.RadioTab
	AITab              *View.AITab
}

func NewActionCenter() *ActionCenter {
	return &ActionCenter{
		NotificationTab: &View.NotificationTab{},
		WifiTab:         &View.WifiTab{},
		ScreenTab:       &View.ScreenTab{},
		RadioTab:        &View.RadioTab{},
		AITab:           &View.AITab{},
	}
}

func (app *ActionCenter) Init() {
	app.notificationServer = NewNotificationServer()
	go app.notificationServer.Init(app)
	if err := app.initWindow(); err != nil {
		return
	}

	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}
	app.container = container
	for _, widget := range Data.WidgetConfs {
		widgetContainer, err := app.createComponent(&widget)
		if err != nil {
			return
		}
		app.container.Add(widgetContainer)
	}
	app.win.Add(app.container)
	app.win.ShowAll()

	// handling signals
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
	gtk.Main()
}

func (app *ActionCenter) GetNotifications() ([]Model.Notification, error) {

	ns, err := app.notificationServer.GetNotifications()
	if err != nil {
		return nil, err
	}
	return ns, nil
}
func (app *ActionCenter) AddNotification(n Model.Notification) error {
	widget := Model.CreateNotificationComponent(n)
	app.NotificationTab.AddNotification(widget)
	app.win.ShowAll()
	return nil
}

func (app *ActionCenter) initWindow() error {
	screen, _ := gdk.ScreenGetDefault()
	visual, _ := screen.GetRGBAVisual()
	display, _ := screen.GetDisplay()
	monitor, _ := display.GetPrimaryMonitor()
	width := monitor.GetGeometry().GetWidth()
	height := monitor.GetGeometry().GetHeight()

	app.win, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	app.win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)
	app.win.SetTitle("action-center-panel")
	app.win.SetDefaultSize(Data.Conf.WINDOW_WIDTH, height-32)
	app.win.Move(width-Data.Conf.WINDOW_WIDTH, 32)
	app.win.SetResizable(false)
	app.win.SetVisual(visual)
	app.win.SetDecorated(false)

	app.win.Connect("configure-event", func(win *gtk.Window, event *gdk.Event) {
		win.Move(width-Data.Conf.WINDOW_WIDTH, 32)
	})
	app.win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	style, err := app.win.GetStyleContext()
	style.AddProvider(Data.StyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	return err
}

func (app *ActionCenter) ToggleVisiblity() {
	app.win.SetVisible(!app.win.GetVisible())
}

func (app *ActionCenter) createComponent(widget *Data.WidgetConfig) (*gtk.Box, error) {
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
		if component, err = app.WifiTab.Create(); err != nil {
			return nil, err
		}
	case "radio":
		if component, err = app.RadioTab.Create(); err != nil {
			return nil, err
		}
	case "ai":
		if component, err = app.AITab.Create(); err != nil {
			return nil, err
		}
	case "notification":
		if component, err = app.NotificationTab.Create(); err != nil {
			return nil, err
		}
	case "capture":
		if component, err = app.ScreenTab.Create(); err != nil {
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

func (app *ActionCenter) createHeaderComponent() (*gtk.Box, error) {
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
	//box.SetSizeRequest(WINDOW_WIDTH, -1)
	notebook, err := gtk.NotebookNew()
	if err != nil {
		return nil, nil, err
	}

	notebook.SetHExpand(true)

	notebook.SetHAlign(gtk.ALIGN_CENTER)

	stylectx, err := notebook.GetStyleContext()
	if err != nil {
		return nil, nil, err
	}
	stylectx.AddClass("tab-viewer")
	stylectx.AddProvider(Data.StyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
	notebook.SetCurrentPage(0)

	box.Add(notebook)
	return box, notebook, nil
}

func (app *ActionCenter) addTab(notebook *gtk.Notebook, tabLabelString string, page *gtk.Box) {
	tabLabel, _ := gtk.LabelNew(tabLabelString)
	tabLabel.SetSizeRequest(50, 50)
	notebook.AppendPage(page, tabLabel)
}
