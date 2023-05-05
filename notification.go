package main

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

type Notification struct {
	app            string
	id             uint32
	icon           string
	summary        string
	message        string
	hints          map[string]dbus.Variant
	expirationTime int32
}

func NewNotification(msg *dbus.Message) (*Notification, error) {
	if msg.Type != dbus.TypeMethodCall {
		return nil, fmt.Errorf("invalid message type: %v", msg.Type)
	}
	if len(msg.Body) < 7 {
		return nil, fmt.Errorf("invalid message body length: %v", len(msg.Body))
	}
	return &Notification{
		app:            msg.Body[0].(string),
		id:             msg.Body[1].(uint32),
		icon:           msg.Body[2].(string),
		summary:        msg.Body[3].(string),
		message:        msg.Body[4].(string),
		hints:          msg.Body[6].(map[string]dbus.Variant),
		expirationTime: msg.Body[7].(int32),
	}, nil
}

/*
body[0]: App name or program name that sent the notification.
body[1]: Notification ID or notification timeout (in milliseconds).
body[2]: Icon path or icon name to display in the notification.
body[3]: Summary or title of the notification.
body[4]: Body or message of the notification.
body[5]: Array of hints or options for the notification.
body[6]: Map of hints or options for the notification.
body[7]: Notification expiration time or timeout (in milliseconds).
*/
