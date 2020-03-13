package signalutils

import (
	"fmt"
	"math"
	"time"
)

//StateTracker state transition tracker
type StateTracker struct {
	onChange                func(*State, *State)
	onUnchanged             func(*State)
	CurrentState            *State
	CandidateState          string
	CandidateCount          int
	lastUnchanged           time.Time
	changeConfirmations     int
	unchangedTimer          time.Duration
	highestLevel            float64
	resetHighestOnunchanged bool
	active                  bool
}

//State event struct
type State struct {
	Name         string
	Start        time.Time
	Stop         *time.Time
	Data         interface{}
	Level        *float64
	HighestLevel *float64
	HighestTime  *time.Time
	HighestData  interface{}
}

//NewStateTracker new state transition tracker instantiation
//initialState - states are simply strings. a different string denotes a new state
//changeConfirmations - number of sequential state samples with a different state before transitioning
//onChange - listener function that will be called on state transition. ex.: func(newState, previousState) {}. nil value disables this
//unchangedTimer - after this time without changing state, 'onUnchanged' func will be invoked recurrently//. current highest sample will be calculated based on this time slice
//onUnchanged - listener function to be invoked if state is not changed after unchangedStateCount. onUnchanged(state). nil value disables this feature
//highestLevelAccordingToUnchangedTimer - calculate highest level according to whole state duration (false) or only during the onChanged recurrent timer
func NewStateTracker(initialState string, changeConfirmations int, onChange func(*State, *State), unchangedTimer time.Duration, onUnchanged func(*State), resetHighestOnunchanged bool) *StateTracker {
	state := State{
		Name:  initialState,
		Start: time.Now(),
	}
	s1 := StateTracker{
		onChange:                onChange,
		CurrentState:            &state,
		lastUnchanged:           time.Now(),
		CandidateState:          "",
		CandidateCount:          0,
		changeConfirmations:     changeConfirmations,
		unchangedTimer:          unchangedTimer,
		onUnchanged:             onUnchanged,
		highestLevel:            -math.MaxFloat64,
		active:                  true,
		resetHighestOnunchanged: resetHighestOnunchanged,
	}
	go s1.verifyUnchanged()
	return &s1
}

//SetTransientState sets a transient state to tracker so that it can find possible transitions if this state gets recurrent
//returns the time this state started
func (s *StateTracker) SetTransientState(stateName string) (*State, error) {
	return s.SetTransientStateWithData(stateName, 0.0, nil)
}

//SetTransientStateWithData sets a transient state to tracker so that it can find possible transitions if this state gets recurrent
//data is any type that will be sent to listener function
//returns current state count
func (s *StateTracker) SetTransientStateWithData(stateName string, level float64, data interface{}) (*State, error) {
	if !s.active {
		return &State{}, fmt.Errorf("State tracker not active")
	}
	// fmt.Printf("setcurrentstate state=%s\n", state)
	if stateName == s.CurrentState.Name {
		s.CandidateState = ""
		s.CandidateCount = 1
		s.CurrentState.Data = data
		s.CurrentState.Level = &level
		if level > s.highestLevel {
			s.highestLevel = level
			s.CurrentState.HighestLevel = &level
			s.CurrentState.HighestData = data
			now := time.Now()
			s.CurrentState.HighestTime = &now
		}

		return s.CurrentState, nil
	}

	// fmt.Printf("Candidate current=%s state=%s candidate=%s count=%d\n", s.CurrentState, state, s.CandidateState, s.CandidateCount)
	//new candidate state
	if s.CandidateState != stateName {
		s.CandidateState = stateName
		s.CandidateCount = 1
		// fmt.Printf("NEW CANDIDATE CC=%d\n", s.CandidateCount)

		//increment candidate confirmations
	} else {
		s.CandidateCount = s.CandidateCount + 1
		// fmt.Printf("INCREMENTED CC=%d\n", s.CandidateCount)
	}

	//state transition. candidate confirmed
	if s.CandidateCount >= s.changeConfirmations {
		// fmt.Printf("Candidate confirm! candidateCount=%d changeConfirmations=%d state=%s\n", s.CandidateCount, s.changeConfirmations, state)
		prevState := s.CurrentState
		now := time.Now()
		prevState.Stop = &now
		s.CurrentState = &State{
			Name:  stateName,
			Start: time.Now(),
			Data:  data,
		}
		if s.onChange != nil {
			onChange := s.onChange
			onChange(prevState, s.CurrentState)
		}
		s.CandidateState = ""
		s.CandidateCount = 0
		s.highestLevel = -math.MaxFloat64
		s.lastUnchanged = time.Now()
	}

	return s.CurrentState, nil
}

//Close closes the internal timers for notifying unchanged
func (s *StateTracker) Close() {
	s.active = false
}

func (s *StateTracker) verifyUnchanged() {
	for ok := s.active; ok; ok = s.active {
		// fmt.Printf(">>> VERIFY UNCHANGED current=%s\n", s.CurrentState)
		elapsed := time.Duration((time.Now().UnixNano() - s.lastUnchanged.UnixNano()))
		if s.onUnchanged != nil && (elapsed.Nanoseconds() > s.unchangedTimer.Nanoseconds()) {
			// fmt.Printf("NOTIFY %s\n", s.CurrentState)
			s.lastUnchanged = time.Now()
			if s.resetHighestOnunchanged {
				s.highestLevel = -math.MaxFloat64
			}
			onUnchanged := s.onUnchanged
			onUnchanged(s.CurrentState)
		}
		time.Sleep(s.unchangedTimer / 2)
	}
}
