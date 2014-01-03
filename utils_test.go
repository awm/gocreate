package gocreate

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestSequence(t *testing.T) {
    values1 := Sequence(0, 10, 1).([]int)
    assert.Equal(t, values1, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

    values2 := Sequence(byte(5), byte(10), byte(1)).([]byte)
    assert.Equal(t, values2, []byte{5, 6, 7, 8, 9})

    values3 := Sequence(uintptr(0), uintptr(10), uintptr(2)).([]uintptr)
    assert.Equal(t, values3, []uintptr{0, 2, 4, 6, 8})

    values4 := Sequence(int16(0), int16(-10), int16(-2)).([]int16)
    assert.Equal(t, values4, []int16{0, -2, -4, -6, -8})

    values5 := Sequence(0, byte(10), 2)
    assert.Nil(t, values5)

    values6 := Sequence(0, 10, -2).([]int)
    assert.Equal(t, values6, []int{})
}
