package main

/*
 * Librarian for the D-Lev Theremin
 * See file "CHANGE_LOG.txt" for details
*/

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"strconv"
)

func main() {

	// init from config file
	port, err := strconv.Atoi(cfg_get("port"))
	if err != nil { 
		cfg_set("port", "0")
		log.Fatalln("> Bad port in config file, setting port to 0!") 
	}

	// print help file if no args
	if len(os.Args) < 2 {
		help(false)  // short help
    } else {
		// parse subcommands
		switch os.Args[1] {
			case "help", "-help", "--help", "-h", "--h": help(true)  // verbose help
			case "reset": reset(port)
			case "ports": port = ports(port)
			case "slots": slots(port)
			case "view": view(port)
			case "ktof": ktof(port)
			case "ftok": ftok(port)
			case "stof": stof(port)
			case "ftos": ftos(port)
			case "stob": stob(port)
			case "btos": btos(port)
			case "dump": dump(port)
			case "pump": pump(port)
			case "hcl": hcl(port)
			case "split": split()
			case "diff": diff(port)
			case "update": update()  // update stuff
			case "dev": dev()  // dev stuff
			default: log.Fatalln("> Unknown command:", os.Args[1])
		}
	}

}  // end of main()


////////////////////
// main functions //
////////////////////

// print version & help info
func help(verbose_f bool) {
	fmt.Print("= D-Lev Librarian version ", VERSION, " =\n") 
	fmt.Print(help_str) 
	if verbose_f { fmt.Print(help_verbose_str) }  // print verbose help
}

// do reset
func reset(port int) {
	sp := sp_open(port)
	rd_str := sp_wr_rd(sp, "0 0xff000000 wr ", false)
	fmt.Print(rd_str)
	fmt.Println(" reset processor")
	sp.Close()
}

// do HCL command
func hcl(port int) {
	if len(os.Args) < 3 { 
		fmt.Println("> Command line is blank!")
	} else {
		sp := sp_open(port)
		wr_str := ""
		for _, cmd := range os.Args[2:] {
			wr_str += cmd + " "
		}
		rd_str := sp_wr_rd(sp, wr_str, false)
		fmt.Print(rd_str)
		fmt.Println(" issued hcl command:", wr_str)
		sp.Close()
	}
}

// list free serial ports / set port
func ports(port int) (int) {
	sub := flag.NewFlagSet("ports", flag.ExitOnError)
	port_str := sub.String("p", "", "`port` number")
	sub.Parse(os.Args[2:])
	//
	ports := sp_list()
	if len(ports) == 0 {
		fmt.Println("> No serial ports found!")
	} else {
		fmt.Println("> available serial ports:")
		for p_num, p_str := range ports { fmt.Printf("  %v : %v\n", p_num, p_str) }
	}
	if *port_str != "" { 
		port_new, err := strconv.Atoi(*port_str)
		if err != nil { log.Fatalln("> Bad port number!") }
		if port_new >= len(ports) { log.Fatalln("> Port number out of range!") }
		cfg_set("port", *port_str)
		port = port_new
		fmt.Println("> set port to:", port)
	} else {
		fmt.Println("> current port:", port)
	}
	return port
}

// view knobs, DLP file, slot
func view(port int) {
	sub := flag.NewFlagSet("view", flag.ExitOnError)
	file := sub.String("f", "", "view `file` name")
	pro := sub.Bool("pro", false, "slot/file profile mode")
	knobs := sub.Bool("k", false, "view knobs")
	slot := sub.String("s", "", "view `slot` number")
	sub.Parse(os.Args[2:])
	//
	mode := "pre"
	if *pro { mode = "pro" }
	if *knobs {  // view current knobs
		rx_str := get_knobs(port)
		fmt.Println(ui_prn_str(knob_ui_strs(rx_str)))
		fmt.Println("> knobs")
	} else if *file != "" {  // view a *.dlp file
		file_bytes := get_dlp(file)
		fmt.Println(ui_prn_str(pre_ui_strs(string(file_bytes), *pro)))
		fmt.Println(">", mode, "file", *file)
	} else if *slot != "" {  // view a slot
		slot_int, err := strconv.Atoi(*slot); if err != nil { log.Fatalln("> Bad slot number!") }
		rx_str := get_slot(port, slot_int, mode)
		fmt.Println(ui_prn_str(pre_ui_strs(rx_str, *pro)))
		fmt.Println(">", mode, "slot", slot_int)
	} else {
		log.Fatalln("> Nothing to do!")
	}
}

