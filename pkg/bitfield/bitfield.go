package bitfield

type BitField struct {
	bitfield []byte
}

func (m *BitField) HasPiece(index uint32) bool {
	if m.bitfield == nil {
		return false
	}
	var byteIndex = int(index / 8)
	if byteIndex >= len(m.bitfield) {
		return false
	}
	var searchMask byte = 0x1 << (7 - index%8)
	return m.bitfield[byteIndex]&searchMask != 0
}

func (m *BitField) SetPiece(index uint32) {
	if m.bitfield == nil || len(m.bitfield) < int(index/8+1) {
		m.SetMaxIndex(index)
	}
	var byteIndex = int(index / 8)
	var pieceMask byte = 0x1 << (7 - index%8)
	m.bitfield[byteIndex] |= pieceMask
}

func (m *BitField) SetMaxIndex(index uint32) {
	if m.bitfield == nil {
		m.bitfield = make([]byte, index/8+1)
	} else if len(m.bitfield) < int(index/8+1) {
		// TODO find the way to set a capacity of a slice
		var n = len(m.bitfield) + int(index/8+1)
		m.bitfield = append(m.bitfield, make([]byte, n)...)
	}
}

func NewBitField(bitfield []byte) BitField {
	return BitField{
		bitfield: bitfield,
	}
}
