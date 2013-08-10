/*
This file is part of Juke MPD client. See juke.go for more details.

This particular file has Juke's concurrent functions (goroutines).
*/

package main

import (
	"code.google.com/p/gompd/mpd"
	"github.com/idealeric/juke/log"
	"github.com/idealeric/juke/ui"
	"strconv"
	"strings"
	"time"
)

// Juke's real, current state
type jukeState uint8

const (
	NOT_CONNECTED jukeState = iota
	CONNECTED_AND_UNKNOWN
	CONNECTED_AND_STOPPED
	CONNECTED_AND_PAUSED
	CONNECTED_AND_PLAYING
)

// A request that one can send to update()
type jukeStateRequest uint8

const (
	POLL_REFREASH jukeStateRequest = iota
	NEXT_TRACK
	PREVIOUS_TRACK
	PLAY_OR_PAUSE
	STOP
	PROGRESS_CHANGE
)

// update() accepts jukeRequests, which consist in a jukeStateRequest and any other
// information that may be required. Members will be added as needed.
type jukeRequest struct {
	state         jukeStateRequest // request type
	progressX     int              // x value of the PROGRESS_CHANGE event request
	progressWidth int              // width progressbar on PROGRESS_CHANGE request
}

// Variable rate at which juke will poll MPD, in ms
const (
	END_POLLING     = 0
	PLAYING_POLLING = 500
	PAUSED_POLLING  = 750
	STOPPED_POLLING = 1000
)

