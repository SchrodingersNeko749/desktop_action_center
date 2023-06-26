package View

import (
	"github.com/actionCenter/Data"
	"github.com/gotk3/gotk3/gtk"
)

func (app *ActionCenterUI) createTabViewerContainer(configWidget Data.WidgetConfig) (*gtk.Box, *gtk.Notebook, error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	box.SetSizeRequest(WINDOW_WIDTH, -1)
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
	stylectx.AddProvider(app.componentStyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
	notebook.SetCurrentPage(0)

	box.Add(notebook)
	return box, notebook, nil
}

func (app *ActionCenterUI) addTab(notebook *gtk.Notebook, tabLabelString string, page *gtk.Box) {
	tabLabel, _ := gtk.LabelNew(tabLabelString)
	tabLabel.SetSizeRequest(50, 50)
	notebook.AppendPage(page, tabLabel)
}
