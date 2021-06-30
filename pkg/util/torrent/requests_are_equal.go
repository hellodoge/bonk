package torrent

import "github.com/hellodoge/bonk/bonk"

func RequestsAreEqual(l bonk.Request, r bonk.Request) bool {
	return l.Index == r.Index && l.Begin == r.Begin
}
