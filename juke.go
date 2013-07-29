/*
Juke is a front-end, GTK+ client for the Music Playing Deamon.

Copyright: Eric Butler 2013
Version:   0.1a

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"fmt"
	"github.com/idealeric/juke/ui"
	"os"
	"os/user"
	"path"
)

// albumArtFilename takes a subdirectory of a song and attempts to string
// it together with a music directory and proper filename.
func albumArtFilename(subDir string) string {

	// TODO - enable music directory to be configured
	usr, err := user.Current()
	if err != nil || subDir == "" {
		fmt.Println("bad", err) // TODO - make this better
		return ui.NO_COVER_ARTWORK
	}

	curSongDir := path.Dir(path.Join(usr.HomeDir, "Music/", subDir))
	// TODO - perhaps supplement this with more filenames and types (config?)
	if _, err := os.Stat(path.Join(curSongDir, "cover.jpg")); err == nil {
		return path.Join(curSongDir, "cover.jpg")
	}
	if _, err := os.Stat(path.Join(curSongDir, "cover.jpeg")); err == nil {
		return path.Join(curSongDir, "cover.jpeg")
	}
	if _, err := os.Stat(path.Join(curSongDir, "cover.png")); err == nil {
		return path.Join(curSongDir, "cover.png")
	}

	return ui.NO_COVER_ARTWORK

} // end albumArtFilename

// Keep main short and sweet!
func main() {

	var (
		updateChannel chan *jukeRequest = make(chan *jukeRequest)
		pollChannel   chan int          = make(chan int)
	)

	ui.InitInterface()

	// Init any concurrent routines:
	go update(updateChannel, pollChannel)

	// For code tidyness, callbacks are defined in a seperate file.
	initCallBacks(updateChannel)

	ui.MainLoop() // This blocks until the GUI is destoryed.

	close(updateChannel) // Tells update to shut off

} // end main
