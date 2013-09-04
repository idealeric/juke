Juke - A Different Kind of MPD Client
================================

Juke is a front-end, GTK+ (v2) client for the Music Playing Deamon. It is being developed for Linux and is written in [Go](http://golang.org/). It has an emphasis on speed and elegance, but also aims for a beautiful interface.

Documentation
-------------------------
Code documenation can be found at [GoDoc.org](http://godoc.org/github.com/idealeric/juke) (as well as in the source code). There is no documenation for input (yet).

Also, huge thanks to the contributors of the Go libraries that Juke uses:
* [go-gtk](https://github.com/mattn/go-gtk)
* [gompd](https://github.com/fhs/gompd)

Installation
-------------------------
You can use `$ go get` to pull the source right into your `$GOPATH`. `$ go install` will compile (and should pull the lib dependencies as well). `$ juke` to run.

Additionally, one will need to set up some symlinks to the images that juke uses:
```
# mkdir /usr/share/pixmaps/juke
# ln -s $GOPATH/src/github.com/idealeric/juke/ui/images/icon.png /usr/share/pixmaps/juke/juke.png
# ln -s $GOPATH/src/github.com/idealeric/juke/ui/images/noCover.png /usr/share/pixmaps/juke/no_cover.png
```

The TODO List (High Priority)
-------------------------

* Icons and actions for secondary controls.
* A library browser.
* Configuration options.
* Internationalization.

The Maybe List
-------------------------

* Drag'n'drop current playlist.
* Inline tag editting.
* Playlist browser/operations.
* GTK3 instead.

The No-Way (Probably-Not) List
-------------------------

* Lyrics support.
* Status icon/tray support/notifications.
* Windows support.
* Last.fm scrobbler (use [mpdas](http://mpd.wikia.com/wiki/Client:Mpdas) or [mpdscribble](http://mpd.wikia.com/wiki/Client:Mpdscribble)).
