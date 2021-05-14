package softdark

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softdark/internal/tools"
	"log"
	"os"
)

type MainForm struct {
	Window         *gtk.ApplicationWindow
	LastAllocation *gtk.Allocation
	Area           *MonitorArea
}

var timeoutContinue = true

// NewMainForm : Creates a new MainForm object
func NewMainForm() *MainForm {
	return new(MainForm)
}

// OpenMainForm : Opens the MainForm window
func (m *MainForm) OpenMainForm(app *gtk.Application) {
	// Initialize gtk
	gtk.Init(&os.Args)

	builder := newSoftBuilder("main.glade")

	// Get the main window from the glade file
	m.Window = builder.getObject("main_window").(*gtk.ApplicationWindow)

	// Set up main window
	m.Window.SetApplication(app)
	m.Window.SetTitle("SoftDark")

	// Hook up the destroy event
	_, err := m.Window.Connect("destroy", m.Window.Close)
	tools.ErrorCheckWithPanic(err, "Failed to connect the mainForm.destroy event")

	// Get fixed area
	monitorArea := builder.getObject("monitor_area").(*gtk.Fixed)
	m.Area = NewSoftDark(monitorArea)

	// TODO : 2021-05-14 : gtk.Fixed have a window if you use the SetHasWindow-function
	// Since gtk.Fixed does not have it's own window
	// you cannot set a background color on it, so we
	// surround it with an EventBox and style the event box
	eventBox := builder.getObject("event_box").(*gtk.EventBox)

	// Quit button
	button := builder.getObject("quit_button").(*gtk.Button)
	_, err = button.Connect("clicked", m.onWindowClose)
	tools.ErrorCheckWithPanic(err, "Failed to connect the quit_button.clicked event")

	// Create CSS provider
	provider, _ := gtk.CssProviderNew()
	if err := provider.LoadFromPath(tools.GetResourcePath("assets", "softdark.css")); err != nil {
		log.Println(err)
	}

	// Set CSS provider
	context, _ := eventBox.GetStyleContext()
	context.AddProvider(provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	// Show the main window
	m.Window.ShowAll()

	m.Area.Init()

	gtk.Main()
}

func (m *MainForm) onWindowClose() {
	// Close and destroy all monitor windows
	for _, monitor := range m.Area.Monitors {
		// Check if the monitor has a window
		if monitor.Form.Window == nil {
			continue
		}

		// Hide the window, if visible
		if monitor.Form.IsVisible {
			monitor.Form.Hide()
		}

		// Destroy the window
		monitor.Form.Window.Destroy()
	}

	// Close main form
	m.Window.Close()
	gtk.MainQuit()
}
