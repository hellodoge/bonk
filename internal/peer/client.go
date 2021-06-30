package peer

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/hellodoge/bonk/bonk"
	"github.com/hellodoge/bonk/internal/controller"
	"github.com/hellodoge/bonk/internal/errors/paired"
	"github.com/hellodoge/bonk/internal/logger"
	"github.com/hellodoge/bonk/pkg/bitfield"
	"github.com/hellodoge/bonk/pkg/handshake"
	"github.com/hellodoge/bonk/pkg/parts"
	"github.com/hellodoge/bonk/pkg/util/torrent"
	"github.com/tevino/abool"
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	state           State
	peer            bonk.Peer
	conn            net.Conn
	config          *bonk.Config
	logger          logger.Logger
	taskQueue       controller.TaskQueue
	receivedQueue   controller.ReceivedQueue
	downloadingInfo *bonk.DownloadingInfo
}

type State struct {
	currentTask    *bonk.PieceInfo
	chokedByPeer   *abool.AtomicBool
	chokingPeer    *abool.AtomicBool
	interested     *abool.AtomicBool
	bitfield       bitfield.BitField
	bitfieldMutex  sync.Mutex
	requested      []bonk.Request
	requestedMutex sync.Mutex
	receivedParts  chan bonk.Part
	receivedPieces map[uint32]*parts.Assembler
	requestedAt    time.Time
	retries        uint
}

var (
	ErrInfoHashDoesNotMatch = errors.New("info hash does not match")
)

func (c *Client) connect(ctx context.Context) error {
	var dialer = net.Dialer{
		Timeout: c.config.PeerTimeout,
	}
	conn, err := dialer.DialContext(ctx, "tcp", c.peer.Addr.String())
	if err != nil {
		return newConnectingToPeerError(err)
	}
	var requestHandshake = handshake.Handshake{
		InfoHash: c.downloadingInfo.Torrent.InfoHash,
		PeerID:   c.config.PeerID,
	}
	_, err = conn.Write(requestHandshake.GetBytes())
	if err != nil {
		err = newWritingToPeerError(err)
		return paired.NewPairedError(err, conn.Close())
	}
	responseHandshake, err := handshake.Read(conn)
	if err != nil {
		err = newGettingHandshakeError(err)
		return paired.NewPairedError(err, conn.Close())
	}
	if responseHandshake.InfoHash != c.downloadingInfo.Torrent.InfoHash {
		return fmt.Errorf("%w: expected %s, got %s", ErrInfoHashDoesNotMatch,
			c.downloadingInfo.Torrent.InfoHash, requestHandshake.InfoHash,
		)
	}
	c.peer.ID = responseHandshake.PeerID
	c.conn = conn
	return nil
}

func NewClient(queue controller.TaskQueue, received controller.ReceivedQueue,
	config *bonk.Config, downloading *bonk.DownloadingInfo) *Client {
	return &Client{
		state: State{
			chokedByPeer:   abool.NewBool(true),
			chokingPeer:    abool.NewBool(true),
			interested:     abool.NewBool(false),
			bitfield:       bitfield.BitField{},
			receivedParts:  make(chan bonk.Part),
			receivedPieces: make(map[uint32]*parts.Assembler),
		},
		peer:            bonk.Peer{},
		config:          config,
		logger:          &log.Logger{},
		taskQueue:       queue,
		receivedQueue:   received,
		downloadingInfo: downloading,
	}
}

func (c *Client) Run(ctx context.Context) error {
	if c.conn == nil {
		err := c.connect(ctx)
		if err != nil {
			return err
		}
		c.log("connected")
	}

	c.state.chokedByPeer = abool.NewBool(true)
	c.state.interested = abool.NewBool(false)
	c.state.chokingPeer = abool.NewBool(true)
	c.state.receivedParts = make(chan bonk.Part)

	go c.handleMessages(ctx)

	for {
		select {
		case <-ctx.Done():
			break
		default:
			err := c.manage(ctx)
			if err != nil {
				c.log(err)
			}
		}
	}
}

func (c *Client) manage(ctx context.Context) error {
	if c.state.currentTask == nil {
		c.state.currentTask = <-c.taskQueue
	}
	select {
	case part := <-c.state.receivedParts:
		c.processPart(part)
	default:
		if c.state.bitfield.HasPiece(c.state.currentTask.Index) {
			// TODO
		}
	}
	return nil
}

func (c *Client) processPart(part bonk.Part) {
	assembler, ok := c.state.receivedPieces[part.Index]
	if !ok {
		assembler = parts.NewAssembler(torrent.CalculatePieceLength(
			c.downloadingInfo.Torrent.Length, c.downloadingInfo.Torrent.PieceLength, part.Index),
		)
		c.state.receivedPieces[part.Index] = assembler
	}
	assembler.AddPart(part)
	piece := assembler.TryToAssemble()
	if piece != nil {
		if sha1.Sum(piece) == c.state.currentTask.Hash {
			c.receivedQueue <- bonk.Piece{
				Index:   part.Index,
				Content: piece,
			}
		} else {
			c.log("hash did not matched")
			c.taskQueue <- c.state.currentTask
		}
		delete(c.state.receivedPieces, part.Index)
		c.state.currentTask = nil
	}
}

func (c *Client) log(args ...interface{}) {
	if c.logger != nil {
		c.logger.Print(c.peer.Addr.String()+" ", args)
	}
}
