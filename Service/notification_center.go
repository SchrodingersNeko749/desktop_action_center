package Service

import (
	"fmt"

	"github.com/actionCenter/Command"
	"github.com/actionCenter/Model"
	"github.com/godbus/dbus/v5"
)

type NotificationCenterService struct {
	Notifications       []Model.Notification
	conn                *dbus.Conn
	obj                 dbus.BusObject
	actionCenterHandler Command.ActionCenterInterface
}

func NewNotificationCenter() *NotificationCenterService {
	return &NotificationCenterService{}
}
func (n *NotificationCenterService) Init(ac Command.ActionCenterInterface) error {
	// Connect to the session bus
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return err
	}
	n.conn = conn
	n.actionCenterHandler = ac

	server := n
	conn.Export(server, "/org/freedesktop/Notifications", "org.freedesktop.Notifications")

	reply, err := conn.RequestName("org.freedesktop.Notifications", dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("Couldnt initialize notification server: dbus name already taken")
	}
	return nil
}
func (n *NotificationCenterService) Run() {
	fmt.Println("Listening...")
	select {}

}
func (n *NotificationCenterService) GetNotifications() ([]Model.Notification, error) {
	return nil, nil
}
