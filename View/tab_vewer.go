package View

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func (app *ActionCenterUI) createTabViewerContainer() error {
	// create a new box for the tabs
	box1, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return err
	}

	// add some buttons to box1
	ntab, _ := app.createTab("Notification")
	box1.Add(ntab)
	atab, _ := app.createTab("Ai")
	box1.Add(atab)
	rtab, _ := app.createTab("Radio")
	box1.Add(rtab)
	testtab, _ := app.createTab("Test")
	box1.Add(testtab)

	// set the max height of the box1 to 200
	box1.SetSizeRequest(123123, 2151515151515100)

	box1s, err := box1.GetStyleContext()
	if err != nil {
		return err
	}
	box1.SetName("tab-container")
	box1s.AddProvider(app.containerStyleProvider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	// create another box for some other widgets
	box2, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return err
	}

	// add some other widgets to box2
	label, err := gtk.LabelNew("Some other widgets")
	if err != nil {
		return err
	}
	box2.Add(label)
	box1.SetHAlign(gtk.ALIGN_CENTER)
	box1.SetHExpand(true)
	// add the boxes to the main container
	app.container.PackStart(box1, false, true, 0)
	app.container.PackStart(box2, true, true, 0)

	return nil
}
func (app *ActionCenterUI) createTab(name string) (*gtk.Button, error) {
	button, err := gtk.ButtonNewWithLabel(name)
	if err != nil {
		return nil, err
	}
	// add the CSS class to the style context
	styleCtx, err := button.GetStyleContext()
	if err != nil {
		log.Fatal("Error getting style context for button:", err)
	}
	styleCtx.AddClass("mytab")
	styleCtx.AddProvider(app.containerStyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	button.Connect("clicked", func() {
		fmt.Println(styleCtx.GetColor(gtk.STATE_FLAG_BACKDROP))
	})
	button.SetSizeRequest(140, 100) // magic number should change

	return button, nil
}
