package main

import (
	"github.com/justmiles/epd/lib/dashboard"
	"log"
)

func main() {
	d, err := dashboard.NewDashboard(dashboard.WithEPD("epd7in5v2"))
	if err != nil {
		log.Panic(err.Error())
	}

	d.EPDService.HardwareInit()
	d.EPDService.Clear()

	err = d.DisplayText("hello world")
	if err != nil {
		panic(err)
	}

	d.EPDService.Sleep()
}
