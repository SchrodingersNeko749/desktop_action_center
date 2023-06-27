package View

import (
	"github.com/actionCenter/Data"
	"github.com/gotk3/gotk3/gtk"
)

type NotificationTab struct {
	win       *gtk.Window
	container *gtk.Box
	listBox   *gtk.ListBox
}

func (app *NotificationTab) Create(window gtk.Window) (*gtk.Box, error) {
	app.win = &window
	container, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	scrollBox, _ := gtk.ScrolledWindowNew(nil, nil)
	container.SetSizeRequest(app.win.GetAllocatedWidth(), app.win.GetAllocatedHeight())
	container.SetVExpand(true)
	container.SetHExpand(true)
	scrollBox.SetHExpand(true)
	scrollBox.SetVExpand(true)

	header, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	label, _ := gtk.LabelNew("Notifications")
	clearBtn, _ := gtk.ButtonNewWithLabel("Clear All")
	header.PackStart(label, false, false, 0)
	header.PackEnd(clearBtn, false, true, 1)

	listBox, _ := gtk.ListBoxNew()
	listBox.SetSelectionMode(gtk.SELECTION_SINGLE)

	app.container = container
	app.listBox = listBox

	container.Add(header)
	container.Add(scrollBox)
	scrollBox.Add(listBox)

	clearBtn.Connect("clicked", func() {
		app.clearNotification()
	})

	listBox.Connect("row-selected", func() {
		selected := listBox.GetSelectedRow()
		if selected != nil {
			listBox.Remove(selected)
		}
	})

	style, _ := container.GetStyleContext()
	style.AddClass("notification-container")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = scrollBox.GetStyleContext()
	style.AddClass("notification-scrollbox")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = header.GetStyleContext()
	style.AddClass("notification-container-header")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = scrollBox.GetStyleContext()
	style.AddClass("notification-scrollbox")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = clearBtn.GetStyleContext()
	style.AddClass("clear-button")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return container, nil
}

func (app *NotificationTab) AddNotification(widget *gtk.ListBoxRow) {
	app.listBox.Insert(widget, 0)
	app.container.ShowAll()
}

func (app *NotificationTab) clearNotification() {
	for app.listBox.GetChildren().Length() > 0 {
		app.listBox.Remove(app.listBox.GetRowAtIndex(0))
	}
}
