package huxley

import (
	"testing"
	"time"
)

// TODO add edge case tests
func TestGetAge(t *testing.T) {

	// ARRANGE ///////////////////////////////////////////////////////////////////////////////////
	tests := []struct {
		scenario    string
		dob         time.Time
		doi         time.Time
		expectedAge Age
	}{
		{
			scenario: "1 x 31 day month",
			dob:      time.Date(2022, 7, 28, 12, 0, 0, 0, time.UTC),
			doi:      time.Date(2022, 8, 28, 12, 0, 0, 0, time.UTC),
			expectedAge: Age{
				Days:   31,
				Weeks:  4,
				Months: 1,
			},
		},
		{
			scenario: "2 x 31 day month",
			dob:      time.Date(2022, 7, 28, 12, 0, 0, 0, time.UTC),
			doi:      time.Date(2022, 9, 28, 12, 0, 0, 0, time.UTC),
			expectedAge: Age{
				Days:   62,
				Weeks:  9,
				Months: 2,
			},
		},
	}

	// ACT and ASSERT  ///////////////////////////////////////////////////////////////////////////////////
	for _, test := range tests {
		if calcAge := getAge(test.dob, test.doi); !areSameAge(calcAge, test.expectedAge) {
			t.Errorf("\n scenario: %v\n got: days %v, weeks %v, months %v\n expected: days %v, weeks %v, months %v",
				test.scenario,
				calcAge.Days, calcAge.Weeks, calcAge.Months,
				test.expectedAge.Days, test.expectedAge.Weeks, test.expectedAge.Months)
		}
	}

}

func areSameAge(calcAge, expectedAge Age) bool {
	if calcAge.Days == expectedAge.Days &&
		calcAge.Months == expectedAge.Months &&
		calcAge.Weeks == expectedAge.Weeks {
		return true
	}
	return false
}
