package main

import (
	"fmt"

	"github.com/flaviostutz/signalutils"
)

func main() {

	fmt.Println(">>>>> signalutils examples <<<<<")

	fmt.Printf("Moving Average\n")
	ma := signalutils.NewMovingAverage(5)
	ma.AddSample(0.00)
	ma.AddSample(99999.00)
	ma.AddSample(1000.00)
	ma.AddSample(2000.00)
	ma.AddSample(3000.00)
	ma.AddSample(4000.00)
	ma.AddSample(5000.00)
	fmt.Printf("Average is %f\n", ma.Average())
}
