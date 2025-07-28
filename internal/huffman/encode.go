package huffman

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log/slog"
)

func (srv *Service) Encode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("Encode: input data is empty")
	}

	// Build the Huffman table
	freq := buildFreqTable(data)
	leaves := buildNodeList(freq)
	huffmanTree := build(leaves)
	table := buildHuffTable(huffmanTree)
	srv.Logger.Info("Encode: Built Huffman tree", slog.Int("symbols", len(freq)))

	// Map input to Huffman codes
	codes := make([]HuffCode, len(data))
	for i, b := range data {
		codes[i] = table[ValueType(b)]
	}
	srv.Logger.Info("Encode: Mapped input to bit patterns")

	// Serialize the freq table
	freqBytes := srv.serializeFreqTable(freq)
	// (serializeFreqTable already emits: [2-byte N][N×(1-byte sym + 4-byte freq)])

	// Pack bits into bytes
	packed, padding := srv.bitsToBytes(codes)
	srv.Logger.Info("Encode: Packed bits", slog.Int("paddingBits", padding))

	// Assemble output with one allocation
	out := make([]byte, 1+len(freqBytes)+len(packed))
	out[0] = byte(padding)
	copy(out[1:], freqBytes)
	copy(out[1+len(freqBytes):], packed)

	return out, nil
}

func (srv *Service) serializeFreqTable(freq map[ValueType]int) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, uint16(len(freq)))

	// Each symbol + its frequency (uint32)
	for sym, count := range freq {
		buf.WriteByte(byte(sym))
		_ = binary.Write(buf, binary.BigEndian, uint32(count))
	}

	return buf.Bytes()
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
