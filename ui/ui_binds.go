package ui

import (
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"unsafe"
)

func callBackCheckandCheckforError(f func() error, cntx *glib.CallbackContext) {
	if err := f(); err != nil {
		fmt.Printf("Juke UI - Error: %v\n", err)
		// TODO - investigate the error with this line
		//fmt.Printf("Juke UI - Error Context: %s\n", cntx.Data().(string))
	}
} // end callBackCheckandCheckforError

// OnExit will bind additional functions to the ui "exit" event.
func OnExit(f func() error) {

	window.Connect("destroy", func(cntx *glib.CallbackContext) {
		callBackCheckandCheckforError(f, cntx)
	})

} // end OnExit

// NextClick will bind to the "release" event on the next button.
func NextClick(f func() error) {

	playBackControls[NEXT_BUTTON].Connect("released", func(cntx *glib.CallbackContext) {
		callBackCheckandCheckforError(f, cntx)
	})

} // end NextClick

// PreviousClick will bind to the "release" event on the previous button.
func PreviousClick(f func() error) {

	playBackControls[PREV_BUTTON].Connect("released", func(cntx *glib.CallbackContext) {
		callBackCheckandCheckforError(f, cntx)
	})

} // end PreviousClick

// PlayPauseClick will bind to the "release" event on the play/pause button.
func PlayPauseClick(f func() error) {

	playBackControls[PLAY_PAUSE_BUTTON].Connect("released", func(cntx *glib.CallbackContext) {
		callBackCheckandCheckforError(f, cntx)
	})

} // end PlayPauseClick

// StopClick will bind to the "release" event on the stop button.
func StopClick(f func() error) {

	playBackControls[STOP_BUTTON].Connect("released", func(cntx *glib.CallbackContext) {
		callBackCheckandCheckforError(f, cntx)
	})

} // end StopClick

// ProgressBarClick will bind to the "click" even on the progress bar.
func ProgressBarClick(f func(int, int) error) {

	progressBarEvent.Connect("button_press_event", func(cntx *glib.CallbackContext) {
		arg := cntx.Args(0)
		eventButton := *(**gdk.EventButton)(unsafe.Pointer(&arg))
		if err := f(int(eventButton.X), progressBar.GetAllocation().Width); err != nil {
			fmt.Printf("Juke UI - Error: %v\n", err)
			// TODO - investigate the error with this line
			//fmt.Printf("Juke UI - Error Context: %s\n", cntx.Data().(string))
		}
	})

} // end ProgressBarClick
