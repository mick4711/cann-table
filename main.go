package main

import "fmt"

func main() {
	standings := getStandings()
	canntable := generateCann(standings)
	fmt.Print(canntable)
}

func getStandings() []byte {
	// TODO http.Get http://api.football-data.org/v4/competitions/PL/standings
	return nil
}

func generateCann(standings []byte) []string {
	// TODO extract points and assign to table
	return []string{"1"}
}
