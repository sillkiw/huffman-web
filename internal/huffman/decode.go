package huffman

import (
	"encoding/binary"
	"fmt"
	"log/slog"
)

func (srv *Service) Decode(data []byte) ([]byte, error) {
	if len(data) < 3 {
		return nil, fmt.Errorf("Decode: data too short")
	}

	padding := int(data[0])

	freqTable, offset, err := srv.unserializeFreqTable(data[:1])
	if err != nil {
		return nil, fmt.Errorf("Decode: failed to read freq table: %w", err)
	}

	payload := data[1+offset:]
	if len(payload) == 0 {
		return nil, fmt.Errorf("Decode: no compressed data found")
	}
	bits := srv.bytesToBits(payload)
	bits = bits[:len(bits)-padding]

	// Rebuild the Huffman tree and decode the bit‑stream
	leaves := buildNodeList(freqTable)
	root := build(leaves)
	decoded := decodeBits(root, bits)

	return decoded, nil
}

// unserializeFreqTable reads:
//
//	[2 bytes: N entries]
//	N × [1 byte sym][4 bytes freq]
//
// from the front of `data`, and returns how many bytes it consumed.
func (srv *Service) unserializeFreqTable(data []byte) (map[ValueType]int, int, error) {
	// if len(data) < 2 {
	// 	return nil, 0, fmt.Errorf("freq data too short")
	// }
	freqLen := int(binary.BigEndian.Uint16(data[0:2]))
	freqTable := make(map[ValueType]int)
	srv.Logger.Info("freqLen", slog.Int("freqLen", freqLen))
	offset := 2
	for i := 0; i < freqLen; i++ {
		sym := data[offset]
		cnt := binary.BigEndian.Uint32(data[offset+1 : offset+5])
		freqTable[ValueType(sym)] = int(cnt)
		offset += 5
	}
	return freqTable, offset, nil
}

// bytesToBits unpacks each byte in data into 8 bits (MSB first).
// e.g. 0xA2 → [1,0,1,0,0,0,1,0]
func (srv *Service) bytesToBits(data []byte) []bool {
	bits := make([]bool, len(data)*8)
	idx := 0
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bits[idx] = ((b >> uint(i)) & 1) == 1
			idx++
		}
	}
	return bits
}

// decodeBits walks the Huffman tree for each bit, emitting a symbol
// when it reaches a leaf, then restarting at the root.
func decodeBits(root *Node, bits []bool) []byte {
	var out []byte
	node := root
	for _, bit := range bits {
		if bit {
			node = node.Right
		} else {
			node = node.Left
		}
		// leaf if both children are nil
		if node.Left == nil && node.Right == nil {
			out = append(out, byte(node.Value))
			node = root
		}
	}
	return out
}
