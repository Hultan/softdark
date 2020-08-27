package main

import (
	"fmt"
	"github.com/hultan/softdark/pkg/monitorInfo"
)

func main() {
	m:=monitorInfo.NewMonitors()
	monitors, err := m.GetMonitors()
	if err!=nil {
		fmt.Println(err.Error())
	}
	fmt.Println(monitors)
}
