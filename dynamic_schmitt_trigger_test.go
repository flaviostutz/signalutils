package signalutils

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDynamicSchmittTrigger1(t *testing.T) {
	dst, err := NewDynamicSchmittTriggerTimeWindow(5*time.Second, 200, 10, 5.0, 0.5, false)
	assert.Nil(t, err)

	for i := 200; i < 500; i++ {
		v := i + rand.Intn(10)
		dst.SetCurrentValue(float64(v))
		// fmt.Printf("v=%d-upper=%v\n", v, dst.IsUpperRange())
		time.Sleep(10 * time.Millisecond)
	}
	assert.True(t, dst.IsUpperRange())

	for i := 500; i > 200; i-- {
		v := i + rand.Intn(10)
		dst.SetCurrentValue(float64(v))
		// fmt.Printf("v=%d-upper=%v\n", v, dst.IsUpperRange())
		time.Sleep(10 * time.Millisecond)
	}
	assert.False(t, dst.IsUpperRange())

	for i := 200; i < 500; i++ {
		v := i + rand.Intn(10)
		dst.SetCurrentValue(float64(v))
		// fmt.Printf("v=%d-upper=%v\n", v, dst.IsUpperRange())
		time.Sleep(10 * time.Millisecond)
	}
	assert.True(t, dst.IsUpperRange())

	for i := 500; i > 200; i-- {
		v := i + rand.Intn(10)
		dst.SetCurrentValue(float64(v))
		// fmt.Printf("v=%d-upper=%v\n", v, dst.IsUpperRange())
		time.Sleep(10 * time.Millisecond)
	}
	assert.False(t, dst.IsUpperRange())

}
