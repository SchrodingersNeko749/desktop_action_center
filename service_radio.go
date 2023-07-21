package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

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

func (radio *RadioTab) AddFoundStation(station Station) error {
	stationRow, _ := gtk.ListBoxRowNew()
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	favicon := ImgDownload(station.Favicon, Conf.ICON_SIZE*2)
	station.FaviconImage = favicon

	hbox.Add(favicon)
	nameLabel, _ := gtk.LabelNew(station.Name)
	vbox.Add(nameLabel)
	hbox.Add(vbox)
	stationRow.Add(hbox)

	radio.listbox.Add(stationRow)
	radio.foundStations = append(radio.foundStations, station)
	radio.listbox.ShowAll()
	return nil
}
