package signalutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRate1(t *testing.T) {
	ts := NewTimeseriesCounterRate(5 * time.Second)

	ts.Inc(20) //20
	time.Sleep(1 * time.Second)

	ts.Inc(30) //50
	time.Sleep(2 * time.Second)

	ts.Inc(150) //200
	time.Sleep(1 * time.Second)

	r, ok := ts.Rate(2 * time.Second)
	assert.True(t, ok)
	assert.InDeltaf(t, float64(75), r, float64(10), "")
}

func TestRate2(t *testing.T) {
	ts := NewTimeseriesCounterRate(3 * time.Second)

	ts.Inc(100000) //100000
	time.Sleep(300 * time.Millisecond)
	ts.Inc(100000) //200000
	time.Sleep(300 * time.Millisecond)
	ts.Inc(200000) //400000
	time.Sleep(300 * time.Millisecond)
	ts.Inc(100000) //500000
	time.Sleep(300 * time.Millisecond)

	r, ok := ts.Rate(750 * time.Millisecond)
	assert.True(t, ok)
	assert.InDeltaf(t, float64(466666), r, float64(10000), "")
}

func TestRate3(t *testing.T) {
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
}

func TestRateRange1(t *testing.T) {
	ts := NewTimeseriesCounterRate(5 * time.Second)

	ts.Inc(20) //20
	time.Sleep(1 * time.Second)

	ts.Inc(30) //50
	time.Sleep(2 * time.Second)

	ts.Inc(150) //200
	time.Sleep(1 * time.Second)

	n1, ok := ts.Timeseries.Last()
	n := n1.Time
	assert.True(t, ok)

	r, ok := ts.RateRange(n.Add(-1*time.Second), n)
	assert.True(t, ok)
	assert.InDeltaf(t, float64(75), r, float64(10), "")

	r, ok = ts.RateRange(n.Add(-2*time.Second), n)
	assert.True(t, ok)
	assert.InDeltaf(t, float64(75), r, float64(10), "")
}

func TestRateOverTime1(t *testing.T) {
	ts := NewTimeseriesCounterRate(5 * time.Second)

	ts.Inc(10) //10
	time.Sleep(500 * time.Millisecond)
	ts.Inc(10) //20
	time.Sleep(500 * time.Millisecond)
	ts.Inc(20) //40
	time.Sleep(500 * time.Millisecond)
	ts.Inc(10) //50
	time.Sleep(500 * time.Millisecond)

	rt, ok := ts.RateOverTime(500*time.Millisecond, 2*time.Second)
	assert.True(t, ok)
	assert.Equal(t, 3, rt.Size())

	// l, ok := rt.Last()
	// assert.True(t, ok)
	// assert.Equal(t, 10.0, l.Value)

	assert.InDeltaf(t, 20.0, rt.Values[0].Value, 1.0, "")
	assert.InDeltaf(t, 40.0, rt.Values[1].Value, 1.0, "")
	assert.InDeltaf(t, 20.0, rt.Values[2].Value, 1.0, "")

}
