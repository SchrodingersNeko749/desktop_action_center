package View

import (
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

func resize(icon *gtk.Image, width int, height int) {
	// Get the current pixbuf from the image
	pixbuf := icon.GetPixbuf()

	// Scale the pixbuf to the new size
	scaledPixbuf, _ := pixbuf.ScaleSimple(width, height, gdk.INTERP_BILINEAR)

	// Update the image with the scaled pixbuf
	icon.SetFromPixbuf(scaledPixbuf)
}
func (app *ActionCenterUI) newNotificationWidget(n Model.Notification) (*NotificationWidget, error) {
	widget := &NotificationWidget{}

	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		return nil, err
	}
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 20)
	if err != nil {
		return nil, err
	}
	// notification icon
	var icon *gtk.Image
	if customImagePath, ok := n.Hints["image-path"].Value().(string); ok {
		icon, err = gtk.ImageNewFromFile(customImagePath)
		resize(icon, 64, 64)
	} else {
		if n.AppIcon == "" {
			icon, err = gtk.ImageNewFromIconName("gtk-dialog-info", gtk.ICON_SIZE_LARGE_TOOLBAR)
			icon.SetPixelSize(64)

		} else {
			icon, err = gtk.ImageNewFromIconName(n.AppIcon, gtk.ICON_SIZE_LARGE_TOOLBAR)
			icon.SetPixelSize(64)

		}
	}

	if err != nil {
		return nil, err
	}

	summaryLabel, err := gtk.LabelNew(n.Summary)
	if err != nil {
		return nil, err
	}
	summaryLabel.SetHAlign(gtk.ALIGN_START)
	summaryLabel.SetLineWrap(true)
	summaryLabel.SetMaxWidthChars(1)
	summaryLabel.SetSizeRequest(WINDOW_WIDTH-200, -1)
	summaryLabel.SetXAlign(0)

	stylectx, err := summaryLabel.GetStyleContext()
	if err != nil {
		return nil, err
	}
	stylectx.AddClass("notification-summary")
	stylectx.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	bodyLabel, err := gtk.LabelNew(n.Body)
	if err != nil {
		return nil, err
	}
	bodyLabel.SetLineWrap(true)
	bodyLabel.SetMaxWidthChars(1)
	bodyLabel.SetSizeRequest(WINDOW_WIDTH-200, -1)
	bodyLabel.SetHAlign(gtk.ALIGN_START)
	bodyLabel.SetXAlign(0)
	stylectx, err = bodyLabel.GetStyleContext()
	if err != nil {
		return nil, err
	}
	stylectx.AddClass("notification-body")
	stylectx.AddProvider(app.componentStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	hbox.PackStart(icon, false, false, 0)
	vbox.PackStart(summaryLabel, false, false, 0)
	vbox.PackStart(bodyLabel, false, false, 0)

	hbox.PackStart(vbox, true, true, 0)

	widget.container = hbox
	return widget, nil
}
func (app *ActionCenterUI) clearNotification() {
	for app.notifications.listBox.GetChildren().Length() > 0 {
		app.notifications.listBox.Remove(app.notifications.listBox.GetRowAtIndex(0))
	}
}
func (app *ActionCenterUI) createNotificationComponent() (*gtk.Box, error) {
	scrollBox, _ := gtk.ScrolledWindowNew(nil, nil)
	container, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)

	scrollBox.SetVExpand(true)
	scrollBox.SetHExpand(false)

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

	//app.ShowNotifications()
	scrollBox.Add(listBox)
	container.Add(scrollBox)
	return container, nil
}
func (app *ActionCenterUI) LoadNotifications() error {

	app.clearNotification()

	return nil
}
func (app *ActionCenterUI) AddNotification(n Model.Notification) error {
	// make notification widget
	widget, err := app.newNotificationWidget(n)
	if err != nil {
		return err
	}
	// make listbox row
	row, err := gtk.ListBoxRowNew()
	if err != nil {
		return err
	}
	row.Add(widget.container)

	app.notifications.listBox.Add(row)
	app.notifications.listBox.ShowAll()
	return nil
}
