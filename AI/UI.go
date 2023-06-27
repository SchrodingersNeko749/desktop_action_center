package AI

import (
	"os/user"

	"github.com/actionCenter/Data"
	"github.com/actionCenter/Model"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type UI struct {
	container *gtk.Box
	listBox   *gtk.ListBox
	Messages  []Model.NotificationWidget
	Service   *Service
}

func (ai *UI) Create() (*gtk.Box, error) {
	ai.Service = &Service{}
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

func (ai *UI) AddMessage(msg string) {
	username, _ := user.Current()
	widget, _ := Model.CreateNotificationComponent(Model.NewNotification("Her.st LLaMa", 0, "", username.Username, msg, nil, nil, 0))
	responseWidget, body := Model.CreateNotificationComponent(Model.NewNotification("Her.st LLaMa", 0, "", "AI", "", nil, nil, 0))

	glib.IdleAdd(func() {
		ai.listBox.Add(widget)
		ai.listBox.Add(responseWidget)
		ai.listBox.ShowAll()
	})

	prompt := ai.Service.GeneratePrompt("instruction", msg, 1024, "guanaco-7B.ggmlv3.q4_0.bin", false, false)
	go ai.Service.RunInference(prompt, body)
}
