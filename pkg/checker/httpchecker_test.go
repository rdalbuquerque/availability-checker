package checker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpChecker_Check(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		url            string
		expectedResult bool
	}{
		{
			name:           "unavailable url",
			url:            "http://unavailable-url.net",
			expectedResult: false,
		},
		{
			name:           "available url",
			url:            "https://google.com",
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare the checker
			checker := HttpChecker{
				URL: tc.url,
			}

			// Call the method under test
			result, _ := checker.Check()

			// Assert the result
			assert.Equal(t, tc.expectedResult, result)

		})
	}
}

func TestHttpChecker_Name(t *testing.T) {
	checker := HttpChecker{
		URL: "test1.io",
	}

	name := checker.Name()
	assert.Equal(t, checker.URL, name)
}

func TestHttpChecker_Fix(t *testing.T) {
	checker := HttpChecker{}
	err := checker.Fix()
	assert.Nil(t, err)
}

func TestHttpChecker_IsFixable(t *testing.T) {
	checker := HttpChecker{}
	fixable := checker.IsFixable()
	assert.False(t, fixable)
}
