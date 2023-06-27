package Command

import "github.com/actionCenter/Model"

type ActionCenterInterface interface {
	GetNotifications() ([]Model.Notification, error)
	AddNotification(Model.Notification)
}
