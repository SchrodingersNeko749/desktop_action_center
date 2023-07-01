package main

import "github.com/gotk3/gotk3/gtk"

type ScreenTab struct {
	container *gtk.Box
	listbox   *gtk.ListBox
}

func (app *ScreenTab) Create() (*gtk.Box, error) {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	label, err := gtk.LabelNew("Wifi")
	if err != nil {
		return nil, err
	}
	container.Add(label)
	return container, nil
}
