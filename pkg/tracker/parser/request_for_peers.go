package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/hellodoge/bonk/bonk"
	"github.com/hellodoge/bonk/pkg/tracker/response"
	bencode "github.com/jackpal/bencode-go"
	"net"
	"time"
)

type Parser struct {}

type trackerResponse struct {
	Interval uint       `json:"interval"`
	Peers    []peerInfo `json:"peers"`
}

type trackerResponseCompact struct {
	Interval uint   `json:"interval"`
	Peers    string `json:"peers"`
}

type peerInfo struct {
	ID   string `bencode:"id"`
	IP   string `bencode:"ip"`
	Port uint16 `bencode:"port"`
}

var (
	ErrNotValidIPFormat   = errors.New("invalid ip address format")
	ErrInvalidPeersLength = errors.New("invalid length of peers")
)

func (p *Parser) ParsePeersResponse(body []byte) (*response.RequestForPeersResponse, error) {
	var decoded trackerResponse
	err := bencode.Unmarshal(bytes.NewBuffer(body), &decoded)
	if err != nil {
		return nil, newParsingResponseError(err)
	}
	return decoded.toRequestForPeersResponse()
}

func (p *Parser) ParsePeersResponseCompact(body []byte) (*response.RequestForPeersResponse, error) {
	var decoded trackerResponseCompact
	err := bencode.Unmarshal(bytes.NewBuffer(body), &decoded)
	if err != nil {
		return nil, newParsingResponseError(err)
	}
	return decoded.toRequestForPeersResponse()
}

func (r *trackerResponse) toRequestForPeersResponse() (*response.RequestForPeersResponse, error) {
	var peers = make([]bonk.Peer, 0, len(r.Peers))
	for _, peer := range r.Peers {
		var addr = net.TCPAddr{
			IP:   net.ParseIP(peer.IP),
			Port: int(peer.Port),
		}
		if addr.IP == nil {
			return nil, newParsingResponseError(ErrNotValidIPFormat)
		}
		peers = append(peers, bonk.Peer{
			ID:   peer.ID,
			Addr: addr,
		})
	}
	return &response.RequestForPeersResponse{
		Interval: intervalToTimeDuration(r.Interval),
		Peers:    peers,
	}, nil
}

func (r *trackerResponseCompact) toRequestForPeersResponse() (*response.RequestForPeersResponse, error) {
	var rawPeers = []byte(r.Peers)
	const portSize = 2 // uint16
	const peerSize = net.IPv4len + portSize
	if len(rawPeers)%peerSize != 0 {
		return nil, newParsingResponseError(ErrInvalidPeersLength)
	}
	var peers = make([]bonk.Peer, 0, len(rawPeers)/peerSize)
	for offset := 0; offset < len(rawPeers); offset += peerSize {
		peers = append(peers, bonk.Peer{
			Addr: net.TCPAddr{
				IP: rawPeers[offset : offset+net.IPv4len],
				Port: int(binary.BigEndian.Uint16(rawPeers[offset+net.IPv4len : offset+net.IPv4len+portSize])),
			},
		})
	}
	return &response.RequestForPeersResponse{
		Interval: intervalToTimeDuration(r.Interval),
		Peers:    peers,
	}, nil
}

func intervalToTimeDuration(interval uint) time.Duration {
	return time.Second * time.Duration(interval)
}
