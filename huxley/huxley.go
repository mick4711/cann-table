// calculate and display Huxley's age in various units
package huxley

import (
	"fmt"
	"log"
	"math"
	"time"

	"html/template"
	"net/http"
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
	<body>
		<h1>{{.Name}} Age</h1>
		<h3>Date: {{.DateOfInterest}}</h3>
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

// write http to http.ResponseWriter, this is like a main() function
func DogStats(w http.ResponseWriter, r *http.Request) {

	dob := time.Date(2022, 7, 28, 12, 0, 0, 0, time.Local)
	age := getAge(dob, time.Now())

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
func getAge(dob time.Time, doi time.Time) Age {
	y, m, d := doi.Date()

	ageDays := doi.Sub(dob).Hours() / 24
	ageWeeks := ageDays / 7
	ageYears := y - dob.Year()

	ageMonths := float64(int(m) - int(dob.Month()) + ageYears*12)
	monthDays := d
	if d < dob.Day() {
		ageMonths--
	} else {
		monthDays -= d
	}
	monthFraction := (math.Round(float64(monthDays) / 7)) / 4
	ageMonths += monthFraction

	return Age{
		DateOfInterest: fmt.Sprint(doi.Format("Mon 02-Jan-2006")), //
		Days:           int(math.Round(ageDays)),
		Weeks:          int(math.Round(ageWeeks)),
		Months:         ageMonths,
	}
}
