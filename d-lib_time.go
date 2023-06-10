package main

/*
 * d-lib support functions
*/

import (
	"fmt"
	"time"
)

// return date string as yyyy-mm-dd
func date() (string) {
	now := time.Now()
	return fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
}

// return time string as hh-mm-ss (24 hour time)
func hms() (string) {
	now := time.Now()
	return fmt.Sprintf("%02d-%02d-%02d", now.Hour(), now.Minute(), now.Second())
}

// return date string as yyyy-mm-dd
func date_hms() (string) {
	now := time.Now()
	return fmt.Sprintf("%d-%02d-%02d_%02d-%02d-%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

// return time for use as random seed
func timeseed() (int) {
	return int(time.Now().UnixNano())
}
