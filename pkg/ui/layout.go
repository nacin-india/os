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
	App        *tview.Application
	MainFlex   *tview.Flex
	Header     *tview.TextView
	Middle     *tview.TextView
	Copyright  *tview.TextView
	StatsPanel *tview.TextView
}

// NewUI creates a new UI instance and sets up the entire UI
func NewUI() *UI {
	app := tview.NewApplication()
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create header with stats panel
	header := tview.NewTextView().
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignLeft)
	header.SetBackgroundColor(tcell.ColorDarkSlateGray)

	statsPanel := tview.NewTextView().
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignRight)
	statsPanel.SetBackgroundColor(tcell.ColorDarkSlateGray)

	headerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	headerFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDarkSlateGray), 2, 0, false)
	headerFlex.AddItem(header, 0, 2, false)
	headerFlex.AddItem(statsPanel, 0, 1, false)
	headerFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDarkSlateGray), 2, 0, false)
	headerFlex.SetBackgroundColor(tcell.ColorDarkSlateGray)

	// Create middle section
	middle := tview.NewTextView().
		SetTextColor(tcell.ColorBlack).
		SetTextAlign(tview.AlignLeft)
	middle.SetBackgroundColor(tcell.ColorYellow)

	middleFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	middleFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorYellow), 2, 0, false)
	middleFlex.AddItem(middle, 0, 1, false)
	middleFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorYellow), 2, 0, false)
	middleFlex.SetBackgroundColor(tcell.ColorYellow)

	// Create copyright
	currentYear := time.Now().Year()
	copyright := tview.NewTextView().
		SetTextColor(tcell.ColorWhite).
		SetText(fmt.Sprintf("© %d Sar Infocom. All rights reserved. ", currentYear)).
		SetTextAlign(tview.AlignCenter)
	copyright.SetBackgroundColor(tcell.ColorDarkSlateGray)

	copyrightFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	copyrightFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDarkSlateGray), 2, 0, false)
	copyrightFlex.AddItem(copyright, 0, 1, false)
	copyrightFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDarkSlateGray), 2, 0, false)
	copyrightFlex.SetBackgroundColor(tcell.ColorDarkSlateGray)

	// Add all sections to the main flex layout
	mainFlex.AddItem(headerFlex, 0, 12, false)
	mainFlex.AddItem(middleFlex, 0, 6, false)
	mainFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorYellow), 0, 18, false) // Empty space
	mainFlex.AddItem(copyrightFlex, 0, 1, false)

	ui := &UI{
		App:        app,
		MainFlex:   mainFlex,
		Header:     header,
		Middle:     middle,
		Copyright:  copyright,
		StatsPanel: statsPanel,
	}

	// Setup key handling
	ui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
			ui.App.Stop()
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
		ui.App.QueueUpdateDraw(func() {
			info := system.GetSystemInfo()

			// Update header text with CPU, GPU, memory, and uptime on the left
			ui.Header.SetText(fmt.Sprintf("\n\nNACIN EXAM SERVER\n\nSar Infocom Virtual Platform\n\n%s\n%s\n%s\n%s\n\n",
				info.CPUInfo,
				info.MemoryInfo,
				info.GPUInfo,
				info.UptimeInfo))

			// Update stats panel with usage and temperature info on the right
			ui.StatsPanel.SetText(fmt.Sprintf("\n\n\n\n\n\n%s\n%s\n%s\n%s\n\n",
				info.CPUUsage,
				info.RAMUsage,
				info.CPUTemp,
				info.GPUTemp))

			// Update middle text with IP addresses
			ipText := "\nDownload tools to manage this host from:\n"
			ipText += strings.Join(info.IPAddresses, "\n") + "\n"
			ui.Middle.SetText(ipText)
		})

		time.Sleep(900 * time.Millisecond)
	}
}

// Run runs the UI application
func (ui *UI) Run() error {
	return ui.App.SetRoot(ui.MainFlex, true).EnableMouse(false).Run()
}
