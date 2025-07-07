package huffman

// Default value
type ValueType int32

// Node in the Huffman tree
type Node struct {
	Parent *Node     // Optional parent node, for fast code read-out
	Left   *Node     // Optional left node
	Right  *Node     // Optional right node
	Count  int       // Relative frequency
	Value  ValueType // Optional value, set if this is a leaf
}

// Binary code
type HuffCode struct {
	r    uint64
	bits byte
}
