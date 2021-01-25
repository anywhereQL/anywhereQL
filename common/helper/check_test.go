package helper

import "testing"

func TestCheckDigit(t *testing.T) {
	testCaass := []struct {
		input    rune
		excepted bool
	}{
		{rune('1'), true},
		{rune('2'), true},
		{rune('3'), true},
		{rune('4'), true},
		{rune('5'), true},
		{rune('6'), true},
		{rune('7'), true},
		{rune('8'), true},
		{rune('9'), true},
		{rune('0'), true},
		{rune('a'), false},
		{rune('私'), false},
	}

	for tn, tc := range testCaass {
		r := IsDigit(tc.input)
		if r != tc.excepted {
			t.Errorf("[%d] Check %c: Expected = %v, but got %v\n", tn, tc.input, tc.excepted, r)
		}
	}
}
