package signalutils

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMovingAverage1(t *testing.T) {
	ma := NewMovingAverage(3)
	assert.Equal(t, ma.Size, 0)
	ma.AddSample(1000)
	assert.Equal(t, ma.Size, 1)
	ma.AddSample(1000)
	assert.Equal(t, ma.Size, 2)
	ma.AddSample(1000)
	assert.Equal(t, ma.Size, 3)
	ma.AddSample(1000)
	assert.Equal(t, ma.Size, 3)
	ma.AddSample(1000)
	assert.Equal(t, ma.Size, 3)
	ma.AddSample(1000)
	assert.Equal(t, 1000.0, ma.Average())
}

func TestMovingAverage2(t *testing.T) {
	ma := NewMovingAverage(5)
	ma.AddSample(10000)
	ma.AddSample(20000)
	ma.AddSample(1000)
	ma.AddSample(2000)
	ma.AddSample(3000)
	ma.AddSample(4000)
	ma.AddSample(5000)
	assert.Equal(t, 3000.0, ma.Average())
}

func TestMovingAverageTimeWindow1(t *testing.T) {
	ma := NewMovingAverageTimeWindow(1*time.Second, 10)
	ma.AddSample(1000)
	ma.AddSample(2000)
	time.Sleep(1100 * time.Millisecond)
	ma.AddSample(3000)
	ma.AddSample(4000)
	time.Sleep(100 * time.Millisecond)
	ma.AddSample(3000)
	ma.AddSample(2000)
	assert.Equal(t, 3000.0, ma.Average())
}

func TestMovingAverageTimeWindow2(t *testing.T) {
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
}

func TestMovingAverageTimeWindow3(t *testing.T) {
	ma := NewMovingAverageTimeWindow(500*time.Millisecond, 5)
	ma.AddSample(10000)
	time.Sleep(200 * time.Millisecond)

	ma.AddSample(1000)
	time.Sleep(105 * time.Millisecond)
	ma.AddSample(2000)
	time.Sleep(105 * time.Millisecond)
	ma.AddSample(3000)
	time.Sleep(105 * time.Millisecond)
	ma.AddSample(4000)
	time.Sleep(105 * time.Millisecond)
	ma.AddSample(5000)
	assert.Equal(t, 3000.0, ma.Average())

	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, 4000.0, ma.Average())

	time.Sleep(400 * time.Millisecond)
	assert.True(t, math.IsNaN(ma.Average()))

	ma.AddSample(5000)
	assert.Equal(t, 5000.0, ma.Average())
	ma.AddSample(10000)
	assert.Equal(t, 5000.0, ma.Average())
}

func TestMovingAverageNearAverage(t *testing.T) {
	ma := NewMovingAverage(5)
	ma.AddSampleIfNearAverage(1000, 1)
	ma.AddSampleIfNearAverage(1000, 1)
	ma.AddSampleIfNearAverage(1000, 1)
	ma.AddSampleIfNearAverage(1000, 1)
	ma.AddSampleIfNearAverage(1000, 1)
	ma.AddSampleIfNearAverage(1000, 1)
	assert.Equal(t, 1000.0, ma.Average())

	ma.AddSampleIfNearAverage(10000, 1) //SKIP
	assert.Equal(t, 1000.0, ma.Average())

	ma.AddSampleIfNearAverage(10000, 1) //SKIP
	assert.Equal(t, 1000.0, ma.Average())

	ma.AddSampleIfNearAverage(1000, 1)
	ma.AddSampleIfNearAverage(2000, 1)
	ma.AddSampleIfNearAverage(4000, 1.1) //SKIP
	assert.Equal(t, 1000.0, ma.Average())

	ma.AddSampleIfNearAverage(-100, 1) //SKIP
	ma.AddSampleIfNearAverage(2000, 1)
	ma.AddSampleIfNearAverage(1000, 1.1)
	assert.Equal(t, 1400.0, ma.Average())
}
