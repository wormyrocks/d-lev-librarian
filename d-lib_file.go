package main

/*
 * d-lev support functions
*/

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"errors"
	"strings"
	"strconv"
)

// check for blank <file> name
func file_blank_chk(file string) {
	if strings.TrimSpace(file) == "" { 
		log.Fatal("> Missing file name!") 
	}
}

// check <file> <ext>, add if missing
func file_ext_chk(file string, ext string) (string) {
    f_ext := filepath.Ext(file)
    if len(f_ext) == 0 { 
		file += ext
	} else if f_ext != ext { 
		log.Fatal("> Wrong file extension: ", f_ext, " (expecting: ", ext, " or none)") 
	}
	return file
}

// check if <file> exists, prompt to overwrite
func file_exists_chk(file string) {
    _, err := os.Stat(file)
    if !errors.Is(err, os.ErrNotExist) {
		fmt.Print("> Overwrite file ", file, " ?  <y|n> ")
		var input string
		fmt.Scanln(&input)
		if input != "y" { log.Fatalln("> Abort, exiting program...") }
	}
}

// get contents of <file>.<ext>, check & add extension if needed
func get_file_str(file, ext string) (string, string) {
	file_blank_chk(file)
	file = file_ext_chk(file, ext)
	file_bytes, err := os.ReadFile(file); if err != nil { log.Fatal(err) }
	return file, string(file_bytes)
}

// get name and contents of all *.<ext> files in <dir>
func get_dir_strs(dir, ext string) ([]string, []string) {
	var name_strs []string
	var data_strs []string
	dir = filepath.Clean(dir)
	files, err := os.ReadDir(dir); if err != nil { log.Fatal(err) }
	for _, file := range files {
		if (filepath.Ext(file.Name()) == ext) && (file.IsDir() == false) {
			file_bytes, err := os.ReadFile(filepath.Join(dir, file.Name())); if err != nil { log.Fatal(err) }
			name_strs = append(name_strs, strings.TrimSuffix(file.Name(), ext))
			data_strs = append(data_strs, string(file_bytes))
		}
    }
    return name_strs, data_strs
}

// split file containers into sub containers
func split_file(file string) {
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
		file_exists_chk(pre_file)
		err = os.WriteFile(pre_file, []byte(pre_str), 0666); if err != nil { log.Fatal(err) }
		file_exists_chk(pro_file)
		err = os.WriteFile(pro_file, []byte(pro_str), 0666); if err != nil { log.Fatal(err) }
		file_exists_chk(spi_file)
		err = os.WriteFile(spi_file, []byte(spi_str), 0666); if err != nil { log.Fatal(err) }
		fmt.Println("> split", file, "to", pre_name, pro_name, spi_name )
	} else {  // pre | pro
		var dlp_str string
		file_num := 0
		for line, str := range str_split {
			dlp_str += str + "\n"
			if line % 64 == 63 { 
				dlp_name := fmt.Sprintf("%03d", file_num) + ".dlp"
				if ext == ".pro" { dlp_name = "pro_" + dlp_name }
				pre_file := filepath.Join(dir, dlp_name)
				err = os.WriteFile(pre_file, []byte(dlp_str), 0666); if err != nil { log.Fatal(err) }
				file_num++
				dlp_str = ""
			}
		}
		fmt.Println("> split", file, "to", file_num, "numbered *.dlp files" )
	}
}

// join sub containers to container
func join_files(file string) {
	file_blank_chk(file)
	dir, file_name := filepath.Split(file)
	dir = filepath.Clean(dir)
	ext := filepath.Ext(file_name)
	switch ext {
		case ".pre", ".pro", ".eeprom" : // these are OK
		default : log.Fatalln("> Unknown file extension", ext)
	}
	file_exists_chk(file)
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
		err = os.WriteFile(file, wr_bytes, 0666); if err != nil { log.Fatal(err) }
		fmt.Println("> merged", pre_name, pro_name, spi_name, "to", file )
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
		err := os.WriteFile(file, wr_bytes, 0666); if err != nil { log.Fatal(err) }
		fmt.Println("> joined", files, "numbered *.dlp files", "to", file)
	}
}

// return a file map for a given directory
// key = file.ext contents (as string)
// value = file name (sans extension)
func map_files(dir string, ext string) (map[string]string) {
	var f_map = make(map[string]string)
	dir = filepath.Clean(dir)
	files, err := os.ReadDir(dir); if err != nil { log.Fatal(err) }
	for _, file := range files {
		if filepath.Ext(file.Name()) == ext && file.IsDir() == false {
			file_bytes, err := os.ReadFile(filepath.Join(dir, file.Name())); if err != nil { log.Fatal(err) }
			f_map[strings.TrimSpace(string(file_bytes))] = strings.TrimSuffix(file.Name(), ext)
		}
    }
    return f_map
}


////////////
// CONFIG //
////////////

// set key value in config file
// create file if it doesn't exist
func cfg_set(key string, value string) {
    _, err := os.Stat(CFG_FILE)
    if errors.Is(err, os.ErrNotExist) {  // missing
		err = os.WriteFile(CFG_FILE, []byte(""), 0666); if err != nil { log.Fatal(err) }
	}
	bytes, err := os.ReadFile(CFG_FILE); if err != nil { log.Fatal(err) }
	lines := (strings.Split(strings.TrimSpace(string(bytes)), "\n"))
	str := ""
	found := false
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			if fields[0] == key {
				str += key + " " + value + "\n"
				found = true
				break
			} else { str += line + "\n"	}
		}
	}
	if !found { str += key + " " + value + "\n" }
	err = os.WriteFile(CFG_FILE, []byte(str), 0666); if err != nil { log.Fatal(err) }
}

// get key value from config file
// create file if it doesn't exist
func cfg_get(key string) (string) {
    _, err := os.Stat(CFG_FILE)
    if errors.Is(err, os.ErrNotExist) {  // missing
		cfg_set(key, strconv.Itoa(CFG_PORT))
	}
	bytes, err := os.ReadFile(CFG_FILE); if err != nil { log.Fatal(err) }
	lines := (strings.Split(strings.TrimSpace(string(bytes)), "\n"))
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == key {
			return fields[1]
		}
	}
	return ""
}

