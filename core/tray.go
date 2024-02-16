package core

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (d *Core) TrayReady() {
	systray.SetIcon(d.Icon)
	systray.SetTitle("Rosa")
	systray.SetTooltip("Beta")

	open := systray.AddMenuItem("Open app", "Openinng app")
	open.SetTooltip("Press to open app")
	go func() {
		for {
			<-open.ClickedCh
			log.Println("Oppening window")
			runtime.WindowShow(d.wailsctx)
		}

	}()

	dht := systray.AddMenuItem("DHT:", "Showing dht status")
	dht.SetTooltip("Showing avalible nodes.")
	network := systray.AddMenuItem("Connected:", "Showing dht status")
	network.SetTooltip("Showing opened connections.")
	go func() {
		for {
			if !d.Started {
				network.SetTitle("Connected: 0")
				time.Sleep(time.Second * 10)
				continue
			}
			if !d.Started {
				dht.SetTitle("DHT: 0")
				time.Sleep(time.Second * 10)
				continue
			}
			dht.SetTitle("DHT: " + strconv.Itoa(len(d.host.H.Peerstore().Peers())))
			network.SetTitle("Connected: " + strconv.Itoa(len(d.host.H.Network().Peers())))
			time.Sleep(time.Second * 10)
		}
	}()

	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuit.SetIcon(d.Icon)
	go func() {
		<-mQuit.ClickedCh
		os.Exit(0)
	}()
}

func (d *Core) TrayExit() {
	log.Println("Exit call")
}
