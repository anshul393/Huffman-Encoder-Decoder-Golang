package huffman

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/anshul393/huffmanCompress/utils"
)

type HuffmanEncoder struct {
	root           HuffmanNode
	fileToEncode   string
	charFreqMap    map[rune]int
	huffmanByteRpr map[rune]string
	paddedBits     int
	separator      []byte // "/n/n/n/n/n"
}

func NewEncoder(fileToEncode string) *HuffmanEncoder {
	mp, err := utils.CharFreqMap(fileToEncode)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	root := BuildHuffmanTree(mp)
	huffmanRepr := TraverseHuffmanTree(root)

	return &HuffmanEncoder{fileToEncode: fileToEncode, charFreqMap: mp, root: root, huffmanByteRpr: huffmanRepr, separator: []byte("\n\n\n\n\n")}

}

// &HuffmanEncoder.Encode encode fileToEncode and writes huffman encoded bytes to "out" file.
// It will create or truncates the "out" file.
func (he *HuffmanEncoder) Encode(out string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}

	fmt.Println(f, he)

	defer f.Close()

	// writing metadata for decoding the compressed(huffman-encoded) file
	err = he.writeJSON(f, he.charFreqMap)
	if err != nil {
		return err
	}

	fmt.Println("here")

	paddedBits := he.paddedBitsCount()
	err = he.writeJSON(f, map[string]int{"bit_count": paddedBits})
	if err != nil {
		return err
	}

	// writing separator to separate the metadata and actual encoded data
	if _, err := f.Write(he.separator); err != nil {
		log.Printf("error: writing separator %q", err)
		return err
	}

	// to write encoded bits to the output file
	he.writeBits(f)

	return nil
}

func (he *HuffmanEncoder) writeBits(w io.Writer) {

	var wg *sync.WaitGroup = &sync.WaitGroup{}

	ch1 := make(chan rune, 2048)
	ch2 := make(chan byte, 2048)

	wg.Add(3)

	go he.read(ch1, wg)
	go he.helper(ch1, ch2, wg)
	go he.write(w, ch2, wg)

	wg.Wait()
}

// &HuffmanEncoder.read reads runes from fileToEncode and writes them to writeChan
func (he *HuffmanEncoder) read(writeChan chan<- rune, wg *sync.WaitGroup) {
	defer wg.Done()

	f, err := os.Open(he.fileToEncode)
	if err != nil {
		log.Printf("error: opening file %q error %q", he.fileToEncode, err)
		return
	}

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanRunes)

	for s.Scan() {
		writeChan <- []rune(s.Text())[0]
	}

	if err := s.Err(); err != nil {
		log.Printf("error: readFromReader %q", err)
	}

	close(writeChan)

}

// &HuffmanEncoder.helper reads the file runes to encode from readChan.
// calculates huffman sequence of bytes,
// writes bytes to writeChan(for &huffmanEncoder.write).
func (he *HuffmanEncoder) helper(readChan <-chan rune, writeChan chan<- byte, wg *sync.WaitGroup) {
	defer wg.Done()

	bitCount := 0
	var bt byte = 0

	for rn := range readChan {
		for _, v := range he.huffmanByteRpr[rn] {
			switch v - '0' {
			case 0:
				bt <<= 1
			case 1:
				bt <<= 1
				bt += 1
			}
			bitCount += 1

			if bitCount == 8 {
				writeChan <- bt
				bitCount = 0
				bt = 0
			}
		}

	}

	// write any left bits
	if bitCount != 0 {
		writeChan <- bt << (8 - bitCount)
	}
	close(writeChan)

}

// &HuffmanEncoder.write reads encoded huffman bytes from readChan
// and writes to the "out" file
func (he *HuffmanEncoder) write(w io.Writer, readChan <-chan byte, wg *sync.WaitGroup) {
	defer wg.Done()

	bufioWriter := bufio.NewWriter(w)

	for bt := range readChan {
		if err := bufioWriter.WriteByte(bt); err != nil {
			log.Printf("error: writeToWriter %q", err)
			return
		}

	}

	err := bufioWriter.Flush()
	if err != nil {
		log.Println("error: flushing the buffered bytes", err)
	}

}

func (he *HuffmanEncoder) writeJSON(w io.Writer, v any) error {
	freqMapBytes, err := json.Marshal(v)
	if err != nil {
		log.Printf("error: marshalling freqMap %q", err)
		return err
	}

	// fmt.Println("here")

	_, err = w.Write(freqMapBytes)
	if err != nil {
		log.Printf("error: writing marshalled bytes of %#v %q", v, err)
	}

	return err
}

// paddedBitsCount return number of 0 to be padded at the end of last huffman sequence of bit to make a  byte.
func (he *HuffmanEncoder) paddedBitsCount() int {
	count := 0

	for cr, ct := range he.charFreqMap {
		count = (count + (len(he.huffmanByteRpr[cr])*ct)%8) % 8
	}
	count = count % 8

	if count == 0{
		he.paddedBits = 0
	}else{
		he.paddedBits = 8 - count
	}


	return he.paddedBits
}
