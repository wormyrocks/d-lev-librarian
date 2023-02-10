package main

/*
 * Librarian for the D-Lev Theremin
 * See file "CHANGE_LOG.TXT" for details
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
	port, _ := strconv.Atoi(cfg_get("port"))

	// print help file if no args
	if len(os.Args) < 2 {
		help(false)  // short help
    } else {
		// parse subcommands
		switch os.Args[1] {
			case "help", "-help", "--help", "-h", "--h": help(true)  // verbose help
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

// list free serial ports / set port
func ports(port int) int {
	ports := sp_list()
	if len(ports) == 0 {
		fmt.Println("> No serial ports found!")
	} else {
		fmt.Println("> Available serial ports:")
		for p_num, p_str := range ports { fmt.Printf("  %v : %v\n", p_num, p_str) }
	}
	if len(os.Args) > 2 { 
		cfg_set("port", os.Args[2])
		port, _ = strconv.Atoi(cfg_get("port"))
		fmt.Println("> Set port to:", port)
	} else {
		fmt.Println("> Current port:", port)
	}
	return port
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
		rd_str := sp_wr_rd(sp, wr_str, false)  // note added space 
		fmt.Println(rd_str)
		sp.Close()
	}
}

// view knobs, DLP file, slot
func view(port int) {
	sub := flag.NewFlagSet("view", flag.ExitOnError)
	file := sub.String("f", "", "view `file` name")
	pro := sub.Bool("pro", false, "slot/file profile mode")
	knobs := sub.Bool("k", false, "view knobs` `")
	slot := sub.Int("s", 0, "view `slot` number (default 0)")
	sub.Parse(os.Args[2:])
	//
	mode := "pre"
	if *pro { mode = "pro" }
	if *knobs {  // view current knobs
		sp := sp_open(port)
		rx_str := sp_wr_rd(sp, "0 159 rk ", false)
		rx_str = decruft_hcl(rx_str)
		fmt.Println(ui_prn_str(knob_ui_strs(rx_str)))
		fmt.Println("> knobs")
		sp.Close()
	} else if *file != "" {  // view a *.dlp file
		*file = file_ext_chk(*file, ".dlp")
		file_bytes, err := os.ReadFile(*file); if err != nil { log.Fatal(err) }
		fmt.Println(ui_prn_str(pre_ui_strs(string(file_bytes), *pro)))
		fmt.Println(">", mode, "file", *file)
	} else {  // view a slot
		sp := sp_open(port)
		addr := spi_slot_addr(*slot, mode)
		rx_str := spi_rd(sp, addr, addr + EE_PG_BYTES - 1, false)
		fmt.Println(ui_prn_str(pre_ui_strs(rx_str, *pro)))
		fmt.Println(">", mode, "slot", *slot)
		sp.Close()
	}
}

// match slot contents w/ DLP files in dir & list
func slots(port int) {
	dir := "."  // default
	if len(os.Args) > 2 { dir = os.Args[2] }
	dir = filepath.Clean(dir)
	file_map := map_files(dir, ".dlp")
	sp := sp_open(port)
	addr, end := spi_bulk_addrs("pre")
	pre_str := spi_rd(sp, addr, end - 1, true)
	fmt.Print(slots_prn_str(map_slots(pre_str, file_map)))
	fmt.Println("> slots matched to *.dlp files in directory", dir)
	sp.Close()
}

// slot names => *.bnk
func stob(port int) {
	sub := flag.NewFlagSet("stob", flag.ExitOnError)
	file := sub.String("f", "", "target `file` name")
	sub.Parse(os.Args[2:])
	//
	*file = file_ext_chk(*file, ".bnk")
	file_exists_chk(*file)
	dir, _ := filepath.Split(*file)
	dir = filepath.Clean(dir)
	file_map := map_files(dir, ".dlp")
	addr, end := spi_bulk_addrs("pre")
	sp := sp_open(port)
	pre_str := spi_rd(sp, addr, end - 1, true)
	slot_names := map_slots(pre_str, file_map)
	str := slots_bnk_str(slot_names)
	err := os.WriteFile(*file, []byte(str), 0666); if err != nil { log.Fatal(err) }
	fmt.Println("> slots matched to *.dlp files written to file", *file)
	sp.Close()
}

//////////////
// download //
//////////////

// dump to file
func dump(port int) {
	sub := flag.NewFlagSet("dump", flag.ExitOnError)
	file := sub.String("f", "", "target `file` name")
	pre := sub.Bool("pre", false, "preset mode")
	pro := sub.Bool("pro", false, "profile mode")
	spi := sub.Bool("spi", false, "spi mode")
	eeprom := sub.Bool("eeprom", false, "eeprom mode")
	sub.Parse(os.Args[2:])
	//
	mode := ""
	if *pre { mode = "pre"
	} else if *pro { mode = "pro"
	} else if *spi { mode = "spi"
	} else if *eeprom { mode = "eeprom" }
	*file = file_ext_chk(*file, "." + mode)
	file_exists_chk(*file)
	addr, end := spi_bulk_addrs(mode)
	sp := sp_open(port)
	rx_str := spi_rd(sp, addr, end - 1, true)
	err := os.WriteFile(*file, []byte(rx_str), 0666); if err != nil { log.Fatal(err) }
	fmt.Println("> downloaded to", mode, "file", *file) 
	sp.Close()
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
	*file = file_ext_chk(*file, ".dlp")
	file_exists_chk(*file)
	sp := sp_open(port)
	rx_str := sp_wr_rd(sp, "0 159 rk ", false)
	rx_str = decruft_hcl(rx_str)
	kints := hexs_to_ints(rx_str, 1)
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
	sp.Close()
}

// slot => *.dlp
func stof(port int) {
	sub := flag.NewFlagSet("stof", flag.ExitOnError)
	slot := sub.Int("s", 0, "source `slot` number (default 0)")
	file := sub.String("f", "", "target `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	mode := "pre"
	if *pro { mode = "pro" }
	*file = file_ext_chk(*file, ".dlp")
	file_exists_chk(*file)
	addr := spi_slot_addr(*slot, mode)
	sp := sp_open(port)
	rx_str := spi_rd(sp, addr, addr + EE_PG_BYTES - 1, false)
	err := os.WriteFile(*file, []byte(rx_str), 0666); if err != nil { log.Fatal(err) }
	fmt.Println("> downloaded", mode, "slot", *slot, "to", mode, "file", *file) 
	sp.Close()
}

////////////
// upload //
////////////

// pump from file
func pump(port int) {
	sub := flag.NewFlagSet("pump", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	pre := sub.Bool("pre", false, "preset mode")
	pro := sub.Bool("pro", false, "profile mode")
	spi := sub.Bool("spi", false, "spi mode")
	eeprom := sub.Bool("eeprom", false, "eeprom mode")
	sub.Parse(os.Args[2:])
	//
	mode := ""
	if *pre { mode = "pre"
	} else if *pro { mode = "pro"
	} else if *spi { mode = "spi"
	} else if *eeprom { mode = "eeprom" }
	*file = file_ext_chk(*file, "." + mode)
	addr, _ := spi_bulk_addrs(mode)
	file_bytes, err := os.ReadFile(*file); if err != nil { log.Fatal(err) }
	sp := sp_open(port)
	spi_wr(sp, addr, string(file_bytes), mode, true)
	fmt.Println("> uploaded", mode, "file", *file)
	sp.Close()
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
	*file = file_ext_chk(*file, ".dlp")
	file_bytes, err := os.ReadFile(*file); if err != nil { log.Fatal(err) }
	pints := hexs_to_ints(string(file_bytes), 4)
	sp := sp_open(port)
	for kidx, kname := range knob_pnames {
		_, _, pidx, pmode := pname_lookup(kname)
		if mode == pmode {
			wr_str := fmt.Sprint(kidx, " ", pints[pidx], " wk ")
			sp_wr_rd(sp, wr_str, false)
		}
	}
	fmt.Println("> uploaded", mode, "file", *file, "to", mode, "knobs") 
	sp.Close()
}

// *.dlp => slot
func ftos(port int) {
	sub := flag.NewFlagSet("ftos", flag.ExitOnError)
	slot := sub.Int("s", 0, "target `slot` number (default 0)")
	file := sub.String("f", "", "source `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	mode := "pre"
	if *pro { mode = "pro" }
	*file = file_ext_chk(*file, ".dlp")
	addr := spi_slot_addr(*slot, mode)
	file_bytes, err := os.ReadFile(*file); if err != nil { log.Fatal(err) }
	sp := sp_open(port)
	spi_wr(sp, addr, string(file_bytes), mode, false)
	fmt.Println("> uploaded", mode, "file", *file, "to", mode, "slot", *slot) 
	sp.Close()
}

// *.bnk => *.dlps => slots
func btos(port int) {
	sub := flag.NewFlagSet("btos", flag.ExitOnError)
	slot := sub.Int("s", 0, "starting `slot` number (default 0)")
	file := sub.String("f", "", "bank `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	mode := "pre"
	if *pro { mode = "pro" }
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
			addr := spi_slot_addr(*slot, mode)
			dlp_bytes, err := os.ReadFile(dlp_file); if err != nil { log.Fatal(err) }
			spi_wr(sp, addr, string(dlp_bytes), mode, false)
			fmt.Println("> uploaded", mode, "file", dlp_file, "to", mode, "slot", *slot)
			if *slot < 0 { *slot--
			} else { *slot++ }
		}
	}
	sp.Close()
}

/////////
// dev //
/////////

func dev() {

//	cfg_set("port", "20")
//	cfg_set("port", "20")
//	fmt.Println(cfg_get("port"))

/*
	// read, update, overwrite all *.dlp in given dir
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	dir := sub.String("dir", ".", "directory")
	pro := sub.Bool("pro", false, "profile flag")
	sub.Parse(os.Args[2:])
	//
	update_dlp(*dir, *pro)
*/

/*
	// generate DLP files with random content for testing
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	dir := sub.String("dir", ".", "directory")
	sub.Parse(os.Args[2:])
	//
	gen_test_dlps(*dir)
*/

/*
	// find DLP files with non-zero PITCH:corr
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	dir := sub.String("dir", ".", "directory")
	sub.Parse(os.Args[2:])
	//
	find_dlp(*dir)
*/
}
