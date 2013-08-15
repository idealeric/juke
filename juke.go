/*
Juke is a front-end, GTK+ client for the Music Playing Deamon.

Copyright: Eric Butler 2013
Version:   0.2a

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
	"github.com/idealeric/juke/ui"
)

// Keep main short and sweet!
func main() {

	var updateChannel chan *jukeRequest = make(chan *jukeRequest)

	ui.InitInterface()

	// Init any concurrent routines:
	go update(updateChannel)

	// For code tidyness, callbacks are defined in a seperate file.
	initCallBacks(updateChannel)

	ui.MainLoop() // This blocks until the GUI is destoryed.

	close(updateChannel) // Tells update to shut off

} // end main
