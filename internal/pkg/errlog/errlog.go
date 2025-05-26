// Package errlog provides functions to log errors to stdout.
package errlog

import "log"

// Print logs error to stdout.
func Print(err error) {
	log.Printf("ERROR: %v", err)
}
