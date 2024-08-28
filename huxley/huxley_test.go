package huxley

import (
	"testing"
	"time"
)

func TestGetData(t *testing.T) {
	// ARRANGE ///////////////////////////////////////////////////////////////////////////////////
	tests := []struct {
		scenario    string
		dob         time.Time
		doi         time.Time
		expectedAge Age
	}{
		{
			scenario: "1.25 month",
			dob:      time.Date(2022, 7, 9, 12, 0, 0, 0, time.UTC),
			doi:      time.Date(2022, 8, 16, 12, 0, 0, 0, time.UTC),
			expectedAge: Age{
				Days:   38,
				Weeks:  5,
				Months: 1.25,
				Years:  0,
			},
		},
		{
			scenario: "1.5 month",
			dob:      time.Date(2022, 7, 9, 12, 0, 0, 0, time.UTC),
			doi:      time.Date(2022, 8, 23, 12, 0, 0, 0, time.UTC),
			expectedAge: Age{
				Days:   45,
				Weeks:  6,
				Months: 1.5,
				Years:  0,
			},
		},
		{
			scenario: "1.75 month",
			dob:      time.Date(2022, 7, 9, 12, 0, 0, 0, time.UTC),
			doi:      time.Date(2022, 8, 30, 12, 0, 0, 0, time.UTC),
			expectedAge: Age{
				Days:   52,
				Weeks:  7,
				Months: 1.75,
				Years:  0,
			},
		},
		{
			scenario: "2 months",
			dob:      time.Date(2022, 7, 9, 12, 0, 0, 0, time.UTC),
			doi:      time.Date(2022, 9, 9, 12, 0, 0, 0, time.UTC),
			expectedAge: Age{
				Days:   62,
				Weeks:  9,
				Months: 2,
				Years:  0,
			},
		},
		{
			scenario: "3 months",
			dob:      time.Date(2022, 7, 9, 12, 0, 0, 0, time.UTC),
			doi:      time.Date(2022, 10, 9, 12, 0, 0, 0, time.UTC),
			expectedAge: Age{
				Days:   92,
				Weeks:  13,
				Months: 3,
				Years:  0.25,
			},
		},
		{
			scenario: "1.0 leap year",
			dob:      time.Date(2023, 7, 9, 12, 0, 0, 0, time.UTC),
			doi:      time.Date(2024, 7, 9, 12, 0, 0, 0, time.UTC),
			expectedAge: Age{
				Days:   366,
				Weeks:  52,
				Months: 12,
				Years:  1,
			},
		},
		{
			scenario: "end of year",
			dob:      time.Date(2023, 12, 9, 12, 0, 0, 0, time.UTC),
			doi:      time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
			expectedAge: Age{
				Days:   30,
				Weeks:  4,
				Months: 1,
				Years:  0,
			},
		},
		{
			scenario: "1.25 years",
			dob:      time.Date(2023, 7, 9, 12, 0, 0, 0, time.UTC),
			doi:      time.Date(2024, 10, 9, 12, 0, 0, 0, time.UTC),
			expectedAge: Age{
				Days:   458,
				Weeks:  65,
				Months: 15,
				Years:  1.25,
			},
		},
		{
			scenario: "1.5 years",
			dob:      time.Date(2023, 7, 9, 12, 0, 0, 0, time.UTC),
			doi:      time.Date(2025, 1, 9, 12, 0, 0, 0, time.UTC),
			expectedAge: Age{
				Days:   550,
				Weeks:  79,
				Months: 18,
				Years:  1.5,
			},
		},
		{
			scenario: "1.75 years",
			dob:      time.Date(2023, 7, 9, 12, 0, 0, 0, time.UTC),
			doi:      time.Date(2025, 4, 9, 12, 0, 0, 0, time.UTC),
			expectedAge: Age{
				Days:   640,
				Weeks:  91,
				Months: 21,
				Years:  1.75,
			},
		},
		{
			scenario: "1 x 31 day month",
			dob:      time.Date(2022, 7, 28, 12, 0, 0, 0, time.UTC),
			doi:      time.Date(2022, 8, 28, 12, 0, 0, 0, time.UTC),
			expectedAge: Age{
				Days:   31,
				Weeks:  4,
				Months: 1,
				Years:  0,
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
				Years:  0,
			},
		},
	}

	// ACT and ASSERT  ///////////////////////////////////////////////////////////////////////////////////
	for _, test := range tests {
		if calcAge := getAge(test.dob, test.doi); !areSameAge(calcAge, test.expectedAge) {
			t.Errorf("\n scenario: %v\n got: days %v, weeks %v, months %v, years %v\n expected: days %v, weeks %v, months %v, years %v",
				test.scenario,
				calcAge.Days, calcAge.Weeks, calcAge.Months, calcAge.Years,
				test.expectedAge.Days, test.expectedAge.Weeks, test.expectedAge.Months, test.expectedAge.Years)
		}
	}
}

func areSameAge(calcAge, expectedAge Age) bool {
	return calcAge.Years == expectedAge.Years &&
		calcAge.Months == expectedAge.Months &&
		calcAge.Weeks == expectedAge.Weeks &&
		calcAge.Days == expectedAge.Days
}
