package signalutils

import (
	"fmt"
	"time"
)

//StateTracker state transition tracker
type StateTracker struct {
	onChange               func(string, string, interface{})
	CurrentState           string
	CurrentStateData       interface{}
	CurrentStateStart      time.Time
	CandidateState         string
	CandidateCount         int
	lastUnchanged          time.Time
	changeConfirmations    int
	unchangedStateDuration time.Duration
	onUnchanged            func(string, time.Duration, interface{})
	active                 bool
}

//NewStateTracker new state transition tracker instantiation
//initialState - states are simply strings. a different string denotes a new state
//changeConfirmations - number of sequential state samples with a different state before transitioning
//onChange - listener function that will be called on state transition. ex.: func(newState, previousState) {}. nil value disables this
//unchangedStateCount - after this number of state samples without changing state, 'onUnchanged' func will be invoked recurrently
//onUnchanged - listener function to be invoked if state is not changed after unchangedStateCount. onUnchanged(state, stateCounter, data). nil value disables this feature
func NewStateTracker(initialState string, changeConfirmations int, onChange func(string, string, interface{}), unchangedStateDuration time.Duration, onUnchanged func(string, time.Duration, interface{})) *StateTracker {
	s1 := StateTracker{
		onChange:               onChange,
		CurrentState:           initialState,
		lastUnchanged:          time.Now(),
		CandidateState:         "",
		CandidateCount:         0,
		changeConfirmations:    changeConfirmations,
		unchangedStateDuration: unchangedStateDuration,
		onUnchanged:            onUnchanged,
		active:                 true,
	}
	go s1.verifyUnchanged()
	return &s1
}

//SetTransientState sets a transient state to tracker so that it can find possible transitions if this state gets recurrent
//returns the time this state started
func (s *StateTracker) SetTransientState(state string) (time.Time, error) {
	return s.SetTransientStateWithData(state, nil)
}

//SetTransientStateWithData sets a transient state to tracker so that it can find possible transitions if this state gets recurrent
//data is any type that will be sent to listener function
//returns current state count
func (s *StateTracker) SetTransientStateWithData(state string, data interface{}) (time.Time, error) {
	if !s.active {
		return time.Time{}, fmt.Errorf("State tracker not active")
	}
	// fmt.Printf("setcurrentstate state=%s\n", state)
	if state == s.CurrentState {
		s.CandidateState = ""
		s.CandidateCount = 1
		s.CurrentStateData = data
		return s.CurrentStateStart, nil
	}

	// fmt.Printf("Candidate current=%s state=%s candidate=%s count=%d\n", s.CurrentState, state, s.CandidateState, s.CandidateCount)
	//new candidate state
	if s.CandidateState != state {
		s.CandidateState = state
		s.CandidateCount = 1
		// fmt.Printf("NEW CANDIDATE CC=%d\n", s.CandidateCount)

		//increment candidate confirmations
	} else {
		s.CandidateCount = s.CandidateCount + 1
		// fmt.Printf("INCREMENTED CC=%d\n", s.CandidateCount)
	}

	//candidate confirmed
	if s.CandidateCount >= s.changeConfirmations {
		// fmt.Printf("Candidate confirm! candidateCount=%d changeConfirmations=%d state=%s\n", s.CandidateCount, s.changeConfirmations, state)
		if s.onChange != nil {
			s.onChange(state, s.CurrentState, data)
		}
		s.CurrentState = state
		s.CurrentStateStart = time.Now()
		s.CandidateState = ""
		s.CandidateCount = 0
		s.lastUnchanged = time.Now()
		// fmt.Printf("CURRENT STATE CONFIRMED %s\n", s.CurrentState)
	}

	return s.CurrentStateStart, nil
}

//Close closes the internal timers for notifying unchanged
func (s *StateTracker) Close() {
	s.active = false
}

func (s *StateTracker) verifyUnchanged() {
	for ok := s.active; ok; ok = s.active {
		// fmt.Printf(">>> VERIFY UNCHANGED current=%s\n", s.CurrentState)
		elapsed := time.Duration((time.Now().UnixNano() - s.lastUnchanged.UnixNano()))
		if s.onUnchanged != nil && (elapsed.Nanoseconds() > s.unchangedStateDuration.Nanoseconds()) {
			// fmt.Printf("NOTIFY %s\n", s.CurrentState)
			s.lastUnchanged = time.Now()
			s.onUnchanged(s.CurrentState, elapsed, s.CurrentStateData)
		}
		time.Sleep(s.unchangedStateDuration / 2)
	}
}
