// Generate a Cann table for the English Premier League. https://en.wikipedia.org/wiki/Cann_table
// A Cann table shows the league positions with gaps to emphasise points differences between teams.
// The standard league table standings are retrieved from api.football-data.org
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/mick4711/moh/huxley"
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

// main entry point - http server
func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/cann", cannHandler)
	http.HandleFunc("/huxley", huxleyHandler)
	http.HandleFunc("/fpl", fplHandler)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// log request details
func logRequest(req *http.Request) {
	log.Println("User-Agent:", req.Header["User-Agent"])
	log.Println("Cf-Ipcountry:", req.Header["Cf-Ipcountry"])
	log.Println("Cf-Connecting-Ip:", req.Header["Cf-Connecting-Ip"])
	log.Println("Sec-Ch-Ua-Platform:", req.Header["Sec-Ch-Ua-Platform"])
	log.Println("Sec-Ch-Ua:", req.Header["Sec-Ch-Ua"])
}

// displays landing page with links to other pages
func homeHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Request on /")
	logRequest(req)

	// generate html output
	homeTemplate := template.Must(template.ParseFiles("HomeTemplate.html"))
	if err := homeTemplate.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}

// displays Huxley's personal details
func huxleyHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Request on /huxley")
	logRequest(req)

	// generate html output
	huxley.DogStats(w , req)
}

// displays FPL league table
func fplHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Request on /fpl")
	logRequest(req)

	type LeagueResponse struct { // response with array of manager entries
		Gameweek  int            `json:"gameweek"`
		Timestamp string         `json:"timestamp"`
	}
	
	leagueResponse := LeagueResponse {1, "Under construction"}

	// convert response to json
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(leagueResponse, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v\n", err)
		return
	}

	// display results
	fmt.Fprintf(w, "%+v\n", string(response))
}

// fetches the standard table standings, generates and outputs the Cann table
func cannHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Request on /cann")
	logRequest(req)

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
	cannTemplate := template.Must(template.ParseFiles("CannTemplate.html"))
	if err := cannTemplate.Execute(w, tbl); err != nil {
		log.Fatal(err)
	}
}
