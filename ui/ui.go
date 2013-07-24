/*
The ui package is responsible for the user interface of Juke. It's objective
is to abstract and de-couple the main program's logic from the user interface.
*/
package ui

import (
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"strconv"
	"strings"
)

// Playback control constant indexes:
const (
	PREV_BUTTON uint8 = iota
	PLAY_PAUSE_BUTTON
	STOP_BUTTON
	NEXT_BUTTON
)

// Left-side control constant indexes:
const (
	SHUFFLE_BUTTON uint8 = iota
	REPEAT_BUTTON
)

// Right-side control constant indexes:
const (
	VOLUME_BUTTON uint8 = iota
	CONNECTION_BUTTON
)

// Constant referances for set program states:
const (
	NOT_CONNECTED_WINDOW_TITLE string = "Not Connected [Juke]"
	NOT_CONNECTED_SONG_LABEL   string = "<span size=\"x-large\" font_weight=\"bold\">Stopped</span>\nNot connected."
	STOPPED_WINDOW_TITLE       string = "Stopped [Juke]"
	STOPPED_SONG_LABEL         string = "<span size=\"x-large\" font_weight=\"bold\">Stopped</span>\nConnected."
	STOPPED_OR_DC_PROGRESS     string = " "
)

// Global referances for all "updating" GUI elements.
var (
	window           *gtk.Window      // Main window
	leftControls     [2]*gtk.Button   // The 2 shuffle/repeat buttons
	playBackControls [4]*gtk.Button   // The 4 playback buttons
	rightControls    [2]*gtk.Button   // The 2 connection/volume buttons
	currentAlbumArt  *gtk.Image       // The current song's album artwork
	currentSongTitle *gtk.Label       // The current song's labeling
	progressBar      *gtk.ProgressBar // Progress bar for song
)

// MainLoop runs the GUI toolkit's main loop.
func MainLoop() {

	gdk.ThreadsEnter()
	gtk.Main()
	gdk.ThreadsLeave()

} // end MainLoop

// SetCurrentSong changes the window title and current song labeling to reflect
// the parameters of song name, artist name, and album name.
func SetCurrentSong(songName, artist, album string) {

	windowTitle, songLabel := "", "<span size=\"x-large\" font_weight=\"bold\">"
	if songName == "" {
		windowTitle += "Unknown by "
		songLabel += "Unknown</span>\nby "
	} else {
		windowTitle += songName + " by "
		songLabel += strings.Replace(songName, "&", "&amp;", -1) + "</span>\nby "
	}

	if artist == "" {
		windowTitle += "Unknown"
		songLabel += "Unknown"
	} else {
		windowTitle += artist
		songLabel += strings.Replace(artist, "&", "&amp;", -1)
	}

	windowTitle += " [Juke]"
	window.SetTitle(windowTitle)

	if album != "" {
		songLabel += " from " + strings.Replace(album, "&", "&amp;", -1)
	}

	currentSongTitle.SetMarkup(songLabel)

} // end SetCurrentSong

// SetCurrentSongNotConnected changes the window title and current song labeling
// to reflect and unconnected client.
func SetCurrentSongNotConnected() {

	window.SetTitle(NOT_CONNECTED_WINDOW_TITLE)
	currentSongTitle.SetMarkup(NOT_CONNECTED_SONG_LABEL)

} // end SetCurrentSongNotConnected

// SetCurrentSongStopped changes the window title and current song labeling to
// reflect a stopped but still connected client.
func SetCurrentSongStopped() {

	window.SetTitle(STOPPED_WINDOW_TITLE)
	currentSongTitle.SetMarkup(STOPPED_SONG_LABEL)

} // end SetCurrentSongStopped

