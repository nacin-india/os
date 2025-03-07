package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/nacin/nacin-os/pkg/system"
	"github.com/rivo/tview"
)

// UI holds all UI components
type UI struct {
	app        *tview.Application
	mainFlex   *tview.Flex
	header     *tview.TextView
	middle     *tview.TextView
	footer     *tview.TextView
	stats      *tview.TextView
	bottomFlex *tview.Flex // Added new field for bottom section
}

// createTextView creates a new text view with specified properties
func createTextView(color tcell.Color, bgColor tcell.Color, align int) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetTextColor(color)
	tv.SetTextAlign(align)
	tv.SetBackgroundColor(bgColor)
	return tv
}

// createPaddedFlex creates a flex container with padding on both sides
func createPaddedFlex(content *tview.TextView, bgColor tcell.Color, padding int) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	flex.AddItem(tview.NewBox().SetBackgroundColor(bgColor), padding, 0, false)
	flex.AddItem(content, 0, 1, false)
	flex.AddItem(tview.NewBox().SetBackgroundColor(bgColor), padding, 0, false)
	flex.SetBackgroundColor(bgColor)
	return flex
}

// NewUI creates a new UI instance and sets up the entire UI
func NewUI() *UI {
	app := tview.NewApplication()
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Define colors
	darkGray := tcell.ColorDarkSlateGray
	yellow := tcell.ColorYellow
	white := tcell.ColorWhite
	black := tcell.ColorBlack

	// Create header
	header := createTextView(white, darkGray, tview.AlignLeft)

	// Create stats panel for system usage
	stats := createTextView(white, yellow, tview.AlignRight)

	// Create header flex without stats
	headerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	headerFlex.AddItem(tview.NewBox().SetBackgroundColor(darkGray), 2, 0, false)
	headerFlex.AddItem(header, 0, 1, false)
	headerFlex.AddItem(tview.NewBox().SetBackgroundColor(darkGray), 2, 0, false)
	headerFlex.SetBackgroundColor(darkGray)

	// Create middle section for IP addresses
	middle := createTextView(black, yellow, tview.AlignLeft)

	// Create bottom section with IP addresses on left and stats on right
	bottomFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	bottomFlex.AddItem(tview.NewBox().SetBackgroundColor(yellow), 2, 0, false)
	bottomFlex.AddItem(middle, 0, 2, false)
	bottomFlex.AddItem(stats, 0, 1, false)
	bottomFlex.AddItem(tview.NewBox().SetBackgroundColor(yellow), 2, 0, false)
	bottomFlex.SetBackgroundColor(yellow)

	// Create footer
	currentYear := time.Now().Year()
	footer := createTextView(white, darkGray, tview.AlignCenter)
	footer.SetText(fmt.Sprintf("© %d Sar Infocom. All rights reserved. ", currentYear))
	footerFlex := createPaddedFlex(footer, darkGray, 2)

	// Add all sections to the main flex layout
	mainFlex.AddItem(headerFlex, 0, 12, false)
	mainFlex.AddItem(bottomFlex, 0, 24, false) // Combined yellow section
	mainFlex.AddItem(footerFlex, 0, 1, false)

	ui := &UI{
		app:        app,
		mainFlex:   mainFlex,
		header:     header,
		middle:     middle,
		footer:     footer,
		stats:      stats,
		bottomFlex: bottomFlex,
	}

	// Setup key handling
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
			ui.app.Stop()
		}
		return event
	})

	// Start updating system information
	go ui.updateSystemInfoPeriodically()

	return ui
}

// updateSystemInfoPeriodically updates the UI with system information every 1 second
func (ui *UI) updateSystemInfoPeriodically() {
	for {
		ui.app.QueueUpdateDraw(func() {
			info := system.GetSystemInfo()

			// Update header text - removed stats info
			ui.header.SetText(fmt.Sprintf("\n\nNACIN EXAM SERVER\n\nSar Infocom Virtual Platform\n\n%s\n%s\n%s\n%s\n\n",
				info.CPUInfo,
				info.MemoryInfo,
				info.GPUInfo,
				info.UptimeInfo))

			// Update stats panel in the bottom yellow section - align vertically with IP addresses
			ui.stats.SetText(fmt.Sprintf("\n%s\n%s\n%s\n%s",
				info.CPUUsage,
				info.RAMUsage,
				info.CPUTemp,
				info.GPUTemp))

			// Update middle text for IP addresses - adjust formatting to align with stats
			ipText := "\nIP addresses:\n"
			ipText += strings.Join(info.IPAddresses, "\n")
			ui.middle.SetText(ipText)
		})

		time.Sleep(900 * time.Millisecond)
	}
}

// Run runs the UI application
func (ui *UI) Run() error {
	return ui.app.SetRoot(ui.mainFlex, true).EnableMouse(false).Run()
}
