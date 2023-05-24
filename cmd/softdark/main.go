package main

import (
	"github.com/hultan/softdark/internal/softdark"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	ApplicationId    = "se.softteam.softdark"
	ApplicationFlags = glib.APPLICATION_FLAGS_NONE
)

func main() {
	// Create a new application
	application, err := gtk.ApplicationNew(ApplicationId, ApplicationFlags)
	if err != nil {
		panic("Failed to create GTK Application")
	}

	mainForm := softdark.NewMainForm()
	// Hook up the activate event handler
	_ = application.Connect("activate", mainForm.OpenMainForm)

	// Start the application (and exit when it is done)
	code := application.Run(nil)
	os.Exit(code)
}
