package softdark

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softdark/internal/tools"
	gtkHelper "github.com/hultan/softteam-tools/pkg/gtk-helper"
	"log"
	"os"
)

type MainForm struct {
	Window         *gtk.ApplicationWindow
	LastAllocation *gtk.Allocation
	SoftDark       *SoftDark
}

// NewMainForm : Creates a new MainForm object
func NewMainForm() *MainForm {
	mainForm := new(MainForm)
	return mainForm
}

// OpenMainForm : Opens the MainForm window
func (m *MainForm) OpenMainForm(app *gtk.Application) {
	// Initialize gtk
	gtk.Init(&os.Args)

	// Create a new gtk helper
	builder, err := gtk.BuilderNewFromFile(tools.GetResourcePath("../assets", "main.glade"))
	tools.ErrorCheckWithPanic(err, "Failed to create builder")
	helper := gtkHelper.GtkHelperNew(builder)

	// Get the main window from the glade file
	window, err := helper.GetApplicationWindow("main_window")
	tools.ErrorCheckWithPanic(err, "Failed to find main_window")

	m.Window = window

	// Set up main window
	window.SetApplication(app)
	window.SetTitle("SoftDark")

	// Hook up the destroy event
	_, err = window.Connect("destroy", window.Close)
	tools.ErrorCheckWithPanic(err, "Failed to connect the mainForm.destroy event")

	// Get fixed area
	monitorArea, err := helper.GetFixed("monitor_area")
	tools.ErrorCheckWithPanic(err, "Failed to get monitor_area")
	m.SoftDark = NewSoftDark(monitorArea)

	// Quit button
	button, err := helper.GetButton("quit_button")
	tools.ErrorCheckWithPanic(err, "Failed to find quit_button")
	_, err = button.Connect("clicked", window.Close)
	tools.ErrorCheckWithPanic(err, "Failed to connect the quit_button.clicked event")

	// Create CSS provider
	provider, _ := gtk.CssProviderNew()
	if err := provider.LoadFromPath(tools.GetResourcePath("../assets", "softdark.css")); err != nil {
		log.Println(err)
	}

	// Set CSS provider
	context, _ := monitorArea.GetStyleContext()
	context.AddProvider(provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	// Show the main window
	window.ShowAll()

	m.SoftDark.Init()
}
