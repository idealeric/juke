/*
The ui package is responsible for the user interface of Juke. It's objective
is to abstract and de-couple the main program's logic from the user interface.
*/
package ui

import (
	"github.com/idealeric/juke/log"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"strconv"
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
	STOPPED_OR_DC_PROGRESS     string = "0:00 / 0:00"
)

// Constant pixmap paths:
const (
	ICON              string = "/usr/share/pixmaps/juke/juke.png"
	NO_COVER_ARTWORK  string = "/usr/share/pixmaps/juke/no_cover.png"
	CUR_PL_ALBUM_SIZE int    = 20
)

const (
	CUR_PL_COL_ID int = iota
	CUR_PL_COL_ART
	CUR_PL_COL_NAME
	CUR_PL_COL_ARTIST
	CUR_PL_COL_ALBUM
	NUM_PL_COLS
)

// CurrentPLRow is an abstraction for other modules to
// work with visible rows in the current playlist.
type CurrentPLRow struct {
	ID          int
	ArtworkPath string
	Name        string
	Artist      string
	Album       string
	Bold        bool
	gref        *gtk.TreeRowReference
}

// Global referances for all "updating" GUI elements.
var (
	window              *gtk.Window                  // Main window
	leftControls        [2]*gtk.Button               // The 2 shuffle/repeat buttons
	playBackControls    [4]*gtk.Button               // The 4 playback buttons
	rightControls       [2]*gtk.Button               // The 2 connection/volume buttons
	controlsSize        int                          // The height of the controls (for current albumart resizing)
	currentAlbumArt     *gtk.Image                   // The current song's album artwork
	currentAlbumArtPath string                       // The current song's album artwork
	currentSongTitle    *gtk.Label                   // The current song's labeling
	currentPause        bool                         // The current state of the play/pause button.
	progressBar         *gtk.ProgressBar             // Progress bar for song
	progressBarEvent    *gtk.EventBox                // Progress bar eventbox (for click events)
	playlistTree        *gtk.TreeView                // Treeview for the current playlist.
	playlistModel       *gtk.ListStore               // Model for the current playlist.
	playlistSortable    *gtk.TreeSortable            // Playlist sortable.
	playlistCols        [3]*gtk.TreeViewColumn       // Playlist columns.
	currentBoldRow      CurrentPLRow                 // Currently bolded row reference.
	currentArtworks     map[string]*gdkpixbuf.Pixbuf // Hash table for fast artwork lookup.
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
		songLabel += escapeHTML(songName) + "</span>\nby "
	}

	if artist == "" {
		windowTitle += "Unknown"
		songLabel += "Unknown"
	} else {
		windowTitle += artist
		songLabel += escapeHTML(artist)
	}

	windowTitle += " [Juke]"
	window.SetTitle(windowTitle)

	if album != "" {
		songLabel += " from " + escapeHTML(album)
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
	window.SetIconFromFile(ICON)
	window.SetSizeRequest(800, -1) // TODO - Remember size
	window.SetTitle(NOT_CONNECTED_WINDOW_TITLE)
	window.SetBorderWidth(8)

	// Ensure we can do icons on buttons.
	var settings *glib.GObject = gtk.SettingsGetDefault().ToGObject()
	settings.Set("gtk-button-images", true)

	// Destory window is fired when the user "exits" the window.
	window.Connect("destroy", gtk.MainQuit)

	mainBox := gtk.NewVBox(false, 8) // Main VBox to glue UI the together

	// Current playlist treeview:
	currentArtworks = make(map[string]*gdkpixbuf.Pixbuf)
	playlistTree = gtk.NewTreeView()
	playlistModel = gtk.NewListStore(gtk.TYPE_INT, gdkpixbuf.GetType(), gtk.TYPE_STRING, gtk.TYPE_STRING, gtk.TYPE_STRING)
	playlistSortable = gtk.NewTreeSortable(playlistModel)
	//playlistTree.SetReorderable(true) // TODO - reordering
	playlistTree.SetModel(playlistModel)
	playlistColNames := []string{"ID", "Art", "Name", "Artist", "Album"}
	var playlistCol *gtk.TreeViewColumn
	for ci := CUR_PL_COL_NAME; ci < NUM_PL_COLS; ci++ {
		if ci == CUR_PL_COL_NAME {
			playlistCol = gtk.NewTreeViewColumn()
			playlistCol.SetSpacing(3)
			playlistCol.SetTitle(playlistColNames[CUR_PL_COL_NAME])
			playlistCol.SetMinWidth(370) // TODO - Remember sizing.
			cellPix := gtk.NewCellRendererPixbuf()
			playlistCol.PackStart(cellPix, false)
			playlistCol.AddAttribute(cellPix, "pixbuf", CUR_PL_COL_ART)
			cellText := gtk.NewCellRendererText()
			playlistCol.PackStart(cellText, true)
			playlistCol.AddAttribute(cellText, "markup", CUR_PL_COL_NAME)
			playlistCol.SetSortColumnId(CUR_PL_COL_NAME)
			playlistSortable.SetSortFunc(CUR_PL_COL_NAME, makeSortFunc(CUR_PL_COL_NAME))
		} else {
			playlistCol = gtk.NewTreeViewColumnWithAttributes(playlistColNames[ci], gtk.NewCellRendererText(), "markup", ci)
			playlistCol.SetMinWidth(190)
			playlistCol.SetSortColumnId(ci)
			playlistSortable.SetSortFunc(ci, makeSortFunc(ci))
		}
		playlistCol.SetResizable(true)
		playlistCol.SetSizing(gtk.TREE_VIEW_COLUMN_FIXED) // TODO - Remember sizing.
		playlistTree.AppendColumn(playlistCol)
		playlistCols[ci-CUR_PL_COL_NAME] = playlistCol
	}
	playlistScroll := gtk.NewScrolledWindow(nil, nil)
	playlistScroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_ALWAYS)
	playlistScroll.Add(playlistTree)
	playlistScroll.SetSizeRequest(-1, 330)
	mainBox.PackStart(playlistScroll, true, true, 0)

	bottomStatusBar := gtk.NewHBox(false, 8) // Main HBox for albumart, controls and other current stuff
	progressAndControls := gtk.NewVBox(false, 5)
	controls := gtk.NewHBox(false, 0)

	// Current album artwork:
	currentAlbumBorder := gtk.NewEventBox()
	currentAlbumBorder.ModifyBG(gtk.STATE_NORMAL, gdk.NewColor("#999"))
	currentAlbumSpace := gtk.NewEventBox()
	currentAlbumSpace.SetBorderWidth(1)
	currentAlbumArt = gtk.NewImage()
	currentAlbumArt.SetNoShowAll(true) // This allows manual display of the album art so that the proper size of the controls can be determined.
	currentAlbumArt.Hide()             // Hidden as well. ^^
	currentAlbumSpace.Add(currentAlbumArt)
	currentAlbumBorder.Add(currentAlbumSpace)
	bottomStatusBar.PackStart(currentAlbumBorder, false, false, 0)

	// Song progress bar:
	progressBar = gtk.NewProgressBar()
	progressBarEvent = gtk.NewEventBox()
	progressBar.SetOrientation(gtk.PROGRESS_LEFT_TO_RIGHT)
	progressBar.SetText(STOPPED_OR_DC_PROGRESS)
	progressBar.SetFraction(0)
	progressBar.SetPulseStep(0.05)
	//progressBar.SetEllipsize(0.05) // TODO - Implement this (maybe)
	progressBarEvent.Add(progressBar)
	progressAndControls.PackStart(progressBarEvent, false, false, 0)

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
	currentPause = false

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

	// Firstly determine the height of the controls (theme dependant)
	// and then show the album art at that size. (Minus 2 for border)
	controlsSize = progressAndControls.GetAllocation().Height - 2
	currentAlbumArtPath = ""
	SetCurrentAlbumArt(NO_COVER_ARTWORK)
	currentAlbumArt.Show()

} // end Init

