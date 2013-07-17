/*
This file is part of Juke MPD client. See juke.go for more details.

This particular file has all of Juke's callback bindings.
*/

package main

import (
	"github.com/idealeric/juke/ui"
)

func initCallBacks(updateChannel chan jukeStateRequest) {

	ui.NextClick(func() error {
		go func() {
			updateChannel <- NEXT_TRACK
		}()
		return nil
	})

	ui.PreviousClick(func() error {
		go func() {
			updateChannel <- PREVIOUS_TRACK
		}()
		return nil
	})

	ui.PlayPauseClick(func() error {
		go func() {
			updateChannel <- PLAY_OR_PAUSE
		}()
		return nil
	})

	ui.StopClick(func() error {
		go func() {
			updateChannel <- STOP
		}()
		return nil
	})

} // end initCallbacks
