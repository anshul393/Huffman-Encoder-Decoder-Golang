package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// CharFreqMap calculates frequency of characters read from r io.Reader. Returns (map[string]int, error)
func CharFreqMap(filepath string) (map[rune]int, error) {

	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	fmp := make(map[rune]int)

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanRunes)

	for s.Scan() {
		// rn := bytes.Runes(s.Bytes())
		// fmt.Println(rn, len(rn))
		fmp[[]rune(s.Text())[0]]++
		// bw.Write(s.Bytes())
	}

	if err := s.Err(); err != nil {
		log.Printf("err: CharFreqMap %q", err)
		return nil, err
	}

	return fmp, nil

}

func PopLastItem(builder *strings.Builder) {
	// Convert the builder's contents to a string
	str := builder.String()

	// Check if the string is non-empty
	if len(str) > 0 {
		// Remove the last character
		str = str[:len(str)-1]

		// Reset the builder and write the modified string back
		builder.Reset()
		builder.WriteString(str)
	}
}