// update blocks waiting for some other thread to tell it to force an update on the UI.
// An update might come from:
//	* A polling update from MPD
//	* The user has interacted with juke in some way as to
//	  force an update (button press, etc)
// The incoming communication is an attempted state change or request
// for general update.
func update(stateRequestChannel chan *jukeRequest, pollChannel chan int) {

	var currentState jukeState = NOT_CONNECTED
	mpdConnection, errDial := mpd.Dial("tcp", "127.0.0.1:6600")

	if errDial != nil {
		// TODO - we need connectionless operation
		log.ErrorOut("update()", "Could not establish MPD connection ("+errDial.Error()+").")
	}

	// On successful connection, init polling.
	go poll(stateRequestChannel, pollChannel)
	// This is a bogus state until we can determine our real state from our first polling.
	currentState = CONNECTED_AND_UNKNOWN

	for request := range stateRequestChannel {

		ui.Lock()

		if currentState > NOT_CONNECTED {
			switch request.state {

			case POLL_REFREASH:

				if status, errStatus := mpdConnection.Status(); errStatus != nil {
					log.ErrorReport("update() POLL_REFREASH", "Could not establish MPD status ("+errStatus.Error()+").")
				} else if status["state"] == "stop" {
					ui.SetPlayPause(false)
					ui.SetCurrentSongStopped()
					ui.SetCurrentAlbumArt(ui.NO_COVER_ARTWORK)
					ui.SetProgressBarTimeStoppedOrDisconnected()
					currentState = CONNECTED_AND_STOPPED
					pollChannel <- STOPPED_POLLING
				} else {

					if status["state"] == "pause" {
						ui.SetPlayPause(false)
						currentState = CONNECTED_AND_PAUSED
						pollChannel <- PAUSED_POLLING
					} else if status["state"] == "play" {
						ui.SetPlayPause(true)
						currentState = CONNECTED_AND_PLAYING
						pollChannel <- PLAYING_POLLING
					}

					// In cases of both pause and play, update the currrent song.
					if curSong, errCurSong := mpdConnection.CurrentSong(); errCurSong != nil {
						log.ErrorReport("update() POLL_REFREASH", "Could not establish MPD current song ("+errCurSong.Error()+").")
					} else {
						ui.SetCurrentSong(curSong["Title"], curSong["Artist"], curSong["Album"])
						ui.SetCurrentAlbumArt(albumArtFilename(curSong["file"]))
						totalTime, errTotalTime := strconv.Atoi(curSong["Time"])
						curTime, errCurTime := strconv.Atoi(strings.SplitN(status["time"], ":", 2)[0])
						if errTotalTime != nil {
							log.ErrorReport("update() POLL_REFREASH", "Could not convert current song total time ("+errTotalTime.Error()+").")
						} else if errCurTime != nil {
							log.ErrorReport("update() POLL_REFREASH", "Could not convert current song time ("+errCurTime.Error()+").")
						} else {
							ui.SetProgressBarTime(curTime, totalTime)
						}
					}

				}

			case NEXT_TRACK, PREVIOUS_TRACK:

				if currentState > CONNECTED_AND_STOPPED {
					if request.state == PREVIOUS_TRACK {
						if errPrev := mpdConnection.Previous(); errPrev != nil {
							log.ErrorReport("update() PREVIOUS_TRACK", "Could not mpd.Previous() ("+errPrev.Error()+").")
						}
					} else { // NEXT_TRACK
						if errNext := mpdConnection.Next(); errNext != nil {
							log.ErrorReport("update() NEXT_TRACK", "Could not mpd.Next() ("+errNext.Error()+").")
						}
					}

					if curSong, errCurSong := mpdConnection.CurrentSong(); errCurSong != nil {
						log.ErrorReport("update() NEXT/PREV_TRACK", "Could not establish current song ("+errCurSong.Error()+").")
					} else {
						ui.SetCurrentSong(curSong["Title"], curSong["Artist"], curSong["Album"])
						ui.SetCurrentAlbumArt(albumArtFilename(curSong["file"]))
						if totalTime, errTotalTime := strconv.Atoi(curSong["Time"]); errTotalTime != nil {
							log.ErrorReport("update() NEXT/PREV_TRACK", "Could not convert current song total time ("+errTotalTime.Error()+").")
						} else {
							ui.SetProgressBarTime(0, totalTime)
						}
					}
				}

			case PLAY_OR_PAUSE:

				if currentState == CONNECTED_AND_PLAYING {
					if errPause := mpdConnection.Pause(true); errPause != nil {
						log.ErrorReport("update() PLAY_OR_PAUSE", "Could not mpd.Pause(true) ("+errPause.Error()+").")
					} else {
						ui.SetPlayPause(false)
						currentState = CONNECTED_AND_PAUSED
					}
				} else if currentState == CONNECTED_AND_PAUSED {
					if errPause := mpdConnection.Pause(false); errPause != nil {
						log.ErrorReport("update() PLAY_OR_PAUSE", "Could not mpd.Pause(false) ("+errPause.Error()+").")
					} else {
						ui.SetPlayPause(true)
						currentState = CONNECTED_AND_PLAYING
					}
				} else if currentState == CONNECTED_AND_STOPPED {
					if errReplay := mpdConnection.PlayId(-1); errReplay != nil {
						log.ErrorReport("update() PLAY_OR_PAUSE", "Could not mpd.PlayId(-1) ("+errReplay.Error()+").")
					} else {

						ui.SetPlayPause(true)
						currentState = CONNECTED_AND_PLAYING

						if curSong, errCurSong := mpdConnection.CurrentSong(); errCurSong != nil {
							log.ErrorReport("update() PLAY_OR_PAUSE", "Could not establish current song ("+errCurSong.Error()+").")
						} else {
							ui.SetCurrentSong(curSong["Title"], curSong["Artist"], curSong["Album"])
							ui.SetCurrentAlbumArt(albumArtFilename(curSong["file"]))
							if totalTime, errTotalTime := strconv.Atoi(curSong["Time"]); errTotalTime != nil {
								log.ErrorReport("update() PLAY_OR_PAUSE", "Could not convert current song total time ("+errTotalTime.Error()+").")
							} else {
								ui.SetProgressBarTime(0, totalTime)
							}
						}

					}
				} // end play/pause state control conditional

			case STOP:

				if errStop := mpdConnection.Stop(); errStop != nil {
					log.ErrorReport("update() STOP", "Could not mpd.Stop() ("+errStop.Error()+").")
				} else {
					ui.SetPlayPause(false)
					ui.SetCurrentSongStopped()
					ui.SetCurrentAlbumArt(ui.NO_COVER_ARTWORK)
					ui.SetProgressBarTimeStoppedOrDisconnected()
					currentState = CONNECTED_AND_STOPPED
				}

			case PROGRESS_CHANGE:

				if currentState > CONNECTED_AND_STOPPED {

					if status, errStatus := mpdConnection.Status(); errStatus != nil {
						log.ErrorReport("update() PROGRESS_CHANGE", "Could not establish MPD status ("+errStatus.Error()+").")
					} else {
						song, intErr1 := strconv.Atoi(status["song"])
						length, intErr2 := strconv.Atoi(strings.SplitN(status["time"], ":", 2)[1])
						if intErr1 != nil {
							log.ErrorReport("update() PROGRESS_CHANGE", "Could not convert song ("+intErr1.Error()+").")
						} else if intErr2 != nil {
							log.ErrorReport("update() PROGRESS_CHANGE", "Could not convert length ("+intErr2.Error()+").")
						} else {
							seektime := int(float64(request.progressX) / float64(request.progressWidth) * float64(length))
							if seekErr := mpdConnection.Seek(song, seektime); seekErr != nil {
								log.ErrorReport("update() PROGRESS_CHANGE", "Could not mpd.Seek() ("+seekErr.Error()+").")
							} else {
								ui.SetProgressBarTime(seektime, length)
							}
						}
					} // end status error check

				} // end is not stopped

			} // end request switch

		} // end if not connected

		ui.Unlock()

	} // end for wait on channel

	pollChannel <- END_POLLING

	// Close the MPD connection, Juke is about to end:
	if currentState > NOT_CONNECTED {
		if errClose := mpdConnection.Close(); errClose != nil {
			log.ErrorReport("update()", "Could not mpd.Close() ("+errClose.Error()+").")
		}
	}

} // end update

func poll(updateChannel chan *jukeRequest, pollChannel chan int) {

	var rate int

	for {
		updateChannel <- &jukeRequest{state: POLL_REFREASH}

		rate = <-pollChannel
		if rate == END_POLLING {
			break
		} else {
			time.Sleep(time.Duration(rate) * time.Millisecond)
		}
	}

} // end poll
