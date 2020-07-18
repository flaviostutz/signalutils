package signalutils

import (
	"fmt"
	"time"
)

//TimeseriesCounterRate this is a utility for storing counter values in time
//while enabling the measurement of rates in various time spans with without
//having to perform average over all points.
//The optimization strategy here is based on the fact that this timeseries contains
//a counter, so that averages between two times are calculated by just
//averaging the first and last points, not all the points between.
//Very useful for metrics monitoring. See more at https://prometheus.io/docs/concepts/metric_types/#counter
type TimeseriesCounterRate struct {
	Timeseries Timeseries
	ccounter   float64
}

//NewTimeseriesCounterRate creates a time timeseries with max time span of timeseriesSpan
func NewTimeseriesCounterRate(timeseriesSpan time.Duration) TimeseriesCounterRate {
	ts := NewTimeseries(timeseriesSpan)
	return TimeseriesCounterRate{
		Timeseries: ts,
	}
}

//Inc increments the last value from the timeseries by 'value' and sets
//add the new point with time.Now() time
func (t *TimeseriesCounterRate) Inc(value float64) error {
	if value < 0 {
		return fmt.Errorf("value cannot be negative")
	}
	t.ccounter = t.ccounter + value
	t.Timeseries.Add(t.ccounter)
	return nil
}

//Set sets the absolute value at time time.Now(). The value cannot be less
//then last value from the timeseries as this must be a counter
func (t *TimeseriesCounterRate) Set(value float64) error {
	if value < t.ccounter {
		return fmt.Errorf("value cannot be less than current counter")
	}
	t.Timeseries.Add(value)
	return nil
}

//Rate calculates the rate of change between the last point in time of this timeseries
//and the time in past, specified by timeSpan
func (t *TimeseriesCounterRate) Rate(timeSpan time.Duration) (float64, bool) {
	if timeSpan > t.Timeseries.TimeseriesSpan {
		return 0, false
	}

	v2, ok := t.Timeseries.Last()
	if !ok {
		return 0, false
	}

	v1, ok := t.Timeseries.Get(v2.Time.Add(-timeSpan))
	if !ok {
		return 0, false
	}

	td := float64(v2.Time.UnixNano()-v1.Time.UnixNano()) / 1000000000
	vd := v2.Value - v1.Value

	return vd / td, true
}
