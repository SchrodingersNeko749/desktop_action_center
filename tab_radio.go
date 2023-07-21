package main

import (
	"fmt"

	"github.com/fhs/gompd/mpd"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type RadioTab struct {
	container         *gtk.Box
	listbox           *gtk.ListBox
	directoryServerIp string
	foundStations     []Station
	currentStation    Station
	mpdClient         mpd.Client
}

func (radio *RadioTab) Create() (*gtk.Box, error) {
	var err error

	if radio.container == nil {
		radio.container, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
		if err != nil {
			return nil, err
		}
	}
	children := radio.container.GetChildren()
	for l := children; l != nil; l = l.Next() {
		child := l.Data().(*gtk.Widget)
		radio.container.Remove(child)
	}
	mpdclient, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		errorLabel, _ := gtk.LabelNew("Error encountered: " + err.Error())
		retryButton, _ := gtk.ButtonNewWithLabel("Retry connection")
		retryButton.Connect("clicked", func() {
			radio.Create()
		})
		radio.container.Add(errorLabel)
		radio.container.Add(retryButton)
		radio.container.ShowAll()
		return radio.container, nil
	}
	radio.mpdClient = *mpdclient
	playerBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	label, _ := gtk.LabelNew("loading")
	go func() {
		glib.TimeoutAdd(uint(1000), func() bool {
			if radio.mpdClient.Ping() != nil {
				radio.Create()
			}
			song, _ := radio.mpdClient.CurrentSong()
			if song != nil {
				label.SetText(song["Title"])
			}
			return true
		})
	}()
	commandBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	commandBox.SetHAlign(gtk.ALIGN_CENTER)
	stopButton, _ := gtk.ButtonNewWithLabel("■")
	stopButton.Connect("clicked", func() {
		mpdclient.Stop()
	})
	playButton, _ := gtk.ButtonNewWithLabel("")
	playButton.Connect("clicked", func() {

	})
	voteButton, _ := gtk.ButtonNewWithLabel("♥")
	voteButton.Connect("clicked", func() {
		fmt.Println("favorited")
	})

	commandBox.Add(stopButton)
	commandBox.Add(playButton)
	commandBox.Add(voteButton)

	playerBox.Add(label)
	stationImg := ImgFromTheme("radio", 128)
	playerBox.Add(stationImg)
	playerBox.Add(commandBox)
	volumeBox, _ := radio.createMpdVolumeComponent()
	playerBox.Add(volumeBox)
	if err != nil {
		return nil, err
	}
	// Advanced search
	inputBox, _ := gtk.EntryNew()
	inputBox.SetPlaceholderText("Search radio stations here")
	inputBox.SetIconFromIconName(gtk.ENTRY_ICON_PRIMARY, "search")
	listBox, _ := gtk.ListBoxNew()
	advancedSearchBoxExpander, _ := gtk.ExpanderNew("Advanced Search")
	advancedSearchBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	inputBoxHintLabel, _ := gtk.LabelNew("Search new stations")
	advancedSearchBox.Add(inputBoxHintLabel)

	listBox.Connect("row-selected", func() {
		selected := listBox.GetSelectedRow()
		if selected != nil {
			radio.mpdClient.Clear()
			radio.currentStation = radio.foundStations[selected.GetIndex()]
			if radio.currentStation.Favicon == "" {
				stationImg.SetFromIconName("radio", 192)
			} else {
				newImg := ImgDownload(radio.currentStation.Favicon, 192)
				stationImg.SetFromPixbuf(newImg.GetPixbuf())
			}
			playerBox.ShowAll()
			label.SetText(radio.currentStation.Name)

			radio.mpdClient.Add(radio.currentStation.URL)
			radio.mpdClient.Play(0)
		}
	})
	inputBox.Connect("activate", func() {
		text, _ := inputBox.GetText()
		inputBox.SetText("")
		radio.foundStations = []Station{}
		stations := radio.AdvancedStationSearch(text, "", 5)
		inputBoxHintLabel.SetText(fmt.Sprintf("Search result for %s, %d Stations found", text, len(stations)))
		for listBox.GetChildren().Length() > 0 {
			listBox.Remove(listBox.GetRowAtIndex(0))
		}

		for _, s := range stations {
			radio.AddFoundStation(s)
		}
	})

	style, _ := inputBox.GetStyleContext()
	style.AddClass("radio-inputbox")
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	advancedSearchBox.Add(listBox)
	advancedSearchBox.Add(inputBox)
	advancedSearchBoxExpander.Add(advancedSearchBox)

	radio.listbox = listBox
	radio.container.Add(playerBox)
	radio.container.PackStart(advancedSearchBoxExpander, false, true, 0)
	radio.container.ShowAll()
	return radio.container, nil
}

func (radio *RadioTab) createMpdVolumeComponent() (*gtk.Box, error) {
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	hbox.SetHAlign(gtk.ALIGN_CENTER)
	hbox.SetVAlign(gtk.ALIGN_END)

	style, _ := hbox.GetStyleContext()
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style.AddClass("scale-box")

	volumeBar, _ := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	volumeBar.SetHExpand(true)
	volumeBar.SetSizeRequest(500, 20)

	volumeBar.SetValue(50)
	radio.mpdClient.SetVolume(50)
	style, _ = volumeBar.GetStyleContext()
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	volumeBar.Connect("value-changed", func() {
		v := volumeBar.GetValue()
		radio.mpdClient.SetVolume(int(v))
	})

	label, _ := gtk.LabelNew("")

	hbox.PackStart(volumeBar, true, true, 0)
	hbox.PackEnd(label, true, true, 0)

	return hbox, nil
}

func (radio *RadioTab) CreateStationWidget() error {
	stationRow, _ := gtk.ListBoxRowNew()
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	favicon := ImgFromTheme("radio", Conf.ICON_SIZE)
	hbox.Add(favicon)
	nameLabel, _ := gtk.LabelNew("Loading ... ")
	vbox.Add(nameLabel)
	hbox.Add(vbox)
	stationRow.Add(hbox)

	radio.listbox.Add(stationRow)
	radio.listbox.ShowAll()
	return nil
}
