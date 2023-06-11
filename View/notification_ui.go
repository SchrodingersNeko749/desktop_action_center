package View

import "github.com/gotk3/gotk3/gtk"

type NotificationCenterUI struct {
	box gtk.Box
}

func (n *NotificationCenterUI) GetComponent() (error, gtk.Box) {
	return nil, n.box
}
