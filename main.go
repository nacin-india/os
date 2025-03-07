package main

import (
	"log"

	"github.com/nacin/nacin-os/pkg/ui"
)

func main() {
	if err := ui.NewUI().Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
