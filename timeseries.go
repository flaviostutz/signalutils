package signalutils

import (
	"time"

	"github.com/gonum/stat"
)

//TimeValue a point in time
type TimeValue struct {
	Time  time.Time
	Value float64
}

//Timeseries utility
type Timeseries struct {
	TimeseriesSpan time.Duration
	Values         []TimeValue
	// gc             int
}

//NewTimeseries create a new timeseries with a limited size in time.
//After that limit older values will be deleted from time to time to
//avoid too much memory usage
func NewTimeseries(maxTimeseriesSpan time.Duration) Timeseries {
	return Timeseries{
		TimeseriesSpan: maxTimeseriesSpan,
		Values:         make([]TimeValue, 0),
	}
}

//Add add a new sample to this timeseries using time.Now()
func (t *Timeseries) Add(value float64) {
	t.Values = append(t.Values, TimeValue{time.Now(), value})
	// t.gc = t.gc + 1
	// if t.gc > 5 {
	i1, _, ok := t.Pos(time.Now().Add(-t.TimeseriesSpan - 1*time.Second))
	if ok && i1 > 1 {
		t.Values = t.Values[i1-1:]
	}
	// t.gc = 0
	// }
}

//Get get value in a specific time in timeseries.
//If time is between two points inside timeseries, the value will
//be interpolated according to the requested time and neighboring values
func (t *Timeseries) Get(time time.Time) (TimeValue, bool) {
	i1, i2, ok := t.Pos(time)
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

//Size current number of elements in this timeseries
func (t *Timeseries) Size() int {
	return len(t.Values)
}

//Pos searches for which two point indexes are between the desired time
//Find the time is exacly the same as a point time, the two returned indexes will be equal
func (t *Timeseries) Pos(time time.Time) (int, int, bool) {
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

//Reset remove all elements from this timeseries
func (t *Timeseries) Reset() {
	t.Values = make([]TimeValue, 0)
}

//Last get last point in time element, the head element
func (t *Timeseries) Last() (TimeValue, bool) {
	l := len(t.Values)
	if l == 0 {
		return TimeValue{}, false
	}
	return t.Values[l-1], true
}

//Avg calculates the average value of points compreended between time 'from' and 'to'
//No interpolation is used here
func (t *Timeseries) Avg(from time.Time, to time.Time) (float64, bool) {
	sum := 0.0
	c := 0
	for _, v := range t.Values {
		if (v.Time == from || v.Time.After(from)) && (v.Time == to || v.Time.Before(to)) {
			sum = sum + v.Value
			c = c + 1
		}
	}
	return sum / float64(c), true
}

//ValuesRange get values in time range
//returns an array of TimeValue and and array with just the float values
func (t *Timeseries) ValuesRange(from time.Time, to time.Time) ([]TimeValue, []float64) {
	vs := make([]TimeValue, 0)
	values := make([]float64, 0)
	for _, v := range t.Values {
		vs = append(vs, v)
		values = append(values, v.Value)
	}
	return vs, values
}

//StdDev calculates the standard deviation and mean for the time range
//returns standard deviation and mean value
func (t *Timeseries) StdDev(from time.Time, to time.Time) (std float64, mean float64) {
	_, values := t.ValuesRange(from, to)
	mean, std = stat.MeanStdDev(values, nil)
	return std, mean
}

//LinearRegression calculates the linear regression coeficients for the time range
//x is in range of time.UnixNano()
//returns alpha and beta as for y = alpha + beta*x and rsquared with fit from 0-1
func (t *Timeseries) LinearRegression(from time.Time, to time.Time) (alpha float64, beta float64, rsquared float64) {
	vs, _ := t.ValuesRange(from, to)
	x := make([]float64, 0)
	y := make([]float64, 0)
	for _, v := range vs {
		// x = append(x, float64(v.Time.UnixNano()-vs[0].Time.UnixNano()))
		x = append(x, float64(v.Time.UnixNano()))
		y = append(y, v.Value)
	}
	alpha, beta = stat.LinearRegression(x, y, nil, false)
	rsquared = stat.RSquared(x, y, nil, alpha, beta)
	return alpha, beta, rsquared
}
