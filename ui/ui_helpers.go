/*
This file is part of Juke MPD client ui package. See juke.go for more details.

This particular file has all of Juke's ui helper functions.
*/

package ui

import (
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
