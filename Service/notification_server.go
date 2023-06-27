package Service

import (
	"fmt"

	"github.com/actionCenter/Model"
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/glib"
)

type NotificationServer struct {
	Notifications []Model.Notification
	conn          *dbus.Conn
	obj           dbus.BusObject
	ActionCenter  *ActionCenter
}

func (n *NotificationServer) Init(ac *ActionCenter) error {
	n.ActionCenter = ac
	// Connect to the session bus
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return err
	}
	n.conn = conn

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

func (n *NotificationServer) GetNotifications() ([]Model.Notification, error) {
	return nil, nil
}

func (n *NotificationServer) Notify(appName string, replacesID uint32, appIcon string, summary string, body string, actions []string, hints map[string]dbus.Variant, expireTimeout int32) (uint32, *dbus.Error) {
	notification := Model.NewNotification(appName, replacesID, appIcon, summary, body, actions, hints, expireTimeout)
	//notification.RemoveHyperLinkFromBody()
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
