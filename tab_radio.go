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
	mpdclient, _ := mpd.Dial("tcp", "localhost:6600")
	radio.mpdClient = *mpdclient
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	playerBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	label, _ := gtk.LabelNew("loading")
	go func() {
		glib.TimeoutAdd(uint(1000), func() bool {
			song, _ := radio.mpdClient.CurrentSong()
			if song != nil {
				/*
					song["XXX"] =
						"file": "http://listen.uturnradio.com/dubstep_32"
						"Title": "Mendum - Forsaken ft. Brenton Mattheus"
						"Name": "Uturn Radio: Dubstep Music"
						"Pos": "0"
						"Id": "8"
				*/
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
	stationImg, _ := gtk.ImageNewFromIconName("radio", gtk.ICON_SIZE_LARGE_TOOLBAR)
	stationImg.SetPixelSize(128)
	playerBox.Add(stationImg)
	playerBox.Add(commandBox)
	volumeBox, _ := radio.createMpdVolumeComponent()
	playerBox.Add(volumeBox)
	if err != nil {
		return nil, err
	}
	// Advanced search
	inputBox, _ := gtk.EntryNew()
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
				stationImg.SetFromIconName("radio", 180)
			} else {
				newImage := radio.currentStation.FaviconImage
				stationImg.SetFromPixbuf(newImage.GetPixbuf())
				Resize(stationImg, 180)
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
	radio.container = container
	container.Add(playerBox)
	container.PackStart(advancedSearchBoxExpander, false, true, 0)

	return container, nil
}
func (radio *RadioTab) createMpdVolumeComponent() (*gtk.Box, error) {
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	hbox.SetHAlign(gtk.ALIGN_CENTER)
	hbox.SetVAlign(gtk.ALIGN_END)

	style, _ := hbox.GetStyleContext()
	style.AddProvider(StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style.AddClass("scale-box")

	volumeBar, err := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
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

	return hbox, err
}

/*	*/
