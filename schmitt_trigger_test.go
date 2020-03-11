package signalutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchmittTrigger1(t *testing.T) {
	st, _ := NewSchmittTrigger(10, 20, false)
	assert.False(t, st.IsUpperRange())
	st.SetCurrentValue(11)
	assert.False(t, st.IsUpperRange())
	st.SetCurrentValue(15)
	assert.False(t, st.IsUpperRange())

	st.SetCurrentValue(21)
	assert.True(t, st.IsUpperRange())
	st.SetCurrentValue(25)
	assert.True(t, st.IsUpperRange())
	st.SetCurrentValue(18)
	assert.True(t, st.IsUpperRange())
	st.SetCurrentValue(12)
	assert.True(t, st.IsUpperRange())

	st.SetCurrentValue(9)
	assert.False(t, st.IsUpperRange())

	st.SetCurrentValue(26)
	assert.True(t, st.IsUpperRange())

	st.SetCurrentValue(-333)
	assert.False(t, st.IsUpperRange())
}
