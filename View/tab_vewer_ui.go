package View

import (
	"fmt"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func (app *ActionCenterUI) createTabViewerContainer() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}
	//box.SetSizeRequest(WINDOW_WIDTH, -1)

	notebook, err := gtk.NotebookNew()
	if err != nil {
		return err
	}

	notebook.SetHExpand(false)
	notebook.SetHAlign(gtk.ALIGN_CENTER)
	//notebook.SetSizeRequest(WINDOW_WIDTH, -1)

	// Add tabs to the notebook
	wtab, _ := gtk.LabelNew("")
	wtab.SetSizeRequest(50, 50)
	wtab.Connect("motion-notify-event", func(widget *gtk.Widget, event *gdk.Event) {
		fmt.Println("test")
	})
	rtab, _ := gtk.LabelNew("")
	rtab.SetSizeRequest(50, 50)

	atab, _ := gtk.LabelNew("🤖")
	atab.SetSizeRequest(50, 50)

	ntab, _ := gtk.LabelNew("")
	ntab.SetSizeRequest(50, 50)

	ctab, _ := gtk.LabelNew("")
	ctab.SetSizeRequest(50, 50)

	w, err := app.createWifiComponent()
	if err != nil {
		return err
	}
	r, err := app.createRadioComponent()
	if err != nil {
		return err
	}
	a, err := app.createAiComponent()
	if err != nil {
		return err
	}
	n, err := app.createNotificationComponent()
	if err != nil {
		return err
	}
	c, err := app.createScreenCaptureComponent()
	if err != nil {
		return err
	}

	stylectx, err := notebook.GetStyleContext()
	if err != nil {
		return nil
	}
	stylectx.AddClass("tab-viewer")
	stylectx.AddProvider(app.containerStyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	notebook.AppendPage(w, wtab)
	notebook.AppendPage(r, rtab)
	notebook.AppendPage(a, atab)
	notebook.AppendPage(n, ntab)
	notebook.AppendPage(c, ctab)

	notebook.SetCurrentPage(0)

	// notebook.Connect("switch-page", func() {
	// 	switch notebook.GetCurrentPage() {
	// 	case 0:
	// 		fmt.Println("wifi")
	// 	case 1:
	// 		fmt.Println("radio")
	// 	case 2:
	// 		fmt.Println("ai")
	// 	case 3:
	// 		fmt.Println("notification")
	// 		app.ShowNotifications()
	// 		app.AddNotification("youtube", "Sample youtube notification", "Check out my video!")

	// 		app.notifications.listBox.ShowAll()

	// 	case 4:
	// 		fmt.Println("capture")
	// 	case -1:
	// 		fmt.Println(-1)
	// 	}

	// })

	box.Add(notebook)
	app.container.Add(box)
	return nil
}
