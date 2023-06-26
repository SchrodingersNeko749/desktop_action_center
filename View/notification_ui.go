package View

import (
	"fmt"

	"github.com/actionCenter/Model"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type NotificationWidget struct {
	container *gtk.Box
	id        int
}
type NotificationList struct {
	container     *gtk.ScrolledWindow
	listBox       *gtk.ListBox
	notifications []NotificationWidget
}

func (app *ActionCenterUI) createNotificationComponent() (*gtk.Box, error) {
	container, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	scrollBox, _ := gtk.ScrolledWindowNew(nil, nil)
	scrollBox.SetHExpand(false)
	scrollBox.SetVExpand(true)
	scrollBox.SetHExpand(false)

	header, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	label, _ := gtk.LabelNew("Notifications")
	clearBtn, _ := gtk.ButtonNewWithLabel("Clear All")
	header.PackStart(label, false, false, 0)
	header.PackEnd(clearBtn, false, true, 1)

	listBox, _ := gtk.ListBoxNew()
	listBox.SetSelectionMode(gtk.SELECTION_SINGLE)

	nlist := NotificationList{
		container: scrollBox,
		listBox:   listBox,
	}
	app.notifications = nlist

	container.Add(header)
	scrollBox.Add(listBox)
	container.Add(scrollBox)

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
	style.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = scrollBox.GetStyleContext()
	style.AddClass("notification-scrollbox")
	style.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = header.GetStyleContext()
	style.AddClass("notification-container-header")
	style.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = scrollBox.GetStyleContext()
	style.AddClass("notification-scrollbox")
	style, _ = listBox.GetStyleContext()
	style.AddClass("notification-scrollbox")
	style.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = clearBtn.GetStyleContext()
	style.AddClass("clear-button")
	style.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return container, nil
}

func (app *ActionCenterUI) AddNotification(n Model.Notification) {
	elementWidth := app.notifications.listBox.GetAllocatedWidth() - ICON_SIZE - HORIZONTAL_SPACING

	widget := NotificationWidget{}
	row, _ := gtk.ListBoxRowNew()
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 20)
	summaryLabel, _ := gtk.LabelNew(n.Summary)
	bodyLabel, _ := gtk.LabelNew(n.Body)
	var icon *gtk.Image = nil

	if _, ok := n.Hints["image-data"]; ok {
		width := n.Hints["image-data"].Value().([]interface{})[0].(int32)
		height := n.Hints["image-data"].Value().([]interface{})[1].(int32)
		rowStride := n.Hints["image-data"].Value().([]interface{})[2].(int32)
		hasAlpha := n.Hints["image-data"].Value().([]interface{})[3].(bool)
		bitsPerSample := n.Hints["image-data"].Value().([]interface{})[4].(int32)

		img := n.Hints["image-data"].Value().([]interface{})[6].([]byte)
		pixbuf, err := gdk.PixbufNewFromData(img, gdk.COLORSPACE_RGB, hasAlpha, int(bitsPerSample), int(width), int(height), int(rowStride))
		icon, err = gtk.ImageNewFromPixbuf(pixbuf)
		if err != nil {
			fmt.Println(err)
		}
	} else if customImagePath, ok := n.Hints["image-path"].Value().(string); ok {
		icon, _ = gtk.ImageNewFromFile(customImagePath)
	} else if n.AppIcon != "" {
		icon, _ = gtk.ImageNewFromIconName(n.AppIcon, gtk.ICON_SIZE_LARGE_TOOLBAR)
	} else {
		icon, _ = gtk.ImageNewFromIconName("gtk-dialog-info", gtk.ICON_SIZE_LARGE_TOOLBAR)
	}
	if icon != nil {
		resize(icon)
		hbox.PackStart(icon, false, false, 0)
	}

	summaryLabel.SetHAlign(gtk.ALIGN_START)
	summaryLabel.SetLineWrap(true)
	summaryLabel.SetMaxWidthChars(1)
	summaryLabel.SetSizeRequest(elementWidth, -1)
	summaryLabel.SetXAlign(0)

	bodyLabel.SetLineWrap(true)
	bodyLabel.SetMaxWidthChars(1)
	bodyLabel.SetSizeRequest(elementWidth, -1)
	bodyLabel.SetHAlign(gtk.ALIGN_START)
	bodyLabel.SetXAlign(0)

	widget.container = hbox
	vbox.PackStart(summaryLabel, false, false, 0)
	vbox.PackStart(bodyLabel, false, false, 0)
	hbox.PackStart(vbox, true, true, 0)
	row.Add(widget.container)
	app.notifications.listBox.Insert(row, 0)

	style, _ := hbox.GetStyleContext()
	style.AddClass("notification-widget")
	style.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = summaryLabel.GetStyleContext()
	style.AddClass("notification-summary")
	style.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = bodyLabel.GetStyleContext()
	style.AddClass("notification-body")
	style.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func (app *ActionCenterUI) clearNotification() {
	for app.notifications.listBox.GetChildren().Length() > 0 {
		app.notifications.listBox.Remove(app.notifications.listBox.GetRowAtIndex(0))
	}
}

func resize(icon *gtk.Image) {
	pixbuf := icon.GetPixbuf()
	if pixbuf == nil {
		theme, _ := gtk.IconThemeGetDefault()
		iconName, _ := icon.GetIconName()
		pixbuf, _ = theme.LoadIconForScale(iconName, ICON_SIZE, 1, gtk.ICON_LOOKUP_FORCE_SIZE)
	}
	scaledPixbuf, _ := pixbuf.ScaleSimple(ICON_SIZE, ICON_SIZE, gdk.INTERP_BILINEAR)
	icon.SetFromPixbuf(scaledPixbuf)
}
