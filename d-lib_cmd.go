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
	fmt.Print("> ", prompt, " <y|ENTER|q>: ")
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
	fmt.Print("> ", prompt, " (q=quit): ")
	input := user_word()
	if input == "q" { log.Fatalln("> Quit, exiting program...") }
	return input
}

// print librarian version & help info
func help_cmd(verbose_f bool) {
	fmt.Print("= D-Lev Librarian version ", VERSION, " =\n") 
	fmt.Print(help_str) 
	if verbose_f { fmt.Print(help_verbose_str) }  // print verbose help
}

// do processor reset
func reset_cmd(port int) {
	sp := sp_open(port)
	sp_wr_rd(sp, "0 0xff000000 wr ", false)
	fmt.Println("> issued processor reset")
	sp.Close()
}

// return D-Lev software version
func get_ver(port int) (string) {
	sp := sp_open(port)
	rd_str := sp_wr_rd(sp, "ver ", false)
	sp.Close()
	return decruft_hcl(rd_str)
}

// check D-Lev software crc
func check_crc(port int) (bool) {
	sp := sp_open(port)
	rd_str := sp_wr_rd(sp, "crc ", false)
	sp.Close()
	if decruft_hcl(rd_str) == CRC { return true }
	return false
}

// get D-Lev software version
func ver_cmd(port int) {
	fmt.Println("> current software version:", get_ver(port))
}

// do ACAL
func acal_cmd(port int) {
	sp := sp_open(port)
	sp_wr_rd(sp, "acal ", false)
	fmt.Println("> issued ACAL")
	sp.Close()
}

