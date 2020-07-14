# signalutils

Event/Signal processing utilities lib for Golang. Online moving averager, linear regression, timed value in Golang etc

## Usage

```go
package main
import (
	"fmt"
	"github.com/flaviostutz/signalutils"
)

func main() {
	fmt.Printf("Moving Average\n")
	ma := signalutils.NewMovingAverage(5)
	ma.AddSample(0.00)
	ma.AddSample(99999.00)
	fmt.Printf("Average is %f\n", ma.Average())
	ma.AddSample(1000.00)
	ma.AddSample(2000.00)
	fmt.Printf("Average is %f\n", ma.Average())
	ma.AddSample(3000.00)
	ma.AddSample(4000.00)
	fmt.Printf("Average is %f\n", ma.Average())
	ma.AddSample(5000.00)
	ma.AddSample(6000.00)
	fmt.Printf("Average is %f\n", ma.Average())
}

```
Results
```
Moving Average
Average is 49999.500000
Average is 25749.750000
Average is 21999.800000
Average is 4000.000000
```

## Utilities

* MovingAverage - add values to an array with a fixed max size and query for the average of values in this fixed size array

* SchmittTrigger - set current values and track current up/down state based on schmitt trigger algorithm

* StateTracker - set state identifications and if state has lots of successive repetitions, perform a state transition. Useful to filter out noises from state changes.

* Timeseries - time/value array with max time span for keeping size at control

* TimeseriesRate - add counter values to a timeseries and query for rate at any time range. Something that ressembles "rate(metric_name[1m])" on Prometheus queries, for example.

