package fieldmask

import "testing"

type TestMask uint8

const (
	One = 1 << iota
	Two
	Three
	Four
	Five
	Six
	Seven
)

func TestEnumerateMaskedFields(t *testing.T) {
	var things []int
	EnumerateFields(Two|Six|Five, &things)
}