// diff DLP files
func diff(port int) {
	sub := flag.NewFlagSet("diff", flag.ExitOnError)
	file := sub.String("f", "", "base `file` name")
	file2 := sub.String("f2", "", "compare `file` name")
	pro := sub.Bool("pro", false, "slot/file profile mode")
	knobs := sub.Bool("k", false, "compare knobs")
	slot := sub.String("s", "", "compare `slot` number")
	sub.Parse(os.Args[2:])
	//
	mode := "pre"
	if *pro { mode = "pro" }
	file_bytes := get_dlp(file)
	if *knobs {  // compare knobs
		rx_str := knob_pre_order(get_knobs(port), *pro)
		fmt.Println(diff_prn_str(diff_strs(string(file_bytes), rx_str, *pro)))
		fmt.Println(">", mode, "file", *file, "vs. knobs" )
	} else if *file2 != "" {  // compare a *.dlp file
		file2_bytes := get_dlp(file2)
		fmt.Println(diff_prn_str(diff_strs(string(file_bytes), string(file2_bytes), *pro)))
		fmt.Println(">", mode, "file", *file, "vs.", *file2 )
	} else if *slot != "" {  // compare a slot
		slot_int, err := strconv.Atoi(*slot); if err != nil { log.Fatalln("> Bad slot number!") }
		rx_str := get_slot(port, slot_int, mode)
		fmt.Println(diff_prn_str(diff_strs(string(file_bytes), rx_str, *pro)))
		fmt.Println(">", mode, "file", *file, "vs. slot", *slot )
	} else {
		log.Fatalln("> Nothing to do!")
	}
}

// match slot contents w/ DLP files in dir & list
func slots(port int) {
	sub := flag.NewFlagSet("slots", flag.ExitOnError)
	dir := sub.String("d", ".", "`directory` name")
	sub.Parse(os.Args[2:])
	//
	*dir = filepath.Clean(*dir)
	file_map := map_files(*dir, ".dlp")
	rx_str := get_slots(port)
	fmt.Print(slots_prn_str(map_slots(rx_str, file_map)))
	fmt.Println("> slots matched to *.dlp files in directory", *dir)
}

// slot names => *.bnk
func stob(port int) {
	sub := flag.NewFlagSet("stob", flag.ExitOnError)
	file := sub.String("f", "", "target `file` name")
	headers := sub.Bool("hdr", false, "include headers")
	sub.Parse(os.Args[2:])
	//
	file_blank_chk(*file)
	*file = file_ext_chk(*file, ".bnk")
	file_exists_chk(*file)
	dir, _ := filepath.Split(*file)
	dir = filepath.Clean(dir)
	file_map := map_files(dir, ".dlp")
	rx_str := get_slots(port)
	str := slots_bnk_str(map_slots(rx_str, file_map), *headers)
	err := os.WriteFile(*file, []byte(str), 0666); if err != nil { log.Fatal(err) }
	fmt.Println("> slots matched to *.dlp files written to file", *file)
}


//////////////////////
// download to file //
//////////////////////

// dump to file
func dump(port int) {
	sub := flag.NewFlagSet("dump", flag.ExitOnError)
	file := sub.String("f", "", "target `file` name")
	sub.Parse(os.Args[2:])
	//
	file_blank_chk(*file)
	file_exists_chk(*file)
	ext := strings.Trim(filepath.Ext(*file), ".")
	switch ext {
		case "pre", "pro", "spi", "eeprom" : // these are OK
		default : log.Fatalln("> Unknown file extension", ext)
	}
	addr, end := spi_bulk_addrs(ext)
	sp := sp_open(port)
	rx_str := spi_rd(sp, addr, end - 1, true)
	sp.Close()
	err := os.WriteFile(*file, []byte(rx_str), 0666); if err != nil { log.Fatal(err) }
	fmt.Println("> dumped to", *file) 
}

