package main

import "encoding/json"

type Station struct {
	ChangeUUID                string  `json:"changeuuid"`
	StationUUID               string  `json:"stationuuid"`
	Name                      string  `json:"name"`
	URL                       string  `json:"url"`
	URLResolved               string  `json:"url_resolved"`
	Homepage                  string  `json:"homepage"`
	Favicon                   string  `json:"favicon"`
	Tags                      string  `json:"tags"`
	Country                   string  `json:"country"`
	CountryCode               string  `json:"countrycode"`
	State                     string  `json:"state"`
	Language                  string  `json:"language"`
	LanguageCodes             string  `json:"languagecodes"`
	Votes                     int     `json:"votes"`
	LastChangeTime            string  `json:"lastchangetime"`
	LastChangeTimeISO8601     string  `json:"lastchangetime_iso8601"`
	Codec                     string  `json:"codec"`
	Bitrate                   int     `json:"bitrate"`
	HLS                       int     `json:"hls"`
	LastCheckOK               int     `json:"lastcheckok"`
	LastCheckTime             string  `json:"lastchecktime"`
	LastCheckTimeISO8601      string  `json:"lastchecktime_iso8601"`
	LastCheckOKTime           string  `json:"lastcheckoktime"`
	LastCheckOKTimeISO8601    string  `json:"lastcheckoktime_iso8601"`
	LastLocalCheckTime        string  `json:"lastlocalchecktime"`
	LastLocalCheckTimeISO8601 string  `json:"lastlocalchecktime_iso8601"`
	ClickTimestamp            string  `json:"clicktimestamp"`
	ClickTimestampISO8601     *string `json:"clicktimestamp_iso8601"`
	ClickCount                int     `json:"clickcount"`
	ClickTrend                int     `json:"clicktrend"`
	SSLError                  int     `json:"ssl_error"`
	GeoLat                    float64 `json:"geo_lat"`
	GeoLong                   float64 `json:"geo_long"`
	HasExtendedInfo           bool    `json:"has_extended_info"`
}
type requestData struct {
	Name        string `json:"name"`
	CountryCode string `json:"countrycode"`
	Limit       int    `json:"limit"`
}

func ParseJsonToRadios(jsonData []byte) ([]Station, error) {
	parsedStation := []Station{}
	err := json.Unmarshal(jsonData, &parsedStation)
	if err != nil {
		panic(err)
	}
	return parsedStation, nil
}
