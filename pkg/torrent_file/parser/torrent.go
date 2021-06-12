package parser

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"github.com/hellodoge/bonk/bonk"
	"github.com/jackpal/bencode-go"
	"io"
)

type torrent struct {
	Announce     string `bencode:"announce"`
	Comment      string `bencode:"comment"`
	CreationDate int64  `bencode:"creation date"`
	Info         struct {
		Length      uint64 `bencode:"length"`
		Name        string `bencode:"name"`
		PieceLength uint   `bencode:"piece length"`
		Pieces      string `bencode:"pieces"`
	} `bencode:"info"`
	buffer []byte // raw bencode; used by torrent.getInfoHash
}

var (
	ErrNoInfo        = errors.New("info not found in torrent file")
	ErrInvalidFormat = errors.New("torrent file has incorrect format")
)

func ParseTorrentFile(r io.Reader) (*bonk.TorrentInfo, error) {
	var torrent torrent
	buffer, err := io.ReadAll(r)
	if err != nil {
		return nil, newParsingError(err)
	}
	torrent.buffer = buffer
	err = bencode.Unmarshal(bytes.NewBuffer(buffer), &torrent)
	if err != nil {
		return nil, newParsingError(err)
	}
	return torrent.getTorrentInfo()
}

func (t *torrent) getTorrentInfo() (*bonk.TorrentInfo, error) {
	var hashes = make([][sha1.Size]byte, len(t.Info.Pieces)/sha1.Size)
	for i := range hashes {
		if copy(hashes[i][:], []byte(t.Info.Pieces)[i*sha1.Size:]) != sha1.Size {
			return nil, ErrInvalidPiecesHashesLength
		}
	}
	infoHash, err := t.getInfoHash()
	return &bonk.TorrentInfo{
		Announce:     t.Announce,
		InfoHash:     infoHash,
		Name:         t.Info.Name,
		Length:       t.Info.Length,
		PieceLength:  t.Info.PieceLength,
		PiecesHashes: hashes,
	}, err
}

func (t *torrent) getInfoHash() ([sha1.Size]byte, error) {
	d, err := bencode.Decode(bytes.NewBuffer(t.buffer))
	if err != nil {
		return [sha1.Size]byte{}, err
	}
	data, ok := d.(map[string]interface{})
	if !ok {
		return [sha1.Size]byte{}, ErrInvalidFormat
	}
	if _, ok := data["info"]; !ok {
		return [sha1.Size]byte{}, newParsingError(ErrNoInfo)
	}
	var buffer bytes.Buffer
	err = bencode.Marshal(&buffer, data["info"])
	if err != nil {
		return [sha1.Size]byte{}, err
	}
	return sha1.Sum(buffer.Bytes()), nil
}
