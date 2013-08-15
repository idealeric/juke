/*
This file is part of Juke MPD client. See juke.go for more details.

This particular file has Juke's helper functions (small operations).
*/

package main

import (
	"github.com/idealeric/juke/log"
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
		log.ErrorReport("albumArtFilename()", "Either sub directory is reported as empty or user.Current() failed.")
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
