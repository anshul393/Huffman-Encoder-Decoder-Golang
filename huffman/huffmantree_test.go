package huffman

import (
	"log"
	"testing"
)

func TestBuildHuffmanTree(t *testing.T) {
	input := map[rune]int{'C': 32, 'D': 42, 'E': 120, 'K': 7, 'L': 42, 'M': 24, 'U': 37, 'Z': 2}

	huffmanNode := BuildHuffmanTree(input)

	root, ok := huffmanNode.(IntermediateNode)
	if !ok {
		log.Fatal("expected IntermediateNode got something else")
	}

	if root.NodeWeight() != 306 {
		log.Fatalf("expected node weight of %d got %d", 306, root.NodeWeight())
	}

	eLeafNode, ok := (*root.left).(LeafNode)
	if !ok {
		log.Fatal("expected LeafNode got something else")
	}

	if eLeafNode.NodeWeight() != 120 {
		log.Fatalf("expected node weight of %d got %d", 120, root.NodeWeight())
	}

	if eLeafNode.char != 'E' {
		log.Fatalf("expected node with char %q got %q", "E", eLeafNode.char)
	}

	iNode, ok := (*root.right).(IntermediateNode)
	if !ok {
		log.Fatal("expected IntermediateNode got something else")
	}
	if iNode.NodeWeight() != 186 {
		log.Fatalf("expected node weight of %d got %d", 186, root.NodeWeight())
	}

}

func TestTraverseHuffmanTree(t *testing.T) {
	input := map[rune]int{'C': 32, 'D': 42, 'E': 120, 'K': 7, 'L': 42, 'M': 24, 'U': 37, 'Z': 2}
	huffmanNode := BuildHuffmanTree(input)

	mp := TraverseHuffmanTree(huffmanNode)

	expected := []struct {
		char    rune
		bitRepr string
	}{
		{'E', "0"},
		{'U', "100"},
		{'D', "101"},
		{'L', "110"},
		{'C', "1110"},
		{'M', "11111"},
		{'K', "111101"},
		{'Z', "111100"},
	}

	for _, tt := range expected {
		if mp[tt.char] != tt.bitRepr {
			log.Fatalf("for %q expected bit representation %q got %q", tt.char, tt.bitRepr, mp[tt.char])
		}
	}

}
