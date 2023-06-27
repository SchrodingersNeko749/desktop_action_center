package Model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/actionCenter/Data"
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type NotificationWidget struct {
	container *gtk.Box
	id        int
}
type Notification struct {
	AppName           string
	Id                uint32
	AppIcon           string
	Summary           string
	Body              string
	Hints             map[string]dbus.Variant
	Actions           []string
	ExpirationTimeOut int32
}

func NewNotification(appName string, replacesID uint32, appIcon string, summary string, body string, actions []string, hints map[string]dbus.Variant, expireTimeout int32) Notification {
	n := Notification{
		AppName:           appName,
		Id:                replacesID,
		AppIcon:           appIcon,
		Summary:           summary,
		Body:              body,
		Hints:             hints,
		Actions:           actions,
		ExpirationTimeOut: expireTimeout,
	}
	return n
}
func NotificationFromVariant(variant map[string]dbus.Variant) Notification {
	notification := Notification{}

	if app, ok := variant["app_name"].Value().(string); ok {
		notification.AppName = app
	}
	if id, ok := variant["id"].Value().(uint32); ok {
		notification.Id = id
	}
	if icon, ok := variant["icon_data"].Value().(string); ok {
		notification.AppIcon = icon
	}
	if summary, ok := variant["summary"].Value().(string); ok {
		notification.Summary = summary
	}
	if body, ok := variant["body"].Value().(string); ok {
		notification.Body = body
	}
	if hints, ok := variant["hints"].Value().(map[string]dbus.Variant); ok {
		notification.Hints = hints
	}
	if expiration, ok := variant["expire_timeout"].Value().(int32); ok {
		notification.ExpirationTimeOut = expiration
	}

	return notification
}
func (n *Notification) RemoveHyperLinkFromBody() {
	re := regexp.MustCompile(`<a.*?>(.*?)</a>`)
	n.Body = re.ReplaceAllString(n.Body, "")

	// Remove empty lines
	lines := strings.Split(n.Body, "\n")
	var filteredLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filteredLines = append(filteredLines, line)
		}
	}
	n.Body = strings.Join(filteredLines, "\n")
}

func CreateNotificationComponent(n Notification) *gtk.ListBoxRow {
	row, _ := gtk.ListBoxRowNew()
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 20)
	summaryLabel, _ := gtk.LabelNew(n.Summary)
	bodyLabel, _ := gtk.LabelNew(n.Body)
	var icon *gtk.Image = nil

	if _, ok := n.Hints["image-data"]; ok {
		width := n.Hints["image-data"].Value().([]interface{})[0].(int32)
		height := n.Hints["image-data"].Value().([]interface{})[1].(int32)
		rowStride := n.Hints["image-data"].Value().([]interface{})[2].(int32)
		hasAlpha := n.Hints["image-data"].Value().([]interface{})[3].(bool)
		bitsPerSample := n.Hints["image-data"].Value().([]interface{})[4].(int32)

		img := n.Hints["image-data"].Value().([]interface{})[6].([]byte)
		pixbuf, err := gdk.PixbufNewFromData(img, gdk.COLORSPACE_RGB, hasAlpha, int(bitsPerSample), int(width), int(height), int(rowStride))
		icon, err = gtk.ImageNewFromPixbuf(pixbuf)
		if err != nil {
			fmt.Println(err)
		}
	} else if customImagePath, ok := n.Hints["image-path"].Value().(string); ok {
		icon, _ = gtk.ImageNewFromFile(customImagePath)
	} else if n.AppIcon != "" {
		icon, _ = gtk.ImageNewFromIconName(n.AppIcon, gtk.ICON_SIZE_LARGE_TOOLBAR)
	} else {
		icon, _ = gtk.ImageNewFromIconName("gtk-dialog-info", gtk.ICON_SIZE_LARGE_TOOLBAR)
	}
	if icon != nil {
		resize(icon)
		hbox.PackStart(icon, false, false, 0)
	}

	summaryLabel.SetXAlign(0)
	summaryLabel.SetHAlign(gtk.ALIGN_START)
	summaryLabel.SetHExpand(true)
	summaryLabel.SetLineWrap(true)
	summaryLabel.SetSelectable(true)
	summaryLabel.SetMaxWidthChars(1)
	summaryLabel.SetSizeRequest(Data.Conf.WINDOW_WIDTH-Data.Conf.HORIZONTAL_SPACING-Data.Conf.ICON_SIZE-128, -1)

	bodyLabel.SetXAlign(0)
	bodyLabel.SetHAlign(gtk.ALIGN_START)
	bodyLabel.SetHExpand(true)
	bodyLabel.SetLineWrap(true)
	bodyLabel.SetSelectable(true)
	bodyLabel.SetMaxWidthChars(1)
	bodyLabel.SetSizeRequest(Data.Conf.WINDOW_WIDTH-Data.Conf.HORIZONTAL_SPACING-Data.Conf.ICON_SIZE-128, -1)

	hbox.SetSizeRequest(Data.Conf.WINDOW_WIDTH-Data.Conf.HORIZONTAL_SPACING-Data.Conf.ICON_SIZE-128, -1)

	vbox.PackStart(summaryLabel, false, false, 0)
	vbox.PackStart(bodyLabel, false, false, 0)
	hbox.PackStart(vbox, false, true, 0)
	container, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	container.PackStart(hbox, false, false, 0)
	container.PackStart(vbox, false, false, 0)
	row.Add(container)

	style, _ := hbox.GetStyleContext()
	style.AddClass("notification-widget")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = summaryLabel.GetStyleContext()
	style.AddClass("notification-summary")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	style, _ = bodyLabel.GetStyleContext()
	style.AddClass("notification-body")
	style.AddProvider(Data.StyleProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return row
}

func resize(icon *gtk.Image) {
	pixbuf := icon.GetPixbuf()
	if pixbuf == nil {
		theme, _ := gtk.IconThemeGetDefault()
		iconName, _ := icon.GetIconName()
		pixbuf, _ = theme.LoadIconForScale(iconName, Data.Conf.ICON_SIZE, 1, gtk.ICON_LOOKUP_FORCE_SIZE)
	}
	scaledPixbuf, _ := pixbuf.ScaleSimple(Data.Conf.ICON_SIZE, Data.Conf.ICON_SIZE, gdk.INTERP_BILINEAR)
	icon.SetFromPixbuf(scaledPixbuf)
}
