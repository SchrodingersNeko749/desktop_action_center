package Service

import (
	"github.com/actionCenter/Model"
	"github.com/godbus/dbus/v5"
)

func (n NotificationCenterService) Notify(appName string, replacesID uint32, appIcon string, summary string, body string, actions []string, hints map[string]dbus.Variant, expireTimeout int32) (uint32, *dbus.Error) {
	notification := Model.NewNotification(appName, replacesID, appIcon, summary, body, actions, hints, expireTimeout)
	notification.RemoveHyperLinkFromBody()
	n.actionCenterHandler.AddNotification(notification)
	return 0, nil
}

func (n NotificationCenterService) GetCapabilities() ([]string, *dbus.Error) {
	return []string{"action-icons", "actions", "body", "body-hyperlinks", "body-images", "body-markup", "icon-multi", "icon-static", "persistence", "sound"}, nil
}

func (n NotificationCenterService) GetServerInformation() (string, string, string, string, *dbus.Error) {
	return "antarctica", "antarctica.com", "1.0", "1.2", nil
}
