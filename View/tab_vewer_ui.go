package View

import (
	"fmt"
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

func (app *ActionCenterUI) createTabViewerContainer(configWidget Widget) (*gtk.Box, *gtk.Notebook, error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, nil, err
	}
	//box.SetSizeRequest(WINDOW_WIDTH, -1)

	for _, child := range configWidget.Children {
		fmt.Println(child.Type)
	}
	notebook, err := gtk.NotebookNew()
	if err != nil {
		return nil, nil, err
	}

	notebook.SetHExpand(false)
	notebook.SetHAlign(gtk.ALIGN_CENTER)

	stylectx, err := notebook.GetStyleContext()
	if err != nil {
		return nil, nil, err
	}
	stylectx.AddClass("tab-viewer")
	stylectx.AddProvider(app.componentStyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
	notebook.SetCurrentPage(0)

	notebook.Connect("switch-page", func() {
		switch notebook.GetCurrentPage() {
		case 0:
			fmt.Println("wifi")
		case 1:
			fmt.Println("radio")
		case 2:
			fmt.Println("ai")
		case 3:
			fmt.Println("notification")
			app.ShowNotifications()
			app.notifications.listBox.ShowAll()

		case 4:
			fmt.Println("capture")
		case -1:
			fmt.Println(-1)
		}

	})

	box.Add(notebook)
	return box, notebook, nil
}
func (app *ActionCenterUI) addTab(notebook *gtk.Notebook, tab *Widget) {
	var page *gtk.Box
	tabLabel, _ := gtk.LabelNew(tab.Properties.Label)
	tabLabel.SetSizeRequest(50, 50)
	var err error

	switch tab.Type {
	case "wifi":
		page, err = app.createWifiComponent()
	case "radio":
		page, _ = app.createRadioComponent()
	case "ai":
		page, _ = app.createAiComponent()
	case "notification":
		page, _ = app.createNotificationComponent()
	case "capture":
		page, _ = app.createScreenCaptureComponent()
	default:
		fmt.Println("invalid tab")
		os.Exit(1)
	}
	if err != nil {
		log.Fatalln(err)
	}
	notebook.AppendPage(page, tabLabel)

}
