package signalutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	notifiedNewState       string
	notifiedPreviousState  string
	notifiedUnchangedState string
)

func TestStateTracker1(t *testing.T) {
	st := NewStateTracker("state1", 0, onChange, 0, nil)
	st.SetTransientState("state2")
	assert.Equal(t, "state2", notifiedNewState)
}

func TestStateTracker2(t *testing.T) {
	st := NewStateTracker("state1", 3, onChange, 0, nil)
	st.SetTransientState("state2")
	st.SetTransientState("state2")
	assert.Equal(t, "state1", st.CurrentState)
	st.SetTransientState("state2")
	assert.Equal(t, "state2", notifiedNewState)
	st.SetTransientState("state3")
	assert.Equal(t, "state2", st.CurrentState)
	st.SetTransientState("state3")
	st.SetTransientState("state2")
	st.SetTransientState("state3")
	st.SetTransientState("state3")
	assert.Equal(t, "state2", st.CurrentState)
	st.SetTransientState("state3")
	st.SetTransientState("state3")
	st.SetTransientState("state3")
	st.SetTransientState("state3")
	st.SetTransientState("state3")
	assert.Equal(t, "state3", st.CurrentState)
}

func TestStateTrackerOnUnchanged(t *testing.T) {
	notifiedUnchangedState = ""
	st := NewStateTracker("state1", 3, onChange, 100*time.Millisecond, onUnchanged)
	st.SetTransientState("state2")
	assert.Equal(t, "", notifiedUnchangedState)
	st.SetTransientState("state2")
	assert.Equal(t, "state1", st.CurrentState)
	st.SetTransientState("state2")
	assert.Equal(t, "state2", st.CurrentState)
	assert.Equal(t, "", notifiedUnchangedState)
	st.SetTransientState("state2")
	time.Sleep(110 * time.Millisecond)
	assert.Equal(t, "state2", notifiedUnchangedState)
	notifiedUnchangedState = ""
	st.SetTransientState("state2")
	assert.Equal(t, "", notifiedUnchangedState)
	st.SetTransientState("state2")
	time.Sleep(110 * time.Millisecond)
	assert.Equal(t, "state2", notifiedUnchangedState)
}

func onChange(news string, previous string, data interface{}) {
	notifiedNewState = news
	notifiedPreviousState = previous
}

func onUnchanged(state string, stateDuration time.Duration, data interface{}) {
	notifiedUnchangedState = state
}
