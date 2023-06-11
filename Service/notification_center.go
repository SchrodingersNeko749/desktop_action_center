package Service

import (
	"fmt"
	"log"

	"github.com/actionCenter/Model"
	"github.com/godbus/dbus/v5"
)

type NotificationCenter struct {
	Notifications []Model.Notification
	conn          *dbus.Conn
	obj           dbus.BusObject
}

func NewNotificationCenter() *NotificationCenter {
	return &NotificationCenter{}
}
func (n *NotificationCenter) Run() error {
	// Connect to the session bus
	conn, err := dbus.SessionBus()
	if err != nil {
		return (err)
	}

	// Get the notification object
	obj := conn.Object("org.freedesktop.Notifications", dbus.ObjectPath("/org/freedesktop/Notifications"))
	n.conn = conn
	n.obj = obj
	fmt.Println("test")
	return nil
}
func (n *NotificationCenter) GetHistory() (error, []Model.Notification) {
	// Call NotificationListHistory method
	call := n.obj.Call("org.dunstproject.cmd0.NotificationListHistory", 0)
	if call.Err != nil {
		log.Fatalln("cant call org.dunstproject.cmd0.NotificationListHistory", call.Err)
	}
	var variants []map[string]dbus.Variant
	err := call.Store(&variants)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d Notifications found: \n", len(variants))
	for _, v := range variants {
		notif := Model.NotificationFromVariant(v)
		n.Notifications = append(n.Notifications, notif)
		fmt.Println("test")
	}
	return nil, n.Notifications
}
