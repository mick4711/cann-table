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
		returnError(err, w)
		return
	}

	cannTable, err := generateCann(standings)
	if err != nil {
		returnError(err, w)
		return
	}

	if err := writeResponse(w, cannTable); err != nil {
		returnError(err, w)
		return
	}
}

// TODO display empty page template with error message
func returnError(err error, w http.ResponseWriter) {
	log.Printf("\n*********** FATAL ERROR ********** [%s]\n", err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, err)
}

// fetch standard table standings
func getStandings() ([]byte, error) {
	// configure request
	url := `http://api.football-data.org/v4/competitions/PL/standings`

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("error creating standings request: %w", err)
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
		return nil, fmt.Errorf("error requesting standings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("standings response status not OK: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading standings response: %w", err)
	}

	return body, nil
}

// generate Cann table from standard standings table
func generateCann(standings []byte) ([]Row, error) {
	// unmarshall json standings into DataResponse slice of TableRows
	var dataResponse DataResponse
	if err := json.Unmarshal(standings, &dataResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling json from standings response:%w", err)
	}

	standingsTable := dataResponse.Standings[0].Table
	maxPoints := standingsTable[0].Points
	minPoints := standingsTable[len(standingsTable)-1].Points

	// generate an empty Cann table with the correct number of rows, set points values
	cannTable := make([]Row, maxPoints-minPoints+1)
	for i := range cannTable {
		cannTable[i].Points = maxPoints - Points(i)
	}

	const rowFormat = "[%d]%s(%d, %+d)"

	// loop thru standard table and assign team names and details to their point values in the Cann table
	for _, row := range standingsTable {
		index := maxPoints - row.Points
		rowData := fmt.Sprintf(rowFormat, row.Position, row.Team.ShortName, row.Played, row.GoalDiff)
		cannTable[index].Teams += fmt.Sprintf(" - %v", rowData)
	}

	return cannTable, nil
}

// write Cann table to response
func writeResponse(w http.ResponseWriter, cannTable []Row) error {
	cannTemplate := template.Must(template.ParseFiles("cann/CannTemplate.html"))
	if err := cannTemplate.Execute(w, cannTable); err != nil {
		return fmt.Errorf("error executing cannTemplate: %w", err)
	}

	return nil
}
