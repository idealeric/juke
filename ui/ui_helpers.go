/*
This file is part of Juke MPD client ui package. See juke.go for more details.

This particular file has all of Juke's ui helper functions.
*/

package ui

import (
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"strings"
)

// escapeHTML removes HTML tokens from a string so that it might be
// rendered without any formating in markup contexts.
func escapeHTML(s string) string {

	ret := strings.Replace(s, "&", "&amp;", -1)
	ret = strings.Replace(ret, "<", "&lt;", -1)
	ret = strings.Replace(ret, ">", "&gt;", -1)
	return strings.Replace(ret, "\"", "&quot;", -1)

} // end escapeHTML

// addBold takes a string and makes it bold.
func addBold(str string) string {

	return "<b>" + str + "</b>"

} // end addBold

// removeBold takes a string and makes it not bold.
func removeBold(str string) string {

	ret := strings.Replace(str, "<b>", "", -1)
	return strings.Replace(ret, "</b>", "", -1)

} // end removeBold

// makeSortFunc creates a sort function for the specified column number.
func makeSortFunc(col int) func(*gtk.TreeModel, *gtk.TreeIter, *gtk.TreeIter) int {

	return func(m *gtk.TreeModel, a *gtk.TreeIter, b *gtk.TreeIter) int {
		var (
			vala glib.GValue
			valb glib.GValue
		)
		playlistModel.GetValue(a, col, &vala)
		playlistModel.GetValue(b, col, &valb)
		stra, strb := removeBold(vala.GetString()), removeBold(valb.GetString())
		if stra == strb {
			return 0
		} else if stra > strb {
			return 1
		} else {
			return -1
		}
	}

} // end makeSortFunc
