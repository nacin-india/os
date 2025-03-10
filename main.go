package main

import (
	"log"

	"server/ui"
)

func main() {
	if err := ui.NewUI().Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
