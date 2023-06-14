package View

import "github.com/gotk3/gotk3/gtk"

func (app *ActionCenterUI) createScreenCaptureComponent() (*gtk.Box, error) {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	label, err := gtk.LabelNew("Wifi")
	if err != nil {
		return nil, err
	}
	container.Add(label)
	return container, nil
}
