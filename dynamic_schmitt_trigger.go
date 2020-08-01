package signalutils

import (
	"sync"
	"time"
)

//DynamicSchmittTrigger adjusts its internal lower and upper limits according to a moving average on observed values
//Only initialize this with NewDynamicSchmittTriggerTimeWindow(..)
type DynamicSchmittTrigger struct {
	minMaxMovAvg                   MovingAverage
	schmittTrigger                 SchmittTrigger
	ignoreSamplesTooDifferentRatio float64
	groupByMinMaxSamples           int
	minMaxUpperLowerRatio          float64
	m                              *sync.Mutex
}

//NewDynamicSchmittTriggerTimeWindow new schmitt trigger creation
//minMaxMovingAverageSamples defines the time window of the moving average that defines min/max values for schmitt trigger
//maxMovingAverageSamples - max number of samples in minmax averager. if a too high rate of samples are set, some may be ignored
//groupByMinMaxSamples - number of signal samples to use to calculate each min/max sampling
//ignoreSamplesTooHighRatio - if SetCurrentValue sets a value that is too high or too low according to min/max moving average, ignore it
//minMaxUpperLowerRatio - 1.0 indicates the lower and upper limits will be placed just like the min/max moving average, which is not too practical. A number between 0.3 and 0.7 is good here.
func NewDynamicSchmittTriggerTimeWindow(minMaxMovingAverageTime time.Duration, maxMovingAverageSamples int, groupByMinMaxSamples int, ignoreSamplesTooDifferentRatio float64, minMaxUpperLowerRatio float64, upperRange bool) (DynamicSchmittTrigger, error) {
	minMaxMovAvg := NewMovingAverageTimeWindow(minMaxMovingAverageTime, maxMovingAverageSamples)
	schmittTrigger, _ := NewSchmittTrigger(0, 0.1, upperRange)
	return DynamicSchmittTrigger{
		minMaxMovAvg:                   minMaxMovAvg,
		schmittTrigger:                 schmittTrigger,
		ignoreSamplesTooDifferentRatio: ignoreSamplesTooDifferentRatio,
		groupByMinMaxSamples:           groupByMinMaxSamples,
		minMaxUpperLowerRatio:          minMaxUpperLowerRatio,
		m:                              &sync.Mutex{},
	}, nil
}

//SetCurrentValue set current value and calculate if it is in upper or lower range
//returns 1-true or false if value was accepted by internal moving averager (rate not too high)
//        2-how much the current value is distant from the lower limit (if it is in 'upperRange' state) or distant from the upper limit (if in 'lowerRange' state) for a new change to occur in trigger. a ratio in relation to max-min range will be returned
func (s *DynamicSchmittTrigger) SetCurrentValue(value float64) (bool, float64) {
	s.m.Lock()
	defer s.m.Unlock()

	b := s.minMaxMovAvg.AddSampleIfNearAverage(value, s.ignoreSamplesTooDifferentRatio)
	min, max := s.minMaxMovAvg.AverageMinMax(s.groupByMinMaxSamples)
	cw := max - min/2
	min2 := min + (cw-min)*(1-s.minMaxUpperLowerRatio)
	max2 := max - (max-cw)*(1-s.minMaxUpperLowerRatio)
	s.schmittTrigger.UpdateLowerUpperLimits(min2, max2)
	s.schmittTrigger.SetCurrentValue(value)

	if s.schmittTrigger.IsUpperRange() {
		return b, (value - s.schmittTrigger.LowerLimit) // / (max - min)
	}
	return b, (value - s.schmittTrigger.UpperLimit) // / (max - min)
}

//IsUpperRange returns if this trigger is in upper or low range
func (s *DynamicSchmittTrigger) IsUpperRange() bool {
	return s.schmittTrigger.IsUpperRange()
}

//GetLowerUpperLimits get current lower and upper limits for this schmtt trigger
func (s *DynamicSchmittTrigger) GetLowerUpperLimits() (float64, float64) {
	return s.schmittTrigger.LowerLimit, s.schmittTrigger.UpperLimit
}
