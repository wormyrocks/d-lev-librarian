package main

/*
 * d-lev constants & helper functions
*/

import (
	"strings"
	"strconv"
)

const (
	VERSION = "2"							// librarian version
	//
	EE_RW_BYTES = 4							// eeprom bytes per read / write cycle
	EE_PG_BYTES = 0x100						// eeprom bytes per page
	EE_SPI_ADDR = 0x10000					// eeprom code start addr
	EE_SPI_SZ = 0x4000						// eeprom code size : 16kB code space
	EE_SPI_END = EE_SPI_ADDR + EE_SPI_SZ	// eeprom code end addr
	EE_PRE_ADDR = 0x0						// eeprom presets start addr
	EE_PRE_SZ = 0x10000						// eeprom presets size : [-120:120] & profile[-7:7] space
	EE_PRE_END = EE_PRE_ADDR + EE_PRE_SZ	// eeprom presets end addr
	EE_PRO_ADDR = 0x7900					// eeprom profiles start addr
	EE_PRO_SZ = 0xF00						// eeprom profiles size : [-7:7] space
	EE_PRO_END = EE_PRO_ADDR + EE_PRO_SZ	// eeprom profiles end addr
	EE_START = EE_PRE_ADDR					// eeprom start addr
	EE_END = EE_SPI_END						// eeprom end addr
	EE_WR_MS = 6							// eeprom write wait time (ms)
	//
	UI_PAGES = 20							// ui pages
	UI_COLS = 2								// ui page columns
	UI_ROWS = 4								// ui page rows
	UI_KNOBS = UI_COLS * UI_ROWS			// ui knobs
	UI_PAGE_KNOB = 7						// ui page selector knob
	UI_PRN_PG_COLS = 4						// ui print pages columns
	UI_PRN_PG_ROWS = 5						// ui print pages rows
	KNOBS = UI_KNOBS * UI_PAGES				// total knobs
	//
	PRE_SLOT_MAX = 120						// highest preset slot
	PRO_SLOT_MAX = 7						// highest profile slot
	SLOTS = 256								// pre + pro slots
	SLOT_BYTES = 256						// bytes per slot
	//
	RX_BUF_BYTES = 512						// serial port rx buffer size
	CHARS_PER_DOT = 4096					// chars for each activity dot printed
	CFG_FILE = "d-lib.cfg"					// config file name
	CFG_PORT = 0							// default port
)

// convert string of multi-byte hex values to slice of ints
// hex string values on separate lines
func hexs_to_ints(hex_str string, bytes int) ([]int) {
	var ints []int
	str_split := (strings.Split(strings.TrimSpace(hex_str), "\n"))
	for _, str := range str_split {
		num, _ := strconv.ParseInt(str, 16, 64)
		num_shr := uint32(num)
		for b:=0; b<bytes; b++ { 
			num_byte := int(uint8(num_shr))
			ints = append(ints, num_byte)
			num_shr >>= 8
		}
	}
	return ints
}

// convert slice of ints to string of multi-byte hex values
// hex string values on separate lines
func ints_to_hexs(ints []int, bytes int) (string) {
	var hex_str string
	for i:=0; i<len(ints); i+=bytes {
		var line_int int64
		for b:=0; b<bytes; b++ { 
			line_int += int64(uint8(ints[i+b])) << (b * 8)
		}
		hex_str += strconv.FormatInt(line_int, 16) + "\n"
	}
	return strings.TrimSpace(hex_str)
}
