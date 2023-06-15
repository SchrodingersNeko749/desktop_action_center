package View

import "github.com/gotk3/gotk3/gtk"

func (app *ActionCenterUI) createAiComponent() (*gtk.Box, error) {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	label, err := gtk.LabelNew("Ai")
	if err != nil {
		return nil, err
	}
	container.Add(label)
	return container, nil
}
