package ui

import (
	"github.com/idealeric/juke/log"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"unsafe"
)

const ROW_BUFFER_SIZE = 50

func callBackCheckandCheckforError(f func() error, cntx *glib.CallbackContext) {
	if err := f(); err != nil {
		log.ErrorReport("UI callBackCheckandCheckforError()", err.Error()+".")
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

// ConnectionClick will bind to the "click" event on the connection button.
func ConnectionClick(f func() error) {

	rightControls[CONNECTION_BUTTON].Connect("released", func(cntx *glib.CallbackContext) {
		callBackCheckandCheckforError(f, cntx)
	})

} // end ConnectionClick

// ProgressBarClick will bind to the "click" event on the progress bar.
func ProgressBarClick(f func(int, int) error) {

	progressBarEvent.Connect("button_press_event", func(cntx *glib.CallbackContext) {
		arg := cntx.Args(0)
		eventButton := *(**gdk.EventButton)(unsafe.Pointer(&arg))
		if err := f(int(eventButton.X), progressBar.GetAllocation().Width); err != nil {
			log.ErrorReport("UI callBackCheckandCheckforError()", err.Error()+".")
		}
	})

} // end ProgressBarClick

// CurrentRowDoubleClick will bind to the "double-click" event on a row in
// the current playlist.
func CurrentRowDoubleClick(f func(*CurrentPLRow) error) {

	playlistTree.Connect("row-activated", func(cntx *glib.CallbackContext) {
		var (
			path *gtk.TreePath
			iter gtk.TreeIter
			val  glib.GValue
			col  *gtk.TreeViewColumn
		)
		playlistTree.GetCursor(&path, &col)
		playlistModel.GetIter(&iter, path)
		playlistModel.GetValue(&iter, CUR_PL_COL_ID, &val)
		if err := f(&CurrentPLRow{ID: val.GetInt(), gref: gtk.NewTreeRowReference(playlistModel, path)}); err != nil {
			log.ErrorReport("UI callBackCheckandCheckforError()", err.Error()+".")
		}
	})

} // end CurrentRowDoubleClick

// CurrentColumnClick will bind to the "column click" event in the
// current playlist.
func CurrentColumnClick(f func(chan *CurrentPLRow) error) {

	for _, c := range playlistCols {
		c.Connect("clicked", func(cntx *glib.CallbackContext) {

			rowsChan := make(chan *CurrentPLRow, ROW_BUFFER_SIZE)

			go func() {
				var iter gtk.TreeIter
				ok := playlistModel.GetIterFirst(&iter)
				for ok {
					var id glib.GValue
					playlistModel.GetValue(&iter, CUR_PL_COL_ID, &id)
					rowsChan <- &CurrentPLRow{ID: id.GetInt()}
					ok = playlistModel.IterNext(&iter)
				}
				close(rowsChan)
			}()

			if err := f(rowsChan); err != nil {
				log.ErrorReport("UI callBackCheckandCheckforError()", err.Error()+".")
			}

		})
	} // end for range of columns

} // end CurrentColumnClick

// CurrentRemoveSongs will bind to the "click" event on the
// remove songs button in the current playlist menu.
func CurrentRemoveSongs(f func(chan *CurrentPLRow) error) {

	playlistMenuRemove.Connect("activate", func(cntx *glib.CallbackContext) {

		rowsChan := make(chan *CurrentPLRow, ROW_BUFFER_SIZE)

		go func() {
			var iter gtk.TreeIter
			ok := playlistModel.GetIterFirst(&iter)
			for ok {

				path := playlistModel.GetPath(&iter)
				defer path.Free()

				if playlistSelection.PathIsSelected(path) {
					var id glib.GValue
					playlistModel.GetValue(&iter, CUR_PL_COL_ID, &id)
					rowsChan <- &CurrentPLRow{ID: id.GetInt(), gref: gtk.NewTreeRowReference(playlistModel, path)}
				}

				ok = playlistModel.IterNext(&iter)

			}
			close(rowsChan)
		}()

		if err := f(rowsChan); err != nil {
			log.ErrorReport("UI callBackCheckandCheckforError()", err.Error()+".")
		}

	})

} // end CurrentRemoveSongs

// CurrentClearSongs will bind the "click" event on the
// clear playlist button in the current playlist menu.
func CurrentClearSongs(f func() error) {

	playlistMenuClear.Connect("activate", func(cntx *glib.CallbackContext) {
		callBackCheckandCheckforError(f, cntx)
	})

} // end CurrentClearSongs
