package tracker

import (
	"errors"
	"github.com/hellodoge/bonk/bonk"
	"github.com/hellodoge/bonk/internal/errors/common"
	"github.com/hellodoge/bonk/pkg/tracker/delivery"
	"github.com/hellodoge/bonk/pkg/tracker/parser"
	"github.com/hellodoge/bonk/pkg/tracker/response"
)

type deliveryService interface {
	RequestForPeers(downloading *bonk.DownloadingInfo, config *bonk.Config) ([]byte, error)
}

type parserService interface {
	ParsePeersResponse([]byte) (*response.RequestForPeersResponse, error)
	ParsePeersResponseCompact([]byte) (*response.RequestForPeersResponse, error)
}

type Tracker struct {
	delivery deliveryService
	parser   parserService
}

func NewTracker() *Tracker {
	return &Tracker{
		delivery: &delivery.Delivery{},
		parser:   &parser.Parser{},
	}
}

func (t *Tracker) GetPeers(downloading *bonk.DownloadingInfo, config *bonk.Config) (*response.RequestForPeersResponse, error) {
	if downloading == nil || config == nil {
		return nil, common.NewInvalidArgumentsError(errors.New("expected value, got nil"))
	}
	content, err := t.delivery.RequestForPeers(downloading, config)
	if err != nil {
		return nil, err
	}
	if config.RequestCompact {
		return t.parser.ParsePeersResponseCompact(content)
	} else {
		return t.parser.ParsePeersResponse(content)
	}
}
