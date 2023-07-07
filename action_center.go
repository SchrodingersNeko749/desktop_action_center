package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type ActionCenter struct {
	StickyShow         bool
	win                *gtk.Window
	container          *gtk.Box
	notificationServer *NotificationServer

	Header          *Header
	TabControl      *gtk.Notebook
	NotificationTab *NotificationTab
	WifiTab         *WifiTab
	ScreenTab       *ScreenTab
	RadioTab        *RadioTab
	AI_Tab          *AiTab
}

func main() {
	gtk.Init(nil)
	LoadConfig()

	app := &ActionCenter{
		notificationServer: &NotificationServer{},
		Header:             &Header{},
		NotificationTab:    &NotificationTab{},
		WifiTab:            &WifiTab{},
		ScreenTab:          &ScreenTab{},
		RadioTab:           &RadioTab{},
		AI_Tab:             &AiTab{},
	}

	go app.HandleSignals()
	go app.notificationServer.Init(app)
	app.RadioTab.SetRadioDirectoryServerIP("all.api.radio-browser.info")

	app.initWindow()
	app.win.ShowAll()
	gtk.Main()
}

func (app *ActionCenter) HandleSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGTERM)
	fmt.Println("Monitoring signals")
	for {
		sig := <-sigs
		fmt.Println(sig)
		switch sig {
		case syscall.SIGUSR1:
			app.ToggleVisiblity()
			app.StickyShow = true
		case syscall.SIGTERM:
			fmt.Println("Closing dbus conn")
			app.notificationServer.Conn.Close()
			os.Exit(0)
		}
	}
}

func (app *ActionCenter) initWindow() {
	screen, _ := gdk.ScreenGetDefault()
	visual, _ := screen.GetRGBAVisual()
	display, _ := screen.GetDisplay()
	monitor, _ := display.GetPrimaryMonitor()
	width := monitor.GetGeometry().GetWidth()
	height := monitor.GetGeometry().GetHeight()

	app.win, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	app.StickyShow = true
	app.win.SetTitle("action-center-panel")
	app.win.SetDefaultSize(Conf.WINDOW_WIDTH, height-32)
	app.win.Move(width-Conf.WINDOW_WIDTH, 32)
	app.win.SetResizable(false)
	app.win.SetVisual(visual)
	app.win.SetDecorated(false)
	app.container, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	for _, widget := range WidgetConfs {
		widgetContainer, _ := app.createComponent(&widget)
		app.container.Add(widgetContainer)
	}
	app.win.Add(app.container)

	app.win.Connect("configure-event", func(win *gtk.Window, event *gdk.Event) {
		win.Move(width-Conf.WINDOW_WIDTH, 32)
	})
	app.win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	app.win.Connect("focus-out-event", func() {
		if !app.StickyShow {
			app.win.Hide()
		}
		app.StickyShow = false
	})

	style, _ := app.win.GetStyleContext()
	style.AddProvider(StyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
}

func (app *ActionCenter) createComponent(widget *WidgetConfig) (*gtk.Box, error) {
	var component *gtk.Box
	var err error

	switch widget.Type {
	case "header":
		component = app.Header.Create()
	case "brightness":
		component, err = app.createBrightnessComponent(widget)
	case "tab-viewer":
		component, app.TabControl, err = app.createTabViewerContainer(widget)
	case "wifi":
		component, err = app.WifiTab.Create()
	case "radio":
		component, err = app.RadioTab.Create()
	case "ai":
		component, err = app.AI_Tab.Create()
	case "notification":
		component, err = app.NotificationTab.Create(*app.win)
	case "capture":
		component, err = app.ScreenTab.Create(app)
	default:
		return nil, fmt.Errorf("unrecognized widget type: %s", widget.Type)
	}

	if err != nil {
		return nil, err
	}

	for _, child := range widget.Children {
		if widget.Type == "tab-viewer" {
			childComponent, err := app.createComponent(child)
			childComponent.SetCanFocus(true)
			if err != nil {
				return nil, err
			}
			tabLabel, _ := gtk.LabelNew(child.Properties.Label)
			tabLabel.SetSizeRequest(Conf.ICON_SIZE, Conf.ICON_SIZE)
			app.TabControl.AppendPage(childComponent, tabLabel)
			app.TabControl.GetNPages()
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

func (app *ActionCenter) createBrightnessComponent(configWidget *WidgetConfig) (*gtk.Box, error) {
	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	hbox.SetHAlign(gtk.ALIGN_CENTER)
	style, _ := hbox.GetStyleContext()
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style.AddClass("notification-container-header")

	brightnessBar, err := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	brightnessBar.SetHExpand(true)
	brightnessBar.SetSizeRequest(500, -1)
	cmd := exec.Command("./assets/getbrightness.sh")
	output, _ := cmd.Output()
	brightnessBar.SetValue(float64(output[0]))
	style, _ = brightnessBar.GetStyleContext()
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	brightnessBar.Connect("value-changed", func() {
		v := brightnessBar.GetValue()
		cmd := exec.Command("./assets/setbrightness.sh", fmt.Sprintf("%d", int(v)))
		output, err := cmd.Output()
		fmt.Println(string(output), err)
	})

	label, _ := gtk.LabelNew(configWidget.Properties.Label)

	hbox.PackStart(brightnessBar, true, true, 0)
	hbox.PackEnd(label, true, true, 0)

	return hbox, err
}

func (app *ActionCenter) createTabViewerContainer(configWidget *WidgetConfig) (*gtk.Box, *gtk.Notebook, error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	box.SetSizeRequest(Conf.WINDOW_WIDTH, -1)
	if err != nil {
		return nil, nil, err
	}

	notebook, err := gtk.NotebookNew()
	notebook.SetHExpand(true)
	notebook.SetHAlign(gtk.ALIGN_FILL)

	stylectx, _ := notebook.GetStyleContext()
	stylectx.AddClass("tab-viewer")
	stylectx.AddProvider(StyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	box.Add(notebook)
	return box, notebook, nil
}

func (app *ActionCenter) AddNotification(n Notification) {
	notifictation, _ := CreateNotificationComponent(n)
	app.NotificationTab.AddNotification(notifictation)
	pageNum := app.TabControl.PageNum(app.NotificationTab.Container)
	app.TabControl.SetCurrentPage(pageNum)
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
