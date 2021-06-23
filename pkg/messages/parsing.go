package messages

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/hellodoge/bonk/bonk"
	wrapped2 "github.com/hellodoge/bonk/internal/errors/wrapped"
	"github.com/hellodoge/bonk/pkg/bitfield"
)

var (
	ErrInvalidMessageID = errors.New("invalid message id")

	ErrParsingHave    = errors.New("error parsing have message")
	ErrParsingRequest = errors.New("error parsing request message")
	ErrParsingPiece   = errors.New("error parsing piece")

	ErrNoPayload          = errors.New("no payload")
	ErrInvalidPieceSize   = errors.New("invalid piece size")
	ErrInvalidPayloadSize = errors.New("invalid payload size")
)

func parse(m *message) (interface{}, error) {
	switch m.id {
	case msgChoke:
		return Choke{}, nil
	case msgUnchoke:
		return Unchoke{}, nil
	case msgInterested:
		return Interested{}, nil
	case msgNotInterested:
		return NotInterested{}, nil
	case msgHave:
		return parseHave(m.payload)
	case msgBitfield:
		return BitField{
			Bitfield: bitfield.NewBitField(m.payload),
		}, nil
	case msgRequest:
		return parseRequest(m.payload)
	case msgPiece:
		return parsePiece(m.payload)
	case msgCancel:
		request, err := parseRequest(m.payload)
		return Cancel{Request: request}, err
	default:
		return nil, newParsingMessageError(
			fmt.Errorf("%w: %d", ErrInvalidMessageID, m.id),
		)
	}
}

func parseHave(payload []byte) (Have, error) {
	if payload == nil {
		return Have{}, newParsingMessageError(
			wrapped2.Error{
				Outer: ErrParsingHave,
				Inner: ErrNoPayload,
			},
		)
	}
	if len(payload) != pieceIndexSize {
		return Have{}, newParsingMessageError(
			wrapped2.Error{
				Outer: ErrParsingHave,
				Inner: ErrInvalidPieceSize,
			},
		)
	}
	return Have{
		Index: binary.BigEndian.Uint32(payload),
	}, nil
}

func parseRequest(payload []byte) (Request, error) {
	if payload == nil {
		return Request{}, newParsingMessageError(
			wrapped2.Error{
				Outer: ErrParsingRequest,
				Inner: ErrNoPayload,
			},
		)
	}
	const uint32size = 4
	const payloadSize = uint32size * 3
	if len(payload) != payloadSize {
		return Request{}, newParsingMessageError(
			wrapped2.Error{
				Outer: ErrParsingRequest,
				Inner: ErrInvalidPayloadSize,
			},
		)
	}
	const beginOffset = uint32size
	const lengthOffset = uint32size * 2
	return Request{
		Request: bonk.Request{
			Index:  binary.BigEndian.Uint32(payload),
			Begin:  binary.BigEndian.Uint32(payload[beginOffset:]),
			Length: binary.BigEndian.Uint32(payload[lengthOffset:]),
		},
	}, nil
}

func parsePiece(payload []byte) (Piece, error) {
	if payload == nil {
		return Piece{}, newParsingMessageError(
			wrapped2.Error{
				Outer: ErrParsingPiece,
				Inner: ErrNoPayload,
			},
		)
	}
	const uint32size = 4
	const minPayloadSize = uint32size * 2
	if len(payload) < minPayloadSize {
		return Piece{}, newParsingMessageError(
			wrapped2.Error{
				Outer: ErrParsingPiece,
				Inner: ErrNoPayload,
			},
		)
	}
	const beginOffset = uint32size
	const blockOffset = uint32size * 2
	return Piece{
		Piece: bonk.Part{
			Index:  binary.BigEndian.Uint32(payload),
			Offset: binary.BigEndian.Uint32(payload[beginOffset:]),
			Block:  payload[:blockOffset],
		},
	}, nil
}
