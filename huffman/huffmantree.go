package huffman

import (
	"container/heap"
	"strings"

	"github.com/anshul393/huffmanCompress/utils"
)

type HuffmanNode interface {
	NodeWeight() int
	LeafCount() int
	LeafLiteralCount() int
}

type priorityQueue []HuffmanNode

func (pq priorityQueue) Len() int { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool {
	if pq[i].NodeWeight() == pq[j].NodeWeight() {
		if pq[i].LeafCount() == pq[j].LeafCount() {
			return pq[i].LeafLiteralCount() < pq[j].LeafLiteralCount()
		}
		return pq[i].LeafCount() < pq[j].LeafCount()
	}
	return pq[i].NodeWeight() < pq[j].NodeWeight()
}
func (pq *priorityQueue) Push(x any) {
	item := x.(HuffmanNode)
	*pq = append(*pq, item)
}
func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

type LeafNode struct {
	weight int
	char   rune
}

func (lf LeafNode) NodeWeight() int       { return lf.weight }
func (lf LeafNode) LeafCount() int        { return 1 }
func (lf LeafNode) LeafLiteralCount() int { return int(lf.char) }

type IntermediateNode struct {
	weight                int
	leafCount             int
	leafNodesLiteralCount int
	left                  *HuffmanNode
	right                 *HuffmanNode
}

func (in IntermediateNode) NodeWeight() int       { return in.weight }
func (lf IntermediateNode) LeafCount() int        { return lf.leafCount }
func (lf IntermediateNode) LeafLiteralCount() int { return lf.leafNodesLiteralCount }

func BuildHuffmanTree(mp map[rune]int) HuffmanNode {
	pq := make(priorityQueue, 0)
	heap.Init(&pq)
	for char, freq := range mp {
		heap.Push(&pq, LeafNode{weight: freq, char: char})
	}

	for len(pq) != 1 {
		n1, n2 := heap.Pop(&pq).(HuffmanNode), heap.Pop(&pq).(HuffmanNode)
		n3 := IntermediateNode{weight: n1.NodeWeight() + n2.NodeWeight(), left: &n1, right: &n2, leafCount: n1.LeafCount() + n2.LeafCount(), leafNodesLiteralCount: n1.LeafLiteralCount() + n2.LeafLiteralCount()}
		heap.Push(&pq, n3)
	}

	huffmanNode := heap.Pop(&pq).(HuffmanNode)

	return huffmanNode
}

func TraverseHuffmanTree(root HuffmanNode) map[rune]string {
	if lnode, ok := root.(LeafNode); ok {
		return map[rune]string{lnode.char: "0"}
	}

	sb := &strings.Builder{}
	mp := map[rune]string{}

	return DFSHuffmanSearch(root, sb, mp)

}

func DFSHuffmanSearch(root HuffmanNode, sb *strings.Builder, mp map[rune]string) map[rune]string {
	if root == nil {
		return mp
	}

	if lnode, ok := root.(LeafNode); ok {
		mp[lnode.char] = sb.String()
		return mp
	}

	inode := root.(IntermediateNode)
	sb.WriteString("0") // left
	DFSHuffmanSearch(*inode.left, sb, mp)
	utils.PopLastItem(sb)
	sb.WriteString("1") // right
	DFSHuffmanSearch(*inode.right, sb, mp)
	utils.PopLastItem(sb)

	return mp
}
