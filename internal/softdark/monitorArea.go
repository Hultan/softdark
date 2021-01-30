package softdark

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softdark/internal/screenShot"
	"github.com/hultan/softdark/internal/softmonitorInfo"
	"github.com/hultan/softdark/internal/tools"
	"log"
	"sort"
)

const buttonImageMargin = 8
const buttonPadding = 6
const buttonPaddingLeft = 5

type MonitorArea struct {
	Area           *gtk.Fixed
	Monitors       []*Monitor
	LastAllocation *gtk.Allocation
}

func NewSoftDark(area *gtk.Fixed) *MonitorArea {
	monitorArea := new(MonitorArea)
	monitorArea.Area = area
	return monitorArea
}

func (s *MonitorArea) Init() {
	// Clear previous buttons from MonitorArea
	s.clearMonitorArea()

	// Refresh SoftDark monitor info
	err := s.refreshMonitorInfo()
	tools.ErrorCheckWithPanic(err, "SoftDark.RefreshMonitorInfo() failed")

	// Calculate scale factor based on window size
	scaleFactor := s.calculateScaleFactor()

	// Sort the monitors based on its left position
	sort.Slice(s.Monitors, func(i, j int) bool {
		if s.Monitors[i].Info.Left < s.Monitors[j].Info.Left {
			return true
		}
		return false
	})

	var padding = 0
	for i := 0; i < len(s.Monitors); i++ {
		currentMonitor := s.Monitors[i]

		// Create a new button
		button, err := gtk.ButtonNew()
		tools.ErrorCheckWithPanic(err, "failed to add button")
		// Store pointer to button
		currentMonitor.Button = button

		// Calculate monitor button size & position
		width := int(float64(currentMonitor.Info.Width) / scaleFactor)
		height := int(float64(currentMonitor.Info.Height) / scaleFactor)
		left := int(float64(currentMonitor.Info.Left) / scaleFactor)
		top := int(float64(currentMonitor.Info.Top) / scaleFactor)

		// Set button size and position on monitor area
		button.SetSizeRequest(width, height)
		s.Area.Put(button, left+padding+buttonPaddingLeft, top)

		// Connect click event
		_, _ = button.Connect("clicked", s.onButtonClicked, currentMonitor)

		// Increase padding for next button
		padding += buttonPadding
	}

	s.updateScreenshots()
	s.Area.ShowAll()
}

// calculateScaleFactor : Calculate the current scale factor
func (s *MonitorArea) calculateScaleFactor() float64 {
	// Get total size for all monitors
	height, width := s.getSize(s.Monitors)
	// Get window size
	allocation := s.Area.GetAllocation()

	//
	heightFactor := float64(height) / float64(allocation.GetHeight())
	widthFactor := float64(width) / float64(allocation.GetWidth())

	factor := widthFactor
	if heightFactor > widthFactor {
		factor = heightFactor
	}
	if factor == 0 {
		return 10
	}
	return factor
}

// getSize : Get the maximum size of all the monitors
func (s *MonitorArea) getSize(monitors []*Monitor) (height, width int) {
	maxWidth, maxHeight := 0, 0

	for _, value := range monitors {
		if value.Info.Top+value.Info.Height > maxHeight {
			maxHeight = value.Info.Top + value.Info.Height
		}
		if value.Info.Left+value.Info.Width > maxWidth {
			maxWidth = value.Info.Left + value.Info.Width
		}
	}

	return maxHeight, maxWidth
}

// refreshMonitorInfo : Refreshes the monitor hardware info
func (s *MonitorArea) refreshMonitorInfo() error {
	// Get monitor hardware info
	monitorInfoTool := softmonitorInfo.NewSoftMonitorInfo()
	monitorInfoDetails, err := monitorInfoTool.GetMonitorInfo()
	if err != nil {
		return err
	}

	s.Monitors = make([]*Monitor, 0)

	for _, info := range monitorInfoDetails {
		monitor := newMonitor(info)
		s.Monitors = append(s.Monitors, monitor)
	}

	return nil
}

// clearMonitorArea : Clears the monitor area (removes the buttons)
func (s *MonitorArea) clearMonitorArea() {
	for _, value := range s.Monitors {
		if value.Button != nil {
			s.Area.Remove(value.Button)
		}
	}
}

func (s *MonitorArea) onButtonClicked(button *gtk.Button, currentMonitor *Monitor) {
	if currentMonitor.Form.IsVisible {
		currentMonitor.Form.Hide()
	} else {
		currentMonitor.Form.Show(currentMonitor.Info)
	}
}

func (s *MonitorArea) updateScreenshots() {
	// Calculate scale factor based on window size
	scaleFactor := s.calculateScaleFactor()

	for i := 0; i < len(s.Monitors); i++ {
		currentMonitor := s.Monitors[i]

		// Calculate monitor button size & position
		width := int(float64(currentMonitor.Info.Width) / scaleFactor)
		height := int(float64(currentMonitor.Info.Height) / scaleFactor)

		// Add a screenshot to the button
		screenshot := screenShot.NewScreenShot()
		image, err := screenshot.GetScreenShot(currentMonitor.Info.Number, width-buttonImageMargin*2, height-buttonImageMargin*2)
		if err != nil {
			log.Println(err)
		} else {
			currentMonitor.Button.SetImage(image)
		}
	}
}