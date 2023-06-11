package Model

import (
	"github.com/godbus/dbus/v5"
)

type Notification struct {
	App            string
	Id             uint32
	Icon           string
	Summary        string
	Message        string
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
	if message, ok := variant["body"].Value().(string); ok {
		notification.Message = message
	}
	if hints, ok := variant["hints"].Value().(map[string]dbus.Variant); ok {
		notification.Hints = hints
	}
	if expiration, ok := variant["expire_timeout"].Value().(int32); ok {
		notification.ExpirationTime = expiration
	}

	return notification
}
