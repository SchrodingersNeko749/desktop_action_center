package main

import (
	"os/user"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type AiTab struct {
	container *gtk.Box
	listBox   *gtk.ListBox
}

func (ai *AiTab) Create() (*gtk.Box, error) {
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
		glib.IdleAdd(func() {
			for listBox.GetChildren().Length() > 0 {
				listBox.Remove(listBox.GetRowAtIndex(0))
			}
		})
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
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = scrollBox.GetStyleContext()
	style.AddClass("ai-scrollbox")
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = header.GetStyleContext()
	style.AddClass("ai-container-header")
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = inputBox.GetStyleContext()
	style.AddClass("ai-inputbox")
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return container, nil
}

func (ai *AiTab) AddMessage(msg string) {
	username, _ := user.Current()
	widget, _ := CreateNotificationComponent(NewNotification("Her.st LLaMa", 0, "", username.Username, msg, nil, nil, 0, nil))
	responseWidget, body := CreateNotificationComponent(NewNotification("Her.st LLaMa", 0, "", "AI", "", nil, nil, 0, nil))

	glib.IdleAdd(func() {
		ai.listBox.Add(widget)
		ai.listBox.Add(responseWidget)
		ai.listBox.ShowAll()
	})

	prompt := GeneratePrompt("instruction", msg, 1024, "nous-hermes-13b.ggmlv3.q4_0.bin", false, false)
	go RunInference(prompt, body)
}
