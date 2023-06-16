package View

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type NotificationWidget struct {
	container *gtk.Box
	id        int
}
type NotificationList struct {
	container     *gtk.Box
	listBox       *gtk.ListBox
	notifications []NotificationWidget
}

func getAppIcon(appIcon string) (*gtk.Image, error) {
	var icon *gtk.Image
	var err error
	if appIcon == "" {
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
	return icon, nil
}

func (app *ActionCenterUI) newNotificationWidget(appIcon string, summary string, body string) (*NotificationWidget, error) {
	if strings.Contains(body, "www.youtube.com") {
		// Set the app icon to the YouTube icon
		appIcon = "youtube"

		// Remove the <a> tag from the body text
		body = strings.Replace(body, "<a href=\"https://www.youtube.com/\">www.youtube.com</a>", "", -1)
		// Remove empty lines from the body text
		body = strings.TrimSpace(body)
	}
	widget := &NotificationWidget{}

	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		return nil, err
	}
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 20)
	if err != nil {
		return nil, err
	}
	vbox.SetHAlign(gtk.ALIGN_START)
	icon, err := getAppIcon(appIcon)
	if err != nil {
		return nil, err
	}
	vbox.SetHAlign(gtk.ALIGN_START)
	summaryLabel, err := gtk.LabelNew(summary)
	if err != nil {
		return nil, err
	}
	summaryLabel.SetHAlign(gtk.ALIGN_START)
	summaryLabel.SetLineWrap(true)
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
	bodyLabel.SetLineWrap(true)
	bodyLabel.SetMaxWidthChars(1)
	bodyLabel.SetSizeRequest(WINDOW_WIDTH-64, -1)
	bodyLabel.SetHAlign(gtk.ALIGN_START)
	bodyLabel.SetXAlign(0)
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

	widget.container = hbox

	return widget, nil
}
func (app *ActionCenterUI) clearNotification() {
	for app.notifications.listBox.GetChildren().Length() > 0 {
		app.notifications.listBox.Remove(app.notifications.listBox.GetRowAtIndex(0))
	}
}
func (app *ActionCenterUI) createNotificationComponent() (*gtk.Box, error) {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
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

	listBox.Connect("row-activated", func() {
		// Get the selected row index
		selected := app.notifications.listBox.GetSelectedRow()
		fmt.Println(selected.GetPreferredWidth())

	})

	nlist := NotificationList{
		container: container,
		listBox:   listBox,
	}
	app.notifications = nlist

	app.ShowNotifications()
	container.Add(listBox)
	return container, nil
}
func (app *ActionCenterUI) ShowNotifications() error {

	app.clearNotification()
	notifications, err := app.actionCenter.GetNotifications()
	if err != nil {
		return err
	}

	// Delay adding the new rows until the GTK event loop has finished updating the user interface
	for _, notification := range notifications {
		err := app.AddNotification(notification.Icon, notification.Summary, notification.Body)
		if err != nil {
			return err
		}
		app.notifications.listBox.ShowAll()

	}
	return nil
}
func (app *ActionCenterUI) AddNotification(icon string, summary string, body string) error {
	// make notification widget
	widget, err := app.newNotificationWidget(icon, summary, body)
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

	return nil
}
