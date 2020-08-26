package main

import (
	"fmt"
	"github.com/hultan/softdark/pkg/monitorInfo"
)

func main() {
	m:=monitorInfo.NewMonitorInfo()
	monitors, err := m.GetMonitorInfo()
	if err!=nil {
		fmt.Println(err.Error())
	}
	fmt.Println(monitors)
}