// InitInterface inits the GUI toolkit and builds most of the base interface.
func InitInterface() {

	// Initialize GTK/GLib.
	// This is a thread-safe mode.
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gtk.Init(nil)
	// Initialize a window.
	window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetIconName("gtk-dialog-info") // TODO - Make an icon and set it.
	window.SetSizeRequest(800, -1)        // TODO - Remember size
	window.SetTitle(NOT_CONNECTED_WINDOW_TITLE)
	window.SetBorderWidth(8)

	// Ensure we can do icons on buttons.
	var settings *glib.GObject = gtk.SettingsGetDefault().ToGObject()
	settings.Set("gtk-button-images", true)

	// Destory window is fired when the user "exits" the window.
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	})

	mainBox := gtk.NewVBox(false, 0)         // Main VBox to glue UI the together
	bottomStatusBar := gtk.NewHBox(false, 8) // Main HBox for albumart, controls and other current stuff
	progressAndControls := gtk.NewVBox(false, 5)
	controls := gtk.NewHBox(false, 0)

	// Current album artwork:
	currentAlbumArt = gtk.NewImageFromFile("/home/chuck/Code/go/src/github.com/idealeric/juke/ui/images/noCover.png") // TODO - Implement changing albumartwork
	bottomStatusBar.PackStart(currentAlbumArt, false, false, 0)

	// Song progress bar:
	progressBar = gtk.NewProgressBar()
	progressBar.SetOrientation(gtk.PROGRESS_LEFT_TO_RIGHT)
	progressBar.SetText(STOPPED_OR_DC_PROGRESS)
	progressBar.SetFraction(0)
	progressBar.SetPulseStep(0.05)
	//progressBar.SetEllipsize(0.05) // TODO - Implement this (maybe)
	progressAndControls.PackStart(progressBar, false, false, 0)

	// Current song labeling:
	currentSongTitleAlign := gtk.NewAlignment(0, 0, 0, 1)
	currentSongTitle = gtk.NewLabel("")
	currentSongTitle.SetJustify(gtk.JUSTIFY_LEFT)
	currentSongTitle.SetMarkup(NOT_CONNECTED_SONG_LABEL)
	currentSongTitleAlign.Add(currentSongTitle)
	progressAndControls.PackStart(currentSongTitleAlign, false, false, 0)

	// Left-hand controls:
	leftControlsBox := gtk.NewHBox(false, 0)

	// TODO - get the right images
	leftControls[SHUFFLE_BUTTON] = gtk.NewButtonFromStock(gtk.STOCK_ZOOM_IN)
	leftControls[SHUFFLE_BUTTON].SetImage(gtk.NewImageFromStock(gtk.STOCK_ZOOM_IN, gtk.ICON_SIZE_DND))

	leftControls[REPEAT_BUTTON] = gtk.NewButtonFromStock(gtk.STOCK_ZOOM_OUT)
	leftControls[REPEAT_BUTTON].SetImage(gtk.NewImageFromStock(gtk.STOCK_ZOOM_OUT, gtk.ICON_SIZE_DND))

	for i := range leftControls {
		leftControls[i].SetCanFocus(false)
		leftControls[i].SetRelief(gtk.RELIEF_HALF)
		leftControls[i].SetLabel("")
		leftControlsBox.PackStart(leftControls[i], false, false, 0)
	}
	controls.PackStart(leftControlsBox, true, true, 0)

	// Playback controls:
	playBox := gtk.NewHBox(false, 0) // This is playback controls only.

	playBackControls[PREV_BUTTON] = gtk.NewButtonFromStock(gtk.STOCK_MEDIA_PREVIOUS)
	playBackControls[PREV_BUTTON].SetImage(gtk.NewImageFromStock(gtk.STOCK_MEDIA_PREVIOUS, gtk.ICON_SIZE_DND))

	playBackControls[PLAY_PAUSE_BUTTON] = gtk.NewButtonFromStock(gtk.STOCK_MEDIA_PLAY)
	playBackControls[PLAY_PAUSE_BUTTON].SetImage(gtk.NewImageFromStock(gtk.STOCK_MEDIA_PLAY, gtk.ICON_SIZE_DND))

	playBackControls[STOP_BUTTON] = gtk.NewButtonFromStock(gtk.STOCK_MEDIA_STOP)
	playBackControls[STOP_BUTTON].SetImage(gtk.NewImageFromStock(gtk.STOCK_MEDIA_STOP, gtk.ICON_SIZE_DND))

	playBackControls[NEXT_BUTTON] = gtk.NewButtonFromStock(gtk.STOCK_MEDIA_NEXT)
	playBackControls[NEXT_BUTTON].SetImage(gtk.NewImageFromStock(gtk.STOCK_MEDIA_NEXT, gtk.ICON_SIZE_DND))

	for i := range playBackControls {
		playBackControls[i].SetCanFocus(false)
		playBackControls[i].SetRelief(gtk.RELIEF_HALF)
		playBackControls[i].SetLabel("")
		playBox.PackStart(playBackControls[i], false, false, 0)
	}
	controls.PackStart(playBox, true, false, 0)

	// Right-hand controls:
	rightControlsBox := gtk.NewHBox(false, 0)

	// TODO - get the right images, or conversely make my own
	rightControls[CONNECTION_BUTTON] = gtk.NewButtonFromStock(gtk.STOCK_CONNECT)
	rightControls[CONNECTION_BUTTON].SetImage(gtk.NewImageFromStock(gtk.STOCK_CONNECT, gtk.ICON_SIZE_DND))

	rightControls[VOLUME_BUTTON] = gtk.NewButtonFromStock(gtk.STOCK_SAVE)
	rightControls[VOLUME_BUTTON].SetImage(gtk.NewImageFromStock(gtk.STOCK_SAVE, gtk.ICON_SIZE_DND))

	for i := range rightControls {
		rightControls[i].SetCanFocus(false)
		rightControls[i].SetRelief(gtk.RELIEF_HALF)
		rightControls[i].SetLabel("")
		rightControlsBox.PackStart(rightControls[i], false, false, 0)
	}
	rightControlsAlign := gtk.NewAlignment(1, 0, 0, 1)
	rightControlsAlign.Add(rightControlsBox)
	controls.PackStart(rightControlsAlign, true, true, 0)

	progressAndControls.PackStart(controls, false, false, 0)
	bottomStatusBar.PackStart(progressAndControls, true, true, 0)
	mainBox.PackStart(bottomStatusBar, false, false, 0)

	window.Add(mainBox)
	window.ShowAll()

} // end Init

