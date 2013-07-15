/*
Juke is a front-end, GTK+ client for the Music Playing Deamon.

Copyright: Eric Butler 2013
Version:   0.1a

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"code.google.com/p/gompd/mpd"
	"fmt"
	"github.com/idealeric/juke/ui"
)

type jukeState uint8

const (
	NOT_CONNECTED jukeState = iota
	CONNECTED_AND_STOPPED
	CONNECTED_AND_PAUSED
	CONNECTED_AND_PLAYING
)

var (
	currentState jukeState = NOT_CONNECTED
)

func main() {

	var songChan chan mpd.Attrs = make(chan mpd.Attrs)

	ui.InitInterface()
	go updateCurrentSong(songChan)

	mpdConnection, err := mpd.Dial("tcp", "127.0.0.1:6600")
	if err != nil {
		fmt.Println("bad", err) // TODO - Make this better
	} else {
		ui.SetPlayPause(true)
		currentState = CONNECTED_AND_PLAYING
		if curSong, erro := mpdConnection.CurrentSong(); erro != nil {
			fmt.Println("bad", err) // TODO - Make this better
		} else {
			songChan <- curSong
		}

		// For code tidyness, callbacks are defined in a seperate file.
		initCallBacks(songChan, mpdConnection)

	}

	ui.MainLoop() // This blocks until the GUI is destoryed.

} // end main
