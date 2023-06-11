package Service

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/actionCenter/View"

	"github.com/gotk3/gotk3/gtk"
)

type ActionCenter struct {
	notifCenter    *NotificationCenterService
	actionCenterUI *View.ActionCenterUI
}

func NewActionCenter() *ActionCenter {
	return &ActionCenter{}
}

func (app *ActionCenter) Init() error {
	app.notifCenter = NewNotificationCenter()
	app.actionCenterUI = &View.ActionCenterUI{}

	log.Println("init()")
	if err := app.actionCenterUI.CreateUI(); err != nil {
		return err
	}

	return nil
}

func (app *ActionCenter) Run() {
	// initializing components
	app.notifCenter.Run()
	app.actionCenterUI.Run()

	// Test: GetNotifications
	app.notifCenter.GetNotifications()

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
				app.notifCenter.conn.Close()
				os.Exit(0)
			}
		}
	}()
	gtk.Main()
}
