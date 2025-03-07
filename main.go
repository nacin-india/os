package main

import (
	"github.com/nacin/nacin-os/pkg/ui"
)

func main() {
	// Create a new UI
	appUI := ui.NewUI()

	// Setup the UI layout
	appUI.SetupLayout()

	// Setup key handling
	appUI.SetupKeyHandling()

	// Start updating system information
	appUI.UpdateSystemInfo()

	// Run the application
	if err := appUI.Run(); err != nil {
		panic(err)
	}
}
