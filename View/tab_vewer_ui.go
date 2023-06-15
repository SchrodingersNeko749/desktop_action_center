package View

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

func (app *ActionCenterUI) createTabViewerContainer() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}
	box.SetHExpand(true)
	box.SetSizeRequest(WINDOW_WIDTH, -1)

	notebook, err := gtk.NotebookNew()
	if err != nil {
		return err
	}
	notebook.SetHAlign(gtk.ALIGN_CENTER)

	// Add tabs to the notebook
	wtab, _ := gtk.LabelNew("")
	wtab.SetSizeRequest(50, 50)

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

	notebook.Connect("switch-page", func() {
		fmt.Println(notebook.GetCurrentPage())
	})

	box.Add(notebook)
	app.container.Add(box)
	return nil
}
