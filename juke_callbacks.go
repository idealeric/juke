/*
This file is part of Juke MPD client. See juke.go for more details.

This particular file has all of Juke's callback bindings.
*/

package main

import (
	"github.com/idealeric/juke/ui"
)

func initCallBacks(updateChannel chan *jukeRequest) {

	ui.NextClick(func() error {
		go func() {
			updateChannel <- &jukeRequest{state: NEXT_TRACK}
		}()
		return nil
	})

	ui.PreviousClick(func() error {
		go func() {
			updateChannel <- &jukeRequest{state: PREVIOUS_TRACK}
		}()
		return nil
	})

	ui.PlayPauseClick(func() error {
		go func() {
			updateChannel <- &jukeRequest{state: PLAY_OR_PAUSE}
		}()
		return nil
	})

	ui.StopClick(func() error {
		go func() {
			updateChannel <- &jukeRequest{state: STOP}
		}()
		return nil
	})

	ui.ProgressBarClick(func(x int, width int) error {
		go func() {
			updateChannel <- &jukeRequest{state: PROGRESS_CHANGE, progressX: x, progressWidth: width}
		}()
		return nil
	})

} // end initCallbacks
