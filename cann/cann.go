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

// A Row contains the points and teams with those points
type Row struct {
	Points Points
	Teams  string
}

// A Table contains the table along with max and min points.
type Table struct {
	Table     map[Points][]string
	MaxPoints Points
	MinPoints Points
}

// A Team contains details for a team.
type Team struct {
	ID        int    `json:"id"`
	ShortName string `json:"shortName"`
}

// A TableRow contains details for a standings table row.
type TableRow struct {
	Team     Team   `json:"team"`
	Position int    `json:"position"`
	Played   int    `json:"playedGames"`
	Points   Points `json:"points"`
	GoalDiff int    `json:"goalDifference"`
}

// A Standings contains a table of Rows, i.e. teams and points.
type Standings struct {
	Table []TableRow `json:"table"`
}

// DataResponse contains the Standings
type DataResponse struct {
	Standings []Standings `json:"standings"`
}

// fetches the standard table standings, generates and outputs the Cann table
func GenerateTable(w http.ResponseWriter, _ *http.Request) {
	standings, err := getStandings()
	if err != nil {
		errMsg := fmt.Sprintf("Unable to read current league standings %s", err)
		log.Printf("\n*********** FATAL ERROR *********************** [%s]  **************\n", errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, errMsg)

		return
	}

	canntable := generateCann(standings)
	setOutput(w, canntable)
}

// fetch standard table standings
func getStandings() ([]byte, error) {
	// configure request
	url := `http://api.football-data.org/v4/competitions/PL/standings`

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("create standings request failure: %w", err)
	}

	// add API token to header
	apiToken, ok := os.LookupEnv("API_TOKEN")
	if !ok {
		return nil, fmt.Errorf("environment variable -API_TOKEN- can not be read")
	}

	req.Header.Add("X-Auth-Token", apiToken)

	// get the response body
	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response failure: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status not OK: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	return body, nil
}

// generate Cann table from standard standings table
func generateCann(standings []byte) Table {
	// unmarshall json standings into table of Row types
	var dataResponse DataResponse
	if err := json.Unmarshal(standings, &dataResponse); err != nil {
		log.Fatalln("error unmarshalling json from response standings:", err)
	}

	standingsTable := dataResponse.Standings[0].Table
	maxPoints := standingsTable[0].Points
	minPoints := standingsTable[len(standingsTable)-1].Points

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

	return Table{cannTable, maxPoints, minPoints}
}

// generate output display
func setOutput(w http.ResponseWriter, cannTable Table) {
	// create slice of Cann rows
	rowsCount := cannTable.MaxPoints - cannTable.MinPoints + 1
	tbl := make([]Row, 0, rowsCount)

	// TODO can this countdown loop be done in the template (maybe using Templ)
	// fill the slice of Cann rows in descending sorted order
	for i := cannTable.MaxPoints; i >= cannTable.MinPoints; i-- {
		teams := ""
		for _, team := range cannTable.Table[i] {
			teams += fmt.Sprintf(" - %v", team)
		}

		// TODO update instead of append with pre allocated slice
		tbl = append(tbl, Row{Points: i, Teams: teams})
	}

	// generate html output
	cannTemplate := template.Must(template.ParseFiles("cann/CannTemplate.html"))
	if err := cannTemplate.Execute(w, tbl); err != nil {
		log.Fatal(err)
	}
}
