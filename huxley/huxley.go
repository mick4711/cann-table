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
}

type Age struct {
	DateOfInterest string
	Days           int
	Weeks          int
	Months         float64
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
		DateOfBirth:    "28 July 2022",
		Breed:          "Golden Retriever",
		DateOfInterest: age.DateOfInterest,
		AgeDays:        age.Days,
		AgeWeeks:       age.Weeks,
		AgeMonths:      age.Months,
	}

	// write result to ResponseWriter using html template
	if err := templ.Execute(w, result); err != nil {
		log.Fatal(err)
	}
}

// TODO work on edge cases
func getAge(dob, doi time.Time) Age {
	y, m, d := doi.Date()

	ageDays := doi.Sub(dob).Hours() / hoursInDay
	ageWeeks := ageDays / daysInWeek
	ageYears := y - dob.Year()
	ageMonths := float64(int(m) - int(dob.Month()) + ageYears*monthsInYear)
	monthDays := d

	if d < dob.Day() {
		ageMonths--
	} else {
		monthDays -= d
	}

	monthFraction := (math.Round(float64(monthDays) / daysInWeek)) / weeksInMonth
	ageMonths += monthFraction

	return Age{
		DateOfInterest: doi.Format("Mon 02-Jan-2006 15:04:05"), //
		Days:           int(math.Round(ageDays)),
		Weeks:          int(math.Round(ageWeeks)),
		Months:         ageMonths,
	}
}