// SetPlayPause changes the image on the play button based on the boolean argument. True will
// display a pause image, while false will display a play image.
func SetPlayPause(pause bool) {

	if pause {
		playBackControls[PLAY_PAUSE_BUTTON].SetImage(gtk.NewImageFromStock(gtk.STOCK_MEDIA_PAUSE, gtk.ICON_SIZE_DND))
	} else {
		playBackControls[PLAY_PAUSE_BUTTON].SetImage(gtk.NewImageFromStock(gtk.STOCK_MEDIA_PLAY, gtk.ICON_SIZE_DND))
	}

} // end SetPlayPause

// SetProgressBarTime takes song progress in the form of "at:total" and updates the
// progress bar to reflect that both textually and visually. The input is in string
// form to consolidate all of the conversion to one place in the codebase.
func SetProgressBarTime(time string) {

	splitTime := strings.SplitN(time, ":", 2)
	at, total := splitTime[0], splitTime[1]

	atNum, errAt := strconv.Atoi(at)
	if errAt != nil {
		return // TODO - make this better
	}
	totalNum, errTotal := strconv.Atoi(total)
	if errTotal != nil {
		return // TODO - make this better
	}

	atSeconds, totalSeconds := atNum%60, totalNum%60
	timeText := strconv.Itoa(atNum/60) + ":"
	if atSeconds < 10 {
		timeText += "0"
	}
	timeText += strconv.Itoa(atSeconds) + " / " + strconv.Itoa(totalNum/60) + ":"
	if totalSeconds < 10 {
		timeText += "0"
	}
	timeText += strconv.Itoa(totalSeconds)
	progressBar.SetText(timeText)
	progressBar.SetFraction(float64(atNum) / float64(totalNum))

} // end SetProgressBarTime

// SetProgressBarTimeStoppedOrDisconnected sets the progress bar to reflect that
// the client is stopped or is disconnected.
func SetProgressBarTimeStoppedOrDisconnected() {

	progressBar.SetText(STOPPED_OR_DC_PROGRESS)
	progressBar.SetFraction(0.0)

} // end SetProgressBarTimeStoppedOrDisconnected

// Lock grabs the ui lock.
func Lock() {

	gdk.ThreadsEnter()

} // end Lock

// Unlock releases any ui locks.
func Unlock() {

	gdk.Flush()
	gdk.ThreadsLeave()

} // end Unlock
