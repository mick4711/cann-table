package main

import (
	"log"
	"os"
	"testing"
)

func TestGenerateCann(t *testing.T) {
	standings, err := os.ReadFile("standings.json")
	if err != nil {
		log.Fatalln(err)
	}

	var tests = []struct {
		input []byte
		want  []string
	}{
		{standings, []string{"1"}},
	}

	for _, test := range tests {
		got := generateCann(test.input)
		for i, elem := range got {
			if elem != test.want[i] {
				t.Errorf("generateCann(%v), got:%v, want:%v", string(test.input), elem, test.want[i])
			}
		}
		// if reflect.DeepEqual(got, test.want) {
		// 	t.Errorf("generateCann(%v), got:%v, want:%v", string(test.input), got, test.want)
		// }
	}
}

/*
want from test standings
map[9:[Sheffield Utd] 10:[] 11:[Burnley] 12:[] 13:[] 14:[] 15:[Luton Town] 16:[Everton] 17:[] 18:[] 19:[Brentford] 20:[Nottingham] 21:[Crystal Palace] 22:[] 23:[] 24:[Fulham] 25:[Bournemouth] 26:[] 27:[] 28:[Chelsea Wolverhampton] 29:[Newcastle] 30:[] 31:[Brighton Hove Man United] 32:[] 33:[] 34:[West Ham] 35:[] 36:[] 37:[] 38:[] 39:[Tottenham] 40:[Man City Arsenal] 41:[] 42:[Aston Villa] 43:[] 44:[] 45:[Liverpool]]
*/
