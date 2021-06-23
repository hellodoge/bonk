package messages

import (
	"encoding/binary"
	"io"
)

// https://wiki.theory.org/BitTorrentSpecification#Messages

const (
	msgChoke         = 0x00
	msgUnchoke       = 0x01
	msgInterested    = 0x02
	msgNotInterested = 0x03
	msgHave          = 0x04
	msgBitfield      = 0x05
	msgRequest       = 0x06
	msgPiece         = 0x07
	msgCancel        = 0x08
)

const (
	lengthPrefixSize = 4 // uint32
	idSize           = 1 // byte

	pieceIndexSize = 4 // uint32
)

type message struct {
	id      byte
	payload []byte
}

func (m *message) getBytes() []byte {
	var buffer = make([]byte, lengthPrefixSize+idSize+len(m.payload))
	binary.BigEndian.PutUint32(buffer, uint32(idSize+len(m.payload)))
	buffer[lengthPrefixSize] = m.id
	copy(buffer[lengthPrefixSize+idSize:], m.payload)
	return buffer
}

// Get returns one of the types defined in types.go file
func Get(r io.Reader) (interface{}, error) {
	message, err := read(r)
	if err != nil {
		return nil, err
	}
	return parse(message)
}

func read(r io.Reader) (*message, error) {
	var lengthBuffer = make([]byte, lengthPrefixSize)
	_, err := io.ReadFull(r, lengthBuffer)
	if err != nil {
		return nil, newReceivingMessageError(err)
	}
	var length = binary.BigEndian.Uint32(lengthBuffer)
	if length == 0 {
		return nil, nil // udp keep-alive
	}

	var idBuffer = make([]byte, idSize)
	_, err = io.ReadFull(r, lengthBuffer)
	if err != nil {
		return nil, newReceivingMessageError(err)
	}
	var message = new(message)
	message.id = idBuffer[0]

	if length > idSize {
		var payloadBuffer = make([]byte, length-idSize)
		_, err = io.ReadFull(r, payloadBuffer)
		if err != nil {
			return nil, newReceivingMessageError(err)
		}
		message.payload = payloadBuffer
	}

	return message, nil
}
