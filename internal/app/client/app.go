package client

import (
	"fmt"
	"github.com/hellodoge/bonk/bonk"
	"github.com/hellodoge/bonk/pkg/torrent_file/parser"
	"github.com/hellodoge/bonk/pkg/tracker"
	"github.com/hellodoge/bonk/pkg/util/torrent"
	"log"
	"os"
	"time"
)

const (
	DefaultPeerId = "bonk_default_peer_id"
	DefaultTrackerTimeout = time.Second * 15
)

func Run() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [.torrent file]", os.Args[0])
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	torrentInfo, err := parser.ParseTorrentFile(file)
	if err != nil {
		log.Fatal(err)
	}

	downloadInfo := bonk.NewDownloadingInfo(torrentInfo)

	peerID, err := torrent.GenerateRandomPeerID()
	if err != nil {
		log.Printf("cannot generate random peer ID (%s), using default", err)
		copy(peerID[:], DefaultPeerId)
	}

	config := bonk.Config{
		PeerID: peerID,
		TrackerTimeout: DefaultTrackerTimeout,
		RequestCompact: true,
	}

	trackerService := tracker.NewTracker()

	response, err := trackerService.GetPeers(downloadInfo, &config)
	if err != nil {
		log.Print(err)
	}

	fmt.Print(response)
}
