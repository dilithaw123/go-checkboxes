package bitset

import "encoding/binary"

type BitSet struct {
	words []uint64
	size  uint64
}

func NewBitSet(size uint64) *BitSet {
	return &BitSet{words: make([]uint64, (size+63)/64), size: size}
}

func (b *BitSet) Set(u uint64) {
	if u >= b.size {
		return
	}
	b.words[u/64] |= 1 << uint(u%64)
}

func (b *BitSet) Clear(u uint64) {
	if u >= b.size {
		return
	}
	b.words[u/64] &^= 1 << uint(u%64)
}

func (b *BitSet) Bytes() []byte {
	bytes := make([]byte, len(b.words)*8)
	for i, word := range b.words {
		binary.LittleEndian.PutUint64(bytes[i*8:], word)
	}
	return bytes
}
