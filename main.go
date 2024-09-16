package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"io"
	"os"

	"github.com/anshul393/huffmanCompress/huffman"
)

func main() {
	encode := flag.Bool("encode", false, "")
	decode := flag.Bool("decode", false, "")
	file := flag.String("filepath", "", "path of file to decode or encode")
	out := flag.String("output", "", "path of output file")

	flag.Parse()

	if encode == decode {
		os.Exit(1)
	}

	if *encode{
		encoder := huffman.NewEncoder(*file)
		encoder.Encode(*out)
	}

	if *decode{
		decoder := huffman.NewDecoder(*file)
		decoder.Decode(*out)
	}

}

// func main() {
// 	f1 := flag.String("f1", "", "")
// 	f2 := flag.String("f2", "", "")
// 	flag.Parse()

// 	// compareFiles(*f1, *f2)

// 	h1, err := HashFile(*f1)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	h2, err := HashFile(*f2)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(h1 == h2, h1, h2)

// 	// f, err := os.Open(*f1)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// defer f.Close()

// 	// ff, err := os.Open(*f2)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// defer ff.Close()

// 	// r1, r2 := bufio.NewScanner(f), bufio.NewScanner(ff)
// 	// r1.Split(bufio.ScanRunes)
// 	// r2.Split(bufio.ScanRunes)

// 	// for {
// 	// 	// fmt.Println("here")
// 	// 	if !r1.Scan() {
// 	// 		break
// 	// 	}
// 	// 	if !r2.Scan() {
// 	// 		break
// 	// 	}

// 	// 	s1, s2 := r1.Text(), r2.Text()

// 	// 	if s1 != s2 {
// 	// 		fmt.Printf("'%v' vs '%v'\n", s1, s2)

// 	// 	}

// 	// }

// 	// fmt.Println(r1.Err(), r2.Err())

// }

func HashFile(filepath string) (string, error) {
	// Open the file
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Create a new MD5 hash
	h := md5.New()

	io.Copy(h,f)

	// Create a scanner to read the file line by line
	// scanner := bufio.NewScanner(f)

	// for scanner.Scan() {
	// 	line := normalizeLineEndings(scanner.Text())
	// 	_, err := io.WriteString(h, line)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// }

	// // Check for any scanning errors
	// if err := scanner.Err(); err != nil {
	// 	return "", err
	// }

	// Get the final hash as a byte slice
	hashInBytes := h.Sum(nil)

	// Convert the hash to a hexadecimal string
	return hex.EncodeToString(hashInBytes), nil
}
