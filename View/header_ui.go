package View

import (
	"time"

	"github.com/actionCenter/Data"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type HeaderUI struct {
	container *gtk.Box
}

func (app *HeaderUI) Create() (*gtk.Box, error) {
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	clockLabel, err := gtk.LabelNew("")
	go func() {
		glib.TimeoutAdd(uint(1000), func() bool {
			clockLabel.SetText(time.Now().Format("Mon 3:04 PM"))
			return true
		})
	}()
	clockLabel.SetName("clock")
	if err != nil {
		return nil, err
	}

	lStyle, err := clockLabel.GetStyleContext()
	if err != nil {
		return nil, err
	}
	lStyle.AddProvider(Data.StyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	vboxStyle, err := vbox.GetStyleContext()
	if err != nil {
		return nil, err
	}
	vboxStyle.AddProvider(Data.StyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
	vbox.SetHAlign(gtk.ALIGN_START)
	vbox.PackStart(clockLabel, true, true, 0)
	return vbox, nil
}
