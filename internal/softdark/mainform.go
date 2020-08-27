package softdark

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softdark/pkg/monitorInfo"
	gtkHelper "github.com/hultan/softteam/gtk"
	"os"
	"path"
	"path/filepath"
)

type MainForm struct {
	Window    *gtk.ApplicationWindow
}

func NewMainForm() *MainForm {
	mainForm := new(MainForm)
	return mainForm
}

func (m *MainForm) OpenMainForm(app *gtk.Application) {
	// Initialize gtk
	gtk.Init(&os.Args)

	// Create a new gtk helper
	resourcePath := path.Join(m.getExecutablePath(),"../assets","main.glade")
	builder, err := gtk.BuilderNewFromFile(resourcePath)
	if err!=nil {
		panic("Failed to create builder")
	}
	helper := gtkHelper.GtkHelperNew(builder)

	// Get the main window from the glade file
	window, err := helper.GetApplicationWindow("main_window")
	if err!=nil {
		panic("Failed to find main_window")
	}

	m.Window = window

	// Set up main window
	window.SetApplication(app)
	title := fmt.Sprintf("SoftDark")
	window.SetTitle(title)
	window.SetDefaultSize(1024, 768)

	// Hook up the destroy event
	_,err = window.Connect("destroy", func() {
		m.CloseMainForm()
	})
	if err!=nil {
		panic("Failed to connect the mainForm.destroy event")
	}

	// Get fixed area
	monitorArea, err := helper.GetFixed("monitor_area")
	if err!=nil {
		panic("Failed to get monitor_area")
	}
	//m.setupMonitors(monitorArea)

	// Quit button
	button, err := helper.GetButton("quit_button")
	if err!=nil {
		panic("Failed to find quit_button")
	}
	_, err =button.Connect("clicked", func() {
		window.Close()
		m.CloseMainForm()
	})
	if err!=nil {
		panic("Failed to connect the quit_button.clicked event")
	}

	// Refresh button
	button, err = helper.GetButton("refresh_button")
	if err!=nil {
		panic("Failed to find refresh_button")
	}
	_, err =button.Connect("clicked", func() {
		m.setupMonitors(monitorArea)
	})
	if err!=nil {
		panic("Failed to connect the quit_button.clicked event")
	}

	m.setupMonitors(monitorArea)

	// Show the main window
	window.ShowAll()
}

func (m *MainForm) CloseMainForm() {

}

func (m *MainForm) setupMonitors(monitorArea *gtk.Fixed) {
	monitors, err := m.getMonitors()
	if err!=nil {
		// TODO : What to do here
	}
	height, width := m.getSize(monitors)
	monitorArea.SetSizeRequest(width/10, height/10)

	for _,value := range monitors {
		//fmt.Println(value.Connection)
		button, err:=gtk.ButtonNewWithLabel(fmt.Sprintf("%d: %s", value.Number, value.Connection))
		if err!=nil {
			// TODO : What to do here
		}
		button.SetSizeRequest(value.Width/10,value.Height/10)
		monitorArea.Put(button, value.Left/10, value.Top/10)
	}

	monitorArea.ShowAll()
}

func (m *MainForm) getMonitors() ([]monitorInfo.Monitor, error){
	info:=monitorInfo.NewMonitorInfo()
	monitors, err := info.GetMonitorInfo()
	return monitors, err
}

// getExecutablePath : Returns the path of the executable
func (m *MainForm) getExecutablePath() string {
	executable, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(executable)
}

func (m *MainForm) getSize(monitors []monitorInfo.Monitor) (height, width int) {
	maxWidth, maxHeight := 0,0

	for _,value := range monitors {
		if value.Top + value.Height > maxHeight {
			maxHeight = value.Top + value.Height
		}
		if value.Left + value.Width > maxWidth {
			maxWidth = value.Left + value.Width
		}
	}

	return maxHeight, maxWidth
}