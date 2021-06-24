package main

import (
	"github.com/justmiles/epd/lib/dashboard"
	"log"
)

func main() {
	d, e := dashboard.NewDashboard(dashboard.WithEPD("epd7in5v2"))
	if e != nil {
		log.Fatalln(e)
	}
	d.EPDService.HardwareInit()

	e = d.DisplayText("hello")
	if e != nil {
		log.Fatalln(e)
	}

}