// SetCurrentAlbumArt sets the current album artwork to the image specified by path.
func SetCurrentAlbumArt(path string) {

	// Since this function is called often (at least as often as MPD is polled)
	// and it is a relatively expensive operation, this safety value prevents
	// it from being run too often. It only needs to run when there is work to
	// be done (as in there is new albumart to load and scale).
	if currentAlbumArtPath == path {
		return
	}
	currentAlbumArtPath = path

	pbuf, pbufErr := gdkpixbuf.NewFromFileAtSize(path, controlsSize, controlsSize)
	if pbufErr != nil {
		log.ErrorReport("SetCurrentAlbumArt()", "Could not load current album art ("+pbufErr.Error()+").")
	} else {
		currentAlbumArt.SetFromPixbuf(pbuf)
	}
	pbuf.Unref()

} // end SetCurrentAlbumArt

// SetPlayPause changes the image on the play button based on the boolean argument. True will
// display a pause image, while false will display a play image.
func SetPlayPause(pause bool) {

	if pause && !currentPause {
		// Attempting to pause. Must also NOT be paused.
		playBackControls[PLAY_PAUSE_BUTTON].SetImage(gtk.NewImageFromStock(gtk.STOCK_MEDIA_PAUSE, gtk.ICON_SIZE_DND))
		currentPause = true
	} else if !pause && currentPause {
		// Attempting to play. Must also NOT be playing.
		playBackControls[PLAY_PAUSE_BUTTON].SetImage(gtk.NewImageFromStock(gtk.STOCK_MEDIA_PLAY, gtk.ICON_SIZE_DND))
		currentPause = false
	}

} // end SetPlayPause

