package main

/*
 * d-lev support functions
*/

import (
	"fmt"
	"log"
	"bytes"
	"strconv"
	"strings"
	"go.bug.st/serial"
)

// return a list of serial ports
func sp_list() ([]string) {
	ports, err := serial.GetPortsList(); if err != nil { log.Fatal(err) }
	return ports
}

// open enumerated serial port
func sp_open() (serial.Port) {
	port := cfg_get("port")
	if port == "" {
		log.Fatal("Current port is not assigned!") 
	}
	// config as 230400bps N81
	mode := &serial.Mode{
		BaudRate: 230400,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	// open port
	sp, err := serial.Open(port, mode); if err != nil { log.Fatal("Error opening port: ", port, " !") }
	return sp
}

// write & read serial port, string i/o, optionally show activity
func sp_wr_rd(sp serial.Port, wr_str string, act_f bool) (string) {
	err := sp.ResetInputBuffer();  if err != nil { log.Fatal(err) }
	err = sp.ResetOutputBuffer(); if err != nil { log.Fatal(err) }
	// write to port
	_, err = sp.Write([]byte(wr_str)); if err != nil { log.Fatal(err) }
	// read port, concat to buffer
	var rd_bytes bytes.Buffer
	var chars int
	for {
		rd_buf := make([]byte, RX_BUF_BYTES)
		n, err := sp.Read(rd_buf); if err != nil { log.Fatal(err) }
		rd_bytes.Write(rd_buf[:n])  // concat
		if bytes.Contains(rd_buf[:n], []byte(">")) { break }  // read until prompt
		chars += n
		if act_f { chars = dots(chars) }
	}
	// done
	if act_f { fmt.Println(" download done") }
	return rd_bytes.String()
}

// get knob data string
func get_knob_str() (string) {
	sp := sp_open()
	rx_str := sp_wr_rd(sp, "0 " + strconv.Itoa(KNOBS-1) + " rk ", false)
	sp.Close()
	rx_str = decruft_hcl(rx_str)
	if strings.Count(rx_str, "\n") != KNOBS-1 { log.Fatalln("> Bad knob info!") }
	return rx_str
}	

// get knob pint data
func get_knob_pints(mode string) ([]int) {
	kints := hexs_to_ints(get_knob_str(), 1)
	return knob_pre_order(kints, mode)
}

// write knob pint data
func put_knob_pints(pints []int, mode string) {
	sp := sp_open()
	for kidx, kname := range knob_pnames {
		_, _, pidx, pmode := pname_lookup(kname)
		if mode == pmode {
			wr_str := fmt.Sprint(kidx, " ", pints[pidx], " wk ")
			sp_wr_rd(sp, wr_str, false)
		}
	}
	sp.Close()
}

