package View

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type NotificationWidget struct {
	*gtk.Box
}

func (app *ActionCenterUI) newNotificationWidget(appIcon string, summary string, body string) (*NotificationWidget, error) {
	widget := &NotificationWidget{}

	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		return nil, err
	}

	var icon *gtk.Image
	if strings.Contains(body, "www.youtube.com") {
		// Set the app icon to the YouTube icon
		appIcon = "youtube"

		// Remove the <a> tag from the body text
		body = strings.Replace(body, "<a href=\"https://www.youtube.com/\">", "", -1)
		body = strings.Replace(body, "</a>", "", -1)
		fmt.Println(body, "string matched")
	}
	if appIcon == "" {
		fmt.Println("empty")
		icon, err = gtk.ImageNewFromIconName("gtk-dialog-info", gtk.ICON_SIZE_LARGE_TOOLBAR)
		if err != nil {
			return nil, err
		}

	} else {
		icon, err = gtk.ImageNewFromIconName(appIcon, gtk.ICON_SIZE_LARGE_TOOLBAR)
		if err != nil {
			return nil, err
		}
	}
	icon.SetPixelSize(64)

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		return nil, err
	}
	vbox.SetHAlign(gtk.ALIGN_START)
	summaryLabel, err := gtk.LabelNew(summary)
	if err != nil {
		return nil, err
	}
	summaryLabel.SetHAlign(gtk.ALIGN_START)

	stylectx, err := summaryLabel.GetStyleContext()
	if err != nil {
		return nil, err
	}
	stylectx.AddClass("notification-summary")
	stylectx.AddProvider(app.containerStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	bodyLabel, err := gtk.LabelNew(body)
	if err != nil {
		return nil, err
	}

	stylectx, err = bodyLabel.GetStyleContext()
	if err != nil {
		return nil, err
	}
	stylectx.AddClass("notification-body")
	stylectx.AddProvider(app.containerStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	hbox.PackStart(icon, false, false, 0)
	vbox.PackStart(summaryLabel, false, false, 0)
	vbox.PackStart(bodyLabel, false, false, 0)

	hbox.PackStart(vbox, true, true, 0)

	widget.Box = hbox

	return widget, nil
}
func (app *ActionCenterUI) createNotificationComponent() (*gtk.Box, error) {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	listBox, err := gtk.ListBoxNew()
	if err != nil {
		return nil, err
	}

	style, err := listBox.GetStyleContext()
	if err != nil {
		return nil, err
	}
	style.AddClass("notification-container")
	style.AddProvider(app.containerStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	// Set selection mode to single
	listBox.SetSelectionMode(gtk.SELECTION_SINGLE)
	notifications, err := app.actionCenter.GetNotifications()
	if err != nil {
		return nil, err
	}

	for _, notification := range notifications {
		widget, err := app.newNotificationWidget(notification.Icon, notification.Summary, notification.Body)
		if err != nil {
			return nil, err
		}

		row, err := gtk.ListBoxRowNew()
		if err != nil {
			return nil, err
		}

		row.Add(widget)
		listBox.Add(row)

		// // Connect the "row-activated" signal to remove the notification from the panel
		// row.Connect("row-activated", func() {
		// 	fmt.Println(row)
		// 	listBox.Remove(row)
		// })
	}
	container.Add(listBox)
	return container, nil
}
