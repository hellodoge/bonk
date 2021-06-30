package handshake

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/hellodoge/bonk/bonk"
	"io"
)

// https://wiki.theory.org/BitTorrentSpecification#Handshake

const (
	protocolIdentifier    = "BitTorrent protocol"
	numberOfReservedBytes = 8

	lengthOfProtocolIdentifierLength = 1 // byte
)

type Handshake struct {
	InfoHash [sha1.Size]byte
	PeerID   bonk.PeerID
}

var (
	ErrInvalidProtocolIdentifierLength = errors.New("invalid protocol identifier length")
	ErrInvalidProtocolIdentifier       = errors.New("invalid protocol")
)

func (h *Handshake) GetBytes() []byte {
	var buffer = make([]byte, lengthOfProtocolIdentifierLength+len(protocolIdentifier)+
		numberOfReservedBytes+sha1.Size+bonk.PeerIdLength)
	buffer[0] = byte(len(protocolIdentifier))
	var cursor = lengthOfProtocolIdentifierLength
	copy(buffer[cursor:], protocolIdentifier)
	cursor += len(protocolIdentifier) + numberOfReservedBytes
	copy(buffer[cursor:], h.InfoHash[:])
	cursor += sha1.Size
	copy(buffer[cursor:], h.PeerID[:])
	return buffer
}

func Read(r io.Reader) (Handshake, error) {
	var bufferProtocolIdentifierLength = make([]byte, lengthOfProtocolIdentifierLength)
	_, err := io.ReadFull(r, bufferProtocolIdentifierLength)
	if err != nil {
		return Handshake{}, newReceivingHandshakeError(err)
	}
	if int(bufferProtocolIdentifierLength[0]) != len(protocolIdentifier) {
		return Handshake{}, newParsingHandshakeError(ErrInvalidProtocolIdentifierLength)
	}

	var protocolIdentifierBuffer = make([]byte, len(protocolIdentifier))
	_, err = io.ReadFull(r, protocolIdentifierBuffer)
	if err != nil {
		return Handshake{}, newReceivingHandshakeError(err)
	}
	if string(protocolIdentifierBuffer) != protocolIdentifier {
		return Handshake{}, newParsingHandshakeError(
			fmt.Errorf("%w: expected: %s got: %s", ErrInvalidProtocolIdentifier,
				protocolIdentifier, string(protocolIdentifierBuffer)),
		)
	}

	var reservedBytesBuffer = make([]byte, numberOfReservedBytes)
	_, err = io.ReadFull(r, reservedBytesBuffer)
	if err != nil {
		return Handshake{}, newReceivingHandshakeError(err)
	}

	var handshake Handshake
	_, err = io.ReadFull(r, handshake.InfoHash[:])
	if err != nil {
		return Handshake{}, newReceivingHandshakeError(err)
	}
	_, err = io.ReadFull(r, handshake.PeerID[:])
	if err != nil {
		return Handshake{}, newReceivingHandshakeError(err)
	}
	return handshake, nil
}
