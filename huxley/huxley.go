// calculate and display Huxley's age in various units
package huxley

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"time"
)

var templ = template.Must(template.New("webpage").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Name}}</title>
		<style>
			h1 {color:blue;}
		</style>
	</head>
	<body style="font-size: xxx-large;">
		<h1>{{.Name}}'s Age</h1>
		<h3>{{.DateOfInterest}}</h3>
		<ul>
			<li>Breed: {{.Breed}}</li>
			<li>Born: {{.DateOfBirth}}</li>
			<li>Years: {{.AgeYears}}</li>
			<li>Months: {{.AgeMonths}}</li>
			<li>Weeks: {{.AgeWeeks}}</li>
			<li>Days: {{.AgeDays}}</li>
		</ul>
	</body>
</html>
`))

type DogStat struct {
	Name           string
	DateOfBirth    string
	Breed          string
	DateOfInterest string
	AgeDays        int
	AgeWeeks       int
	AgeMonths      float64
	AgeYears       float64
}

type Age struct {
	DateOfInterest string
	Days           int
	Weeks          int
	Months         float64
	Years          float64
}

const (
	hoursInDay   = 24
	daysInWeek   = 7
	weeksInMonth = 4
	monthsInYear = 12
)

// write http to http.ResponseWriter, this is like a main() function
func DogStats(w http.ResponseWriter, _ *http.Request) {
	dob := time.Date(2022, 7, 28, 12, 0, 0, 0, time.Local)

	loc, err := time.LoadLocation("Europe/Dublin")
	if err != nil {
		loc = time.UTC
	}

	age := getAge(dob, time.Now().In(loc))

	result := DogStat{
		Name:           "Huxley",
		DateOfBirth:    dob.Format("2 January 2006"),
		Breed:          "Golden Retriever",
		DateOfInterest: age.DateOfInterest,
		AgeDays:        age.Days,
		AgeWeeks:       age.Weeks,
		AgeMonths:      age.Months,
		AgeYears:       age.Years,
	}

	// write result to ResponseWriter using html template
	if err := templ.Execute(w, result); err != nil {
		log.Fatal(err)
	}
}

func getAge(dob, doi time.Time) Age {
	y, m, d := doi.Date()
	ageDays := doi.Sub(dob).Hours() / hoursInDay
	ageWeeks := ageDays / daysInWeek
	ageYears := float64(y - dob.Year())

	ageMonths := getAgeMonthsFractional(m, d, dob, ageYears)

	ageYears = getAgeYearsFractional(ageMonths, ageYears)

	return Age{
		DateOfInterest: doi.Format(time.RFC1123),
		Days:           int(math.Round(ageDays)),
		Weeks:          int(math.Round(ageWeeks)),
		Months:         ageMonths,
		Years:          float64(ageYears),
	}
}

// calculate fractional age in months to nearest quarter month
func getAgeMonthsFractional(m time.Month, d int, dob time.Time, ageYears float64) float64 {
	var partialMonthDays int

	ageMonths := float64(int(m)-int(dob.Month())) + ageYears*monthsInYear
	daysInBirthMonth := daysInMonth(dob.Month(), dob.Year())

	if d < dob.Day() {
		ageMonths--
		partialMonthDays = d + daysInBirthMonth - dob.Day()
	} else {
		partialMonthDays = d - dob.Day()
	}

	monthFraction := (math.Round(float64(partialMonthDays) / daysInWeek)) / weeksInMonth
	ageMonths += monthFraction

	return ageMonths
}

// trick to get the number of days in a month. Add 1 month to day zero (the day before the 1st of the month)
func daysInMonth(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// get the fractional age in years to the nearest quarter
func getAgeYearsFractional(ageMonths, ageYears float64) float64 {
	yr, frac := math.Modf(ageMonths / monthsInYear)

	//nolint:gomnd //readability, values easier to read than named constants
	switch {
	case frac < 0.25:
		ageYears = yr
	case frac < 0.5:
		ageYears = yr + 0.25
	case frac < 0.75:
		ageYears = yr + 0.5
	case frac < 1:
		ageYears = yr + 0.75
	}

	return ageYears
}
