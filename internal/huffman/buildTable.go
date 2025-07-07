package huffman

import "sort"

// sortNodes implements sort.Interface, order defined by Node.Count.
type sortNodes []*Node

func (sn sortNodes) Len() int           { return len(sn) }
func (sn sortNodes) Less(i, j int) bool { return sn[i].Count < sn[j].Count }
func (sn sortNodes) Swap(i, j int)      { sn[i], sn[j] = sn[j], sn[i] }

// buildNodeList builds initiate leaves of the Huffman tree
func buildNodeList(data []byte) []*Node {

	freq := make(map[ValueType]int)

	for _, el := range data {
		freq[ValueType(el)]++
	}

	leaves := make([]*Node, 0)
	for key, value := range freq {
		leaves = append(leaves, &Node{Value: key, Count: value})
	}

	return leaves

}

// build builds a Huffman tree from the specified leaves.
// The content of the passed slice is modified, if this is unwanted, pass a copy.
// Guaranteed that the same input slice will result in the same Huffman tree.
func build(leaves []*Node) *Node {
	// We sort once and use binary insertion later on
	sort.Stable(sortNodes(leaves)) // Note: stable sort for deterministic output!

	return buildSorted(leaves)
}

// buildSorted builds a Huffman tree from the specified leaves which must be sorted by Node.Count.
// The content of the passed slice is modified, if this is unwanted, pass a copy.
// Guaranteed that the same input slice will result in the same Huffman tree.
func buildSorted(leaves []*Node) *Node {
	if len(leaves) == 0 {
		return nil
	}

	for len(leaves) > 1 {
		left, right := leaves[0], leaves[1]
		parentCount := left.Count + right.Count
		parent := &Node{Left: left, Right: right, Count: parentCount}
		left.Parent = parent
		right.Parent = parent

		// Where to insert parent in order to remain sorted?
		ls := leaves[2:]
		idx := sort.Search(len(ls), func(i int) bool { return ls[i].Count >= parentCount })
		idx += 2

		// Insert
		copy(leaves[1:], leaves[2:idx])
		leaves[idx-1] = parent
		leaves = leaves[1:]
	}

	return leaves[0]
}

// buildHuffTable traverse huffman tree and fill Table
func buildHuffTable(root *Node) map[ValueType]HuffCode {
	// traverse traverses a subtree from the given node,
	// using the prefix code leading to this node, having the number of bits specified.
	table := make(map[ValueType]HuffCode)
	var traverse func(n *Node, code uint64, bits byte)

	traverse = func(n *Node, code uint64, bits byte) {
		if n.Left == nil && n.Right == nil {
			table[n.Value] = HuffCode{code, bits}
			return
		}
		bits++
		traverse(n.Left, code<<1, bits)
		traverse(n.Right, code<<1+1, bits)
	}

	traverse(root, 0, 0)
	return table
}
