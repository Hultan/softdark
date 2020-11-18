package screenShot

import (
	"bytes"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/kbinani/screenshot"
	"image/png"
)

type ScreenShot struct {

}

func NewScreenShot() *ScreenShot {
	return new(ScreenShot)
}

// getScreenShot : Get a screenshot of a monitor, with the specified width/height
func (s *ScreenShot) GetScreenShot(monitor, width, height int) (*gtk.Image, error) {
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
	scaledPixbuf, err := imagePixBuf.ScaleSimple(width, height, gdk.INTERP_HYPER)
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
