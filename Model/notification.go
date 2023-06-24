package Model

import (
	"regexp"
	"strings"

	"github.com/godbus/dbus/v5"
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
func NotificationFromVariant(variant map[string]dbus.Variant) Notification {
	notification := Notification{}

	if app, ok := variant["app_name"].Value().(string); ok {
		notification.AppName = app
	}
	if id, ok := variant["id"].Value().(uint32); ok {
		notification.Id = id
	}
	if icon, ok := variant["icon_data"].Value().(string); ok {
		notification.AppIcon = icon
	}
	if summary, ok := variant["summary"].Value().(string); ok {
		notification.Summary = summary
	}
	if body, ok := variant["body"].Value().(string); ok {
		notification.Body = body
	}
	if hints, ok := variant["hints"].Value().(map[string]dbus.Variant); ok {
		notification.Hints = hints
	}
	if expiration, ok := variant["expire_timeout"].Value().(int32); ok {
		notification.ExpirationTimeOut = expiration
	}

	return notification
}
func (n *Notification) RemoveHyperLinkFromBody() {
	re := regexp.MustCompile(`<a.*?>(.*?)</a>`)
	n.Body = re.ReplaceAllString(n.Body, "")

	// Remove empty lines
	lines := strings.Split(n.Body, "\n")
	var filteredLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filteredLines = append(filteredLines, line)
		}
	}
	n.Body = strings.Join(filteredLines, "\n")
}
