package signalutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	notifiedNewState       *State
	notifiedPreviousState  *State
	notifiedUnchangedState *State
)

func TestStateTracker1(t *testing.T) {
	st := NewStateTracker("state1", 0, onChange, 0, nil, true)
	st.SetTransientState("state2")
	assert.Equal(t, "state2", notifiedNewState.Name)
}

func TestStateTracker2(t *testing.T) {
	st := NewStateTracker("state1", 3, onChange, 0, nil, true)
	st.SetTransientState("state2")
	st.SetTransientState("state2")
	assert.Equal(t, "state1", st.CurrentState.Name)
	st.SetTransientState("state2")
	assert.Equal(t, "state2", notifiedNewState.Name)
	st.SetTransientState("state3")
	assert.Equal(t, "state2", st.CurrentState.Name)
	st.SetTransientState("state3")
	st.SetTransientState("state2")
	st.SetTransientState("state3")
	st.SetTransientState("state3")
	assert.Equal(t, "state2", st.CurrentState.Name)
	st.SetTransientState("state3")
	st.SetTransientState("state3")
	st.SetTransientState("state3")
	st.SetTransientState("state3")
	st.SetTransientState("state3")
	assert.Equal(t, "state3", st.CurrentState.Name)
}

func TestStateTrackerOnUnchanged(t *testing.T) {
	notifiedUnchangedState = nil
	st := NewStateTracker("state1", 3, onChange, 100*time.Millisecond, onUnchanged, true)
	st.SetTransientState("state2")
	assert.Nil(t, notifiedUnchangedState)
	st.SetTransientState("state2")
	assert.Equal(t, "state1", st.CurrentState.Name)
	st.SetTransientState("state2")
	assert.Equal(t, "state2", st.CurrentState.Name)
	assert.Nil(t, notifiedUnchangedState)
	st.SetTransientState("state2")
	time.Sleep(110 * time.Millisecond)
	assert.Equal(t, "state2", notifiedUnchangedState.Name)
	notifiedUnchangedState = nil
	st.SetTransientState("state2")
	assert.Nil(t, notifiedUnchangedState)
	st.SetTransientState("state2")
	time.Sleep(110 * time.Millisecond)
	assert.Equal(t, "state2", notifiedUnchangedState.Name)
}

func TestStateTrackerHighest(t *testing.T) {
	st := NewStateTracker("state1", 3, onChange, 100*time.Millisecond, onUnchanged, true)
	st.SetTransientState("state2")
	st.SetTransientState("state2")
	st.SetTransientStateWithData("state2", 10.0, 10.0)
	st.SetTransientStateWithData("state2", 20.0, 20.0)
	st.SetTransientStateWithData("state2", 5.0, 5.0)
	assert.Equal(t, 20.0, *st.CurrentState.HighestLevel)
	notifiedUnchangedState = nil
	time.Sleep(110 * time.Millisecond)
	st.SetTransientStateWithData("state2", 15.0, 15.0)
	assert.Equal(t, 15.0, *st.CurrentState.HighestLevel)
	st.SetTransientStateWithData("state2", 40.0, 40.0)
	st.SetTransientStateWithData("state2", 41.0, 41.0)
	st.SetTransientStateWithData("state2", 15.0, 15.0)
	assert.Equal(t, 41.0, st.CurrentState.HighestData.(float64))
	notifiedUnchangedState = nil
	time.Sleep(110 * time.Millisecond)
	assert.Equal(t, 41.0, *st.CurrentState.HighestLevel)
	assert.NotNil(t, notifiedUnchangedState)
	notifiedUnchangedState = nil
	st.SetTransientStateWithData("state2", 30.0, 30.0)
	assert.Equal(t, 30.0, *st.CurrentState.HighestLevel)
	assert.Nil(t, notifiedUnchangedState)
}

func TestStateTrackerHighest2(t *testing.T) {
	st := NewStateTracker("state2", 3, onChange, 100*time.Millisecond, onUnchanged, false)
	st.SetTransientState("state2")
	st.SetTransientState("state2")
	st.SetTransientStateWithData("state2", 10.0, 10.0)
	st.SetTransientStateWithData("state2", 20.0, 20.0)
	st.SetTransientStateWithData("state2", 5.0, 5.0)
	assert.Equal(t, 20.0, *st.CurrentState.HighestLevel)
	notifiedUnchangedState = nil
	time.Sleep(110 * time.Millisecond)
	st.SetTransientStateWithData("state2", 15.0, 15.0)
	assert.Equal(t, 20.0, *st.CurrentState.HighestLevel)
	st.SetTransientStateWithData("state2", 40.0, 40.0)
	st.SetTransientStateWithData("state2", 41.0, 41.0)
	st.SetTransientStateWithData("state2", 15.0, 15.0)
	assert.Equal(t, 41.0, st.CurrentState.HighestData.(float64))
	notifiedUnchangedState = nil
	time.Sleep(110 * time.Millisecond)
	assert.Equal(t, 41.0, *st.CurrentState.HighestLevel)
	assert.NotNil(t, notifiedUnchangedState)
	notifiedUnchangedState = nil
	st.SetTransientStateWithData("state2", 30.0, 30.0)
	assert.Equal(t, 41.0, *st.CurrentState.HighestLevel)
	assert.Nil(t, notifiedUnchangedState)
}

func onChange(prevState *State, curState *State) {
	notifiedPreviousState = prevState
	notifiedNewState = curState
}

func onUnchanged(curState *State) {
	notifiedUnchangedState = curState
}
