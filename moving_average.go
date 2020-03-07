package signalutils

//MovingAverage running moving averager
type MovingAverage struct {
	Size            int
	Samples         []float64
	lastResultValid bool
	lastResult      float64
}

//NewMovingAverage creates a new moving averager with a fixed size
func NewMovingAverage(size int) MovingAverage {
	return MovingAverage{
		Samples:         make([]float64, size),
		lastResultValid: false,
	}
}

//Average computes average with current samples in fixed length list
func (m *MovingAverage) Average() float64 {
	if m.Size == 0 {
		return 0
	}
	if !m.lastResultValid {
		sum := 0.0
		for i := 0; i < m.Size; i++ {
			sum = sum + m.Samples[i]
		}
		m.lastResult = sum / float64(m.Size)
		m.lastResultValid = true
	}
	return m.lastResult
}

//AddSample adds a new sample to the moving average. If there is more than 'size' samples, the oldest sample will be removed
func (m *MovingAverage) AddSample(value float64) {
	if m.Size < len(m.Samples) {
		m.Size = m.Size + 1
	} else {
		//put new sample in tail
		for i := 0; i < len(m.Samples)-1; i++ {
			m.Samples[i] = m.Samples[i+1]
		}
	}
	m.Samples[m.Size-1] = value
	m.lastResultValid = false
}

//Reset internal samples
func (m *MovingAverage) Reset() {
	m.Samples = make([]float64, len(m.Samples))
	m.Size = 0
}
