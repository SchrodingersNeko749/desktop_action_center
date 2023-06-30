package Radio

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

	"github.com/actionCenter/Data"
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

	stationWidget, _ := gtk.ListBoxRowNew()
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	fmt.Print("start")
	img, err := GetFavicon(station.Favicon)
	if err != nil {
		return err
	}
	fmt.Print("finish")

	hbox.Add(img)

	hbox.Add(vbox)
	stationWidget.Add(hbox)

	radio.listbox.Add(stationWidget)
	radio.listbox.ShowAll()
	return nil
}
func Play() {

}
func Stop() {

}
func GetFavicon(faviconUrl string) (*gtk.Image, error) {
	response, err := http.Get(faviconUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	img, err := png.Decode(response.Body)
	if err != nil {
		return nil, err
	}
	// Convert the image to a GTK pixbuf
	pixbuf, err := gdk.PixbufNewFromBytes(buf.Bytes(), gdk.COLORSPACE_RGB, false, 8, img.Bounds().Size().X, img.Bounds().Size().Y, 0)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	gtkimg, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		return nil, err
	}
	return gtkimg, nil
}

// Duplicate code
func Resize(icon *gtk.Image) {
	pixbuf := icon.GetPixbuf()
	if pixbuf == nil {
		theme, _ := gtk.IconThemeGetDefault()
		iconName, _ := icon.GetIconName()
		pixbuf, _ = theme.LoadIconForScale(iconName, Data.Conf.ICON_SIZE, 1, gtk.ICON_LOOKUP_FORCE_SIZE)
	}
	scaledPixbuf, _ := pixbuf.ScaleSimple(Data.Conf.ICON_SIZE, Data.Conf.ICON_SIZE, gdk.INTERP_BILINEAR)
	icon.SetFromPixbuf(scaledPixbuf)
}
