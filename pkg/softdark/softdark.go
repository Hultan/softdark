package softdark

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softdark/pkg/softmonitorInfo"
	"github.com/hultan/softdark/pkg/tools"
)

type SoftDark struct {
	MonitorArea *gtk.Fixed
	// TODO : Replace with map???
	Monitors       []Monitor
	LastAllocation *gtk.Allocation
}

func NewSoftDark(monitorArea *gtk.Fixed) *SoftDark {
	softDark := new(SoftDark)
	softDark.MonitorArea = monitorArea
	_, err := softDark.MonitorArea.Connect("size_allocate", softDark.monitorAreaResized)
	tools.ErrorCheckWithPanic(err, "Failed to connect size_allocate signal")
	return softDark
}

func (s *SoftDark) Init() {
	// Clear previous buttons from MonitorArea
	s.clearMonitorArea()

	// Refresh SoftDark monitor info
	err := s.RefreshMonitorInfo()
	tools.ErrorCheckWithPanic(err, "SoftDark.RefreshMonitorInfo() failed")

	scaleFactor := s.calculateScaleFactor()

	for i := 0; i < len(s.Monitors); i++ {
		currentMonitor := s.Monitors[i]

		// Create a new button
		button, err := gtk.ButtonNewWithLabel(fmt.Sprintf("%d: %s", currentMonitor.Info.Number, currentMonitor.Info.Connection))
		tools.ErrorCheckWithPanic(err, "failed to add button")
		// Store pointer to button
		currentMonitor.Button = button

		// Set button size
		button.SetSizeRequest(
			int(float64(currentMonitor.Info.Width)/scaleFactor),
			int(float64(currentMonitor.Info.Height)/scaleFactor))
		// Place button on MonitorArea
		s.MonitorArea.Put(button,
			int(float64(currentMonitor.Info.Left)/scaleFactor),
			int(float64(currentMonitor.Info.Top)/scaleFactor))
	}

	s.MonitorArea.ShowAll()
}

func (s *SoftDark) calculateScaleFactor() float64 {
	height, width := s.getSize(s.Monitors)
	allocation := s.MonitorArea.GetAllocation()
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
func (s *SoftDark) getSize(monitors []Monitor) (height, width int) {
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

func (s *SoftDark) RefreshMonitorInfo() error {
	// Get monitor hardware info
	monitorInfoTool := softmonitorInfo.NewSoftMonitorInfo()
	monitorInfoDetails, err := monitorInfoTool.GetMonitorInfo()
	if err != nil {
		return err
	}

	s.Monitors = make([]Monitor, 0)

	for _, info := range monitorInfoDetails {
		monitor := newMonitor(info)
		s.Monitors = append(s.Monitors, *monitor)
	}

	return nil
}

func (s *SoftDark) clearMonitorArea() {
	for _, value := range s.Monitors {
		if value.Button != nil {
			s.MonitorArea.Remove(value.Button)
		}
	}
}

func (s *SoftDark) monitorAreaResized(monitorArea *gtk.Fixed) {
	allocation := monitorArea.GetAllocation()
	if s.allocationHasChanged(allocation) {
		scaleFactor := s.calculateScaleFactor()

		//fmt.Println("New size : ", allocation.GetHeight(), allocation.GetWidth())
		fmt.Println("Scalefactor : ", scaleFactor)

	}
}

func (s *SoftDark) allocationHasChanged(allocation *gtk.Allocation) bool {
	if s.LastAllocation == nil {
		s.LastAllocation = allocation
		return false
	}
	if allocation.GetHeight() != s.LastAllocation.GetHeight() || allocation.GetWidth() != s.LastAllocation.GetWidth() {
		s.LastAllocation = allocation
		return true
	}
	return false
}
