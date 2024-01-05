package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// A Team contains details for a team.
type CannRow struct {
	Points int
	Teams  []string
}

// A Team contains details for a team.
type Team struct {
	ID        int    `json:"id"`
	ShortName string `json:"shortName"`
}

// A Row contains details for a table row.
type Row struct {
	Team     Team `json:"team"`
	Position int  `json:"position"`
	Played   int  `json:"playedGames"`
	Points   int  `json:"points"`
}

// A Standings contains a table of Rows.
type Standings struct {
	Table []Row `json:"table"`
}

// DataResponse contains the Standings
type DataResponse struct {
	Standings []Standings `json:"standings"`
}

func main() {
	standings := getStandings()
	// fmt.Println(string(standings))
	canntable := generateCann(standings)
	fmt.Print(canntable)
}

func getStandings() []byte {
	url := `http://api.football-data.org/v4/competitions/PL/standings`
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln("error getting standings:", err)
	}

	apiToken := os.Getenv("API_TOKEN")
	if len(apiToken) == 0 {
		log.Fatalln("API_TOKEN env var not set")
	}
	log.Println("API_TOKEN: ", len(apiToken))
	req.Header.Add("X-Auth-Token", apiToken)

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

func generateCann(standings []byte) []string {
	var dataResponse DataResponse
	if err := json.Unmarshal(standings, &dataResponse); err != nil {
		log.Fatalln("error unmarshalling json from response standings:", err)
	}
	standingsTable := dataResponse.Standings[0].Table
	maxPoints := standingsTable[0].Points
	minPoints := standingsTable[len(standingsTable)-1].Points

	cannTable := make(map[int][]string)
	for i := maxPoints; i >= minPoints; i-- {
		cannTable[i] = []string{}
	}

	for _, row := range standingsTable {
		points := row.Points
		cannTable[points] = append(cannTable[points], row.Team.ShortName)
	}

	for i := maxPoints; i >= minPoints; i-- {
		teams := ""
		for _, team := range cannTable[i] {
			teams += fmt.Sprintf(" - %v", team)
		}
		fmt.Printf("%02d %v\n", i, teams)
	}

	// fmt.Println(cannTable)
	return []string{"1"}
}
