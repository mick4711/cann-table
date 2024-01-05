package main

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

	cannTable := CannTable{map[int][]string{
		45: {"Liverpool"}, 42: {"Aston Villa"}, 40: {"Man City", "Arsenal"},
		39: {"Tottenham"}, 44: {}, 43: {}, 41: {},
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
			t.Errorf("generateCann(%v) Table, got:%v, want:%v", string(test.input), got.Table, test.want.Table)
		}
		if got.MaxPoints != test.want.MaxPoints {
			t.Errorf("generateCann(%v) MaxPoints, got:%v, want:%v", string(test.input), got.MaxPoints, test.want.MaxPoints)
		}
		if got.MinPoints != test.want.MinPoints {
			t.Errorf("generateCann(%v) MinPoints, got:%v, want:%v", string(test.input), got.MinPoints, test.want.MinPoints)
		}
	}
}
