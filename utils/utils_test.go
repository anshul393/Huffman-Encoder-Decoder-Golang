package utils

import (
	"fmt"
	"testing"
)

func TestCharFreqMap(t *testing.T) {
	filepath := `D:\coding-challenges\HuffmanEncodingDecoding\test.txt`

	mp, err := CharFreqMap(filepath)
	if err != nil {
		t.Fatalf("error in CharFreqMap err %q", err)
	}

	tests := map[rune]int{
		'X': 333,
		't': 223000,
	}

	for ch, val := range mp {
		fmt.Printf("%q:::%d", ch, val)
	}

	for char, expectedFreq := range tests {
		if expectedFreq != mp[char] {
			t.Fatalf("char:%q expected frequency %d got %d", char, expectedFreq, mp[char])
		}
	}
}
