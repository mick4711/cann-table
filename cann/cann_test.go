package cann

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func TestGenerateCann(t *testing.T) {
	standings, err := os.ReadFile("standings_test.json")
	if err != nil {
		log.Fatalln(err)
	}

	cannTable := CannTable{map[Points][]string{
		45: {"[1]Liverpool(20, -25)"}, 
		42: {"[2]Aston Villa(20, +16)"}, 
		40: {"[3]Man City(19, +24)", "[4]Arsenal(20, +17)"},
		39: {"[5]Tottenham(20, +13)"}, 
		44: {}, 
		43: {}, 
		41: {},
	}, 45, 39}

	var tests = []struct {
		input []byte
		want  CannTable
	}{
		{standings, cannTable},
	}

	for _, test := range tests {
		got := generateCann(test.input)
		if !reflect.DeepEqual(got.Table, test.want.Table) {
			t.Errorf("generateCann(%v) Table, \ngot :%v, \nwant:%v", string(test.input), got.Table, test.want.Table)
		}
		if got.MaxPoints != test.want.MaxPoints {
			t.Errorf("generateCann(%v) MaxPoints, got:%v, want:%v", string(test.input), got.MaxPoints, test.want.MaxPoints)
		}
		if got.MinPoints != test.want.MinPoints {
			t.Errorf("generateCann(%v) MinPoints, got:%v, want:%v", string(test.input), got.MinPoints, test.want.MinPoints)
		}
	}
}