// do HCL command
func hcl_cmd(port int) {
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
func loop_cmd(port int) {
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
func ports_cmd(port int, port_str string) (int) {
	ports := sp_list()
	if len(ports) == 0 {
		fmt.Println("> No serial ports found!")
	} else {
		fmt.Println("> Available serial ports:")
		for p_num, p_str := range ports { fmt.Printf("  %v : %v\n", p_num, p_str) }
	}
	if port_str != "" { 
		port_new, err := strconv.Atoi(port_str)
		if err != nil { 
			fmt.Println("> Current port:", port)
			log.Fatalln("> Bad port number!") 
		} else if port_new >= len(ports) { 
			fmt.Println("> Current port:", port)
			log.Fatalln("> Port number out of range!") 
		} else {
			cfg_set("port", port_str)
			port = port_new
			fmt.Println("> Set port to:", port)
		}
	} else if len(os.Args) > 2 { 
		fmt.Println("> Current port:", port)
		log.Fatalln("> Use the -p flag to set the port!")
	} else {
		fmt.Println("> Current port:", port)
	}
	return port
}

// view knobs, DLP file, slot
func view_cmd(port int, file string, pro, knobs bool, slot string) {
	mode := "pre"
	if pro { mode = "pro" }
	if knobs {  // view current knobs
		knob_str := get_knob_str(port)
		fmt.Println(ui_prn_str(knob_ui_strs(knob_str)))
		fmt.Println("> knobs")
	} else if file != "" {  // view a *.dlp file
		var file_str string
		file, file_str = get_file_str(file, ".dlp")
		fmt.Println(ui_prn_str(pre_ui_strs(file_str, pro)))
		fmt.Println(">", mode, "file", file)
	} else if slot != "" {  // view a slot
		slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
		slot_str := get_slot_str(port, slot_int, mode)
		fmt.Println(ui_prn_str(pre_ui_strs(slot_str, pro)))
		fmt.Println(">", mode, "slot", slot_int)
	} else {
		log.Fatalln("> Nothing to do!")
	}
}

// twiddle knob
func knob_cmd(port int, knob, offset, val string) {
	str_split := (strings.Split(strings.TrimSpace(knob), ":"))
	if len(str_split) < 2  { log.Fatalln("> Bad knob value!") }
	knob_int, err := strconv.Atoi(str_split[1]); if err != nil { log.Fatalln("> Bad knob value!") }
	if knob_int < 0 || knob_int > UI_PG_KNOBS - 2 { log.Fatalln("> Bad knob value!") }
	pg_name, pg_idx := page_lookup(str_split[0])
	if pg_idx < 0 { log.Fatalln("> Bad page name!") }
	knob_idx := knob_int + pg_idx * UI_PG_KNOBS
	ptype, plabel, _, _ := pname_lookup(knob_pnames[knob_idx])
	sp := sp_open(port)
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
		sp := sp_open(port)
		sp_wr_rd(sp, strconv.Itoa(knob_idx) + " " + strconv.Itoa(rw_int) + " wk ", false)
		sp.Close()
		fmt.Print("=>[", strings.TrimSpace(enc_disp(rw_int, ptype)), "]")
	}
	fmt.Println("")
}

// diff DLP file(s) / slot(s) / knobs
func diff_cmd(port int, file, file2 string, pro, knobs bool, slot, slot2 string) {
	mode := "pre"
	if pro { mode = "pro" }
	if file != "" {  // compare to a *.dlp file
		file_str := ""
		file, file_str = get_file_str(file, ".dlp")
		if knobs {  // file vs. knobs
			knob_str := knob_pre_str(get_knob_str(port), pro)
			fmt.Println(diff_prn_str(diff_pres(file_str, knob_str, pro)))
			fmt.Println(">", mode, "file", file, "vs. knobs" )
		} else if file2 != "" {  // file vs. file2
			file2_str := ""
			file2, file2_str = get_file_str(file2, ".dlp")
			fmt.Println(diff_prn_str(diff_pres(file_str, file2_str, pro)))
			fmt.Println(">", mode, "file", file, "vs.", file2 )
		} else if slot != "" {  // file vs. slot
			slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
			slot_str := get_slot_str(port, slot_int, mode)
			fmt.Println(diff_prn_str(diff_pres(file_str, slot_str, pro)))
			fmt.Println(">", mode, "file", file, "vs. slot", slot )
		} else {
			log.Fatalln("> Nothing to do!")
		}
	} else if slot != "" {  // compare to a slot
		slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
		slot_str := get_slot_str(port, slot_int, mode)
		if knobs {  // slot vs. knobs
			knob_str := knob_pre_str(get_knob_str(port), pro)
			fmt.Println(diff_prn_str(diff_pres(slot_str, knob_str, pro)))
			fmt.Println(">", mode, "slot", slot, "vs. knobs" )
		} else if slot2 != "" {  // slot vs. slot2
			slot2_int, err := strconv.Atoi(slot2); if err != nil { log.Fatalln("> Bad slot number!") }
			slot2_str := get_slot_str(port, slot2_int, mode)
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
func match_cmd(port int, dir, dir2 string, pro, hdr, guess, slots bool) {
	mode := "pre"
	if pro { mode = "pro" }
	name_strs, data_strs := get_dir_strs(dir, ".dlp")
	if len(data_strs) == 0 {  log.Fatalln("> No", mode, "files in", dir) }
	if slots {
		slots_strs := get_slots_strs(port)
		fmt.Print(slots_prn_str(comp_file_data(slots_strs, name_strs, data_strs, pro, guess), pro, hdr))
		fmt.Println("> matched", mode, "slots to", mode, "files in", dir)
	} else {
		name2_strs, data2_strs := get_dir_strs(dir2, ".dlp")
		if len(data2_strs) == 0 {  log.Fatalln("> No", mode, "files in", dir2) }
		fmt.Print(files_prn_str(name2_strs, comp_file_data(data2_strs, name_strs, data_strs, pro, guess)))
		fmt.Println("> matched", mode, "files in", dir2, "to", mode, "files in", dir)
	}
}

func dump_cmd(port int, file string, yes bool) {
	file_blank_chk(file)
	ext := strings.Trim(filepath.Ext(file), ".")
	switch ext {
		case "pre", "pro", "spi", "eeprom" : // these are OK
		default : log.Fatalln("> Unknown file extension", ext)
	}
	addr, end := spi_bulk_addrs(ext)
	sp := sp_open(port)
	rx_str := spi_rd(sp, addr, end - 1, true)
	sp.Close()
	if file_write(file, []byte(rx_str), yes) {
		fmt.Println("> dumped to", file) 
	}
}

// knobs => *.dlp
func ktof_cmd(port int, file string, pro, yes bool) {
	mode := "pre"
	if pro { mode = "pro" }
	file_blank_chk(file)
	file = file_ext_chk(file, ".dlp")
	pints := get_knob_pints(port, mode)
	if file_write(file, []byte(ints_to_hexs(pints, 4)), yes) {
		fmt.Println("> downloaded", mode, "knobs to", mode, "file", file) 
	}
}

// slot => *.dlp
func stof_cmd(port int, slot, file string, pro, yes bool) {
	slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
	mode := "pre"
	if pro { mode = "pro" }
	file_blank_chk(file)
	file = file_ext_chk(file, ".dlp")
	if file_write(file, []byte(get_slot_str(port, slot_int, mode)), yes) {
		fmt.Println("> downloaded", mode, "slot", slot_int, "to", mode, "file", file) 
	}
}

// pump from file
func pump_cmd(port int, file string) {
	file_blank_chk(file)
	ext := strings.Trim(filepath.Ext(file), ".")
	switch ext {
		case "pre", "pro", "spi", "eeprom" : // these are OK
		default : log.Fatalln("> Unknown file extension", ext)
	}
	addr, _ := spi_bulk_addrs(ext)
	file_bytes, err := os.ReadFile(file); if err != nil { log.Fatal(err) }
	sp := sp_open(port)
	spi_wr(sp, addr, string(file_bytes), ext, true)
	sp.Close()
	fmt.Println("> pumped from", file)
	if ext == "spi" || ext == "eeprom" { reset_cmd(port) }
}

// *.dlp => knobs
func ftok_cmd(port int, file string, pro bool) {
	mode := "pre"
	if pro { mode = "pro" }
	var file_str string
	file, file_str = get_file_str(file, ".dlp")
	pints := hexs_to_ints(file_str, 4)
	if len(pints) < SLOT_BYTES { log.Fatalln("> Bad file info!") }
	put_knob_pints(port, pints, mode)
	fmt.Println("> uploaded", mode, "file", file, "to", mode, "knobs") 
}

// *.dlp => slot
func ftos_cmd(port int, slot, file string, pro bool) {
	slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
	mode := "pre"
	if pro { mode = "pro" }
	var file_str string
	file, file_str = get_file_str(file, ".dlp")
	addr := spi_slot_addr(slot_int, mode)
	sp := sp_open(port)
	spi_wr(sp, addr, file_str, mode, false)
	sp.Close()
	fmt.Println("> uploaded", mode, "file", file, "to", mode, "slot", slot_int) 
}

// *.bnk => *.dlps => slots
func btos_cmd(port int, slot, file string, pro bool) {
	slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
	mode := "pre"
	if pro { mode = "pro" }
	file_blank_chk(file)
	file = file_ext_chk(file, ".bnk")
	dir, bnk_file := filepath.Split(file)
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

// split file containers into sub containers
func split_cmd(file string, yes bool) {
	file_blank_chk(file)
	dir, file_name := filepath.Split(file)
	dir = filepath.Clean(dir)
	ext := filepath.Ext(file_name)
	base_name := strings.TrimSuffix(file_name, ext)
	switch ext {
		case ".pre", ".pro", ".eeprom" : // these are OK
		default : log.Fatalln("> Unknown file extension", ext)
	}
	file_bytes, err := os.ReadFile(file); if err != nil { log.Fatal(err) }
	str_split := (strings.Split(strings.TrimSpace(string(file_bytes)), "\n"))
	if ext == ".eeprom" {
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
		pre_name := base_name + ".pre"
		pro_name := base_name + ".pro"
		spi_name := base_name + ".spi"
		//
		pre_file := filepath.Join(dir, pre_name)
		pro_file := filepath.Join(dir, pro_name)
		spi_file := filepath.Join(dir, spi_name)
		file_write(pre_file, []byte(pre_str), yes)
		file_write(pro_file, []byte(pro_str), yes)
		file_write(spi_file, []byte(spi_str), yes)
		fmt.Println("> split", file, "to", pre_name, pro_name, spi_name )
	} else {  // pre | pro
		var dlp_str string
		file_num := 0
		for line, str := range str_split {
			dlp_str += str + "\n"
			if line % 64 == 63 { 
				dlp_name := fmt.Sprintf("%03d", file_num) + ".dlp"
				if ext == ".pro" { dlp_name = "pro_" + dlp_name }
				dlp_file := filepath.Join(dir, dlp_name)
				file_write(dlp_file, []byte(dlp_str), yes)
				file_num++
				dlp_str = ""
			}
		}
		fmt.Println("> split", file, "to", file_num, "numbered *.dlp files" )
	}
}

// join sub containers to container
func join_cmd(file string, yes bool) {
	file_blank_chk(file)
	dir, file_name := filepath.Split(file)
	dir = filepath.Clean(dir)
	ext := filepath.Ext(file_name)
	switch ext {
		case ".pre", ".pro", ".eeprom" : // these are OK
		default : log.Fatalln("> Unknown file extension", ext)
	}
	base_name := strings.TrimSuffix(file_name, ext)
	if ext == ".eeprom" {
		pre_name := base_name + ".pre"
		pro_name := base_name + ".pro"
		spi_name := base_name + ".spi"
		pre_file := filepath.Join(dir, pre_name)
		pro_file := filepath.Join(dir, pro_name)
		spi_file := filepath.Join(dir, spi_name)
		pre_bytes, err := os.ReadFile(pre_file); if err != nil { log.Fatal(err) }
		pro_bytes, err := os.ReadFile(pro_file); if err != nil { log.Fatal(err) }
		spi_bytes, err := os.ReadFile(spi_file); if err != nil { log.Fatal(err) }
		wr_bytes := append(pre_bytes, pro_bytes...)
		wr_bytes = append(wr_bytes, spi_bytes...)
		if file_write(file, wr_bytes, yes) {
			fmt.Println("> merged", pre_name, pro_name, spi_name, "to", file )
		}
	} else {  // pre | pro
		var wr_bytes []byte
		files := PRE_SLOTS
		if ext == ".pro" { files = PRO_SLOTS }
		for file_num := 0; file_num < files; file_num++ {
			dlp_name := fmt.Sprintf("%03d", file_num) + ".dlp"
			if ext == ".pro" { dlp_name = "pro_" + dlp_name }
			rd_file := filepath.Join(dir, dlp_name)
			rd_bytes, err := os.ReadFile(rd_file); if err != nil { log.Fatal(err) }
			wr_bytes = append(wr_bytes, rd_bytes...)
		}
		if file_write(file, wr_bytes, yes) {
			fmt.Println("> joined", files, "numbered *.dlp files", "to", file)
		}
	}
}

func morph_cmd(port int, file string, knobs bool, slot string, seed int, mo, mn, me, mf, mr int) {
	rand.Seed(int64(seed))
	prn_str := ""
	var pints []int
	if mo | mn | me | mf | mr == 0 {
		log.Fatalln("> Nothing to do!")
	} else if knobs {  // morph current knobs
		pints = pints_signed(get_knob_pints(port, "pre"), false)
		prn_str = fmt.Sprint("> morphed knobs")
	} else if file != "" {  // morph a *.dlp file
		var file_str string
		file, file_str = get_file_str(file, ".dlp")
		pints = pints_signed(hexs_to_ints(file_str, 4), false)
		prn_str = fmt.Sprint("> morphed file ", file)
	} else if slot != "" {  // morph a slot
		slot_int, err := strconv.Atoi(slot); if err != nil { log.Fatalln("> Bad slot number!") }
		slot_str := get_slot_str(port, slot_int, "pre")
		pints = pints_signed(hexs_to_ints(slot_str, 4), false)
		prn_str = fmt.Sprint("> morphed slot ", slot_int)
	} else {
		log.Fatalln("> Nothing to do!")
	}
	prn_str += fmt.Sprint(" (-i=", seed, ")")
	pints = morph_pints(pints, mo, mn, me, mf, mr)
	put_knob_pints(port, pints, "pre")
	fmt.Println(prn_str)
}


// do a bunch of update stuff via interactive menu
func update_cmd(port int, dir_all, dir_work string) {
	ver_old := "27c263bf"
	ver_new := "73c6c3d7"
	date_new := "2023-05-24"
	file_spi := filepath.Join(dir_all, ver_new + ".spi")
	file_factory := filepath.Join(dir_all, date_new + ".eeprom")
	file_bank := filepath.Join(dir_all, date_new + ".bnk")
	file_bank_new := filepath.Join(dir_all, date_new + "_new.bnk")
	file_pre_dl := filepath.Join(dir_work, "download.pre")
	file_pre_ul := filepath.Join(dir_work, "upload.pre")
	prompt := false
	for {
		if prompt {
			fmt.Println()
			user_input("<ENTER> to return to main MENU")
		}
		prompt = false
		fmt.Println("\n")
		fmt.Println("** D-LEV SOFTWARE & PRESETS UPDATE MENU **")
		fmt.Println(" 0. README!")
		fmt.Println(" 1. Serial port setup & check.")
		fmt.Println(" 2. Backup your system to an EEPROM file.")
		fmt.Println(" 3. Check & update the D-Lev software.")
		fmt.Println(" 4. Download all D-Lev preset slots to the", dir_work, "directory.")
		fmt.Println(" 5. Update all preset files in the", dir_work, "directory.")
		fmt.Println(" 6. Overwrite all D-Lev preset slots with files from the", dir_work, "directory.")
		fmt.Println(" 7. Upload the latest new preset files.")
		fmt.Println(" 8. Convert all preset files in the", dir_work, "directory to MONO.")
		fmt.Println(" 9. Overwrite all D-Lev preset slots with presets from the", dir_all, "directory.")
		fmt.Println("10. Factory Reset: Overwrite EVERYTHING with the latest factory EEPROM file.")
		fmt.Println()
		menu_sel := user_input("Menu select")
		fmt.Print("\n")
		switch {
		case menu_sel == "0" :
			prompt = true
			fmt.Println()
			fmt.Println("** README README README **")
			fmt.Println(" - To UPDATE the software and UPDATE ALL of the preset SLOTS: Do 1 thru 7.")
			fmt.Println(" - To UPDATE the software and OVERWRITE ALL of the preset SLOTS: Do 1 thru 3, then 9.")
			fmt.Println(" - TO UPDATE & OVERWRITE ABSOLUTELY EVERYTHING INCLUDING PROFILE SLOTS: Do 1, 2, 10.")
			fmt.Println(" - To CONVERT ALL of the preset SLOTS to MONO: Do 1, 2, 4, 8, 6.")
			fmt.Println(" - If you run into trouble, quit and pump the backup EEPROM file created in step 2.")
			fmt.Println(" - Valid prompt responses: y=yes, ENTER=no, q=quit the program.")
			fmt.Println(" - If unresponsive, do a CTRL-C (hold down the CONTROL key and press the C key).")
		case menu_sel == "1" :
			ports_cmd(port, "")
			fmt.Println()
			prompt = user_prompt("Do you want to CHANGE current port " + strconv.Itoa(port) + "?", false)
			if prompt {
				port_new := user_input("Input PORT number")
				if port_new == "" {
					prompt = false
					fmt.Println("> -CANCEL-")
				} else {
					port = ports_cmd(port, port_new)
				}
			}
			fmt.Println()
			prompt = user_prompt("Do you want to TEST port " + strconv.Itoa(port) + "" + " (do a CTRL-C if it hangs)?", false)
			if prompt {
				get_ver(port)
				fmt.Println("> Port", port, "seems to be OK!")
			}
		case menu_sel == "2" :
			backup_file := date_hms() + ".eeprom"
			prompt = user_prompt("Do you want to BACKUP your ENTIRE D-Lev to the FILE: " + backup_file + "?", false)
			if prompt {
				dump_cmd(port, backup_file, false)
			}
		case menu_sel == "3" :
			ver_current := get_ver(port)
			crc_ok := check_crc(port)
			fmt.Println("> Your D-Lev current software VERSION is:", ver_current)
			if (ver_current == ver_new) && crc_ok { 
				fmt.Println("> Your D-Lev software is UP-TO-DATE and PASSED the CRC check!") 
				prompt = true
			} else {
				if !crc_ok { 
					fmt.Println("> Your SOFTWARE FAILED THE CRC CHECK!  You need to RE-UPLOAD or UPDATE your software.") 
				}
				if (ver_current != ver_new) {
					fmt.Println("> Your SOFTWARE is OLD, you may want to UPDATE it.")
					if (ver_current != ver_old)  {
						fmt.Print("> Your PRESETS are TOO OLD TO UPDATE using this version of the librarian!\n") 
						fmt.Print("  The preset update process only works correctly for ", ver_old, " => ", ver_new, ".\n")
						fmt.Print("  - You can replace your presets with the latest ones using Menu #9.\n")
						fmt.Print("  - Or you can contact Eric for further options.\n")
					}
				}
				fmt.Println()
				prompt = user_prompt("Do you want to UPDATE your D-Lev SOFTWARE with the FILE: "+ file_spi + "?", false)
				if prompt {
					pump_cmd(port, file_spi)
				}
			}
		case menu_sel == "4" :
			prompt = user_prompt("Do you want to DOWNLOAD ALL D-Lev preset SLOTS to FILES in " + dir_work +"?", false)
			if prompt {
				dump_cmd(port, file_pre_dl, true)
				split_cmd(file_pre_dl, true)
			}
		case menu_sel == "5" :
			prompt = user_prompt("Do you want to UPDATE ALL preset FILES in " + dir_work + "?", false)
			if prompt {
				process_dlps(dir_work, dir_work, false, false, true, false, true)
			}
		case menu_sel == "6" :
			prompt = user_prompt("Do you want to OVERWRITE ALL of your D-Lev preset SLOTS with the FILES in "+ dir_work + "?", false)
			if prompt {
				join_cmd(file_pre_ul, true)
				pump_cmd(port, file_pre_ul)
			}
		case menu_sel == "7" :
			fmt.Println("> Here is a LIST of the latest NEW PRESETS:\n")
			_, file_str := get_file_str(file_bank_new, ".bnk")
			fmt.Println(file_str)
			prompt = user_prompt("Do you want to EXAMINE your current D-Lev preset SLOTS?", false)
			if prompt { 
				match_cmd(port, dir_all, "", false, false, true, true)
			}
			fmt.Println("")
			prompt = user_prompt("Do you want to UPLOAD the latest NEW PRESETS?", false)
			if prompt {
				fmt.Println("")
				slot := user_input("What SLOT do you want to START the upload?")
				if slot == "" {
					prompt = false
				} else {
					btos_cmd(port, slot, file_bank_new, false)
				}
			}
		case menu_sel == "8" :
			prompt = user_prompt("Do you want to CONVERT ALL of the preset FILES in " + dir_work + " to MONO?", false) 
			if prompt {
				process_dlps(dir_work, dir_work, false, true, false, false, true)
			}
		case menu_sel == "9" :
			prompt = user_prompt("Do you want to OVERWRITE ALL of your D-Lev preset SLOTS with the FILES in " + dir_all + "?", false)
			if prompt {
				btos_cmd(port, "0", file_bank, false)
			}
		case menu_sel == "10" :
			prompt = user_prompt("Do you want to OVERWRITE ABSOLUTELY EVERYTHING in your D-Lev with the FILE "+ file_factory + "?", false)
			if prompt {
				pump_cmd(port, file_factory)
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

	fmt.Println(date())
	fmt.Println(hms())

}
