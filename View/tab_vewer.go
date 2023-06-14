package View

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

func (app *ActionCenterUI) createTabViewerContainer() error {
	notebook, err := gtk.NotebookNew()
	if err != nil {
		return err
	}
	notebook.SetHAlign(gtk.ALIGN_CENTER)
	notebook.SetHExpand(false)
	notebook.SetVExpand(false)
	notebook.SetSizeRequest(140, 140)
	// Add tabs to the notebook
	tab1, _ := gtk.LabelNew("")
	tab2, _ := gtk.LabelNew("")
	tab3, _ := gtk.LabelNew("🤖")
	tab4, _ := gtk.LabelNew("")

	label, err := gtk.LabelNew("Wifi")
	if err != nil {
		return err
	}
	label2, err := gtk.LabelNew("Radio")
	if err != nil {
		return err
	}
	label3, err := gtk.LabelNew("Ai")
	if err != nil {
		return err
	}
	label4, err := gtk.LabelNew("Notification")
	if err != nil {
		return err
	}

	stylectx, err := notebook.GetStyleContext()
	if err != nil {
		return nil
	}
	stylectx.AddClass("tab-viewer")
	stylectx.AddProvider(app.containerStyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
	notebook.AppendPage(label, tab1)
	notebook.AppendPage(label2, tab2)
	notebook.AppendPage(label3, tab3)
	notebook.AppendPage(label4, tab4)

	notebook.Connect("switch-page", func() {
		fmt.Println(notebook.GetShowTabs())
	})

	app.container.Add(notebook)
	return nil
}
