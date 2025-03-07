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

// NewUI creates a new UI instance
func NewUI() *UI {
	app := tview.NewApplication()

	// Create main flex layout
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create header
	header := tview.NewTextView().
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignLeft)
	header.SetBackgroundColor(tcell.ColorDarkSlateGray)

	// Create middle section
	middle := tview.NewTextView().
		SetTextColor(tcell.ColorBlack).
		SetTextAlign(tview.AlignLeft)
	middle.SetBackgroundColor(tcell.ColorYellow)

	// Create copyright
	currentYear := time.Now().Year()
	copyright := tview.NewTextView().
		SetTextColor(tcell.ColorWhite).
		SetText(fmt.Sprintf("© %d Sar Infocom. All rights reserved. ", currentYear)).
		SetTextAlign(tview.AlignCenter)
	copyright.SetBackgroundColor(tcell.ColorDarkSlateGray)

	// Create stats panel
	statsPanel := tview.NewTextView().
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignRight)
	statsPanel.SetBackgroundColor(tcell.ColorDarkSlateGray)

	return &UI{
		App:        app,
		MainFlex:   mainFlex,
		Header:     header,
		Middle:     middle,
		Copyright:  copyright,
		StatsPanel: statsPanel,
	}
}

// SetupLayout sets up the UI layout
func (ui *UI) SetupLayout() {
	// Add padding to header
	headerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	headerFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDarkSlateGray), 2, 0, false)
	headerFlex.AddItem(ui.Header, 0, 2, false)
	headerFlex.AddItem(ui.StatsPanel, 0, 1, false)
	headerFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDarkSlateGray), 2, 0, false)
	headerFlex.SetBackgroundColor(tcell.ColorDarkSlateGray)

	// Add padding to middle section
	middleFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	middleFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorYellow), 2, 0, false)
	middleFlex.AddItem(ui.Middle, 0, 1, false)
	middleFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorYellow), 2, 0, false)
	middleFlex.SetBackgroundColor(tcell.ColorYellow)

	// Add padding to copyright
	copyrightFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	copyrightFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDarkSlateGray), 2, 0, false)
	copyrightFlex.AddItem(ui.Copyright, 0, 1, false)
	copyrightFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDarkSlateGray), 2, 0, false)
	copyrightFlex.SetBackgroundColor(tcell.ColorDarkSlateGray)

	// Add all sections to the main flex layout with appropriate spacing
	ui.MainFlex.AddItem(headerFlex, 0, 12, false)
	ui.MainFlex.AddItem(middleFlex, 0, 6, false)
	ui.MainFlex.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorYellow), 0, 18, false) // Empty space
	ui.MainFlex.AddItem(copyrightFlex, 0, 1, false)
}

// SetupKeyHandling sets up key handling
func (ui *UI) SetupKeyHandling() {
	ui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
			ui.App.Stop()
		}
		return event
	})
}

// UpdateSystemInfo updates the UI with system information
func (ui *UI) UpdateSystemInfo() {
	go func() {
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

			time.Sleep(2 * time.Second)
		}
	}()
}

// Run runs the UI application
func (ui *UI) Run() error {
	return ui.App.SetRoot(ui.MainFlex, true).EnableMouse(false).Run()
}
