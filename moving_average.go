package signalutils

import (
	"math"
	"time"
)

//MovingAverage running moving averager
type MovingAverage struct {
	Size                      int
	Samples                   []float64
	lastResultValid           bool
	lastResult                float64
	samplesTimeUnixNano       []int64
	samplesDurationNano       int64
	lastSampleTimeUnixNano    int64
	minTimeNanoBetweenSamples int64
	lastResultTime            int64
}

//NewMovingAverage creates a new moving averager with a fixed size
func NewMovingAverage(size int) MovingAverage {
	return MovingAverage{
		Samples:             make([]float64, size),
		lastResultValid:     false,
		samplesDurationNano: -1,
	}
}

//NewMovingAverageTimeWindow creates a new moving averager that will average samples no older than 'samplesDuration', limiting the number of samples to 'maxSamples' in time window. If two consecutive samples are added to the averager in a period less than duration/maxSamples, it will be ignored.
func NewMovingAverageTimeWindow(samplesDuration time.Duration, maxSamples int) MovingAverage {
	minTime := samplesDuration.Nanoseconds() / int64(maxSamples)
	return MovingAverage{
		Samples:                   make([]float64, maxSamples),
		samplesDurationNano:       samplesDuration.Nanoseconds(),
		samplesTimeUnixNano:       make([]int64, maxSamples),
		lastResultValid:           false,
		minTimeNanoBetweenSamples: minTime,
		lastSampleTimeUnixNano:    0,
	}
}

//AddSample adds a new sample to the moving average. If there is more than 'size' samples, the oldest sample will be removed. If this is a timed window averager and the last sample was added in less than sampleDurate/maxSamples time, it will be ignored.
func (m *MovingAverage) AddSample(value float64) bool {
	if m.samplesDurationNano != -1 {
		if (time.Now().UnixNano() - m.lastSampleTimeUnixNano) < m.minTimeNanoBetweenSamples {
			return false
		}
	}

	m.lastSampleTimeUnixNano = time.Now().UnixNano()

	if m.Size < len(m.Samples) {
		m.Size = m.Size + 1
	} else {
		//put new sample in tail
		for i := 0; i < len(m.Samples)-1; i++ {
			m.Samples[i] = m.Samples[i+1]
			if m.samplesDurationNano != -1 {
				m.samplesTimeUnixNano[i] = m.samplesTimeUnixNano[i+1]
			}
		}
	}
	m.Samples[m.Size-1] = value

	//time window
	if m.samplesDurationNano != -1 {
		m.samplesTimeUnixNano[m.Size-1] = time.Now().UnixNano()
	}
	m.lastResultValid = false
	return true
}

//AddSampleIfNearAverage Add sample only if its value is near current average to avoid espurious samples to be added to the average.
//avgDiff 1 means samples between [-currentAvg, +currentAvg] will be accepted.
//Returns true if sample was accepted
func (m *MovingAverage) AddSampleIfNearAverage(value float64, avgDiff float64) bool {
	avg := m.Average()
	if math.IsNaN(avg) || (math.Abs(avg-value) <= (avg * avgDiff)) {
		return m.AddSample(value)
	}
	return false
}

//Average computes average with current samples in fixed length list
func (m *MovingAverage) Average() float64 {
	if m.Size == 0 {
		return math.NaN()
	}

	//invalidate cache if using timed window
	if m.samplesDurationNano != -1 && m.lastResultValid {
		if (time.Now().UnixNano() - m.lastResultTime) > m.minTimeNanoBetweenSamples {
			m.lastResultValid = false
		}
	}

	// fmt.Printf("CACHE %v\n", m.lastResultValid)
	n := 0
	if !m.lastResultValid {
		sum := 0.0
		for i := 0; i < m.Size; i++ {
			if m.samplesDurationNano != -1 {
				//skip this sample if too old
				if (time.Now().UnixNano() - m.samplesTimeUnixNano[i]) > m.samplesDurationNano {
					// fmt.Printf("SKIP OLD %f i=%d\n", m.Samples[i], i)
					continue
				}
			}
			sum = sum + m.Samples[i]
			n = n + 1
		}

		m.lastResult = math.NaN()
		if n > 0 {
			m.lastResult = sum / float64(n)
			// fmt.Printf("%f/%f=%f\n", sum, float64(n), m.lastResult)
			m.lastResultValid = true
			m.lastResultTime = time.Now().UnixNano()
		}
	}
	return m.lastResult
}

//Reset internal samples
func (m *MovingAverage) Reset() {
	m.Samples = make([]float64, len(m.Samples))
	m.Size = 0
	m.lastResultValid = false
}
