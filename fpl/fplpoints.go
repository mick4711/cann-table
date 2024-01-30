// reads a list of comma seperated FPL manager ids from environment variable "managers"
// and retrieves the current gameweek scores for the managers.
package fpl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type FplResponse struct { // fields retrieved from FPL API
	CurrentEvent         int    `json:"current_event"`
	ID                   int    `json:"id"`
	ManagerFirstName     string `json:"player_first_name"`
	ManagerLastName      string `json:"player_last_name"`
	Name                 string `json:"name"`
	SummaryOverallPoints int    `json:"summary_overall_points"`
	SummaryOverallRank   int    `json:"summary_overall_rank"`
	SummaryEventPoints   int    `json:"summary_event_points"`
	SummaryEventRank     int    `json:"summary_event_rank"`
}
type ManagerEntry struct { // stats for a manager for current gameweek
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Team     string `json:"team"`
	Points   int    `json:"points"`
	Rank     int    `json:"rank"`
	GwPoints int    `json:"gw_points"`
	GwRank   int    `json:"gw_rank"`
	Link     string `json:"link"`
}
type ManagerEntryResult struct { // result wrapper for ManagerEntry, Gameweek, Error
	Gameweek          int
	ManagerEntryValue ManagerEntry
	Error             error
}
type LeagueResponse struct { // response with array of manager entries
	Gameweek  int            `json:"gameweek"`
	Timestamp string         `json:"timestamp"`
	League    []ManagerEntry `json:"league"`
}

var fplURL = "https://fantasy.premierleague.com/api/entry/%v/"

// var fplURL = "http://MIKE-DEV.local:3001/api/entry/%v/"
// var fplURL = "http://MIKE-ALT.local:3001/api/entry/%v/"

func FplPoints(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	managers, ok := os.LookupEnv("managers")
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "LookupEnv managers not ok")
		return
	}

	// retrieve and filter data from FPL for the list of manager ids
	leagueResponse, err := getData(managers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v\n", err)
		return
	}

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

func getData(managers string) (LeagueResponse, error) {

	// initialise
	managerList := strings.Split(managers, ",")
	league := []ManagerEntry{}                        // slice of manager gameweek entries
	chManagerEntries := make(chan ManagerEntryResult) // channel to gather manager entries
	var gameweek int                                  // var to hold the gameweek value

	// loop thru manager list, fire off goroutines to get entries for each manager, results sent to channels
	for _, manager := range managerList {
		go getManagerEntries(strings.TrimSpace(manager), chManagerEntries)
	}

	// receive results from channels, set gameweek once and build up league table
	for range managerList {
		managerEntries := <-chManagerEntries
		if managerEntries.Error != nil {
			return LeagueResponse{}, managerEntries.Error
		}

		gameweekResponse := managerEntries.Gameweek
		if gameweek == 0 && gameweekResponse > 0 {
			gameweek = gameweekResponse
		}
		league = append(league, managerEntries.ManagerEntryValue)
	}

	// construct response
	leagueResponse := LeagueResponse{
		Gameweek:  gameweek,
		Timestamp: time.Now().Format("Mon Jan _2 15:04:05 MST 2006"),
		League:    league,
	}

	return leagueResponse, nil
}

func getManagerEntries(entry string, chManagerEntries chan<- ManagerEntryResult) {
	url := fmt.Sprintf(fplURL, entry)
	resp, err := http.Get(url)
	if err != nil {
		chManagerEntries <- ManagerEntryResult{Error: err}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		chManagerEntries <- ManagerEntryResult{Error: err}
		return
	}

	if resp.StatusCode == http.StatusNotFound {
		chManagerEntries <- ManagerEntryResult{
			Gameweek:          -1,
			ManagerEntryValue: ManagerEntry{Name: fmt.Sprintf("ID %v Not Found (404)", entry)},
		}
		return
	}

	if resp.StatusCode != http.StatusOK {
		chManagerEntries <- ManagerEntryResult{Error: fmt.Errorf("get manager ID %v not OK, Status: %v", entry, resp.Status)}
		return
	}

	var fplResponse FplResponse
	if err := json.Unmarshal(body, &fplResponse); err != nil {
		chManagerEntries <- ManagerEntryResult{Error: err}
		return
	}

	gw := fplResponse.CurrentEvent
	managerEntryResult := ManagerEntryResult{
		Gameweek: gw,
		ManagerEntryValue: ManagerEntry{
			ID:       fplResponse.ID,
			Name:     fmt.Sprintf("%v %v", fplResponse.ManagerFirstName, fplResponse.ManagerLastName),
			Team:     fplResponse.Name,
			Points:   fplResponse.SummaryOverallPoints,
			Rank:     fplResponse.SummaryOverallRank,
			GwPoints: fplResponse.SummaryEventPoints,
			GwRank:   fplResponse.SummaryEventRank,
			Link:     fmt.Sprintf("https://fantasy.premierleague.com/entry/%v/event/%d", entry, gw),
		},
	}
	chManagerEntries <- managerEntryResult
}
