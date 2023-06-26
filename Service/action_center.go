package Service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/actionCenter/Model"
	"github.com/actionCenter/View"

	"github.com/gotk3/gotk3/gtk"
)

type ActionCenter struct {
	notificationServer *NotificationServer
	actionCenterUI     *View.ActionCenterUI
}

func NewActionCenter() *ActionCenter {
	return &ActionCenter{}
}

func (app *ActionCenter) Init() error {
	app.notificationServer = NewNotificationServer()
	go app.notificationServer.Init(app)

	app.actionCenterUI = &View.ActionCenterUI{}
	if err := app.actionCenterUI.CreateUI(app); err != nil {
		return err
	}

	return nil
}

func (app *ActionCenter) GetNotifications() ([]Model.Notification, error) {

	ns, err := app.notificationServer.GetNotifications()
	if err != nil {
		return nil, err
	}
	return ns, nil
}
func (app *ActionCenter) AddNotification(n Model.Notification) error {
	app.actionCenterUI.AddNotification(n)

	app.actionCenterUI.ShowAll()
	return nil
}

func (app *ActionCenter) Run() {
	app.actionCenterUI.Run()

	// handling signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGTERM)

	go func() {
		fmt.Println("Monitoring signals")
		for {
			sig := <-sigs
			fmt.Println(sig)
			switch sig {
			case syscall.SIGUSR1:
				app.actionCenterUI.ToggleVisiblity()
			case syscall.SIGTERM:
				fmt.Println("Closing dbus conn")
				app.notificationServer.conn.Close()
				os.Exit(0)
			}
		}
	}()
	gtk.Main()
}
