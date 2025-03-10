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
	stats      *tview.TextView
	bottomFlex *tview.Flex
}

// createTextView creates a new text view with specified properties
func createTextView(color, bgColor tcell.Color, align int) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetTextColor(color)
	tv.SetTextAlign(align)
	tv.SetBackgroundColor(bgColor)
	tv.SetDynamicColors(true)
	return tv
}

// NewUI creates a new UI instance and sets up the entire UI
func NewUI() *UI {
	// Define colors
	darkGray := tcell.ColorBlack
	yellow := tcell.ColorYellow
	white := tcell.ColorWhite
	black := tcell.ColorBlack

	// Create app and main layout
	app := tview.NewApplication()
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create text components
	header := createTextView(white, darkGray, tview.AlignLeft)
	copyright := createTextView(white, darkGray, tview.AlignRight)
	stats := createTextView(black, yellow, tview.AlignRight)
	stats.SetWordWrap(true)
	middle := createTextView(black, yellow, tview.AlignLeft)
	middle.SetWordWrap(true)

	// Create header flex
	headerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	headerFlex.AddItem(tview.NewBox().SetBackgroundColor(darkGray), 2, 0, false)
	headerFlex.AddItem(header, 0, 1, false)
	headerFlex.AddItem(copyright, 0, 1, false)
	headerFlex.AddItem(tview.NewBox().SetBackgroundColor(darkGray), 2, 0, false)
	headerFlex.SetBackgroundColor(darkGray)

	// Create bottom section
	bottomFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	bottomFlex.AddItem(tview.NewBox().SetBackgroundColor(yellow), 2, 0, false)
	bottomFlex.AddItem(middle, 0, 2, false)
	bottomFlex.AddItem(stats, 0, 1, false)
	bottomFlex.AddItem(tview.NewBox().SetBackgroundColor(yellow), 2, 0, false)
	bottomFlex.SetBackgroundColor(yellow)

	// Add sections to main layout
	mainFlex.AddItem(headerFlex, 0, 12, false)
	mainFlex.AddItem(bottomFlex, 0, 14, false)

	// Create UI instance
	ui := &UI{
		app:        app,
		mainFlex:   mainFlex,
		header:     header,
		copyright:  copyright,
		middle:     middle,
		stats:      stats,
		bottomFlex: bottomFlex,
	}

	// Setup key handling for exit
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

// createEnhancedTitle creates a title with box drawing characters
func createEnhancedTitle(title string) string {
	width := len(title) + 2
	return fmt.Sprintf("┌%s┐\n│ %s │\n└%s┘",
		strings.Repeat("─", width),
		title,
		strings.Repeat("─", width))
}

// updateSystemInfoPeriodically updates the UI with system information every 1 second
func (ui *UI) updateSystemInfoPeriodically() {
	for {
		ui.app.QueueUpdateDraw(func() {
			info := system.GetSystemInfo()
			title := createEnhancedTitle("NACIN EXAM SERVER")

			// Update header info
			ui.header.SetText(fmt.Sprintf("\n%s\n\n[::b]%s[::]\n[::b]%s[::]\n[::b]%s[::]\n[::b]%s[::]\n\n",
				title, info.CPUInfo, info.MemoryInfo, info.GPUInfo, info.UptimeInfo))

			// Update copyright
			ui.copyright.SetText("\n\n[::b]by Sar Infocom[::]\n\n\n\n\n\n")

			// Update stats
			ui.stats.SetText(fmt.Sprintf("\n[::b]%s[::]\n[::b]%s[::]\n",
				info.CPUUsage, info.RAMUsage))

			// Update IP addresses
			ipText := "\n[::b]IP addresses:[::]\n"
			for _, ip := range info.IPAddresses {
				ipText += fmt.Sprintf("[::b]%s[::]\n", ip)
			}
			ui.middle.SetText(ipText)
		})

		time.Sleep(900 * time.Millisecond)
	}
}

// Run runs the UI application
func (ui *UI) Run() error {
	return ui.app.SetRoot(ui.mainFlex, true).EnableMouse(true).Run()
}
