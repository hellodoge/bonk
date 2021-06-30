package peer

import (
	"context"
	"fmt"
	"github.com/hellodoge/bonk/pkg/messages"
	"github.com/hellodoge/bonk/pkg/util/torrent"
)

func (c *Client) handleMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.handleMessage()
		}
	}
}

func (c *Client) handleMessage() {
	message, err := messages.Get(c.conn)
	if err != nil {
		c.log(err)
		return
	}
	switch message := message.(type) {
	case messages.Choke:
		c.state.chokedByPeer.Set()
	case messages.Unchoke:
		c.state.chokedByPeer.UnSet()
	case messages.Interested:
		c.state.interested.Set()
	case messages.NotInterested:
		c.state.interested.UnSet()
	case messages.Have:
		c.state.bitfieldMutex.Lock()
		c.state.bitfield.SetPiece(message.Index)
		c.state.bitfieldMutex.Unlock()
	case messages.BitField:
		c.state.bitfieldMutex.Lock()
		c.state.bitfield = message.Bitfield
		c.state.bitfieldMutex.Unlock()
	case messages.Request:
		c.state.requestedMutex.Lock()
		c.state.requested = append(c.state.requested, message.Request)
		c.state.requestedMutex.Unlock()
	case messages.Piece:
		c.state.receivedParts <- message.Piece
	case messages.Cancel:
		c.state.requestedMutex.Lock()
		for i := range c.state.requested {
			if torrent.RequestsAreEqual(c.state.requested[i], message.Request.Request) {
				c.state.requested[i].Cancelled = true
				break
			}
		}
		c.state.requestedMutex.Unlock()
	default:
		c.log(fmt.Errorf("invalid message type: %T", message))
	}
}
