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
