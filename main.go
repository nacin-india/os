package main

import (
	"github.com/nacin/nacin-os/pkg/ui"
)

func main() {
	// Create and run the UI application
	appUI := ui.NewUI()

	if err := appUI.Run(); err != nil {
		panic(err)
	}
}
