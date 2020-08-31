# signalutils

[![Build Status](https://travis-ci.org/flaviostutz/signalutils.svg?branch=master)](https://travis-ci.org/flaviostutz/signalutils)

Event/Signal processing utilities lib for Golang. Online moving averager, linear regression, timed value, worker frequency control in Golang etc

See API documentation at https://pkg.go.dev/github.com/flaviostutz/signalutils?tab=doc

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

```golang
	ma := NewMovingAverageTimeWindow(1*time.Second, 10)
	ma.AddSample(1000)
	ma.AddSample(2000)
	ma.AddSample(3000)
	ma.AddSample(4000)
	ma.AddSample(3000)
	ma.AddSample(2000)
	ma.AddSample(3000)
	ma.AddSample(2000)
	assert.Equal(t, 1000.0, ma.Average())
	time.Sleep(100 * time.Millisecond)
	ma.AddSample(3000)
	assert.Equal(t, 2000.0, ma.Average())
	ma.AddSample(2000)
	assert.Equal(t, 2000.0, ma.Average())
```

* SchmittTrigger - set current values and track current up/down state based on schmitt trigger algorithm

```golang
	st, _ := NewSchmittTrigger(10, 20, false)
	assert.False(t, st.IsUpperRange())
	st.SetCurrentValue(11)
	assert.False(t, st.IsUpperRange())
	st.SetCurrentValue(15)
	assert.False(t, st.IsUpperRange())

	st.SetCurrentValue(21)
	assert.True(t, st.IsUpperRange())
	st.SetCurrentValue(25)
	assert.True(t, st.IsUpperRange())
	st.SetCurrentValue(18)
	assert.True(t, st.IsUpperRange())
```

* DynamicSchmittTrigger - set current values and track current up/down state based on schmitt trigger algorithm. Trigger points are set dynamically set in a timelly manner.

* StateTracker - set state identifications and if state has lots of successive repetitions, perform a state transition. Useful to filter out noises from state changes.

```golang
	st := NewStateTracker("state1", 3, onChange, 0, nil, true)
	st.SetTransientState("state2")
	st.SetTransientState("state2")
	assert.Equal(t, "state1", st.CurrentState.Name)
	st.SetTransientState("state2")
	assert.Equal(t, "state2", notifiedNewState.Name)
	st.SetTransientState("state3")
	assert.Equal(t, "state2", st.CurrentState.Name)
	st.SetTransientState("state3")
	st.SetTransientState("state2")
	st.SetTransientState("state3")
	st.SetTransientState("state3")
	assert.Equal(t, "state2", st.CurrentState.Name)
```

* Timeseries - time/value array with max time span for keeping size at control. If you try to get values between time points, interpolation will occur.

```golang
	ts := NewTimeseries(1000 * time.Millisecond)

	ts.AddSample(-100)
	time.Sleep(500 * time.Millisecond)
	ts.AddSample(-1000)

	nv, ok := ts.GetValue(time.Now().Add(-250 * time.Millisecond))
	assert.True(t, ok)
	assert.InDeltaf(t, float64(-555), nv.Value, float64(20), "")
```

* TimeseriesCounterRate - add counter values to a timeseries and query for rate at any time range. Something that ressembles "rate(metric_name[1m])" on Prometheus queries, for example.

```golang
	ts := NewTimeseriesCounterRate(1 * time.Second)

	_, ok := ts.Rate(1 * time.Second)
	assert.False(t, ok)

	ts.Inc(100000) //100000
	time.Sleep(300 * time.Millisecond)
	ts.Inc(100000) //200000
	time.Sleep(300 * time.Millisecond)
	ts.Inc(200000) //400000
	time.Sleep(300 * time.Millisecond)
	ts.Inc(100000) //500000
	time.Sleep(300 * time.Millisecond)

	_, ok = ts.Rate(2 * time.Second)
	assert.False(t, ok)
```

* Worker - useful for workloads that works on a "while true" loop. It launches a Go routine with a function, limits the loop frequency, measures actual frequency and alerts if frequency is outside desired limits.

```golang
	w := StartWorker(context.Background(), "test1", func() error {
		//do some real work here
		time.Sleep(15 * time.Millisecond)
		return nil
	}, 1, 5, true)
	time.Sleep(200 * time.Millisecond)
	assert.True(t, w.active)
	time.Sleep(2000 * time.Millisecond)
	assert.InDeltaf(t, 5, w.CurrentFreq, 2, "")
	assert.InDeltaf(t, 15, w.CurrentStepTime.Milliseconds(), 5, "")
	w.Stop()
	time.Sleep(300 * time.Millisecond)
	assert.False(t, w.active)
```
