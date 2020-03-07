package signalutils

import (
	"testing"

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
	assert.Equal(t, ma.Average(), 1000.0)
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
	assert.Equal(t, ma.Average(), 3000.0)
}
