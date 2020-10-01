package softdark

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softdark/pkg/softmonitorInfo"
)

type Monitor struct {
	Button   *gtk.Button
	Info     softmonitorInfo.MonitorInfo
	Form     DarkForm
}

func newMonitor(info softmonitorInfo.MonitorInfo) *Monitor {
	monitor := new(Monitor)
	monitor.Info = info
	return monitor
}
