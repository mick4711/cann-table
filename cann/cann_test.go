package cann

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func TestGenerateCann(t *testing.T) {
	validStandings, err := os.ReadFile("standings_test.json")
	if err != nil {
		log.Fatalln(err)
	}

	validCannTable := []Row{
		{45, " - [1]Liverpool(20, -25)"},
		{44, ""},
		{43, ""},
		{42, " - [2]Aston Villa(20, +16)"},
		{41, ""},
		{40, " - [3]Man City(19, +24) - [4]Arsenal(20, +17)"},
		{39, " - [5]Tottenham(20, +13)"},
	}

	tests := []struct {
		input []byte
		want  []Row
		hasError bool
	}{
		{validStandings, validCannTable, false},
		{[]byte{}, []Row(nil), true},
	}

	for _, test := range tests {
		got, err := generateCann(test.input)
		if hasError := err != nil; hasError != test.hasError {
			t.Errorf("generateCann()\n got err:%v, \nwant hasError:%v", err, test.hasError)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("generateCann()\ngot :%#v, \nwant:%#v", got, test.want)
		}
	}
}
