package signalutils

import (
	"fmt"
	"math"
)

//SchmittTrigger utility
type SchmittTrigger struct {
	LowerLimit float64
	UpperLimit float64
	UpperRange bool
}

//NewSchmittTrigger new schmitt trigger creation
func NewSchmittTrigger(lowerLimit float64, upperLimit float64, upperRange bool) (SchmittTrigger, error) {
	if upperLimit <= lowerLimit {
		return SchmittTrigger{}, fmt.Errorf("upperLimit cannot be less than lowerLimit")
	}
	return SchmittTrigger{
		LowerLimit: lowerLimit,
		UpperLimit: upperLimit,
		UpperRange: upperRange,
	}, nil
}

//SetCurrentValue set current value and calculate if it is in upper or lower range
func (s *SchmittTrigger) SetCurrentValue(value float64) {
	if s.UpperRange {
		if value < s.LowerLimit {
			s.UpperRange = false
		}
	} else {
		if value > s.UpperLimit {
			s.UpperRange = true
		}
	}
}

//IsUpperRange returns whatever it is in upper range or not
func (s *SchmittTrigger) IsUpperRange() bool {
	return s.UpperRange
}

//UpdateLowerUpperLimits changes current lower/upper limits for schmitt trigger
func (s *SchmittTrigger) UpdateLowerUpperLimits(lowerLimit float64, upperLimit float64) {
	// fmt.Printf("updateLimits %f-%f\n", lowerLimit, upperLimit)
	if !math.IsNaN(lowerLimit) {
		s.LowerLimit = lowerLimit
	}
	if !math.IsNaN(upperLimit) {
		s.UpperLimit = upperLimit
	}
}
