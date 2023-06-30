package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type NotificationTab struct {
	win       *gtk.Window
	Container *gtk.Box
	ListBox   *gtk.ListBox
}
type NotificationWidget struct {
	container *gtk.Box
	id        int
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

	app.Container = container
	app.ListBox = listBox

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
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = scrollBox.GetStyleContext()
	style.AddClass("notification-scrollbox")
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = header.GetStyleContext()
	style.AddClass("notification-container-header")
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = scrollBox.GetStyleContext()
	style.AddClass("notification-scrollbox")
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = clearBtn.GetStyleContext()
	style.AddClass("clear-button")
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return container, nil
}

func (app *NotificationTab) AddNotification(widget *gtk.ListBoxRow) {
	app.ListBox.Insert(widget, 0)
	if app.win.GetVisible() {
		app.Container.ShowAll()
	}
}

func (app *NotificationTab) clearNotification() {
	for app.ListBox.GetChildren().Length() > 0 {
		app.ListBox.Remove(app.ListBox.GetRowAtIndex(0))
	}
}

func (n *Notification) RemoveHyperLinkFromBody() {
	re := regexp.MustCompile(`<a.*?>(.*?)</a>`)
	n.Body = re.ReplaceAllString(n.Body, "")

	// Remove empty lines
	lines := strings.Split(n.Body, "\n")
	var filteredLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filteredLines = append(filteredLines, line)
		}
	}
	n.Body = strings.Join(filteredLines, "\n")
}

func CreateNotificationComponent(n Notification) (*gtk.ListBoxRow, *gtk.Label) {
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
	summaryLabel.SetSelectable(true)

	bodyLabel.SetHAlign(gtk.ALIGN_START)
	bodyLabel.SetLineWrap(true)
	bodyLabel.SetSelectable(true)

	vbox.PackStart(summaryLabel, false, false, 0)
	vbox.PackStart(bodyLabel, false, false, 0)
	hbox.PackStart(vbox, false, true, 0)

	container, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	container.PackStart(hbox, false, false, 0)
	container.PackStart(vbox, false, false, 0)

	row.SetHAlign(gtk.ALIGN_FILL)
	row.Add(container)

	style, _ := hbox.GetStyleContext()
	style.AddClass("notification-widget")
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = summaryLabel.GetStyleContext()
	style.AddClass("notification-summary")
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = bodyLabel.GetStyleContext()
	style.AddClass("notification-body")
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return row, bodyLabel
}

func resize(icon *gtk.Image) {
	pixbuf := icon.GetPixbuf()
	if pixbuf == nil {
		theme, _ := gtk.IconThemeGetDefault()
		iconName, _ := icon.GetIconName()
		pixbuf, _ = theme.LoadIconForScale(iconName, Conf.ICON_SIZE, 1, gtk.ICON_LOOKUP_FORCE_SIZE)
	}
	scaledPixbuf, _ := pixbuf.ScaleSimple(Conf.ICON_SIZE, Conf.ICON_SIZE, gdk.INTERP_BILINEAR)
	icon.SetFromPixbuf(scaledPixbuf)
}
