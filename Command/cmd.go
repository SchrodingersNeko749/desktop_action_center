package Command

import "github.com/actionCenter/Model"

type ActionCenterInterface interface {
	Hello() string
	GetNotifications() ([]Model.Notification, error)
}
