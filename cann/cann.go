// Generate a Cann table for the English Premier League. https://en.wikipedia.org/wiki/Cann_table
// A Cann table shows the league positions with gaps to emphasise points differences between teams.
// The standard league table standings are retrieved from api.football-data.org
package cann

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type Points int

// A CannRow contains the points and teams with those points
type CannRow struct {
	Points Points
	Teams  string
}

// A CannTable contains the table along with max and min points.
type CannTable struct {
	Table     map[Points][]string
	MaxPoints Points
	MinPoints Points
}

// A Team contains details for a team.
type Team struct {
	ID        int    `json:"id"`
	ShortName string `json:"shortName"`
}

// A Row contains details for a standings table row.
type Row struct {
	Team     Team   `json:"team"`
	Position int    `json:"position"`
	Played   int    `json:"playedGames"`
	Points   Points `json:"points"`
	GoalDiff int    `json:"goalDifference"`
}

// A Standings contains a table of Rows, i.e. teams and points.
type Standings struct {
	Table []Row `json:"table"`
}

// DataResponse contains the Standings
type DataResponse struct {
	Standings []Standings `json:"standings"`
}

// fetches the standard table standings, generates and outputs the Cann table
func GenerateTable(w http.ResponseWriter, req *http.Request) {
	standings := getStandings()
	canntable := generateCann(standings)
	setOutput(w, canntable)
}

// fetch standard table standings
func getStandings() []byte {
	//configure request
	url := `http://api.football-data.org/v4/competitions/PL/standings`
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln("error getting standings:", err)
	}

	// add API token to header
	apiToken := os.Getenv("API_TOKEN")
	if len(apiToken) == 0 {
		log.Fatalln("API_TOKEN env var not set")
	}
	req.Header.Add("X-Auth-Token", apiToken)

	// get the response body
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("error getting standings:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalln("error getting standings, status not OK:", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("error getting standings, cannot read response:", err)
	}

	return body
}

// generate Cann table from standard standings table
func generateCann(standings []byte) CannTable {
	// unmarshall json standings into table of Row types
	var dataResponse DataResponse
	if err := json.Unmarshal(standings, &dataResponse); err != nil {
		log.Fatalln("error unmarshalling json from response standings:", err)
	}
	standingsTable := dataResponse.Standings[0].Table
	maxPoints := Points(standingsTable[0].Points)
	minPoints := Points(standingsTable[len(standingsTable)-1].Points)

	// TODO refactor to populate slice of CannRows directly
	// generate an empty Cann table with the correct number of empty rows with points values as keys
	cannTable := make(map[Points][]string)
	for i := maxPoints; i >= minPoints; i-- {
		cannTable[i] = []string{}
	}

	// loop thru standard table and assign team names and details to their point values in the Cann table
	const rowFormat = "[%d]%s(%d, %+d)"
	for _, row := range standingsTable {
		points := row.Points
		rowData := fmt.Sprintf(rowFormat, row.Position, row.Team.ShortName, row.Played, row.GoalDiff)
		cannTable[points] = append(cannTable[points], rowData)
	}

	return CannTable{cannTable, maxPoints, minPoints}
}

// generate output display
func setOutput(w http.ResponseWriter, cannTable CannTable) {
	// create slice of Cann rows
	rowsCount := cannTable.MaxPoints - cannTable.MinPoints + 1
	tbl := make([]CannRow, 0, rowsCount)

	// fill the slice of Cann rows in descending sorted order
	for i := cannTable.MaxPoints; i >= cannTable.MinPoints; i-- {
		teams := ""
		for _, team := range cannTable.Table[i] {
			teams += fmt.Sprintf(" - %v", team)
		}
		tbl = append(tbl, CannRow{Points: i, Teams: teams})
	}

	// generate html output
	cannTemplate := template.Must(template.ParseFiles("cann/CannTemplate.html"))
	if err := cannTemplate.Execute(w, tbl); err != nil {
		log.Fatal(err)
	}
}
