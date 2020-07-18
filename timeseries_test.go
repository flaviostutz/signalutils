package signalutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTSMaxSize(t *testing.T) {
	ts := NewTimeseries(500 * time.Millisecond)

	ts.Add(1)
	ts.Add(1)
	ts.Add(1)
	ts.Add(1)
	ts.Add(1)
	ts.Add(1)
	ts.Add(1)
	time.Sleep(300 * time.Millisecond)

	ts.Add(1)
	ts.Add(1)
	ts.Add(1)
	ts.Add(1)
	time.Sleep(650 * time.Millisecond)

	assert.Equal(t, 11, ts.Size())

	time.Sleep(300 * time.Millisecond)
	ts.Add(1)
	ts.Add(1)

	time.Sleep(600 * time.Millisecond)
	ts.Add(1)
	ts.Add(1)
	ts.Add(1)

	time.Sleep(300 * time.Millisecond)
	ts.Add(1)
	ts.Add(1)
	ts.Add(1)
	ts.Add(1)
	ts.Add(1)

	time.Sleep(300 * time.Millisecond)
	ts.Add(1)
	ts.Add(1)

	time.Sleep(300 * time.Millisecond)
	ts.Add(1)
	ts.Add(1)
	ts.Add(1)

	time.Sleep(1000 * time.Millisecond)
	ts.Add(1)
	ts.Add(1)
	ts.Add(1)

	assert.Equal(t, 10, ts.Size())
}

func TestTSGetInterpolated1(t *testing.T) {
	ts := NewTimeseries(5000 * time.Millisecond)

	ts.Add(10)
	time.Sleep(100 * time.Millisecond)
	ts.Add(20)
	time.Sleep(30 * time.Millisecond)
	ts.Add(10)
	time.Sleep(20 * time.Millisecond)
	ts.Add(20)
	time.Sleep(40 * time.Millisecond)
	ts.Add(10)
	time.Sleep(30 * time.Millisecond)

	ts.Add(10)
	time.Sleep(100 * time.Millisecond)
	ts.Add(30)

	nv, ok := ts.Get(time.Now().Add(-50 * time.Millisecond))
	assert.True(t, ok)
	assert.InDeltaf(t, float64(20), nv.Value, float64(5), "")
}

func TestTSGetInterpolated2(t *testing.T) {
	ts := NewTimeseries(5000 * time.Millisecond)

	ts.Add(100)
	time.Sleep(500 * time.Millisecond)
	ts.Add(30)

	nv, ok := ts.Get(time.Now().Add(-200 * time.Millisecond))
	assert.True(t, ok)
	assert.InDeltaf(t, float64(58), nv.Value, float64(5), "")
}

func TestTSGetInterpolated3(t *testing.T) {
	ts := NewTimeseries(1000 * time.Millisecond)

	ts.Add(-100)
	time.Sleep(500 * time.Millisecond)
	ts.Add(100)

	nv, ok := ts.Get(time.Now().Add(-250 * time.Millisecond))
	assert.True(t, ok)
	assert.InDeltaf(t, float64(0), nv.Value, float64(5), "")
}

func TestTSGetInterpolated4(t *testing.T) {
	ts := NewTimeseries(1000 * time.Millisecond)

	ts.Add(-100)
	time.Sleep(500 * time.Millisecond)
	ts.Add(-1000)

	nv, ok := ts.Get(time.Now().Add(-250 * time.Millisecond))
	assert.True(t, ok)
	assert.InDeltaf(t, float64(-555), nv.Value, float64(50), "")
}

func TestTSReset(t *testing.T) {
	ts := NewTimeseries(1000 * time.Millisecond)

	ts.Add(-100)
	time.Sleep(500 * time.Millisecond)
	ts.Add(-1000)

	assert.Equal(t, 2, ts.Size())

	ts.Reset()
	_, ok := ts.Get(time.Now().Add(-250 * time.Millisecond))
	assert.False(t, ok)

	assert.Equal(t, 0, ts.Size())
}

func TestTSLastValue(t *testing.T) {
	ts := NewTimeseries(1000 * time.Millisecond)

	_, ok := ts.Last()
	assert.False(t, ok)

	ts.Add(-100)
	time.Sleep(500 * time.Millisecond)
	ts.Add(-1000)

	v, ok := ts.Last()
	assert.True(t, ok)
	assert.Equal(t, float64(-1000), v.Value)
}
