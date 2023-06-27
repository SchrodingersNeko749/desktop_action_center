package View

import (
	"github.com/actionCenter/Data"
	"github.com/actionCenter/Model"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type AiWidget struct {
	container *gtk.Box
	id        int
}
type AITab struct {
	container *gtk.Box
	listBox   *gtk.ListBox
	Messages  []Model.NotificationWidget
}

func (ai *AITab) Create() (*gtk.Box, error) {
	container, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	scrollBox, _ := gtk.ScrolledWindowNew(nil, nil)
	scrollBox.SetHExpand(true)
	scrollBox.SetVExpand(true)

	header, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	label, _ := gtk.LabelNew("Her.st LLaMa API")
	clearBtn, _ := gtk.ButtonNewWithLabel("Clear All")
	listBox, _ := gtk.ListBoxNew()
	inputBox, _ := gtk.EntryNew()

	listBox.SetSelectionMode(gtk.SELECTION_SINGLE)
	header.PackStart(label, false, false, 0)
	header.PackEnd(clearBtn, false, true, 1)

	ai.container = container
	ai.listBox = listBox

	container.Add(header)
	scrollBox.Add(listBox)
	container.Add(scrollBox)
	container.Add(inputBox)

	clearBtn.Connect("clicked", func() {
	})
	listBox.Connect("row-selected", func() {
		selected := listBox.GetSelectedRow()
		if selected == nil {
			return
		}
		glib.IdleAdd(func() {
			listBox.Remove(selected)
		})
	})
	inputBox.Connect("activate", func() {
		text, _ := inputBox.GetText()
		inputBox.SetText("")
		ai.AddMessage(text)
	})

	style, _ := container.GetStyleContext()
	style.AddClass("ai-container")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = scrollBox.GetStyleContext()
	style.AddClass("ai-scrollbox")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = header.GetStyleContext()
	style.AddClass("ai-container-header")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = inputBox.GetStyleContext()
	style.AddClass("ai-inputbox")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return container, nil
}

func (ai *AITab) AddMessage(msg string) {
	elementWidth := ai.listBox.GetAllocatedWidth() - Data.Conf.ICON_SIZE - Data.Conf.HORIZONTAL_SPACING
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 20)
	row, _ := gtk.ListBoxRowNew()
	summaryLabel, _ := gtk.LabelNew(msg)
	bodyLabel, _ := gtk.LabelNew(msg)

	summaryLabel.SetHAlign(gtk.ALIGN_START)
	summaryLabel.SetLineWrap(true)
	summaryLabel.SetMaxWidthChars(1)
	summaryLabel.SetSizeRequest(elementWidth, -1)
	summaryLabel.SetXAlign(0)

	bodyLabel.SetLineWrap(true)
	bodyLabel.SetMaxWidthChars(1)
	bodyLabel.SetSizeRequest(elementWidth, -1)
	bodyLabel.SetHAlign(gtk.ALIGN_START)
	bodyLabel.SetXAlign(0)

	row.Add(hbox)
	vbox.PackStart(summaryLabel, false, false, 0)
	vbox.PackStart(bodyLabel, false, false, 0)
	hbox.PackStart(vbox, true, true, 0)
	ai.listBox.Insert(row, 0)

	style, _ := hbox.GetStyleContext()
	style.AddClass("notification-widget")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	stylectx, _ := summaryLabel.GetStyleContext()
	stylectx.AddClass("ai-summary")
	stylectx.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	stylectx, _ = bodyLabel.GetStyleContext()
	stylectx.AddClass("ai-body")
	stylectx.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	ai.listBox.ShowAll()
}
