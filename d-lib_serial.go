package main

/*
 * d-lev support functions
*/

import (
	"fmt"
	"log"
	"bytes"
	"go.bug.st/serial"
)

// return a list of serial ports
func sp_list() ([]string) {
	ports, err := serial.GetPortsList(); if err != nil { log.Fatal(err) }
	return ports
}

// open enumerated serial port
func sp_open(port int) (serial.Port) {
	// check port
	ports := sp_list()
	if port >= len(ports) || port < 0 { 
		log.Fatalln("> Invalid port:", port) 
	}
	// config as 230400bps N81
	mode := &serial.Mode{
		BaudRate: 230400,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	// open port
	sp, err := serial.Open(ports[port], mode); if err != nil { log.Fatal(err) }
	return sp
}


// write & read serial port, string i/o, optionally show activity
func sp_wr_rd(sp serial.Port, wr_str string, act_f bool) (string) {
	// write to port
	_, err := sp.Write([]byte(wr_str)); if err != nil { log.Fatal(err) }
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
	if act_f { fmt.Println(" done!") }
	return rd_bytes.String()
}
