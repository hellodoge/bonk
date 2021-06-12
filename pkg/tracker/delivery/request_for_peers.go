package delivery

import (
	"github.com/hellodoge/bonk/bonk"
	"github.com/hellodoge/bonk/pkg/logger"
	"github.com/hellodoge/bonk/pkg/util"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Delivery struct {
	logger logger.Logger
}

func buildURL(downloading *bonk.DownloadingInfo, config *bonk.Config) (*url.URL, error) {
	addr, err := url.Parse(downloading.Torrent.Announce)
	if err != nil {
		return nil, newParsingTrackerURLError(err)
	}
	// https://wiki.theory.org/BitTorrent_Tracker_Protocol
	params := url.Values{
		"info_hash":  []string{string(downloading.Torrent.InfoHash[:])},
		"peer_id":    []string{string(config.PeerID[:])},
		"port":       []string{strconv.Itoa(int(config.Port))},
		"uploaded":   []string{strconv.FormatUint(downloading.Uploaded, 10)},
		"downloaded": []string{strconv.FormatUint(downloading.Downloaded, 10)},
		"left":       []string{strconv.FormatUint(downloading.Left, 10)},
		"compact":    []string{strconv.Itoa(util.BoolToInt(config.RequestCompact))},
	}
	addr.RawQuery = params.Encode()
	return addr, nil
}

func (d *Delivery) RequestForPeers(downloading *bonk.DownloadingInfo, config *bonk.Config) ([]byte, error) {
	addr, err := buildURL(downloading, config)
	if err != nil {
		return nil, err
	}
	var client = &http.Client{
		Timeout: config.TrackerTimeout,
	}
	conn, err := client.Get(addr.String())
	if err != nil {
		if err, ok := err.(*url.Error); ok {
			if err.Timeout() {
				return nil, newTimeoutError(err)
			}
		}
		return nil, newMakingRequestError(err)
	}
	defer func() {
		err := conn.Body.Close()
		if d.logger != nil && err != nil {
			d.logger.Print(err)
		}
	}()
	response, err := io.ReadAll(conn.Body)
	if err != nil {
		return nil, newMakingRequestError(err)
	}
	return response, nil
}