package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func (radio *RadioTab) SetRadioDirectoryServerIP(host string) error {
	addrs, err := net.LookupHost(host)
	var servers []string

	if err != nil {
		return err
	}

	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip != nil && ip.To4() != nil {
			servers = append(servers, ip.String())
		}
	}
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(len(servers))
	radio.directoryServerIp = servers[num]
	return nil
}
func (radio *RadioTab) AdvancedStationSearch(name string, countryCode string, limit int) []Station {
	if limit == -1 {
		limit = 10
	}
	url := "http://" + radio.directoryServerIp + "/json/stations/search"
	data := requestData{Name: name, CountryCode: countryCode, Limit: limit}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/plain")

	// send the HTTP request and print the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}
	stations, err := ParseJsonToRadios(body)
	if err != nil {
		panic(err)
	}
	return stations

}
func (radio *RadioTab) AddStation(station Station) error {
	stationRow, _ := gtk.ListBoxRowNew()
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	var img *gtk.Image
	var err error
	if station.Favicon != "" {
		img, err = DownloadFavicon(station.Favicon)
		if err != nil {
			img, _ = gtk.ImageNewFromIconName("radio", gtk.ICON_SIZE_LARGE_TOOLBAR)
			img.SetPixelSize(64)
		}
		if img != nil {
			hbox.Add(img)
		}

	} else {
		img, _ = gtk.ImageNewFromIconName("radio", gtk.ICON_SIZE_LARGE_TOOLBAR)
		img.SetPixelSize(64)
	}
	station.FaviconImage = img
	Resize(img, Conf.ICON_SIZE)
	hbox.Add(img)
	nameLabel, _ := gtk.LabelNew(station.Name)
	vbox.Add(nameLabel)
	hbox.Add(vbox)
	stationRow.Add(hbox)

	radio.listbox.Add(stationRow)
	radio.foundStations = append(radio.foundStations, station)
	radio.listbox.ShowAll()
	return nil
}
func Play() {

}

func DownloadFavicon(faviconUrl string) (*gtk.Image, error) {
	fmt.Println(faviconUrl)
	response, err := http.Get(faviconUrl)
	if err != nil || response.StatusCode != 200 {
		fmt.Println("errors found")
		return nil, err
	}
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, response.Body)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(buf)
	if err != nil {
		return nil, err
	}

	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y
	stride := width * 4
	pixels := make([]byte, height*stride)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			p := (y*width + x) * 4
			pixels[p] = byte(r)
			pixels[p+1] = byte(g)
			pixels[p+2] = byte(b)
			pixels[p+3] = byte(a)
		}
	}

	pixbuf, err := gdk.PixbufNewFromData(pixels, gdk.COLORSPACE_RGB, true, 8, width, height, stride)
	if err != nil {
		return nil, err
	}

	gtkimg, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		return nil, err
	}
	return gtkimg, nil
}
