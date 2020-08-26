package monitorInfo

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type MonitorInfo struct {
}

type Monitor struct {
	Number int
	Main bool
	Connection string
	Width int
	Height int
}

func NewMonitorInfo() *MonitorInfo {
	return new(MonitorInfo)
}

func (m *MonitorInfo) GetMonitorInfo() ([]Monitor, error) {
	// Call xrandr to get monitor info
	monitorInfo, err := getMonitorInfo()
	if err!=nil {
		return nil, err
	}

	// Parse monitor info
	monitors, err := parseMonitorsInfo(monitorInfo)

	return monitors, err
}

func getMonitorInfo() (string, error) {
	// Call xrandr to get monitor info
	out, err := exec.Command("xrandr", "--listmonitors").Output()
	output := string(out[:])
	return output, err
}

func parseMonitorsInfo(monitorInfo string) ([]Monitor,error) {
	// Split xrandr result into lines
	lines := strings.Split(monitorInfo, "\n")
	var monitors []Monitor
	// Start parsing at row 1, since row 0
	// contains the monitor count, which we
	// don't need
	for i:=1;i<len(lines);i++ {
		line := lines[i]
		// Ignore empty lines (there is usually one at the end)
		if line == "" {
			continue
		}
		// Parse the monitor line
		monitor, err := getMonitor(line)
		if err!=nil {
			return nil, err
		}
		// Append the monitor to the monitors slice
		monitors = append(monitors, *monitor)
	}
	// Return the monitors
	return monitors, nil
}

func getMonitor(monitorInfoString string) (*Monitor,error) {
	// Regular expression to parse the output of xrandr --listmonitors
	regexString := `\s(?P<Number>\d+): \+(?P<Main>\*?)(?P<Connection>.+) (?P<Width>\d+)\/.+x(?P<Height>\d+)\/.+\s`
	r := regexp.MustCompile(regexString)
	matches := r.FindStringSubmatch(monitorInfoString)

	if len(matches)>1 {
		m:=Monitor{}
		// Get monitor number
		number, err := strconv.Atoi(matches[1])
		if err!=nil {
			return nil, err
		}
		m.Number = number
		// Get main monitor
		if matches[2]=="*" {
			m.Main = true
		}
		// Get monitor connection
		m.Connection = matches[3]
		// Get monitor width
		width, err := strconv.Atoi(matches[4])
		if err!=nil {
			return nil, err
		}
		m.Width = width

		// Get monitor height
		height, err := strconv.Atoi(matches[5])
		if err!=nil {
			return nil, err
		}
		m.Height = height

		return &m, nil
	} else {
		return nil, errors.New("failed to parse monitor")
	}

}
