package main

import (
	"fmt"
	"html"

	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/glib"
	strip "github.com/grokify/html-strip-tags-go"
)

type NotificationServer struct {
	Notifications []Notification
	Conn          *dbus.Conn
	obj           dbus.BusObject
	ActionCenter  *ActionCenter
}

func (n *NotificationServer) Init(ac *ActionCenter) error {
	n.ActionCenter = ac
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return err
	}
	n.Conn = conn

	server := n
	conn.Export(server, "/org/freedesktop/Notifications", "org.freedesktop.Notifications")

	reply, err := conn.RequestName("org.freedesktop.Notifications", dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("Couldnt initialize notification server: dbus name already taken")
	}
	fmt.Println("Listening...")
	select {}
}

func (n *NotificationServer) Notify(appName string, replacesID uint32, appIcon string, summary string, body string, actions []string, hints map[string]dbus.Variant, expireTimeout int32) (uint32, *dbus.Error) {
	unescaped := strip.StripTags(body)
	unescaped = html.UnescapeString(unescaped)
	notification := NewNotification(appName, replacesID, appIcon, summary, unescaped, actions, hints, expireTimeout, nil)

	glib.IdleAdd(func() {
		n.ActionCenter.AddNotification(notification)
	})
	return 0, nil
}

func (n *NotificationServer) GetCapabilities() ([]string, *dbus.Error) {
	return []string{"action-icons", "actions", "body", "body-hyperlinks", "body-images", "body-markup", "icon-multi", "icon-static", "persistence", "sound"}, nil
}

func (n *NotificationServer) GetServerInformation() (string, string, string, string, *dbus.Error) {
	return "antarctica", "antarctica.com", "1.0", "1.2", nil
}
