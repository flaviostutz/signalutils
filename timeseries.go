package signalutils

import (
	"time"
)

type TimeValue struct {
	Time  time.Time
	Value float64
}

type Timeseries struct {
	TimeseriesSpan time.Duration
	Values         []TimeValue
	// gc             int
}

func NewTimeseries(maxTimeseriesSpan time.Duration) Timeseries {
	return Timeseries{
		TimeseriesSpan: maxTimeseriesSpan,
		Values:         make([]TimeValue, 0),
	}
}

func (t *Timeseries) AddSample(value float64) {
	t.Values = append(t.Values, TimeValue{time.Now(), value})
	// t.gc = t.gc + 1
	// if t.gc > 5 {
	i1, _, ok := t.FindPos(time.Now().Add(-t.TimeseriesSpan - 1*time.Second))
	if ok && i1 > 1 {
		t.Values = t.Values[i1-1:]
	}
	// t.gc = 0
	// }
}

func (t *Timeseries) GetValue(time time.Time) (TimeValue, bool) {
	i1, i2, ok := t.FindPos(time)
	if !ok {
		return TimeValue{}, false
	}
	if i1 == i2 {
		return t.Values[i1], true
	}
	v1 := t.Values[i1]
	v2 := t.Values[i2]
	// fmt.Printf("%f %f", v1.value, v2.value)
	td := float64(v2.Time.UnixNano() - v1.Time.UnixNano())
	vd := v2.Value - v1.Value
	vdr := v1.Value + ((vd / td) * float64(time.UnixNano()-v1.Time.UnixNano()))
	return TimeValue{time, vdr}, true
}

func (t *Timeseries) Size() int {
	return len(t.Values)
}

func (t *Timeseries) FindPos(time time.Time) (int, int, bool) {
	for i1, v1 := range t.Values {
		if v1.Time == time {
			return i1, i1, true
		}
		i2 := i1 + 1
		if i2 < len(t.Values) {
			v2 := t.Values[i2]
			if v2.Time == time {
				return i2, i2, true
			}
			if time.After(v1.Time) && time.Before(v2.Time) {
				return i1, i2, true
			}
		}
	}
	return -1, -1, false
}

func (t *Timeseries) Reset() {
	t.Values = make([]TimeValue, 0)
}

func (t *Timeseries) GetLastValue() (TimeValue, bool) {
	l := len(t.Values)
	if l == 0 {
		return TimeValue{}, false
	}
	return t.Values[l-1], true
}
