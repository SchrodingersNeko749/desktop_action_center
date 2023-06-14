package View

import "github.com/gotk3/gotk3/gtk"

type NotificationWidget struct {
	*gtk.Box
}

func NewNotificationWidget(appIcon string, summary string, body string) (*NotificationWidget, error) {
	widget := &NotificationWidget{}

	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		return nil, err
	}

	icon, err := gtk.ImageNewFromIconName(appIcon, gtk.ICON_SIZE_LARGE_TOOLBAR)
	if err != nil {
		return nil, err
	}

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		return nil, err
	}
	summaryLabel, err := gtk.LabelNew(summary)
	if err != nil {
		return nil, err
	}

	bodyLabel, err := gtk.LabelNew(body)
	if err != nil {
		return nil, err
	}

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
	notifications := []struct {
		appIcon string
		summary string
		body    string
	}{
		{"chromium", "New tab opened", "https://www.example.com"},
		{"telegram", "New message", "John: Hello! How are you?"},
		{"email", "New email", "From: Jane Doe <jane@example.com> Subject: Hello World"},
		// ...
	}

	for _, notification := range notifications {
		widget, err := NewNotificationWidget(notification.appIcon, notification.summary, notification.body)
		if err != nil {
			return nil, err
		}

		row, err := gtk.ListBoxRowNew()
		if err != nil {
			return nil, err
		}

		row.Add(widget)
		listBox.Add(row)

		// Connect the "row-activated" signal to remove the notification from the panel
		row.Connect("activate", func() {
			listBox.Remove(row)
		})
	}
	container.Add(listBox)
	return container, nil
}
