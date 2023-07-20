package main

/*
 * d-lev support functions
*/

import (
	"fmt"
	"log"
	"strings"
	"strconv"
	"go.bug.st/serial"
	"time"
)

// show activity via printed dots
func dots(chars int) (int) {
	if chars > 0 {
		chars -= CHARS_PER_DOT
		fmt.Print(".") 
	}
	return chars
}

// read SPI port to string, trim cruft, optionally show activity
func spi_rd(sp serial.Port, addr int, addr_end int, act_f bool) (string) {
	rd_buf := sp_wr_rd(sp, strconv.Itoa(addr) + " " + strconv.Itoa(addr_end) + " rs ", act_f)
	rd_str := decruft_hcl(string(rd_buf))
	if len(strings.Split(rd_str, "\n")) != 1 + (addr_end - addr) / 4 { log.Fatalln("> Bad SPI read!") }
	return rd_str
}

// SPI write enable
func spi_wr_en(sp serial.Port) {
	sp_wr_rd(sp, "6 6 wr ", false)
	sp_wr_rd(sp, "6 0x100 wr ", false)  // csn hi
}

// SPI write & wait
func spi_wr_wait(sp serial.Port) {
	sp_wr_rd(sp, "6 0x100 wr ", false)  // csn hi
	time.Sleep(EE_WR_MS * time.Millisecond)
}

// SPI write protect & unprotect
func spi_wr_prot(sp serial.Port, prot_f bool) {
	spi_wr_en(sp)
	sp_wr_rd(sp, "6 1 wr ", false)  // wrsr reg
	if prot_f { sp_wr_rd(sp, "6 0xc wr ", false)
	} else { sp_wr_rd(sp, "6 0 wr ", false)	}
	spi_wr_wait(sp)
}

// write string to SPI port, optionally show activity
func spi_wr(sp serial.Port, addr int, wr_str string, act_f bool) {
	spi_wr_prot(sp, false)
	split_strs := (strings.Split(strings.TrimSpace(wr_str), "\n"))
	var chars int
	for _, line_str := range split_strs {
		var cmd string
		line_str := strings.TrimSpace(line_str)
		if addr % EE_PG_BYTES == 0 {  // page boundary
			spi_wr_wait(sp)
			spi_wr_en(sp)
			cmd = strconv.Itoa(addr) + " "
		}
		if line_str != "0" { cmd += "0x" }  // no 0x for zero data
		cmd += line_str + " ws "
		sp_wr_rd(sp, cmd, false)
		chars += len(cmd)
		addr += EE_RW_BYTES;
		if act_f { chars = dots(chars) }
	}
	// done
	spi_wr_wait(sp);
	spi_wr_prot(sp, true);
	if act_f { fmt.Println(" upload done") }
}

// return spi bulk addresses
func spi_bulk_addrs(mode string) (addr int, end int) {
	switch mode {
		case ".pre" :
			addr = EE_PRE_ADDR
			end = EE_PRE_END
		case ".pro" :
			addr = EE_PRO_ADDR
			end = EE_PRO_END
		case ".spi" :
			addr = EE_SPI_ADDR
			end = EE_SPI_END
		case ".eeprom" :
			addr = EE_START
			end = EE_END
		default :
			log.Fatalln("> Unknown mode:", mode)
	}
	return
}

// return spi slot addr
func spi_slot_addr(slot int, mode string) (int) {
	switch mode {
		case "pre" :
			if slot < 0 || slot >= PRE_SLOTS { log.Fatalln("- Slot out of range:", slot) }
		case "pro" :
			if slot < 0 || slot >= PRO_SLOTS { log.Fatalln("- Slot out of range:", slot) }
			slot += PRE_SLOTS
		default :
			log.Fatalln("> Unknown mode:", mode)
	}
	return slot * EE_PG_BYTES
}

// trim command, address, and prompt cruft from hcl read string
func decruft_hcl(str_i string) (string) {
	lines_i := strings.Split(strings.TrimSpace(str_i), "\n")
	lines_o := ""
	for idx, line := range lines_i {
		if (idx != 0) && (idx != len(lines_i) - 1) {
			line := strings.TrimSpace(line)
			addr_end := strings.Index(line, "]")
			lines_o += line[addr_end+1:] + "\n" 
		}
	}
	return strings.TrimSpace(lines_o)
}

// get single slot data string
func get_slot_str(slot int, mode string) (string) {
	addr := spi_slot_addr(slot, mode)
	sp := sp_open()
	rx_str := spi_rd(sp, addr, addr + EE_PG_BYTES - 1, false)
	sp.Close()
	return rx_str
}

// get all slots data strings
func get_slots_strs() ([]string) {
	addr, _ := spi_bulk_addrs(".pre")
	_, end := spi_bulk_addrs(".pro")
	sp := sp_open()
	rx_str := spi_rd(sp, addr, end - 1, true)
	sp.Close()
	split_strs := strings.Split(rx_str, "\n")
	if len(split_strs) < SLOTS * SLOT_BYTES/4 { log.Fatalln("> Bad slots info!") }
	var strs []string
	for s:=0; s<SLOTS; s++ {
		pre_str := ""
		for i:=s*SLOTS/4; i<(s+1)*SLOTS/4; i++ {
			pre_str += split_strs[i] + "\n"
		}
		strs = append(strs, pre_str)
	}
	return strs
}
