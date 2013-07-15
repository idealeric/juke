/*
This file is part of Juke MPD client. See juke.go for more details.

This particular file has Juke's concurrent functions (goroutines).
*/

package main

import (
	"code.google.com/p/gompd/mpd"
	"github.com/idealeric/juke/ui"
)

// updateCurrentSong blocks waiting for an update the current song.
func updateCurrentSong(updateChannel chan mpd.Attrs) {

	for curSong := range updateChannel {
		ui.SetCurrentSong(curSong["Title"], curSong["Artist"], curSong["Album"])
	}

} // end updateCurrentSong
