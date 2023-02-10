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
	"math/rand"
	"time"
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
			case "hcl": hcl(port)
			case "loop":loop(port)
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
			case "split": split()
			case "join": join()
			case "diff": diff(port)
			case "morph": morph(port)
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
		wr_str := ""
		for _, cmd := range os.Args[2:] {
			wr_str += cmd + " "
		}
		sp := sp_open(port)
		rd_str := sp_wr_rd(sp, wr_str, false)
		sp.Close()
		fmt.Print(rd_str)
		fmt.Println(" issued hcl command:", wr_str)
	}
}

// do loop command
func loop(port int) {
	if len(os.Args) < 3 { 
		fmt.Println("> Loop text is blank!")
	} else {
		wr_str := ""
		for _, arg := range os.Args[2:] {
			wr_str += arg + " "
		}
		wr_str = strings.TrimSpace(wr_str)
		sp := sp_open(port)
		rd_str := sp_wr_rd(sp, wr_str + ">", false)
		sp.Close()
		fmt.Println("> tx:", wr_str)
		fmt.Println("> rx:", strings.TrimSuffix(strings.TrimSpace(rd_str), ">"))
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
		if err != nil { 
			fmt.Println("> current port:", port)
			log.Fatalln("> Bad port number!") 
		} else if port_new >= len(ports) { 
			fmt.Println("> current port:", port)
			log.Fatalln("> Port number out of range!") 
		} else {
			cfg_set("port", *port_str)
			port = port_new
			fmt.Println("> set port to:", port)
		}
	} else if len(os.Args) > 2 { 
		fmt.Println("> current port:", port)
		log.Fatalln("> Use the -p flag to set the port!")
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
		knob_str := get_knob_str(port)
		fmt.Println(ui_prn_str(knob_ui_strs(knob_str)))
		fmt.Println("> knobs")
	} else if *file != "" {  // view a *.dlp file
		var file_str string
		*file, file_str = get_file_str(*file, ".dlp")
		fmt.Println(ui_prn_str(pre_ui_strs(file_str, *pro)))
		fmt.Println(">", mode, "file", *file)
	} else if *slot != "" {  // view a slot
		slot_int, err := strconv.Atoi(*slot); if err != nil { log.Fatalln("> Bad slot number!") }
		slot_str := get_slot_str(port, slot_int, mode)
		fmt.Println(ui_prn_str(pre_ui_strs(slot_str, *pro)))
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
	var file_str string
	*file, file_str = get_file_str(*file, ".dlp")
	if *knobs {  // compare knobs
		knob_str := knob_pre_str(get_knob_str(port), *pro)
		fmt.Println(diff_prn_str(diff_pres(file_str, knob_str, *pro)))
		fmt.Println(">", mode, "file", *file, "vs. knobs" )
	} else if *file2 != "" {  // compare a *.dlp file
		var file2_str string
		*file2, file2_str = get_file_str(*file2, ".dlp")
		fmt.Println(diff_prn_str(diff_pres(file_str, file2_str, *pro)))
		fmt.Println(">", mode, "file", *file, "vs.", *file2 )
	} else if *slot != "" {  // compare a slot
		slot_int, err := strconv.Atoi(*slot); if err != nil { log.Fatalln("> Bad slot number!") }
		slot_str := get_slot_str(port, slot_int, mode)
		fmt.Println(diff_prn_str(diff_pres(file_str, slot_str, *pro)))
		fmt.Println(">", mode, "file", *file, "vs. slot", *slot )
	} else {
		log.Fatalln("> Nothing to do!")
	}
}

// compare slots contents w/ DLP files in dir & list
func slots(port int) {
	sub := flag.NewFlagSet("slots", flag.ExitOnError)
	dir := sub.String("d", ".", "`directory` name")
	inf := sub.Bool("inf", false, "infer best guess")
	sub.Parse(os.Args[2:])
	//
	name_strs, data_strs := get_dir_strs(*dir, ".dlp")
	slots_strs := get_slots_strs(port)
	if *inf {
		fmt.Print(slots_prn_str(comp_slots(slots_strs, name_strs, data_strs, *inf)))
	} else {
		file_map := map_files(*dir, ".dlp")
		fmt.Print(slots_prn_str(map_slots(slots_strs, file_map)))
	}
	fmt.Println("> slots matched to *.dlp files in directory", *dir)
}

// slot names => *.bnk
func stob(port int) {
	sub := flag.NewFlagSet("stob", flag.ExitOnError)
	file := sub.String("f", "", "bank `file` name")
	headers := sub.Bool("hdr", false, "include headers")
	inf := sub.Bool("inf", false, "infer best guess")
	sub.Parse(os.Args[2:])
	//
	file_blank_chk(*file)
	*file = file_ext_chk(*file, ".bnk")
	file_exists_chk(*file)
	dir, _ := filepath.Split(*file)
	name_strs, data_strs := get_dir_strs(dir, ".dlp")
	slots_strs := get_slots_strs(port)
	str := slots_bnk_str(comp_slots(slots_strs, name_strs, data_strs, *inf), *headers)
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
	pints := get_knob_pints(port, mode)
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
	slot_str := get_slot_str(port, slot_int, mode)
	err = os.WriteFile(*file, []byte(slot_str), 0666); if err != nil { log.Fatal(err) }
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
	var file_str string
	*file, file_str = get_file_str(*file, ".dlp")
	pints := hexs_to_ints(file_str, 4)
	if len(pints) < SLOT_BYTES { log.Fatalln("> Bad file info!") }
	put_knob_pints(port, pints, mode)
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
	var file_str string
	*file, file_str = get_file_str(*file, ".dlp")
	addr := spi_slot_addr(slot_int, mode)
	sp := sp_open(port)
	spi_wr(sp, addr, file_str, mode, false)
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
	split_file(*file)
}


// join bulk files
func join() {
	sub := flag.NewFlagSet("join", flag.ExitOnError)
	file := sub.String("f", "", "target `file` name")
	sub.Parse(os.Args[2:])
	//
	join_files(*file)
}


///////////
// morph //
///////////

func morph(port int) {
	sub := flag.NewFlagSet("morph", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	knobs := sub.Bool("k", false, "source knobs")
	slot := sub.String("s", "", "source `slot` number")
	seed := sub.Int("seed", int(time.Now().UnixNano()), "random seed")
	mo := sub.Int("mo", 0, "oscillator `mult`iplier")
	mn := sub.Int("mn", 0, "noise `mult`iplier")
	me := sub.Int("me", 0, "eq (bass & treble) `mult`iplier")
	mf := sub.Int("mf", 0, "filter `mult`iplier")
	mr := sub.Int("mr", 0, "resonator `mult`iplier")
	sub.Parse(os.Args[2:])
	rand.Seed(int64(*seed))
	prn_str := ""
	var pints []int
	if *mo | *mn | *me | *mf | *mr == 0 {
		log.Fatalln("> Nothing to do!")
	} else if *knobs {  // morph current knobs
		pints = pints_signed(get_knob_pints(port, "pre"), false)
		prn_str = fmt.Sprint("> morphed knobs")
	} else if *file != "" {  // morph a *.dlp file
		var file_str string
		*file, file_str = get_file_str(*file, ".dlp")
		pints = pints_signed(hexs_to_ints(file_str, 4), false)
		prn_str = fmt.Sprint("> morphed file ", *file)
	} else if *slot != "" {  // morph a slot
		slot_int, err := strconv.Atoi(*slot); if err != nil { log.Fatalln("> Bad slot number!") }
		slot_str := get_slot_str(port, slot_int, "pre")
		pints = pints_signed(hexs_to_ints(slot_str, 4), false)
		prn_str = fmt.Sprint("> morphed slot ", slot_int)
	} else {
		log.Fatalln("> Nothing to do!")
	}
	prn_str += fmt.Sprint(" (-seed=", *seed, ")")
	pints = morph_pints(pints, *mo, *mn, *me, *mf, *mr)
	put_knob_pints(port, pints, "pre")
	fmt.Println(prn_str)
}


////////////
// update //
////////////

// read, update, overwrite all *.dlp in given dir
func update() {
	sub := flag.NewFlagSet("update", flag.ExitOnError)
	dir := sub.String("d", ".", "`directory` name")
	pro := sub.Bool("pro", false, "profile mode")
	dry := sub.Bool("dry", true, "dry-run mode")
	rob := sub.Bool("rob", false, "schwimmer mode")
	sub.Parse(os.Args[2:])
	//
	update_dlp(*dir, *pro, *dry, *rob)
}


/////////
// dev //
/////////

func dev() {

	// find DLP files with certain characteristics
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	dir := sub.String("d", ".", "`directory` name")
	sub.Parse(os.Args[2:])
	//
	find_dlp(*dir)

}
