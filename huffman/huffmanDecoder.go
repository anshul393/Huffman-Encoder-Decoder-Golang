package huffman

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sync"
)

type HuffmanDecoder struct {
	fileToDecode string
	separator    []byte //"\n\n\n\n\n"
	root         HuffmanNode
	currNode     HuffmanNode
	paddedBits   int
}

func NewDecoder(fileToDecode string) *HuffmanDecoder {
	return &HuffmanDecoder{fileToDecode: fileToDecode, separator: []byte("\n\n\n\n\n")}
}

// &HuffmanDecoder.Decode decodes the huffman encoded bytes from filetoDecode
// and writes decoded runes to the "out" file
func (hd *HuffmanDecoder) Decode(out string) error {
	f, err := os.Open(hd.fileToDecode)
	if err != nil {
		log.Printf("error: opening file %q error %q", hd.fileToDecode, err)
		return err
	}

	charFreqMap := make(map[rune]int)

	dec := json.NewDecoder(f)
	if err := dec.Decode(&charFreqMap); err != nil {
		log.Printf("error: json decoding file %q error %q", hd.fileToDecode, err)
		return err
	}

	paddedBits := make(map[string]int)
	if err := dec.Decode(&paddedBits); err != nil {
		log.Printf("error: json decoding file %q error %q", hd.fileToDecode, err)
		return err
	}
	f.Close()

	root := BuildHuffmanTree(charFreqMap)
	hd.root = root
	hd.currNode = root
	hd.paddedBits = paddedBits["bit_count"]

	ch1 := make(chan byte, 2048)
	ch2 := make(chan rune, 2048)

	var wg *sync.WaitGroup = &sync.WaitGroup{}
	wg.Add(3)

	go hd.read(ch1, wg)
	go hd.helper(ch1, ch2, wg)
	go hd.write(out, ch2, wg)

	wg.Wait()

	return nil

}

// &HuffmanDecoder.read reads the bytes from the fileToDecode,
// it skips all the bits until the separator and writes bytes to decode into writeChan.
func (hd *HuffmanDecoder) read(writeChan chan<- byte, wg *sync.WaitGroup) {

	defer wg.Done()

	f, err := os.Open(hd.fileToDecode)
	if err != nil {
		log.Printf("error: opening file %q err %q", hd.fileToDecode, err)
		return
	}
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanBytes)

	sep := ""

	for {
		if sep == string(hd.separator) {
			break
		}
		s.Scan()
		if len(sep) == 5 {
			sep = sep[1:] + s.Text()
		} else {
			sep += s.Text()
		}

	}

	for s.Scan() {
		writeChan <- s.Bytes()[0]
	}

	close(writeChan)
}

// &HuffmanDecoder.helper reads bytes to decode from readChan,
// decodes the bytes using huffmanTree and writes decoded runes to the writeChan.
func (hd *HuffmanDecoder) helper(readChan <-chan byte, writeChan chan<- rune, wg *sync.WaitGroup) {

	defer wg.Done()

	b1 := <-readChan

	for {
		b2, ok := <-readChan
		if !ok {
			break
		}
		for i := 7; i >= 0; i-- {
			step := (b1 >> i) & 1

			if lNode, ok := hd.currNode.(LeafNode); ok {
				writeChan <- lNode.char
				hd.currNode = hd.root
			}

			iNode := hd.currNode.(IntermediateNode)

			if step == 0 {
				hd.currNode = *iNode.left
			} else {
				hd.currNode = *iNode.right
			}
		}

		b1 = b2
	}

	for i := 7; i >= hd.paddedBits; i-- {
		step := (b1 >> i) & 1

		if lNode, ok := hd.currNode.(LeafNode); ok {
			writeChan <- lNode.char
			hd.currNode = hd.root
		}

		iNode := hd.currNode.(IntermediateNode)

		if step == 0 {
			hd.currNode = *iNode.left
		} else {
			hd.currNode = *iNode.right
		}
	}

	if lNode, ok := hd.currNode.(LeafNode); ok {
		writeChan <- lNode.char
		hd.currNode = hd.root
	}

	close(writeChan)
}

// &HuffmanDecoder.write reads from readChan decoded runes and writes them to the "out" file.
func (hd *HuffmanDecoder) write(out string, readChan <-chan rune, wg *sync.WaitGroup) {
	defer wg.Done()

	f, err := os.Create(out)
	if err != nil {
		log.Printf("error %q", err)
		return
	}
	bufWriter := bufio.NewWriter(f)

	for bt := range readChan {
		bufWriter.WriteRune(bt)
	}

	if err := bufWriter.Flush(); err != nil {
		log.Printf("error %q", err)
	}
}
