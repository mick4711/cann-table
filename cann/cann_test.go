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
		err   error
	}{
		{validStandings, validCannTable, nil},
		// TODO custom error {[]byte{}, []Row{}, nil},
	}

	for _, test := range tests {
		got, err := generateCann(test.input)
		if err != test.err {
			t.Errorf("generateCann(%v) err:%v, want:%v", string(test.input), err, test.err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("generateCann()\ngot :%v, \nwant:%v", got, test.want)
		}
	}
}
