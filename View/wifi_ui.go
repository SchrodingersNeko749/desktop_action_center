package View

import (
	"github.com/gotk3/gotk3/gtk"
)

func (app *ActionCenterUI) createWifiComponent() (*gtk.Box, error) {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	container.SetHExpand(false)
	if err != nil {
		return nil, err
	}
	label, err := gtk.LabelNew("wifi")

	if err != nil {
		return nil, err
	}
	container.Add(label)
	return container, nil
}