// SetProgressBarTime takes song progress and updates the progress bar to
// reflect that both textually and visually.
func SetProgressBarTime(at, total int) {

	atSeconds, totalSeconds := at%60, total%60
	timeText := strconv.Itoa(at/60) + ":"
	if atSeconds < 10 {
		timeText += "0"
	}
	timeText += strconv.Itoa(atSeconds) + " / " + strconv.Itoa(total/60) + ":"
	if totalSeconds < 10 {
		timeText += "0"
	}
	timeText += strconv.Itoa(totalSeconds)
	progressBar.SetText(timeText)
	progressBar.SetFraction(float64(at) / float64(total))

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

// AddManyRowstoCurrentPlaylist adds mulitple rows at once.
// It is designed to be more efficient than adding one row at a time.
func AddManyRowstoCurrentPlaylist(rows []*CurrentPLRow) {

	playlistTree.SetModel(nil)
	for _, row := range rows {
		AddRowtoCurrentPlaylist(row)
	}
	playlistTree.SetModel(playlistModel)

} // end AddManyRowstoCurrentPlaylist

// AddRowtoCurrentPlaylist adds a row to the current playlist view.
func AddRowtoCurrentPlaylist(row *CurrentPLRow) {

	var iter gtk.TreeIter
	playlistModel.Append(&iter)

	if val, exists := currentArtworks[row.ArtworkPath]; exists {
		playlistModel.Set(&iter, row.ID, val.GPixbuf, escapeHTML(row.Name), escapeHTML(row.Artist), escapeHTML(row.Album))
	} else {
		pbuf, pbufErr := gdkpixbuf.NewFromFileAtSize(row.ArtworkPath, CUR_PL_ALBUM_SIZE, CUR_PL_ALBUM_SIZE)
		if pbufErr != nil {
			log.ErrorReport("AddRowtoCurrentPlaylist()", "Could not load artwork ("+pbufErr.Error()+").")
		} else {
			currentArtworks[row.ArtworkPath] = pbuf
			playlistModel.Set(&iter, row.ID, currentArtworks[row.ArtworkPath].GPixbuf, escapeHTML(row.Name), escapeHTML(row.Artist), escapeHTML(row.Album))
		}
	}

	if row.Bold {
		path := playlistModel.GetPath(&iter)
		defer path.Free()
		BoldRowByReference(&CurrentPLRow{ID: row.ID, gref: gtk.NewTreeRowReference(playlistModel, path)})
	}

} // end AddRowtoCurrentPlaylist

// BoldRowByReference makes a row in the current playlist bold.
// Only one row can be bold at a time.
func BoldRowByReference(row *CurrentPLRow) {

	var iter gtk.TreeIter
	var strv string

	// First, unbold the old row if there is one
	if currentBoldRow.gref != nil && currentBoldRow.gref.Valid() {
		path := currentBoldRow.gref.GetPath()
		defer path.Free()
		playlistModel.GetIter(&iter, path)
		for vali := CUR_PL_COL_NAME; vali < NUM_PL_COLS; vali++ {
			var val glib.GValue
			playlistModel.GetValue(&iter, vali, &val)
			strv = val.GetString()
			playlistModel.SetValue(&iter, vali, removeBold(strv))
		}
		currentBoldRow.gref.Free()
	}
	currentBoldRow.gref = row.gref
	path := currentBoldRow.gref.GetPath()
	defer path.Free()
	var id glib.GValue
	playlistModel.GetIter(&iter, path)
	playlistModel.GetValue(&iter, CUR_PL_COL_ID, &id)
	currentBoldRow.ID = id.GetInt()
	for vali := CUR_PL_COL_NAME; vali < NUM_PL_COLS; vali++ {
		var val glib.GValue
		playlistModel.GetValue(&iter, vali, &val)
		strv = val.GetString()
		playlistModel.SetValue(&iter, vali, addBold(strv))
	}

} // end BoldRowByReference

// BoldRowById makes a row in the current playlist by the ID
// of the song. Only one row can be bold at a time.
func BoldRowById(rowId int) {

	// Since this operation is O(n), only
	// do it if need be.
	if currentBoldRow.ID == rowId {
		return
	}

	var iter gtk.TreeIter
	ok := playlistModel.GetIterFirst(&iter)
	for ok {
		var id glib.GValue
		playlistModel.GetValue(&iter, CUR_PL_COL_ID, &id)
		if rowId == id.GetInt() {
			path := playlistModel.GetPath(&iter)
			defer path.Free()
			// This will set currentBoldRow.ID properly:
			BoldRowByReference(&CurrentPLRow{ID: id.GetInt(), gref: gtk.NewTreeRowReference(playlistModel, path)})
			return
		}
		ok = playlistModel.IterNext(&iter)
	}

	log.ErrorReport("BoldRowById()", "Never found row with ID "+string(rowId))

} // end BoldRowById

// ClearCurrentPlaylist clears the current playlist view.
func ClearCurrentPlaylist() {

	currentBoldRow.gref = nil
	currentBoldRow.ID = -1
	playlistModel.Clear()
	playlistSortable.SetSortColumnId(gtk.TREE_SORTABLE_UNSORTED_SORT_COLUMN_ID, gtk.SORT_ASCENDING)

	// In addition to cleaning up the model and all its
	// references, the hashmap references need to be
	// freed.
	for key, value := range currentArtworks {
		value.Unref() // Let Gobject clean up after us
		delete(currentArtworks, key)
	}

} // end ClearCurrentPlaylist
