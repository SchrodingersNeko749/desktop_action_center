package main

import (
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Header struct {
	container *gtk.Box
}

func (app *Header) Create() *gtk.Box {
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	clockLabel, _ := gtk.LabelNew(time.Now().Format("Mon 3:04 PM | Jun 2"))

	go func() {
		glib.TimeoutAdd(uint(30000), func() bool {
			clockLabel.SetText(time.Now().Format("Mon 3:04 PM | Jun 2"))
			return true
		})
	}()

	clockLabel.SetName("clock")
	vbox.SetHAlign(gtk.ALIGN_START)
	vbox.PackStart(clockLabel, true, true, 0)

	lStyle, _ := clockLabel.GetStyleContext()
	lStyle.AddProvider(StyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
	vboxStyle, _ := vbox.GetStyleContext()
	vboxStyle.AddProvider(StyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	return vbox
}
