package softdark

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softdark/internal/tools"
	"github.com/hultan/softdark/pkg/softmonitorInfo"
	gtkHelper "github.com/hultan/softteam-tools/pkg/gtk-helper"
	"log"
)

type DarkForm struct {
	Window    *gtk.Window
	IsVisible bool
}

func (d *DarkForm) init(info softmonitorInfo.MonitorInfo) {
	// Create a new gtk helper
	builder, err := gtk.BuilderNewFromFile(tools.GetResourcePath("../assets", "main.glade"))
	tools.ErrorCheckWithPanic(err, "Failed to create builder")
	helper := gtkHelper.GtkHelperNew(builder)

	// Get the main window from the glade file
	window, err := helper.GetWindow("dark_window")
	tools.ErrorCheckWithPanic(err, "Failed to find dark_window")

	d.Window = window

	window.Move(info.Left, info.Top)
	window.SetSizeRequest(info.Width, info.Height)
	window.Fullscreen()

	// Hook up the destroy event
	_, err = window.Connect("button-press-event", d.Hide)
	tools.ErrorCheckWithPanic(err, "Failed to connect the dark_window.Close() event")

	// Create CSS provider
	provider, _ := gtk.CssProviderNew()
	if err := provider.LoadFromPath(tools.GetResourcePath("../assets", "softdark.css")); err != nil {
		log.Println(err)
	}

	// Set CSS provider
	context, _ := window.GetStyleContext()
	context.AddProvider(provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func (d *DarkForm) Show(info softmonitorInfo.MonitorInfo) {
	if d.Window == nil {
		d.init(info)
	}
	d.Window.Show()
	d.IsVisible = true
}

func (d *DarkForm) Hide() {
	d.Window.Hide()
	d.IsVisible = false
}