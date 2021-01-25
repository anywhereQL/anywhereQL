package value

import (
	"reflect"
	"testing"
)

func TestConvert(t *testing.T) {
	testCases := []struct {
		input    string
		isError  bool
		expected Value
	}{}

	for tn, tc := range testCases {
		r, e := Convert(tc.input)
		if tc.isError && e == nil {
			t.Errorf("[%d] Expected Error, but no error", tn)
		}
		if !tc.isError && e != nil {
			t.Errorf("[%d] Expected No Error, but error", tn)
		}
		if !tc.isError && !reflect.DeepEqual(tc.expected, r) {
			t.Errorf("[%d] Expected Value Struct not equal return Value struct", tn)
		}
	}
}
