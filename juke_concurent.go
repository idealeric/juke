/*
This file is part of Juke MPD client. See juke.go for more details.

This particular file has Juke's concurrent functions (goroutines).
*/

package main

import (
	"code.google.com/p/gompd/mpd"
	"fmt"
	"github.com/idealeric/juke/ui"
)

// Juke's real, current state
type jukeState uint8

const (
	NOT_CONNECTED jukeState = iota
	CONNECTED_AND_STOPPED
	CONNECTED_AND_PAUSED
	CONNECTED_AND_PLAYING
)

// A request that one can send to update()
type jukeStateRequest uint8

const (
	NEXT_TRACK jukeStateRequest = iota
	PREVIOUS_TRACK
	PLAY_OR_PAUSE
	STOP
)

// update blocks waiting for some other thread to tell it to force an update on the UI.
// An update might come from:
//	* A polling update from MPD
//	* The user has interacted with juke in some way as to
//	  force an update (button press, etc)
// The incoming communication is an attempted state change or request
// for general update.
func update(stateRequestChannel chan jukeStateRequest) {

	var currentState jukeState = NOT_CONNECTED
	mpdConnection, err := mpd.Dial("tcp", "127.0.0.1:6600")

	if err != nil {
		fmt.Println("Could not establish MPD connection.", err) // TODO - Make this better
		return // TODO - we need connectionless operation
	}

	currentState = CONNECTED_AND_PLAYING // TODO, we need to not assume this and get it "naturally"

	for requestedState := range stateRequestChannel {

		if currentState > NOT_CONNECTED {
			switch requestedState {

			case NEXT_TRACK, PREVIOUS_TRACK:

				if requestedState == PREVIOUS_TRACK {
					if err := mpdConnection.Previous(); err != nil {
						fmt.Println("bad", err) // TODO - Make this better
					}
				} else { // NEXT_TRACK
					if err := mpdConnection.Next(); err != nil {
						fmt.Println("bad", err) // TODO - Make this better
					}
				}

				if curSong, erro := mpdConnection.CurrentSong(); erro != nil {
					fmt.Println("bad", erro) // TODO - Make this better
				} else {
					ui.SetCurrentSong(curSong["Title"], curSong["Artist"], curSong["Album"])
				}

			case PLAY_OR_PAUSE:

				if currentState == CONNECTED_AND_PLAYING {
					if err := mpdConnection.Pause(true); err != nil {
						fmt.Println("bad", err) // TODO - Make this better
					} else {
						ui.SetPlayPause(false)
						currentState = CONNECTED_AND_PAUSED
					}
				} else if currentState == CONNECTED_AND_PAUSED {
					if err := mpdConnection.Pause(false); err != nil {
						fmt.Println("bad", err) // TODO - Make this better
					} else {
						ui.SetPlayPause(true)
						currentState = CONNECTED_AND_PLAYING
					}
				} else if currentState == CONNECTED_AND_STOPPED {
					if err := mpdConnection.PlayId(-1); err != nil {
						fmt.Println("bad", err) // TODO - Make this better
					} else {

						ui.SetPlayPause(true)
						currentState = CONNECTED_AND_PLAYING

						if curSong, erro := mpdConnection.CurrentSong(); erro != nil {
							fmt.Println("bad", erro) // TODO - Make this better
						} else {
							ui.SetCurrentSong(curSong["Title"], curSong["Artist"], curSong["Album"])
						}

					}
				} // end play/pause state control conditional

			case STOP:

				if err := mpdConnection.Stop(); err != nil {
					fmt.Println("bad", err) // TODO - Make this better
				} else {
					ui.SetPlayPause(false)
					ui.SetCurrentSongStopped()
					currentState = CONNECTED_AND_STOPPED
				}

			} // end requestedState switch

		} // end if not connected

	} // end for wait on channel

	// Close the MPD connection, Juke is about to end:
	if currentState > NOT_CONNECTED {
		if err := mpdConnection.Close(); err != nil {
			fmt.Println("bad", err) // TODO - Make this better
		}
	}

} // end update
