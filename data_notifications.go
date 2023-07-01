package main

import "github.com/godbus/dbus/v5"

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
