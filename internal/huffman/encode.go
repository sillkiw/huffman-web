package huffman

// Encode takse data from the input file and encode it using Huffman table
func (srv *Service) Encode(data []byte) []byte {
	// Result
	var bitStream []HuffCode
	// Build huffman table for encoding
	leaves := buildNodeList(data)
	huffmanTree := build(leaves)
	table := buildHuffTable(huffmanTree)
	srv.Logger.Info("Build Huffman Tree")

	// Encoding
	for _, el := range data {
		bitStream = append(bitStream, table[ValueType(el)])
	}
	srv.Logger.Info("Data successfully encoded")

	packed, padding := srv.bitsToBytes(bitStream)

	// Return packed data with padding
	// TODO: serialize freq table
	return append([]byte{byte(padding)}, packed...)
}

// BitsToBytes packs bitStream into seq of bytes
func (srv *Service) bitsToBytes(bitStream []HuffCode) ([]byte, int) {
	var result []byte
	var buffer uint64 = 0
	var bitCount byte = 0

	for _, c := range bitStream {
		buffer = (buffer << c.bits) | c.r
		bitCount += c.bits
		for bitCount >= 8 {
			shift := bitCount - 8
			b := byte(buffer >> shift)
			result = append(result, b)
			buffer &= (1 << shift) - 1
			bitCount -= 8
		}
	}

	padding := 0
	if bitCount > 0 {
		padding = int(8 - bitCount)
		buffer <<= padding
		result = append(result, byte(buffer))
	}
	return result, padding
}

// WriteToFile writes seq of bytes packed from encoded data to a file
// func (srv *Service) writeToFile(bitStream []HuffCode) {
// 	packed, padding := srv.bitsToBytes(bitStream)
// 	err := os.WriteFile("output.bin", append([]byte{byte(padding)}, packed...), 0644)
// 	if err != nil {

// 	}
// 	fmt.Println("Сохранено в output.bin")
// }
