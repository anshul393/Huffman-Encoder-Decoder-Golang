package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// func main() {
// 	filename := flag.String("filename", "", "path to file to compress")
// 	output := flag.String("output", "", "compressed file path")
// 	flag.Parse()

// 	// encoder := huffman.NewEncoder(*filename)
// 	// if err := encoder.Encode(*output); err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	decoder := huffman.NewDecoder(*filename)
// 	decoder.Decode(*output)

// }

func main() {
	f1 := flag.String("f1", "", "")
	f2 := flag.String("f2", "", "")
	flag.Parse()

	// compareFiles(*f1, *f2)

	h1, err := HashFile(*f1)
	if err != nil {
		log.Fatal(err)
	}

	h2, err := HashFile(*f2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(h1 == h2, h1, h2)

	// f, err := os.Open(*f1)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer f.Close()

	// ff, err := os.Open(*f2)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer ff.Close()

	// r1, r2 := bufio.NewScanner(f), bufio.NewScanner(ff)
	// r1.Split(bufio.ScanRunes)
	// r2.Split(bufio.ScanRunes)

	// for {
	// 	// fmt.Println("here")
	// 	if !r1.Scan() {
	// 		break
	// 	}
	// 	if !r2.Scan() {
	// 		break
	// 	}

	// 	s1, s2 := r1.Text(), r2.Text()

	// 	if s1 != s2 {
	// 		fmt.Printf("'%v' vs '%v'\n", s1, s2)

	// 	}

	// }

	// fmt.Println(r1.Err(), r2.Err())

}
func normalizeLineEndings(input string) string {
	return strings.ReplaceAll(input, "\r\n", "\n")
}

func HashFile(filepath string) (string, error) {
	// Open the file
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Create a new MD5 hash
	h := md5.New()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := normalizeLineEndings(scanner.Text())
		_, err := io.WriteString(h, line)
		if err != nil {
			return "", err
		}
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// Get the final hash as a byte slice
	hashInBytes := h.Sum(nil)

	// Convert the hash to a hexadecimal string
	return hex.EncodeToString(hashInBytes), nil
}

func compareFiles(f1, f2 string) error {
	// Open the files
	file1, err := os.Open(f1)
	if err != nil {
		return err
	}
	defer file1.Close()

	file2, err := os.Open(f2)
	if err != nil {
		return err
	}
	defer file2.Close()

	// Create scanners for both files
	r1 := bufio.NewScanner(file1)
	r2 := bufio.NewScanner(file2)
	r1.Split(bufio.ScanLines)
	r2.Split(bufio.ScanLines)

	// Compare files rune by rune
	for {
		hasNext1 := r1.Scan()
		hasNext2 := r2.Scan()

		// If one file ends before the other
		if hasNext1 != hasNext2 {
			fmt.Printf("Difference found: %q vs %q\n", r1.Text(), r2.Text())
			fmt.Println("Files have different lengths")
			break
		}

		// Break if both are done scanning
		if !hasNext1 {
			break
		}

		// Compare runes
		s1, s2 := r1.Text(), r2.Text()
		if s1 != s2 {
			fmt.Printf("Difference found: %q vs %q\n", s1, s2)
			break
		}
	}

	for r1.Scan() {
		fmt.Printf("%q", r1.Text())
	}

	for r2.Scan() {
		fmt.Printf("%q", r2.Text())
	}

	if err := r1.Err(); err != nil {
		return fmt.Errorf("error reading file1: %v", err)
	}
	if err := r2.Err(); err != nil {
		return fmt.Errorf("error reading file2: %v", err)
	}

	fmt.Println("Files are identical")
	return nil
}
