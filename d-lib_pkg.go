package main

/*
 * d-lev constants & helper functions
*/

import (
	"strings"
	"strconv"
)

const (
	LIB_VER = "8"										// current librarian version
	SW_DATE = "2023-06-20"								// current sw date (for pre, pro, spi, eeprom file names)
	SW_V8 = "7bbb846b"									// sw ver 2023-06-20
	SW_V7 = "73c6c3d7"									// sw ver 2023-05-24
	SW_V6 = "27c263bf"									// sw ver 2023-01-31
	SW_V5 = "2d58f653"									// sw ver 2023-01-01
	SW_V2 = "add46826"									// sw ver 2022-10-06
	SW_OV129 = "7bc1bd55"								// sw ver 2022-07-05
	SW_OV128 = "93152c8b"								// sw ver 2022-05-10
	SW_OV127 = "d202d35"								// sw ver 2022-05-04
	SW_OV126 = "af3f63c4"								// sw ver 2022-04-30
	SW_OV125 = "67517a97"								// sw ver 2022-04-17
	SW_OV124 = "5ba55477"								// sw ver 2022-03-17
	SW_OV121 = "7b6a0484"								// sw ver 2022-01-01
	SW_OV120 = "84f7f31c"								// sw ver 2021-12-18
	SW_OV119 = "52fe7d"									// sw ver 2021-12-04
	SW_OV115 = "240b1e68"								// sw ver 2021-10-31
	//
	SLOTS = 256											// pre + pro slots
	SLOT_BYTES = 256									// bytes per slot
	PRO_SLOTS = 6										// profile[0:5]
	PRE_SLOTS = SLOTS - PRO_SLOTS						// preset[0:249]
	//
	EE_RW_BYTES = 4										// eeprom bytes per read / write cycle
	EE_PG_BYTES = 256									// eeprom bytes per page
	//
	EE_PRE_ADDR = 0x0									// eeprom pre start addr
	EE_PRE_END = EE_PRE_ADDR + (PRE_SLOTS * SLOT_BYTES)	// eeprom pre end addr
	//
	EE_PRO_ADDR = EE_PRE_END							// eeprom pro start addr
	EE_PRO_END = EE_PRO_ADDR + (PRO_SLOTS * SLOT_BYTES)	// eeprom pro end addr
	//
	EE_SPI_ADDR = EE_PRO_END							// eeprom code start addr
	EE_SPI_SZ = 0x4000									// eeprom code size : 16kB code space
	EE_SPI_END = EE_SPI_ADDR + EE_SPI_SZ				// eeprom code end addr
	//
	EE_START = EE_PRE_ADDR								// eeprom start addr
	EE_END = EE_SPI_END									// eeprom end addr
	EE_WR_MS = 6										// eeprom write wait time (ms)
	//
	UI_PAGES = 20										// ui pages
	UI_PG_COLS = 2										// ui page columns
	UI_PG_ROWS = 4										// ui page rows
	UI_PG_KNOBS = UI_PG_COLS * UI_PG_ROWS				// ui page knobs
	UI_PAGE_KNOB = 7									// ui page selector knob
	UI_PRN_PG_COLS = 4									// ui print pages columns
	UI_PRN_PG_ROWS = 5									// ui print pages rows
	KNOBS = UI_PG_KNOBS * UI_PAGES						// total knobs
	//
	RX_BUF_BYTES = 512									// serial port rx buffer size
	CHARS_PER_DOT = 4096								// chars for each activity dot printed
	CFG_FILE = "d-lib.cfg"								// config file name
	WORK_DIR = "_WORK_"									// work scratch dir
	PRESETS_DIR = "_ALL_"								// presets dir
	CRC = "debb20e3"									// good CRC
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
	return hex_str
}

// check for hexness
func str_is_hex(str string) bool {
	if len(str) == 0 { return false }
	for _, ch := range str {
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) { return false }
	}
	return true
}

// return index of string in slice, else -1
func str_exists(strs []string, str string) (int) {
	for idx, entry := range strs { if str == entry { return idx } }
	return -1
}
