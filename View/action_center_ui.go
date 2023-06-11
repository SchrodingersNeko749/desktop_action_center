package View

import (
	"log"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

func (app *ActionCenter) addLabel() error {
	// Create the label widget
	l, err := gtk.LabelNew(time.Now().Format("Mon 3:04 PM"))
	l.SetName("header")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	l2, err := gtk.LabelNew("Action Center")
	l2.SetName("header")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	// Create the vertical box container and add the label widget to it
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}
	vbox.Add(l)
	vbox.Add(l2)

	// Create a new CSS provider and load the CSS file
	provider, err := gtk.CssProviderNew()
	if err != nil {
		return err
	}

	if err := provider.LoadFromPath("/home/neko/Projects/programming/go/desktop_action_center/window.css"); err != nil {
		return err
	}

	// Apply the CSS provider to the label widget's style context
	lStyle, err := l.GetStyleContext()
	if err != nil {
		return err
	}
	lStyle.AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	// Apply the CSS provider to the box container's style context
	vboxStyle, err := vbox.GetStyleContext()
	if err != nil {
		return err
	}
	vboxStyle.AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
	vbox.SetHAlign(gtk.ALIGN_START)
	// Add the box container to the window
	app.win.Add(vbox)

	return nil
}
