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

	ui.ConnectionClick(func() error {
		go func() {
			updateChannel <- &jukeRequest{state: CONNECTION_REFREASH}
		}()
		return nil
	})

	ui.CurrentRowDoubleClick(func(row *ui.CurrentPLRow) error {
		go func() {
			updateChannel <- &jukeRequest{state: CHANGE_TRACK, clickedRow: row}
		}()
		return nil
	})

	ui.CurrentColumnClick(func(rc chan *ui.CurrentPLRow) error {
		go func() {
			updateChannel <- &jukeRequest{state: SORT_PLAYLIST, playlistChan: rc}
		}()
		return nil
	})

} // end initCallbacks
