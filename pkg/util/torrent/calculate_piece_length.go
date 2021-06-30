package torrent

func CalculatePieceLength(length uint64, pieceLength uint, index uint32) uint {
	var incompleteLength = uint(length % uint64(pieceLength))
	if incompleteLength == 0 {
		return pieceLength
	} else if uint64(index) == length/uint64(pieceLength) {
		return incompleteLength
	} else {
		return pieceLength
	}
}
