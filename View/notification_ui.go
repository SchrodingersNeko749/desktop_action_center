package View

import "github.com/gotk3/gotk3/gtk"

type notification_ui struct {
	box gtk.Box
}

func (n *notification_ui) GetComponent() (error, gtk.Box) {
	return nil, n.box
}
