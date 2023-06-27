package View

import "github.com/gotk3/gotk3/gtk"

type RadioTab struct {
	container *gtk.Box
	listbox   *gtk.ListBox
}

func (app *RadioTab) Create() (*gtk.Box, error) {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	label, err := gtk.LabelNew("Radio")
	if err != nil {
		return nil, err
	}
	container.Add(label)
	return container, nil
}
