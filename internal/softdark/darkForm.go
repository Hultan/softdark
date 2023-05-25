package softdark

import (
	_ "embed"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softdark/internal/monitorInfo"
	"github.com/hultan/softdark/internal/tools"
	gtkHelper "github.com/hultan/softteam-tools/pkg/gtk-helper"
	"log"
)

//go:embed assets/main.glade
var mainGlade string

type DarkForm struct {
	Window    *gtk.Window
	IsVisible bool
}

func (d *DarkForm) init(info monitorInfo.MonitorInfo) {
	// Create a new gtk helper
	builder, err := gtk.BuilderNewFromString(mainGlade)
	if err != nil {
		log.Fatal("Failed to create builder : " + err.Error())
	}
	helper := gtkHelper.GtkHelperNew(builder)

	// Get the main window from the glade file
	window, err := helper.GetWindow("dark_window")
	if err != nil {
		log.Fatal("Failed to find dark_window : " + err.Error())
	}

	d.Window = window

	window.Move(info.Left, info.Top)
	window.SetSizeRequest(info.Width, info.Height)
	window.Fullscreen()

	// Hook up the destroy event
	_ = window.Connect("button-press-event", d.Hide)

	// Create CSS provider
	provider, _ := gtk.CssProviderNew()
	if err := provider.LoadFromPath(tools.GetResourcePath("assets", "softdark.css")); err != nil {
		log.Println(err)
	}

	// Set CSS provider
	context, _ := window.GetStyleContext()
	context.AddProvider(provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func (d *DarkForm) Show(info monitorInfo.MonitorInfo) {
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
