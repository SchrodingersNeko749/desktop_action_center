package View

import (
	"github.com/gotk3/gotk3/gtk"
)

type AiWidget struct {
	container *gtk.Box
	id        int
}
type AiChatList struct {
	container *gtk.ScrolledWindow
	listBox   *gtk.ListBox
	Messages  []NotificationWidget
}

func (app *ActionCenterUI) Create() (*gtk.Box, error) {
	scrollBox, _ := gtk.ScrolledWindowNew(nil, nil)
	container, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)

	container.SetCanFocus(true)
	scrollBox.SetCanFocus(true)
	scrollBox.SetVExpand(true)
	scrollBox.SetHExpand(true)

	label, err := gtk.LabelNew("Ai")
	if err != nil {
		return nil, err
	}
	container.Add(label)

	clearBtn, err := gtk.ButtonNewWithLabel("Clear")
	if err != nil {
		return nil, err
	}
	clearBtn.Connect("clicked", func() {
		app.clearNotification()
	})
	container.Add(clearBtn)

	listBox, _ := gtk.ListBoxNew()
	style, err := listBox.GetStyleContext()
	if err != nil {
		return nil, err
	}
	style.AddClass("notification-container")
	style.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	listBox.SetSelectionMode(gtk.SELECTION_NONE)

	nlist := NotificationList{
		container: scrollBox,
		listBox:   listBox,
	}
	app.notifications = nlist
	scrollBox.Add(listBox)
	container.Add(scrollBox)

	inputBox, err := gtk.EntryNew()
	inputBox.SetEditable(true)
	inputBox.SetCanFocus(true)
	container.Add(inputBox)

	return container, nil
}
