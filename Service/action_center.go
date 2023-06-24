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
	notificationCenter *NotificationCenterService
	actionCenterUI     *View.ActionCenterUI
}

func NewActionCenter() *ActionCenter {
	return &ActionCenter{}
}

func (app *ActionCenter) Init() error {
	app.notificationCenter = NewNotificationCenter()
	if err := app.notificationCenter.Init(app); err != nil {
		return err
	}
	go app.notificationCenter.Run()

	app.actionCenterUI = &View.ActionCenterUI{}
	if err := app.actionCenterUI.CreateUI(app, "test.json"); err != nil {
		return err
	}

	return nil
}
func (app *ActionCenter) GetNotifications() ([]Model.Notification, error) {

	ns, err := app.notificationCenter.GetNotifications()
	if err != nil {
		return nil, err
	}
	return ns, nil
}
func (app *ActionCenter) AddNotification(n Model.Notification) error {
	err := app.actionCenterUI.AddNotification(n)
	if err != nil {
		return err
	}
	app.actionCenterUI.ShowAll()
	return nil
}
func (app *ActionCenter) Hello(name string) string {
	return fmt.Sprintf("Hello %s", name)
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
				// Perform any necessary actions for SIGUSR1
				app.actionCenterUI.ToggleVisiblity()
			case syscall.SIGTERM:
				fmt.Println("Closing dbus conn")
				app.notificationCenter.conn.Close()
				os.Exit(0)
			}
		}
	}()
	gtk.Main()
}