// knobs => *.dlp
func ktof(port int) {
	sub := flag.NewFlagSet("ktof", flag.ExitOnError)
	file := sub.String("f", "", "target `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	mode := "pre"
	if *pro { mode = "pro" }
	file_blank_chk(*file)
	*file = file_ext_chk(*file, ".dlp")
	file_exists_chk(*file)
	rx_str := get_knobs(port)
	kints := hexs_to_ints(rx_str, 1)
	if len(kints) != KNOBS { log.Fatalln("> Bad knob info!") }
	pints := make([]int, SLOT_BYTES)
	for kidx, kname := range knob_pnames {
		_, _, pidx, pmode := pname_lookup(kname)
		if mode == pmode {
			pints[pidx] = kints[kidx]
		}
	}
	hexs := ints_to_hexs(pints, 4)
	err := os.WriteFile(*file, []byte(hexs), 0666); if err != nil { log.Fatal(err) }
	fmt.Println("> downloaded", mode, "knobs to", mode, "file", *file) 
}

// slot => *.dlp
func stof(port int) {
	sub := flag.NewFlagSet("stof", flag.ExitOnError)
	slot := sub.String("s", "", "source `slot` number")
	file := sub.String("f", "", "target `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	slot_int, err := strconv.Atoi(*slot); if err != nil { log.Fatalln("> Bad slot number!") }
	mode := "pre"
	if *pro { mode = "pro" }
	file_blank_chk(*file)
	*file = file_ext_chk(*file, ".dlp")
	file_exists_chk(*file)
	rx_str := get_slot(port, slot_int, mode)
	err = os.WriteFile(*file, []byte(rx_str), 0666); if err != nil { log.Fatal(err) }
	fmt.Println("> downloaded", mode, "slot", slot_int, "to", mode, "file", *file) 
}


//////////////////////
// upload from file //
//////////////////////

// pump from file
func pump(port int) {
	sub := flag.NewFlagSet("pump", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	sub.Parse(os.Args[2:])
	//
	file_blank_chk(*file)
	ext := strings.Trim(filepath.Ext(*file), ".")
	switch ext {
		case "pre", "pro", "spi", "eeprom" : // these are OK
		default : log.Fatalln("> Unknown file extension", ext)
	}
	addr, _ := spi_bulk_addrs(ext)
	file_bytes, err := os.ReadFile(*file); if err != nil { log.Fatal(err) }
	sp := sp_open(port)
	spi_wr(sp, addr, string(file_bytes), ext, true)
	sp.Close()
	fmt.Println("> pumped from", *file)
	if ext == "spi" || ext == "eeprom" { reset(port) }
}

// *.dlp => knobs
func ftok(port int) {
	sub := flag.NewFlagSet("ftok", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	mode := "pre"
	if *pro { mode = "pro" }
	file_bytes := get_dlp(file)
	pints := hexs_to_ints(string(file_bytes), 4)
	if len(pints) < SLOT_BYTES { log.Fatalln("> Bad file info!") }
	sp := sp_open(port)
	for kidx, kname := range knob_pnames {
		_, _, pidx, pmode := pname_lookup(kname)
		if mode == pmode {
			wr_str := fmt.Sprint(kidx, " ", pints[pidx], " wk ")
			sp_wr_rd(sp, wr_str, false)
		}
	}
	sp.Close()
	fmt.Println("> uploaded", mode, "file", *file, "to", mode, "knobs") 
}

// *.dlp => slot
func ftos(port int) {
	sub := flag.NewFlagSet("ftos", flag.ExitOnError)
	slot := sub.String("s", "", "target `slot` number")
	file := sub.String("f", "", "source `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	slot_int, err := strconv.Atoi(*slot); if err != nil { log.Fatalln("> Bad slot number!") }
	mode := "pre"
	if *pro { mode = "pro" }
	file_bytes := get_dlp(file)
	addr := spi_slot_addr(slot_int, mode)
	sp := sp_open(port)
	spi_wr(sp, addr, string(file_bytes), mode, false)
	sp.Close()
	fmt.Println("> uploaded", mode, "file", *file, "to", mode, "slot", slot_int) 
}

// *.bnk => *.dlps => slots
func btos(port int) {
	sub := flag.NewFlagSet("btos", flag.ExitOnError)
	slot := sub.String("s", "", "starting `slot` number")
	file := sub.String("f", "", "bank `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	slot_int, err := strconv.Atoi(*slot); if err != nil { log.Fatalln("> Bad slot number!") }
	mode := "pre"
	if *pro { mode = "pro" }
	file_blank_chk(*file)
	*file = file_ext_chk(*file, ".bnk")
	dir, bnk_file := filepath.Split(*file)
	dir = filepath.Clean(dir)
	bnk_bytes, err := os.ReadFile(filepath.Join(dir, bnk_file)); if err != nil { log.Fatal(err) }
	bnk_split := strings.Split(strings.TrimSpace(string(bnk_bytes)), "\n")
	sp := sp_open(port)
	for _, line := range bnk_split {
		line_str := strings.TrimSpace(string(line));
		if !strings.HasPrefix(line_str, "//") {  // skip commented lines
			dlp_file := file_ext_chk(filepath.Join(dir, line_str), ".dlp")
			addr := spi_slot_addr(slot_int, mode)
			dlp_bytes, err := os.ReadFile(dlp_file); if err != nil { log.Fatal(err) }
			spi_wr(sp, addr, string(dlp_bytes), mode, false)
			fmt.Println("> uploaded", mode, "file", dlp_file, "to", mode, "slot", slot_int)
			slot_int++
		}
	}
	sp.Close()
}

// split bulk files
func split() {
	sub := flag.NewFlagSet("split", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	sub.Parse(os.Args[2:])
	//
	file_blank_chk(*file)
	ext := strings.Trim(filepath.Ext(*file), ".")
	switch ext {
		case "pre", "pro", "eeprom" : // these are OK
		default : log.Fatalln("> Unknown file extension", ext)
	}
	file_bytes, err := os.ReadFile(*file); if err != nil { log.Fatal(err) }
	str_split := (strings.Split(strings.TrimSpace(string(file_bytes)), "\n"))
	if ext == "eeprom" {
		var pre_str string
		var pro_str string
		var spi_str string
		for line, str := range str_split {
			if line < PRE_SLOTS*SLOT_BYTES/4 { 
				pre_str += str + "\n"
			} else if line < SLOTS*SLOT_BYTES/4 { 
				pro_str += str + "\n"
			} else { 
				spi_str += str + "\n"
			}
		}
		base_file := strings.TrimSuffix(*file, filepath.Ext(*file))
		pre_file := base_file + ".pre"
		pro_file := base_file + ".pro"
		spi_file := base_file + ".spi"
		file_exists_chk(pre_file)
		err = os.WriteFile(pre_file, []byte(pre_str), 0666); if err != nil { log.Fatal(err) }
		file_exists_chk(pro_file)
		err = os.WriteFile(pro_file, []byte(pro_str), 0666); if err != nil { log.Fatal(err) }
		file_exists_chk(spi_file)
		err = os.WriteFile(spi_file, []byte(spi_str), 0666); if err != nil { log.Fatal(err) }
		fmt.Println("> split", *file, "to", pre_file, pro_file, spi_file )
	} else {  // pre | pro
		var dlp_str string
		file_num := 0
		for line, str := range str_split {
			dlp_str += str + "\n"
			if line % 64 == 63 { 
				pre_file := fmt.Sprintf("%03d", file_num) + ".dlp"
				if ext == "pro" { pre_file = "pro_" + pre_file }
				err = os.WriteFile(pre_file, []byte(dlp_str), 0666); if err != nil { log.Fatal(err) }
				file_num++
				dlp_str = ""
			}
		}
		fmt.Println("> split", *file, "to", file_num, "numbered *.dlp files" )
	}
}


////////////
// update //
////////////

// read, update, overwrite all *.dlp in given dir
func update() {
	sub := flag.NewFlagSet("slots", flag.ExitOnError)
	dir := sub.String("d", ".", "`directory` name")
	pro := sub.Bool("pro", false, "profile mode")
	dry := sub.Bool("dry", true, "dry-run mode")
	sub.Parse(os.Args[2:])
	//
	update_dlp(*dir, *pro, *dry)
}


/////////
// dev //
/////////

func dev() {

/*
	// generate DLP files with random content for testing
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	dir := sub.String("d", ".", "`directory` name")
	sub.Parse(os.Args[2:])
	//
	gen_test_dlps(*dir)
*/


	// find DLP files with certain characteristics
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	dir := sub.String("d", ".", "`directory` name")
	sub.Parse(os.Args[2:])
	//
	find_dlp(*dir)

}
