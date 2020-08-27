package test

import (
	"github.com/hultan/softdark/pkg/monitorInfo"
	"testing"
)

func TestMonitors(t *testing.T) {
	m := monitorInfo.NewMonitors()
	got, err := m.GetMonitors()
	if err!=nil {
		t.Error(err)
	}
	if len(got)!=3 {
		t.Errorf("incorrect number of monitors")
	}
	if got[0].Main != true && got[1].Main != false && got[2].Main != false {
		t.Errorf("first monitor should be main")
	}
	for key,value := range got {
		if value.Height!=1080 {
			t.Errorf("wrong height (%d) for monitor {%d}", value.Height, key)
		}
		if value.Width!=1920 {
			t.Errorf("wrong width (%d) for monitor {%d}", value.Width, key)
		}
	}
	if got[0].Connection != "DisplayPort-0" {
		t.Errorf("wrong connection (%s) for monitor 0", got[0].Connection)
	}
	if got[1].Connection != "DisplayPort-1" {
		t.Errorf("wrong connection (%s) for monitor 1", got[1].Connection)
	}
	if got[2].Connection != "DVI-D-0" {
		t.Errorf("wrong connection (%s) for monitor 2", got[2].Connection)
	}
}