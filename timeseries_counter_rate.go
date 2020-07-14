package signalutils

import (
	"fmt"
	"time"
)

type TimeseriesCounterRate struct {
	Timeseries Timeseries
	ccounter   float64
}

func NewTimeseriesCounterRate(timeseriesSpan time.Duration) TimeseriesCounterRate {
	ts := NewTimeseries(timeseriesSpan)
	return TimeseriesCounterRate{
		Timeseries: ts,
	}
}

func (t *TimeseriesCounterRate) Inc(value float64) error {
	if value < 0 {
		return fmt.Errorf("value cannot be negative")
	}
	t.ccounter = t.ccounter + value
	t.Timeseries.AddSample(t.ccounter)
	return nil
}

func (t *TimeseriesCounterRate) Rate(timeSpan time.Duration) (float64, bool) {
	if timeSpan > t.Timeseries.TimeseriesSpan {
		return 0, false
	}

	v2, ok := t.Timeseries.GetLastValue()
	if !ok {
		return 0, false
	}

	v1, ok := t.Timeseries.GetValue(v2.Time.Add(-timeSpan))
	if !ok {
		return 0, false
	}

	td := float64(v2.Time.UnixNano()-v1.Time.UnixNano()) / 1000000000
	vd := v2.Value - v1.Value

	return vd / td, true
}
