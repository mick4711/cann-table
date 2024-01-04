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
