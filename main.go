package main

import (
	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady)
}

func onReady() {
	systray.SetTitle("Awesome app")
	systray.SetTooltip("Awesomeeeee tooltip")
	systray.AddMenuItem("Quit", "Quit it!")
}
