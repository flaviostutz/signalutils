package signalutils

//StateTracker state transition tracker
type StateTracker struct {
	onChange            func(string, string, interface{})
	CurrentState        string
	CurrentStateCounter int
	CandidateState      string
	CandidateCount      int
	changeConfirmations int
	unchangedStateCount int
	onUnchanged         func(string, int, interface{})
}

//NewStateTracker new state transition tracker instantiation
//initialState - states are simply strings. a different string denotes a new state
//changeConfirmations - number of sequential state samples with a different state before transitioning
//onChange - listener function that will be called on state transition. ex.: func(newState, previousState) {}. nil value disables this
//unchangedStateCount - after this number of state samples without changing state, 'onUnchanged' func will be invoked recurrently
//onUnchanged - listener function to be invoked if state is not changed after unchangedStateCount. onUnchanged(state, stateCounter, data). nil value disables this feature
func NewStateTracker(initialState string, changeConfirmations int, onChange func(string, string, interface{}), unchangedStateCount int, onUnchanged func(string, int, interface{})) StateTracker {
	return StateTracker{
		onChange:            onChange,
		CurrentState:        initialState,
		CurrentStateCounter: 1,
		CandidateState:      "",
		CandidateCount:      0,
		changeConfirmations: changeConfirmations,
		unchangedStateCount: unchangedStateCount,
		onUnchanged:         onUnchanged,
	}
}

//SetTransientState sets a transient state to tracker so that it can find possible transitions if this state gets recurrent
func (s *StateTracker) SetTransientState(state string) int {
	return s.SetTransientStateWithData(state, nil)
}

//SetTransientStateWithData sets a transient state to tracker so that it can find possible transitions if this state gets recurrent
//data is any type that will be sent to listener function
//returns current state count
func (s *StateTracker) SetTransientStateWithData(state string, data interface{}) int {
	// fmt.Printf("setcurrentstate state=%s\n", state)
	if state == s.CurrentState {
		s.CurrentStateCounter = s.CurrentStateCounter + 1
		s.CandidateState = ""
		s.CandidateCount = 1
		if s.onUnchanged != nil && s.CurrentStateCounter%s.unchangedStateCount == 0 {
			s.onUnchanged(s.CurrentState, s.CurrentStateCounter, data)
		}
		return s.CurrentStateCounter
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
		s.CurrentStateCounter = 1
		s.CandidateState = ""
		s.CandidateCount = 0
	}

	return s.CurrentStateCounter
}
