package fpl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TEST DATA ///////////////////////////////////////////////////////////////////////////////////
// const Gameweek = 99
// const ContentType = "Content-Type"
// const ApplicationJSON = "application/json"
const (
	Gameweek         = 99
	ContentType      = "Content-Type"
	ApplicationJSON  = "application/json"
	EntryPlaceholder = "/%v"
)

var mockFplResponse = []Response{
	{
		CurrentEvent:         Gameweek,
		ID:                   1,
		ManagerFirstName:     "first1",
		ManagerLastName:      "last1",
		Name:                 "team1",
		SummaryOverallPoints: 77,
		SummaryOverallRank:   66,
		SummaryEventPoints:   55,
		SummaryEventRank:     44,
	},
	{
		CurrentEvent:         Gameweek,
		ID:                   2,
		ManagerFirstName:     "first2",
		ManagerLastName:      "last2",
		Name:                 "team2",
		SummaryOverallPoints: 177,
		SummaryOverallRank:   166,
		SummaryEventPoints:   155,
		SummaryEventRank:     144,
	},
}

func TestGetData(t *testing.T) {
	// ARRANGE ///////////////////////////////////////////////////////////////////////////////////
	expectedManagersResponse := setExpectedManagersResponse()

	// httptest server to serve up mock json response
	ts := setTestServer()
	defer ts.Close()

	// overwrite fplURL to use httptest URL
	fplURL = ts.URL + EntryPlaceholder

	// ACT //////////////////////////////////////////////////////////////////////////////////////////////
	testResponse, err := getData("1, 2")
	// ASSERT ///////////////////////////////////////////////////////////////////////////////////////////
	// check err
	if err != nil {
		t.Errorf(`getData("1, 2") err = (%v), want: nil err`, err)
	}

	// correct gameweek, current event
	if testResponse.Gameweek != Gameweek {
		t.Errorf(`getData Gameweek = %v, want (%v)`, testResponse.Gameweek, Gameweek)
	}

	// correct number of league entries
	if len(testResponse.League) != len(mockFplResponse) {
		t.Errorf(`getData League manager entries count = %v, want (%v)`, len(testResponse.League), len(mockFplResponse))
	}

	// correct values for each league entry
	checkValidResponse(t, testResponse, expectedManagersResponse)
}

func setExpectedManagersResponse() []ManagerEntry {
	var expectedManagersResponse []ManagerEntry

	for _, fplResponse := range mockFplResponse {
		managerEntry := ManagerEntry{
			ID:       fplResponse.ID,
			Name:     fmt.Sprintf("%v %v", fplResponse.ManagerFirstName, fplResponse.ManagerLastName),
			Team:     fplResponse.Name,
			Points:   fplResponse.SummaryOverallPoints,
			Rank:     fplResponse.SummaryOverallRank,
			GwPoints: fplResponse.SummaryEventPoints,
			GwRank:   fplResponse.SummaryEventRank,
			Link:     fmt.Sprintf("https://fantasy.premierleague.com/entry/%v/event/%d", fplResponse.ID, Gameweek),
		}
		expectedManagersResponse = append(expectedManagersResponse, managerEntry)
	}

	return expectedManagersResponse
}

func setTestServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(ContentType, ApplicationJSON)

		switch r.URL.Path {
		case "/1":
			mockJSONResponse, err := json.Marshal(mockFplResponse[0])
			if err != nil {
				panic(err)
			}

			fmt.Fprintln(w, string(mockJSONResponse))
		case "/2":
			mockJSONResponse, err := json.Marshal(mockFplResponse[1])
			if err != nil {
				panic(err)
			}

			fmt.Fprintln(w, string(mockJSONResponse))
		}
	}))

	return ts
}

func checkValidResponse(t *testing.T, testResponse LeagueResponse, expectedManagersResponse []ManagerEntry) {
	t.Helper()

	for _, leagueEntry := range testResponse.League {
		switch leagueEntry.ID {
		case 1:
			if !((leagueEntry.Name == expectedManagersResponse[0].Name) &&
				(leagueEntry.Team == expectedManagersResponse[0].Team) &&
				(leagueEntry.Points == expectedManagersResponse[0].Points) &&
				(leagueEntry.Rank == expectedManagersResponse[0].Rank) &&
				(leagueEntry.GwPoints == expectedManagersResponse[0].GwPoints) &&
				(leagueEntry.GwRank == expectedManagersResponse[0].GwRank) &&
				(leagueEntry.Link == expectedManagersResponse[0].Link)) {
				t.Errorf(`getData League manager entry = (%v), want (%v)`, leagueEntry, expectedManagersResponse[0])
			}
		case 2:
			if !((leagueEntry.Name == expectedManagersResponse[1].Name) &&
				(leagueEntry.Team == expectedManagersResponse[1].Team) &&
				(leagueEntry.Points == expectedManagersResponse[1].Points) &&
				(leagueEntry.Rank == expectedManagersResponse[1].Rank) &&
				(leagueEntry.GwPoints == expectedManagersResponse[1].GwPoints) &&
				(leagueEntry.GwRank == expectedManagersResponse[1].GwRank) &&
				(leagueEntry.Link == expectedManagersResponse[1].Link)) {
				t.Errorf(`getData returned (%v) not   matching expected value (%v)`, leagueEntry, expectedManagersResponse[1])
			}
		}
	}
}

func TestFpl404(t *testing.T) {
	// httptest server to serve up mock json response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(ContentType, ApplicationJSON)

		switch r.URL.Path {
		case "/1":
			w.WriteHeader(http.StatusNotFound)
		case "/2":
			mockJSONResponse, err := json.Marshal(mockFplResponse[1])
			if err != nil {
				panic(err)
			}

			fmt.Fprintln(w, string(mockJSONResponse))
		}
	}))
	defer ts.Close()

	// overwrite fplURL to use httptest URL
	fplURL = ts.URL + EntryPlaceholder

	// ACT //////////////////////////////////////////////////////////////////////////////////////////////
	testResponse, err := getData("1, 2")
	// ASSERT ///////////////////////////////////////////////////////////////////////////////////////////
	// check err
	if err != nil {
		t.Errorf(`getData("1, 2") err = (%v), want: nil err`, err)
	}

	// 404 in name of unfound league entry
	for _, leagueEntry := range testResponse.League {
		if leagueEntry.ID != 2 {
			if !strings.Contains(leagueEntry.Name, "404") {
				t.Errorf(`getData manager entry = (%v), want: contains 404 in name`, leagueEntry.Name)
			}
		}
	}
}

func TestFpl5XX(t *testing.T) {
	// httptest server to serve up mock json response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(ContentType, ApplicationJSON)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	// overwrite fplURL to use httptest URL
	fplURL = ts.URL + EntryPlaceholder

	// ACT //////////////////////////////////////////////////////////////////////////////////////////////
	_, err := getData("1, 2")

	// ASSERT ///////////////////////////////////////////////////////////////////////////////////////////
	if err == nil {
		t.Error(`getData server error is nil, want: err not nil`)
	} else if !strings.Contains(err.Error(), "not OK, Status:") {
		t.Errorf(`getData server error (%v), want: "...not OK, Status:..."`, err)
	}
}
