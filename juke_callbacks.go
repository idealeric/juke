/*
This file is part of Juke MPD client. See juke.go for more details.

This particular file has all of Juke's callback bindings.
*/

package main

import (
	"code.google.com/p/gompd/mpd"
	"fmt"
	"github.com/idealeric/juke/ui"
)

func initCallBacks(songChan chan mpd.Attrs, mpdConnection *mpd.Client) {

	ui.OnExit(func() error {

		// Close any channels so that we clean up any running goroutines.
		close(songChan)

		if currentState > NOT_CONNECTED {
			if err := mpdConnection.Close(); err != nil {
				return err
			}
		}

		return nil

	})

	ui.NextClick(func() error {

		if currentState > NOT_CONNECTED {

			if err := mpdConnection.Next(); err != nil {
				return err
			}

			if curSong, erro := mpdConnection.CurrentSong(); erro != nil {
				fmt.Println("bad", erro) // TODO - Make this better
			} else {
				songChan <- curSong
			}

		}

		return nil

	})

	ui.PreviousClick(func() error {

		if currentState > NOT_CONNECTED {

			if err := mpdConnection.Previous(); err != nil {
				return err
			}

			if curSong, erro := mpdConnection.CurrentSong(); erro != nil {
				fmt.Println("bad", erro) // TODO - Make this better
			} else {
				songChan <- curSong
			}

		}

		return nil

	})

	ui.PlayPauseClick(func() error {

		if currentState == CONNECTED_AND_PLAYING {
			if err := mpdConnection.Pause(true); err != nil {
				return err
			} else {
				ui.SetPlayPause(false)
				currentState = CONNECTED_AND_PAUSED
			}
		} else if currentState == CONNECTED_AND_PAUSED {
			if err := mpdConnection.Pause(false); err != nil {
				return err
			} else {
				ui.SetPlayPause(true)
				currentState = CONNECTED_AND_PLAYING
			}
		} else if currentState == CONNECTED_AND_STOPPED {
			if err := mpdConnection.PlayId(-1); err != nil {
				return err
			} else {

				ui.SetPlayPause(true)
				currentState = CONNECTED_AND_PLAYING

				if curSong, erro := mpdConnection.CurrentSong(); erro != nil {
					fmt.Println("bad", err) // TODO - Make this better
				} else {
					songChan <- curSong
				}

			}
		} // end state control conditional

		return nil

	})

	ui.StopClick(func() error {

		if currentState > NOT_CONNECTED {
			if err := mpdConnection.Stop(); err != nil {
				return err
			} else {
				ui.SetPlayPause(false)
				currentState = CONNECTED_AND_STOPPED
			}

			ui.SetCurrentSongStopped()

		}

		return nil

	})

} // end initCallbacks
