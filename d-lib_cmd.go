package main

/*
 * d-lib support functions
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"math/rand"
)


// return first word from user input line
func user_word() (string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	fields := strings.Fields(scanner.Text())
	if len(fields) != 0 { return fields[0] }
	return ""
}

// pause and ask user yes | no | quit question
func user_prompt(prompt string, yes bool) (bool) {
	if yes { return true }
	fmt.Print("\n> ", prompt, " <y|ENTER|q>: ")
	input := user_word()
	switch input {
		case "q": log.Fatalln("> Quit, exiting program...")
		case "y": return true 
		default: fmt.Println("> -CANCELED-")

	}
	return false	
}

// pause and ask user text | quit question
func user_input(prompt string) (string) {
	fmt.Print("\n> ", prompt, " (q=quit): ")
	input := user_word()
	if input == "q" { log.Fatalln("> Quit, exiting program...") }
	return input
}

// print librarian version & help info
func help_cmd(verbose_f bool) {
	fmt.Print("= D-Lev Librarian version ", LIB_VER, " =\n") 
	fmt.Print(help_str) 
	if verbose_f { fmt.Print(help_verbose_str) }  // print verbose help
}

// do processor reset
func reset_cmd() {
	sp := sp_open()
	sp_wr_rd(sp, "0 0xff000000 wr ", false)
	fmt.Println("> Issued processor reset")
	sp.Close()
}

// return D-Lev software version
func get_ver() (string) {
	sp := sp_open()
	rd_str := sp_wr_rd(sp, "ver ", false)
	sp.Close()
	rd_str = decruft_hcl(rd_str)
	if !str_is_hex(rd_str) { log.Fatalln("> Something went wrong reading your VERSION:", rd_str) }
	return rd_str
}

// check D-Lev software crc
func check_crc() (bool) {
	sp := sp_open()
	rd_str := sp_wr_rd(sp, "crc ", false)
	sp.Close()
	rd_str = decruft_hcl(rd_str)
	if !str_is_hex(rd_str) { log.Fatalln("> Something went wrong reading your CRC:", rd_str) }
	if rd_str == CRC { return true }
	return false
}

// get versions
func ver_cmd() {
	fmt.Println("> Librarian version:", LIB_VER)
	sw_pre_chk(false)
}

// do ACAL
func acal_cmd() {
	sp := sp_open()
	sp_wr_rd(sp, "acal ", false)
	fmt.Println(" Issued ACAL")
	sp.Close()
}

// do HCL command
func hcl_cmd() {
	if len(os.Args) < 3 { 
		fmt.Println("> Command line is blank!")
	} else {
		wr_str := ""
		for _, cmd := range os.Args[2:] {
			wr_str += cmd + " "
		}
		sp := sp_open()
		rd_str := sp_wr_rd(sp, wr_str, false)
		sp.Close()
		fmt.Print(rd_str)
		fmt.Println(" Issued hcl command:", wr_str)
	}
}

// do loop command
func loop_cmd() {
	if len(os.Args) < 3 { 
		fmt.Println("> Loop text is blank!")
	} else {
		wr_str := ""
		for _, arg := range os.Args[2:] {
			wr_str += arg + " "
		}
		wr_str = strings.TrimSpace(wr_str)
		sp := sp_open()
		rd_str := sp_wr_rd(sp, wr_str + ">", false)
		sp.Close()
		fmt.Println("> TX:", wr_str)
		fmt.Println("> RX:", strings.TrimSuffix(strings.TrimSpace(rd_str), ">"))
	}
}

// list free serial ports / set port
func ports_cmd(port_new string) {
	port := cfg_get("port")
	port_list := sp_list()
	port_idx := str_exists(port_list, port)
	if len(port_list) == 0 {
		fmt.Println("> No serial ports found!")
	} else {
		fmt.Println("> Available serial ports:")
		for p_num, p_str := range port_list { fmt.Printf(" [%v] %v\n", p_num, p_str) }
	}
	if port_new != "" { 
		port_num, err := strconv.Atoi(port_new)
		if err != nil || port_num < 0 || port_num >= len(port_list) { 
			log.Fatalln("> Bad port number!") 
		} else {
			port = port_list[port_num]
			cfg_set("port", port)
			fmt.Print("> Set port to: [", port_num, "] ", port, "\n")
		}
	} else if len(os.Args) > 2 { 
		log.Fatalln("> Use the -p flag to set the port!")
	} else if port == "" {
		fmt.Println("> Current port is not assigned!")
	} else if port_idx < 0 {
		fmt.Println("> Current port:", port, "doesn't exist!")
	} else {
		fmt.Print("> Current port: [", port_idx, "] ", port, "\n")
	}
}

// view knobs, DLP file, slot
func view_cmd(file string, pro, knobs bool, slot string) {
	mode := "pre"
	if pro { mode = "pro" }
	if knobs {  // view current knobs
		knob_str := get_knob_str()
		fmt.Println(ui_prn_str(knob_ui_strs(knob_str)))
		fmt.Println("> knobs")
	} else if file != "" {  // view a *.dlp file
		file = file_read_chk(file, ".dlp")
		file_str := file_read_str(file)
		fmt.Println(ui_prn_str(pre_ui_strs(file_str, pro)))
		fmt.Println(">", mode, "file", file)
	} else if slot != "" {  // view a slot
		slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
		slot_str := get_slot_str(slot_int, mode)
		fmt.Println(ui_prn_str(pre_ui_strs(slot_str, pro)))
		fmt.Println(">", mode, "slot", slot_int)
	} else {
		log.Fatalln("> Nothing to do!")
	}
}

// twiddle knob
func knob_cmd(knob, offset, val string) {
	str_split := (strings.Split(strings.TrimSpace(knob), ":"))
	if len(str_split) < 2  { log.Fatalln("> Bad knob value!") }
	knob_int, err := strconv.Atoi(str_split[1]); if err != nil { log.Fatalln("> Bad knob value!") }
	if knob_int < 0 || knob_int > UI_PG_KNOBS - 2 { log.Fatalln("> Bad knob value!") }
	pg_name, pg_idx := page_lookup(str_split[0])
	if pg_idx < 0 { log.Fatalln("> Bad page name!") }
	knob_idx := knob_int + pg_idx * UI_PG_KNOBS
	ptype, plabel, _, _ := pname_lookup(knob_pnames[knob_idx])
	sp := sp_open()
	rd_str := sp_wr_rd(sp, strconv.Itoa(knob_idx) + " rk ", false)
	sp.Close()
	rd_uint, _ := strconv.ParseInt(decruft_hcl(rd_str), 16, 32)
	fmt.Print("> ", pg_name, ":", strings.TrimSpace(plabel), "[", strings.TrimSpace(enc_disp(int(rd_uint), ptype)), "]")
	if offset != "" || val != "" {
		rw_int := ptype_signed(ptype, int(rd_uint))
		min := ptype_min(ptype)
		max := ptype_max(ptype)
		if offset != "" {
			offset_int, err := strconv.Atoi(offset); if err != nil { log.Fatalln("> Bad offset value!") }
			rw_int += offset_int
			if rw_int > max { rw_int = max }
			if rw_int < min { rw_int = min }
		} else {
			val_int, err := strconv.Atoi(val); if err != nil { log.Fatalln("> Bad set value!") }
			rw_int = val_int
			if rw_int > max { rw_int = max }
			if rw_int < min { rw_int = min }
		}
		sp := sp_open()
		sp_wr_rd(sp, strconv.Itoa(knob_idx) + " " + strconv.Itoa(rw_int) + " wk ", false)
		sp.Close()
		fmt.Print("=>[", strings.TrimSpace(enc_disp(rw_int, ptype)), "]")
	}
	fmt.Println("")
}

// diff DLP file(s) / slot(s) / knobs
func diff_cmd(file, file2 string, pro, knobs bool, slot, slot2 string) {
	mode := "pre"
	if pro { mode = "pro" }
	if file != "" {  // compare to a *.dlp file
		file = file_read_chk(file, ".dlp")
		file_str := file_read_str(file)
		if knobs {  // file vs. knobs
			knob_str := knob_pre_str(get_knob_str(), pro)
			fmt.Println(diff_prn_str(diff_pres(file_str, knob_str, pro)))
			fmt.Println(">", mode, "file", file, "vs. knobs" )
		} else if file2 != "" {  // file vs. file2
			file2 = file_read_chk(file2, ".dlp")
			file2_str := file_read_str(file2)
			fmt.Println(diff_prn_str(diff_pres(file_str, file2_str, pro)))
			fmt.Println(">", mode, "file", file, "vs.", file2 )
		} else if slot != "" {  // file vs. slot
			slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
			slot_str := get_slot_str(slot_int, mode)
			fmt.Println(diff_prn_str(diff_pres(file_str, slot_str, pro)))
			fmt.Println(">", mode, "file", file, "vs. slot", slot )
		} else {
			log.Fatalln("> Nothing to do!")
		}
	} else if slot != "" {  // compare to a slot
		slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
		slot_str := get_slot_str(slot_int, mode)
		if knobs {  // slot vs. knobs
			knob_str := knob_pre_str(get_knob_str(), pro)
			fmt.Println(diff_prn_str(diff_pres(slot_str, knob_str, pro)))
			fmt.Println(">", mode, "slot", slot, "vs. knobs" )
		} else if slot2 != "" {  // slot vs. slot2
			slot2_int, err := strconv.Atoi(slot2); if err != nil { log.Fatalln("> Bad slot number!") }
			slot2_str := get_slot_str(slot2_int, mode)
			fmt.Println(diff_prn_str(diff_pres(slot_str, slot2_str, pro)))
			fmt.Println(">", mode, "slot", slot, "vs. slot", slot2 )
		} else {
			log.Fatalln("> Nothing to do!")
		}
	} else {
		log.Fatalln("> Nothing to do!")
	}
}

// match slots / DLP files w/ DLP files & list
func match_cmd(dir, dir2 string, pro, hdr, guess, slots bool) {
	name_strs, data_strs := get_dir_strs(dir, ".dlp")
	mode := "pre"
	if pro { mode = "pro" }
	if len(data_strs) == 0 {  log.Fatalln("> No", mode, "files in", dir) }
	if slots {
		slots_strs := get_slots_strs()
		fmt.Print(slots_prn_str(comp_file_data(slots_strs, name_strs, data_strs, pro, guess), pro, hdr))
		fmt.Println("> matched", mode, "slots to", mode, "files in", dir)
	} else {
		name2_strs, data2_strs := get_dir_strs(dir2, ".dlp")
		if len(data2_strs) == 0 {  log.Fatalln("> No", mode, "files in", dir2) }
		fmt.Print(files_prn_str(name2_strs, comp_file_data(data2_strs, name_strs, data_strs, pro, guess)))
		fmt.Println("> matched", mode, "files in", dir2, "to", mode, "files in", dir)
	}
}

func dump_cmd(file string, yes bool) {
	file_blank_chk(file)
	ext := filepath.Ext(file)
	switch ext {
		case ".pre", ".pro", ".spi", ".eeprom" : // these are OK
		default : log.Fatalln("> Unknown file extension:", ext)
	}
	addr, end := spi_bulk_addrs(ext)
	sp := sp_open()
	rx_str := spi_rd(sp, addr, end - 1, true)
	sp.Close()
	if file_write_str(file, rx_str, yes) {
		fmt.Println("> dumped to", file) 
	}
}

// knobs => *.dlp
func ktof_cmd(file string, pro, yes bool) {
	file_blank_chk(file)
	file = file_ext_chk(file, ".dlp")
	mode := "pre"
	if pro { mode = "pro" }
	pints := get_knob_pints(mode)
	if file_write_str(file, ints_to_hexs(pints, 4), yes) {
		fmt.Println("> downloaded", mode, "knobs to", mode, "file", file) 
	}
}

// slot => *.dlp
func stof_cmd(slot, file string, pro, yes bool) {
	slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
	file_blank_chk(file)
	file = file_ext_chk(file, ".dlp")
	mode := "pre"
	if pro { mode = "pro" }
	if file_write_str(file, get_slot_str(slot_int, mode), yes) {
		fmt.Println("> downloaded", mode, "slot", slot_int, "to", mode, "file", file) 
	}
}

// pump from file
func pump_cmd(file string) {
	ext := filepath.Ext(file)
	file = file_read_chk(file, ext)
	switch ext {
		case ".pre", ".pro", ".spi", ".eeprom" : // these are OK
		default : log.Fatalln("> Unknown file extension:", ext)
	}
	file_str := file_read_str(file)
	addr, _ := spi_bulk_addrs(ext)
	sp := sp_open()
	spi_wr(sp, addr, file_str, true)
	sp.Close()
	fmt.Println("> pumped from", file)
	if ext == ".spi" || ext == ".eeprom" { reset_cmd() }
}

// *.dlp => knobs
func ftok_cmd(file string, pro bool) {
	file = file_read_chk(file, ".dlp")
	file_str := file_read_str(file)
	pints := hexs_to_ints(file_str, 4)
	if len(pints) < SLOT_BYTES { log.Fatalln("> Bad file info!") }
	mode := "pre"
	if pro { mode = "pro" }
	put_knob_pints(pints, mode)
	fmt.Println("> uploaded", mode, "file", file, "to", mode, "knobs") 
}

// *.dlp => slot
func ftos_cmd(slot, file string, pro bool) {
	slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
	file = file_read_chk(file, ".dlp")
	file_str := file_read_str(file)
	mode := "pre"
	if pro { mode = "pro" }
	addr := spi_slot_addr(slot_int, mode)
	sp := sp_open()
	spi_wr(sp, addr, file_str, false)
	sp.Close()
	fmt.Println("> uploaded", mode, "file", file, "to", mode, "slot", slot_int) 
}

// *.bnk => *.dlps => slots
func btos_cmd(slot, file string, pro bool) {
	slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
	file = file_read_chk(file, ".bnk")
	bnk_str := file_read_str(file)
	bnk_split := strings.Split(bnk_str, "\n")
	dir, _ := filepath.Split(file)
	dir = filepath.Clean(dir)
	mode := "pre"
	if pro { mode = "pro" }
	sp := sp_open()
	for _, line := range bnk_split {
		line_str := strings.TrimSpace(string(line));
		if !strings.HasPrefix(line_str, "//") {  // skip commented lines
			addr := spi_slot_addr(slot_int, mode)
			dlp_file := file_read_chk(filepath.Join(dir, line_str), ".dlp")
			dlp_str := file_read_str(dlp_file)
			spi_wr(sp, addr, dlp_str, false)
			fmt.Println("> uploaded", mode, "file", line_str, "to", mode, "slot", slot_int)
			slot_int++
		}
	}
	sp.Close()
}

// split file containers into sub containers
func split_cmd(file string, yes bool) {
	ext := filepath.Ext(file)
	file_read_chk(file, ext)
	dir, base := filepath.Split(file)
	dir = filepath.Clean(dir)
	base = strings.TrimSuffix(base, ext)
	file_str := file_read_str(file)
	str_split := strings.Split(file_str, "\n")
	switch ext {
		case ".eeprom" :
			pre_str := ""
			pro_str := ""
			spi_str := ""
			for line, str := range str_split {
				if line < PRE_SLOTS*SLOT_BYTES/4 { 
					pre_str += str + "\n"
				} else if line < SLOTS*SLOT_BYTES/4 { 
					pro_str += str + "\n"
				} else { 
					spi_str += str + "\n"
				}
			}
			pre_file := base + ".pre"
			pro_file := base + ".pro"
			spi_file := base + ".spi"
			pre_path := filepath.Join(dir, pre_file)
			pro_path := filepath.Join(dir, pro_file)
			spi_path := filepath.Join(dir, spi_file)
			file_write_str(pre_path, pre_str, yes)
			file_write_str(pro_path, pro_str, yes)
			file_write_str(spi_path, spi_str, yes)
			fmt.Println("> split", file, "to", pre_file, pro_file, spi_file )
		case ".pre", ".pro" :
			var dlp_str string
			file_num := 0
			for line, str := range str_split {
				dlp_str += str + "\n"
				if line % 64 == 63 { 
					dlp_name := fmt.Sprintf("%03d", file_num) + ".dlp"
					if ext == ".pro" { dlp_name = "pro_" + dlp_name }
					dlp_file := filepath.Join(dir, dlp_name)
					file_write_str(dlp_file, dlp_str, yes)
					file_num++
					dlp_str = ""
				}
			}
			fmt.Println("> split", file, "to", file_num, "numbered *.dlp files" )
		default : log.Fatalln("> Unknown file extension:", ext)
	}
}

// join sub containers to container
func join_cmd(file string, yes bool) {
	ext := filepath.Ext(file)
	dir, base := filepath.Split(file)
	dir = filepath.Clean(dir)
	base = strings.TrimSuffix(base, ext)
	switch ext {
		case ".eeprom" :
			base_path := filepath.Join(dir, base)
			pre_path := file_read_chk(base_path, ".pre")
			pre_str := file_read_str(pre_path)
			pro_path := file_read_chk(base_path, ".pro")
			pro_str := file_read_str(pro_path)
			spi_path := file_read_chk(base_path, ".spi")
			spi_str := file_read_str(spi_path)
			file_str := pre_str + "\n"
			file_str += pro_str + "\n"
			file_str += spi_str
			if file_write_str(file, file_str, yes) {
				fmt.Println("> joined", pre_path, pro_path, spi_path, "to", file )
			}
		case ".pre", ".pro" :
			file_str := ""
			files := PRE_SLOTS
			if ext == ".pro" { files = PRO_SLOTS }
			for file_num := 0; file_num < files; file_num++ {
				dlp_name := fmt.Sprintf("%03d", file_num)
				if ext == ".pro" { dlp_name = "pro_" + dlp_name }
				dlp_path := file_read_chk(filepath.Join(dir, dlp_name), ".dlp")
				dlp_str := file_read_str(dlp_path)
				file_str += dlp_str + "\n"
			}
			if file_write_str(file, file_str, yes) {
				fmt.Println("> joined", files, "numbered *.dlp files", "to", file)
			}
		default : log.Fatalln("> Unknown file extension:", ext)
	}
}

func morph_cmd(file string, knobs bool, slot string, seed int, mo, mn, me, mf, mr int) {
	rand.Seed(int64(seed))
	prn_str := ""
	var pints []int
	if mo | mn | me | mf | mr == 0 {
		log.Fatalln("> Nothing to do!")
	} else if knobs {  // morph current knobs
		pints = pints_signed(get_knob_pints("pre"), false)
		prn_str = fmt.Sprint("> morphed knobs")
	} else if file != "" {  // morph a *.dlp file
		file = file_read_chk(file, ".dlp")
		file_str := file_read_str(file)
		pints = pints_signed(hexs_to_ints(file_str, 4), false)
		prn_str = fmt.Sprint("> morphed file ", file)
	} else if slot != "" {  // morph a slot
		slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
		slot_str := get_slot_str(slot_int, "pre")
		pints = pints_signed(hexs_to_ints(slot_str, 4), false)
		prn_str = fmt.Sprint("> morphed slot ", slot_int)
	} else {
		log.Fatalln("> Nothing to do!")
	}
	prn_str += fmt.Sprint(" (-i=", seed, ")")
	pints = morph_pints(pints, mo, mn, me, mf, mr)
	put_knob_pints(pints, "pre")
	fmt.Println(prn_str)
}

// check the software
func sw_pre_chk(pre_chk bool) (bool) {
	sw_upd := false
	sw_ver := get_ver()
	fmt.Println("> Software version:", sw_ver)
	switch sw_ver {
	case SW_V8 :
		fmt.Println("> Software version is CURRENT.") 
	case SW_V7, SW_V6, SW_V5, SW_V2 :
		fmt.Println("> Software version is OLD, you may want to UPDATE it.") 
		sw_upd = true
	default :
		fmt.Println("> Software version is UNKNOWN, you may want to UPDATE it.") 
		sw_upd = true
	}
	if check_crc() { 
		fmt.Println("> Software PASSED the CRC check.") 
	} else {  
		fmt.Println("> Software FAILED the CRC check!") 
		fmt.Println("> You may need to RE-UPLOAD or UPDATE your software.") 
		sw_upd = true
	}
	if pre_chk {
		switch sw_ver {
		case SW_V8, SW_V7 : 
			fmt.Println("> Presets should be OK.") 
		case SW_V6 :
			fmt.Println("> Presets can be UPDATED with this version of the librarian.") 
		default :
			fmt.Println("> Presets cannot be UPDATED using this version of the librarian,") 
			fmt.Println("> You can REPLACE your presets, or contact Eric for further options.")
		}
	}
	return sw_upd
}

// do a bunch of update stuff via interactive menu
func update_cmd(dir_work string) {
	dir_work = filepath.Clean(dir_work)
	path_exe, err := os.Executable(); if err != nil { log.Fatal(err) }
	dir_exe := filepath.Dir(path_exe)
	dir_all := filepath.Join(dir_exe, PRESETS_DIR)
	//
	file_spi := filepath.Join(PRESETS_DIR, SW_DATE + ".spi")
	file_factory := filepath.Join(PRESETS_DIR, SW_DATE + ".eeprom")
	file_bank := filepath.Join(PRESETS_DIR, SW_DATE + ".bnk")
	file_bank_new := filepath.Join(PRESETS_DIR, SW_DATE + "_new.bnk")
	//
	path_spi := filepath.Join(dir_exe, file_spi)
	path_factory := filepath.Join(dir_exe, file_factory)
	path_bank := filepath.Join(dir_exe, file_bank)
	path_bank_new := filepath.Join(dir_exe, file_bank_new)	
	//
	path_pre_dl := filepath.Join(dir_work, "download.pre")
	path_pre_ul := filepath.Join(dir_work, "upload.pre")
	//
	prompt := false
	for {
		if prompt {
			user_input("Please press <ENTER> to return to the MENU")
		}
		prompt = false
		fmt.Println()
		fmt.Println()
		fmt.Println(" ---------------------------------")
		fmt.Println(" |  D-LEV LIBRARIAN - VERSION", LIB_VER, " |")
		fmt.Println(" | SOFTWARE & PRESET UPDATE MENU |")
		fmt.Println(" ---------------------------------")
		fmt.Println("  0. README!")
		fmt.Println("  1. Serial port setup & check.")
		fmt.Println("  2. Backup your system to an EEPROM file.")
		fmt.Println("  3. Check & update the D-Lev software.")
		fmt.Println("  4. Download all D-Lev presets to the", dir_work, "directory.")
		fmt.Println("  5. Update all presets in the", dir_work, "directory.")
		fmt.Println("  6. Upload all D-Lev presets from the", dir_work, "directory.")
		fmt.Println("  7. Upload the latest new presets.")
		fmt.Println("  8. Convert all presets in the", dir_work, "directory to MONO.")
		fmt.Println("  9. Overwrite all D-Lev preset slots with presets from the", PRESETS_DIR, "directory.")
		fmt.Println(" 10. Factory Reset: Overwrite EVERYTHING with the latest factory EEPROM file.")
		menu_sel := user_input("Please select a MENU option")
		switch {
		case menu_sel == "0" :
			prompt = true
			fmt.Println()
			fmt.Println()
			fmt.Println(" ////////////")
			fmt.Println(" // README //")
			fmt.Println(" ////////////")
			fmt.Println(" - To UPDATE the software and UPDATE ALL of the preset SLOTS: Do 1 thru 7.")
			fmt.Println(" - To UPDATE the software and OVERWRITE ALL of the preset SLOTS: Do 1 thru 3, then 9.")
			fmt.Println(" - TO UPDATE & OVERWRITE ABSOLUTELY EVERYTHING INCLUDING PROFILE SLOTS: Do 1, 2, 10.")
			fmt.Println(" - To CONVERT ALL of the preset SLOTS to MONO: Do 1, 2, 4, 8, 6.")
			fmt.Println(" - If you run into trouble, quit and pump the backup EEPROM file created in step 2.")
			fmt.Println(" - Valid prompt responses: y=yes, ENTER=no, q=quit the program.")
			fmt.Println(" - If unresponsive, do a CTRL-C (hold down the CONTROL key and press the C key).")
			fmt.Println(" - DO NOT turn or press any D-Lev knobs during the upload / download process!")
		case menu_sel == "1" :
			fmt.Println()
			ports_cmd("")
			prompt = user_prompt("Do you want to CHANGE current port?", false)
			if prompt {
				port_new := user_input("Please input PORT number")
				if port_new == "" {
					prompt = false
					fmt.Println("> -CANCEL-")
				} else {
					fmt.Println()
					ports_cmd(port_new)
				}
			}
			prompt = user_prompt("Do you want to TEST the port (do a CTRL-C if it hangs)?", false)
			if prompt {
				get_ver()
				fmt.Println("> Port seems to be OK!")
			}
		case menu_sel == "2" :
			file_backup := date_hms() + ".eeprom"
			prompt = user_prompt("Do you want to BACKUP your ENTIRE D-Lev to the FILE: " + file_backup + "?", false)
			if prompt {
				dump_cmd(file_backup, false)
			}
		case menu_sel == "3" :
			fmt.Println()
			sw_upd := sw_pre_chk(true)
			prompt = true
			if (sw_upd)  {
				prompt = user_prompt("Do you want to UPDATE your D-Lev SOFTWARE with the FILE: "+ file_spi + "?", false)
				if prompt {
					pump_cmd(path_spi)
					fmt.Println()
					sw_pre_chk(false)
				}
			}
		case menu_sel == "4" :
			prompt = user_prompt("Do you want to DOWNLOAD your D-Lev presets to " + dir_work +"?", false)
			if prompt {
				dump_cmd(path_pre_dl, true)
				split_cmd(path_pre_dl, true)
			}
		case menu_sel == "5" :
			prompt = user_prompt("Do you want to UPDATE the presets in " + dir_work + "?", false)
			if prompt {
				process_dlps(dir_work, dir_work, false, false, true, false, true)
			}
		case menu_sel == "6" :
			prompt = user_prompt("Do you want to UPLOAD the presets in "+ dir_work + "?", false)
			if prompt {
				join_cmd(path_pre_ul, true)
				pump_cmd(path_pre_ul)
			}
		case menu_sel == "7" :
			fmt.Println("> Here is a LIST of the latest NEW presets:")
			fmt.Println()
			file_read_chk(path_bank_new, ".bnk")
			file_str := file_read_str(path_bank_new)
			fmt.Println(file_str)
			prompt = user_prompt("Do you want to EXAMINE your current D-Lev presets?", false)
			if prompt { match_cmd(dir_all, "", false, false, true, true) }
			prompt = user_prompt("Do you want to UPLOAD the latest NEW presets?", false)
			if prompt {
				slot := user_input("What SLOT do you want to START the upload?")
				if slot == "" {
					prompt = false
				} else {
					btos_cmd(slot, path_bank_new, false)
				}
				prompt = user_prompt("Do you want to EXAMINE your current D-Lev presets?", false)
				if prompt { match_cmd(dir_all, "", false, false, true, true) }
			}
		case menu_sel == "8" :
			prompt = user_prompt("Do you want to CONVERT all of the presets in " + dir_work + " to MONO?", false) 
			if prompt {
				process_dlps(dir_work, dir_work, false, true, false, false, true)
			}
		case menu_sel == "9" :
			prompt = user_prompt("Do you want to OVERWRITE all D-Lev preset slots with presets in " + PRESETS_DIR + "?", false)
			if prompt {
				btos_cmd("0", path_bank, false)
				prompt = user_prompt("Do you want to EXAMINE your current D-Lev presets?", false)
				if prompt { match_cmd(dir_all, "", false, false, true, true) }
			}
		case menu_sel == "10" :
			prompt = user_prompt("Do you want to OVERWRITE ABSOLUTELY EVERYTHING in your D-Lev with the latest EEPROM file?", false)
			if prompt {
				pump_cmd(path_factory)
				fmt.Println()
				sw_pre_chk(false)
			}
		case menu_sel == "" :
			// do nothing
		default:
			prompt = true
			fmt.Println("> Invalid menu selection!")
		}
	}
}

func dev_cmd() {

	mydir, myerr := os.Getwd(); if myerr != nil { log.Fatal(myerr) }
	exe, exeerr := os.Executable(); if exeerr != nil { log.Fatal(exeerr) }
	exedir := filepath.Dir(exe)

	reldir, rerr := filepath.Rel(mydir, exedir); if rerr != nil { log.Fatal(rerr) }
 
	fmt.Println(mydir)
	fmt.Println(exedir)
	fmt.Println(reldir)
}
