package Service

import (
	"fmt"

	"github.com/actionCenter/Model"
	"github.com/godbus/dbus/v5"
)

type NotificationCenterService struct {
	Notifications []Model.Notification
	conn          *dbus.Conn
	obj           dbus.BusObject
}

func NewNotificationCenter() *NotificationCenterService {
	return &NotificationCenterService{}
}
func (n *NotificationCenterService) Run() error {
	// Connect to the session bus
	conn, err := dbus.SessionBus()
	if err != nil {
		return (err)
	}

	// Get the notification object
	obj := conn.Object("org.freedesktop.Notifications", dbus.ObjectPath("/org/freedesktop/Notifications"))
	n.conn = conn
	n.obj = obj

	return nil
}
func (n *NotificationCenterService) GetNotifications() ([]Model.Notification, error) {
	call := n.obj.Call("org.dunstproject.cmd0.NotificationListHistory", 0)
	if call.Err != nil {
		return nil, fmt.Errorf("error calling NotificationListHistory: %w", call.Err)
	}

	var variants []map[string]dbus.Variant
	if err := call.Store(&variants); err != nil {
		return nil, fmt.Errorf("error decoding notification variants: %w", err)
	}

	notifications := make([]Model.Notification, len(variants))
	for i, v := range variants {
		notifications[i] = Model.NotificationFromVariant(v)
	}

	return notifications, nil
}
