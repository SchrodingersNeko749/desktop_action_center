package Model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/godbus/dbus/v5"
)

type Notification struct {
	App            string
	Id             uint32
	Icon           string
	Summary        string
	Body           string
	Hints          map[string]dbus.Variant
	ExpirationTime int32
}

func NotificationFromVariant(variant map[string]dbus.Variant) Notification {
	notification := Notification{}

	if app, ok := variant["app_name"].Value().(string); ok {
		notification.App = app
	}
	if id, ok := variant["id"].Value().(uint32); ok {
		notification.Id = id
	}
	if icon, ok := variant["icon_data"].Value().(string); ok {
		notification.Icon = icon
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
		notification.ExpirationTime = expiration
	}

	return notification
}
func (n *Notification) FixEmptyIcon() {
	//gtk-dialog-info
	if strings.Contains(n.Body, "<a href=") {
		re := regexp.MustCompile(`<a\s+href="https?://(?:[a-zA-Z0-9-]+\.)?([a-zA-Z0-9-]+)\.[a-zA-Z]{2,}(?:/[a-zA-Z0-9-._]*)?"[^>]*>(.*?)</a>`)
		match := re.FindStringSubmatch(n.Body)
		if len(match) > 1 {
			//domain := match[1]
			n.Icon = match[1]
			n.Body = strings.Replace(n.Body, match[0], "", -1)
			n.Body = strings.TrimSpace(n.Body)
		}

	} else {
		//fmt.Println("warning: no match found for app icon (notification.go)")
		fmt.Println(n.Body)
	}

}
