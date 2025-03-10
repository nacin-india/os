package ui

import (
	"fmt"
	"strings"
	"time"

	"server/system"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI holds all UI components
type UI struct {
	app        *tview.Application
	mainFlex   *tview.Flex
	header     *tview.TextView
	copyright  *tview.TextView
	middle     *tview.TextView
	footer     *tview.TextView
	stats      *tview.TextView
	bottomFlex *tview.Flex // Added new field for bottom section
}

// createTextView creates a new text view with specified properties
func createTextView(color tcell.Color, bgColor tcell.Color, align int) *tview.TextView {
	tv := tview.NewTextView()
	if color != tcell.ColorWhite && color != tcell.ColorBlack {
		color = tcell.ColorBlack
	}
	tv.SetTextColor(color)
	tv.SetTextAlign(align)
	tv.SetBackgroundColor(bgColor)
	tv.SetDynamicColors(true)
	return tv
}

// NewUI creates a new UI instance and sets up the entire UI
func NewUI() *UI {
	app := tview.NewApplication()
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Define colors - using high contrast colors only
	darkGray := tcell.ColorBlack
	yellow := tcell.ColorYellow
	white := tcell.ColorWhite
	black := tcell.ColorBlack

	// Create header
	header := createTextView(white, darkGray, tview.AlignLeft)

	// Create copyright text view - right aligned
	copyright := createTextView(white, darkGray, tview.AlignRight)

	// Create stats panel for system usage - use black text on yellow for contrast
	stats := createTextView(black, yellow, tview.AlignRight)
	stats.SetWordWrap(true) // Enable word wrap for better space utilization

	// Create header flex with title on left and copyright on right
	headerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	headerFlex.AddItem(tview.NewBox().SetBackgroundColor(darkGray), 2, 0, false)
	headerFlex.AddItem(header, 0, 1, false)
	headerFlex.AddItem(copyright, 0, 1, false)
	headerFlex.AddItem(tview.NewBox().SetBackgroundColor(darkGray), 2, 0, false)
	headerFlex.SetBackgroundColor(darkGray)

	// Create middle section for IP addresses
	middle := createTextView(black, yellow, tview.AlignLeft)
	middle.SetWordWrap(true) // Enable word wrap for better space utilization

	// Create bottom section with IP addresses on left and stats on right
	bottomFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	bottomFlex.AddItem(tview.NewBox().SetBackgroundColor(yellow), 2, 0, false) // Restored original padding
	bottomFlex.AddItem(middle, 0, 2, false)
	bottomFlex.AddItem(stats, 0, 1, false)
	bottomFlex.AddItem(tview.NewBox().SetBackgroundColor(yellow), 2, 0, false) // Restored original padding
	bottomFlex.SetBackgroundColor(yellow)

	// Add all sections to the main flex layout
	mainFlex.AddItem(headerFlex, 0, 12, false)
	mainFlex.AddItem(bottomFlex, 0, 14, false) // Further reduced from 18 to 14 to make the yellow section even shorter

	ui := &UI{
		app:        app,
		mainFlex:   mainFlex,
		header:     header,
		copyright:  copyright,
		middle:     middle,
		footer:     nil, // No footer needed anymore
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

// createEnhancedTitle creates a slightly larger title using box drawing characters
func createEnhancedTitle(title string) string {
	// Create a box around the title to make it stand out
	topBorder := "┌" + strings.Repeat("─", len(title)+2) + "┐"
	titleLine := fmt.Sprintf("│ %s │", title)
	bottomBorder := "└" + strings.Repeat("─", len(title)+2) + "┘"

	// Build the enhanced title with spacing for better visibility
	enhancedTitle := fmt.Sprintf("%s\n%s\n%s", topBorder, titleLine, bottomBorder)

	return enhancedTitle
}

// updateSystemInfoPeriodically updates the UI with system information every 1 second
func (ui *UI) updateSystemInfoPeriodically() {
	for {
		ui.app.QueueUpdateDraw(func() {
			info := system.GetSystemInfo()
			copyrightText := "by Sar Infocom"
			enhancedTitle := createEnhancedTitle("NACIN EXAM SERVER")

			ui.header.SetText(fmt.Sprintf("\n%s\n\n[::b]%s[::]\n[::b]%s[::]\n[::b]%s[::]\n[::b]%s[::]\n\n",
				enhancedTitle,
				info.CPUInfo,
				info.MemoryInfo,
				info.GPUInfo,
				info.UptimeInfo))

			ui.copyright.SetText(fmt.Sprintf("\n\n[::b]%s[::]\n\n\n\n\n\n", copyrightText))

			ui.stats.SetText(fmt.Sprintf("\n[::b]%s[::]\n[::b]%s[::]\n",
				info.CPUUsage,
				info.RAMUsage))

			middleText := "\n[::b]IP addresses:[::]\n"
			for _, ip := range info.IPAddresses {
				middleText += fmt.Sprintf("[::b]%s[::]\n", ip)
			}

			ui.middle.SetText(middleText)
		})

		time.Sleep(900 * time.Millisecond)
	}
}

// Run runs the UI application
func (ui *UI) Run() error {
	return ui.app.SetRoot(ui.mainFlex, true).EnableMouse(true).Run()
}
