package bonk

import "crypto/sha1"

type TorrentInfo struct {
	Announce     string
	InfoHash     [sha1.Size]byte
	Name         string
	Length       uint64
	PieceLength  uint
	PiecesHashes [][sha1.Size]byte
}

type DownloadingInfo struct {
	Torrent    *TorrentInfo
	Uploaded   uint64
	Downloaded uint64
	Left       uint64
}

func NewDownloadingInfo(torrent *TorrentInfo) *DownloadingInfo {
	return &DownloadingInfo{
		Torrent:    torrent,
		Uploaded:   0,
		Downloaded: 0,
		Left:       torrent.Length,
	}
}