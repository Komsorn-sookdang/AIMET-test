package usecases_test

import (
	"aimet-test/internal/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatTimeTo12Hour(t *testing.T) {
	type testCase struct {
		name string
		time24Str string
		expected string
	}

	testCases := []testCase{
		{
			name: "Test AM 1",
			time24Str: "09:10",
			expected: "9:10 AM",
		},
		{
			name: "Test AM 2",
			time24Str: "00:10",
			expected: "12:10 AM",
		},
		{
			name: "Test AM 3",
		time24Str: "11:00",
			expected: "11 AM",
		},
		{
			name: "Test PM 1",
			time24Str: "13:10",
			expected: "1:10 PM",
		},
		{
			name: "Test PM 2",
			time24Str: "12:10",
			expected: "12:10 PM",
		},
		{
			name: "Test PM 3",
			time24Str: "23:00",
			expected: "11 PM",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := usecases.FormatTimeTo12Hour(tc.time24Str)

			assert.Equal(t, tc.expected, actual)
		})
	}
}	
