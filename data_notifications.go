package main

import (
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/gtk"
)

type Notification struct {
	AppName           string
	Id                uint32
	AppIcon           string
	Summary           string
	Body              string
	Hints             map[string]dbus.Variant
	Actions           []string
	ExpirationTimeOut int32
	Image             *gtk.Image
}

func NewNotification(appName string, replacesID uint32, appIcon string, summary string, body string, actions []string, hints map[string]dbus.Variant, expireTimeout int32, image *gtk.Image) Notification {
	n := Notification{
		AppName:           appName,
		Id:                replacesID,
		AppIcon:           appIcon,
		Summary:           summary,
		Body:              body,
		Hints:             hints,
		Actions:           actions,
		ExpirationTimeOut: expireTimeout,
		Image:             image,
	}
	return n
}
