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
	// tx SPI read command, rx result
	rd_buf := sp_wr_rd(sp, strconv.Itoa(addr) + " " + strconv.Itoa(addr_end) + " rs ", act_f)
	// trim cruft
	return decruft_hcl(string(rd_buf))
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
// confine writes to mode sections
func spi_wr(sp serial.Port, addr int, wr_str string, mode string, act_f bool) {
	spi_wr_prot(sp, false)
	split_strs := (strings.Split(strings.TrimSpace(wr_str), "\n"))
	var chars int
	for _, line_str := range split_strs {
		var cmd string
		if spi_addr_chk(addr, mode) {  // do write
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
		} else {  // don't write
			chars += 20
		}
		addr += EE_RW_BYTES;
		if act_f { chars = dots(chars) }
	}
	// done
	spi_wr_wait(sp);
	spi_wr_prot(sp, true);
	if act_f { fmt.Println(" done!") }
}

// return spi bulk addresses
func spi_bulk_addrs(mode string) (addr int, end int) {
	switch mode {
		case "pre" :
			addr = EE_PRE_ADDR
			end = EE_PRE_END
		case "pro" :
			addr = EE_PRO_ADDR
			end = EE_PRO_END
		case "spi" :
			addr = EE_SPI_ADDR
			end = EE_SPI_END
		case "eeprom" :
			addr = EE_START
			end = EE_END
		default :
			log.Fatalln("> Unknown mode:", mode)
	}
	return
}

// confine writes to mode sections
func spi_addr_chk(addr int, mode string) (bool) {
	mode_addr, mode_end := spi_bulk_addrs(mode)
	mode_f := addr >= mode_addr && addr < mode_end
	if mode == "pre" {  // pro hole in pre section
		pro_addr, pro_end := spi_bulk_addrs("pro")
		pro_f := addr >= pro_addr && addr < pro_end
		mode_f = mode_f && !pro_f
	}
	return mode_f
}

// return spi slot addr
func spi_slot_addr(slot int, mode string) (int) {
	switch mode {
		case "pre" :
			if slot < -PRE_SLOT_MAX || slot > PRE_SLOT_MAX { log.Fatalln("- Slot out of range:", slot) }
			if slot < 0 { slot += SLOTS }
		case "pro" :
			if slot < -PRO_SLOT_MAX || slot > PRO_SLOT_MAX { log.Fatalln("- Slot out of range:", slot) }
			slot += SLOTS / 2
		default :
			log.Fatalln("> Unknown mode:", mode)
	}
	return slot * EE_PG_BYTES
}

// trim command, address, and prompt cruft from hcl read string
func decruft_hcl(str_i string) (string) {
	split_strs := strings.Split(strings.TrimSpace(str_i), "\n")
	var str_all string
	for _, line_str := range split_strs {
		line_str := strings.TrimSpace(line_str)
		idx := strings.Index(line_str, "]")
		if idx >= 0 { str_all += line_str[idx+1:] + "\n" }
	}
	return str_all
}
