/*
The log package is responsible for reporting errors, warnings, and messages
to the user. It also handles "graceful" crashes.
*/
package log

import (
	"fmt"
	"time"
)

// ErrorReport prints an error message to stdout (colorized) with a where,
// a when, and a what.
func ErrorReport(where, what string) {

	when := time.Now()
	fmt.Printf("[\033[31mJuke Error\033[0m ")
	fmt.Printf("@ \033[33m%s\033[0m] ", when.Format("3:04pm"))
	fmt.Printf("\033[34m%s\033[0m: %s\n", where, what)

} // end ErrorReport

// ErrorOut is ErrorReport with the addition of a call to panic()
func ErrorOut(where, what string) {

	ErrorReport(where, what)
	panic("Juke errored out. :(")

} // end ErrorOut

// MessageReport prints a message to stdout (colorized) with a where,
// a when, and a what.
func MessageReport(where, what string) {

	when := time.Now()
	fmt.Printf("[\033[32mJuke Msg\033[0m ")
	fmt.Printf("@ \033[33m%s\033[0m] ", when.Format("3:04pm"))
	fmt.Printf("\033[34m%s\033[0m: %s\n", where, what)

} // end MessageReport
