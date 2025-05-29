package main

import (
	"os"
)

func main() {

	data, err := os.ReadFile("example.txt")

	if err != nil {
		panic(err)
	}

	leaves := buildNodeList(data)
	huffman_tree := Build(leaves)

	Table = make(map[ValueType]HuffCode)
	BuildHuffTable(huffman_tree)

	Encode(data)

	Decode()

}
