package softdark

import (
	"bytes"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softdark/internal/tools"
	"github.com/hultan/softdark/pkg/softmonitorInfo"
	"github.com/kbinani/screenshot"
	"image/png"
	"log"
)

const buttonImageMargin = 8

type SoftDark struct {
	MonitorArea *gtk.Fixed
	// TODO : Replace with map???
	Monitors       []*Monitor
	LastAllocation *gtk.Allocation
}

func NewSoftDark(monitorArea *gtk.Fixed) *SoftDark {
	softDark := new(SoftDark)
	softDark.MonitorArea = monitorArea
	return softDark
}

func (s *SoftDark) Init() {
	// Clear previous buttons from MonitorArea
	s.clearMonitorArea()

	// Refresh SoftDark monitor info
	err := s.refreshMonitorInfo()
	tools.ErrorCheckWithPanic(err, "SoftDark.RefreshMonitorInfo() failed")

	scaleFactor := s.calculateScaleFactor()

	for i := 0; i < len(s.Monitors); i++ {
		currentMonitor := s.Monitors[i]

		// Create a new button
		button, err := gtk.ButtonNew()
		tools.ErrorCheckWithPanic(err, "failed to add button")
		// Store pointer to button
		currentMonitor.Button = button

		// Calculate monitor button size & position
		width := int(float64(currentMonitor.Info.Width)/scaleFactor)
		height := int(float64(currentMonitor.Info.Height)/scaleFactor)
		left := int(float64(currentMonitor.Info.Left)/scaleFactor)
		top := int(float64(currentMonitor.Info.Top)/scaleFactor)
		//fmt.Printf("Placing button at (%v,%v), size (%v,%v)\n", top, left, height, width)

		// Set button size
		button.SetSizeRequest(width, height)
		// Place button on MonitorArea
		s.MonitorArea.Put(button, left, top)

		// Add a screenshot to the button
		image, err := s.getScreenShot(i,width, height)
		if err != nil {
			log.Println(err)
		} else {
			button.SetImage(image)
		}

		_,_ = button.Connect("clicked", s.onButtonClicked, currentMonitor)
	}

	s.MonitorArea.ShowAll()
}

// getScreenShot : Get a screenshot of a monitor, with the specified width/height
func (s *SoftDark) getScreenShot(monitor, width, height int) (*gtk.Image, error) {
	// Get screenshot of monitor
	screenImage, err := screenshot.CaptureDisplay(monitor)
	if err != nil {
		return nil, err
	}
	// Convert screenshot to byte array
	var b bytes.Buffer
	err = png.Encode(&b, screenImage)
	if err != nil {
		return nil, err
	}
	// Create a PixBufLoader
	loader, err := gdk.PixbufLoaderNew()
	if err != nil {
		return nil, err
	}
	// Write byte array to PixBufLoader
	imagePixBuf, err := loader.WriteAndReturnPixbuf(b.Bytes())
	if err != nil {
		return nil, err
	}
	// Scale image down to a reasonable size
	scaledPixbuf, err := imagePixBuf.ScaleSimple(width-buttonImageMargin*2, height-buttonImageMargin*2, gdk.INTERP_HYPER)
	if err != nil {
		return nil, err
	}
	// Create an gtk.Image from the PixBuf
	image, err := gtk.ImageNewFromPixbuf(scaledPixbuf)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// calculateScaleFactor : Calculate the current scale factor
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
func (s *SoftDark) getSize(monitors []*Monitor) (height, width int) {
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
func (s *SoftDark) refreshMonitorInfo() error {
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
func (s *SoftDark) clearMonitorArea() {
	for _, value := range s.Monitors {
		if value.Button != nil {
			s.MonitorArea.Remove(value.Button)
		}
	}
}

func (s *SoftDark) onButtonClicked(button *gtk.Button, currentMonitor *Monitor) {
	if currentMonitor.Form.IsVisible {
		currentMonitor.Form.Hide()
	} else {
		currentMonitor.Form.Show(currentMonitor.Info)
	}
}