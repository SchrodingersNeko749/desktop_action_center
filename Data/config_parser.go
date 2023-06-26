package Data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"

	"github.com/gotk3/gotk3/gtk"
)

var Conf *Config
var WidgetConfs []WidgetConfig
var StyleProvider *gtk.CssProvider

type Config struct {
	PATH               string
	WINDOW_WIDTH       int    `json:"window_width"`
	HORIZONTAL_SPACING int    `json:"horizontal_spacing"`
	VERTICAL_SPACING   int    `json:"vertical_spacing"`
	ICON_SIZE          int    `json:"icon_size"`
	CSS_THEME_FILE     string `json:"css_theme_file"`
}

type WidgetConfig struct {
	Type       string           `json:"type"`
	Properties WidgetProperties `json:"properties"`
	Children   []*WidgetConfig  `json:"children"`
	Action     string           `json:"action"`
	Script     string           `json:"script"`
}
type WidgetProperties struct {
	Orientation string `json:"orientation"`
	Spacing     int    `json:"spacing"`
	Label       string `json:"label"`
}

func LoadConfig() {
	user, _ := user.Current()
	path := "/home/" + user.Username + "/.config/actionCenter/"
	configPath := path + "config.json"
	widgetsPath := path + "widgets.json"

	cfgFile, cer := ioutil.ReadFile(configPath)
	wcfgFile, wer := ioutil.ReadFile(widgetsPath)

	if cer != nil {
		fmt.Println("Error reading config file:", cer)
	}
	if wer != nil {
		fmt.Println("Error reading widget file:", wer)
	}

	var config Config
	var widgets []WidgetConfig

	cer = json.Unmarshal(cfgFile, &config)
	wer = json.Unmarshal(wcfgFile, &widgets)

	if cer != nil {
		fmt.Println("Error parsing config file:", cer)
	}
	if wer != nil {
		fmt.Println("Error parsing widget file:", wer)
	}

	config.PATH = path
	Conf = &config
	WidgetConfs = widgets

	provider, err := gtk.CssProviderNew()
	if err != nil {
		fmt.Println("Error creating css provider: ", err)
		panic(err)
	}

	err = provider.LoadFromPath(Conf.PATH + Conf.CSS_THEME_FILE)

	if err != nil {
		fmt.Println("Error loading "+Conf.PATH+Conf.CSS_THEME_FILE+": ", err)
		panic(err)
	}

	StyleProvider = provider
}
