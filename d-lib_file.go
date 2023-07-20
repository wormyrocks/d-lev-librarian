package main

/*
 * d-lib support functions
*/

import (
	"log"
	"os"
	"path/filepath"
	"errors"
	"strings"
)

// check if dir or file exists
func path_exists_chk(path string) (bool) {
	path = filepath.Clean(path)
    _, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// create directory for file if directory does not exist.
func file_make_dir(file string) {
	dir, _ := filepath.Split(file)
	dir = filepath.Clean(dir)
	if dir == "" { return }
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {	log.Println(err) }
}

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

// check if <file><ext> exists, check & add extension if needed
func file_read_chk(file, ext string) (string) {
	file_blank_chk(file)
	file = file_ext_chk(file, ext)
    if !path_exists_chk(file) { log.Fatal("> File ", file, " does not exist!") }
	return file
}

// check if <file> exists, prompt to overwrite
func file_write_chk(file string, yes bool) (bool) {
	file_blank_chk(file)
    if path_exists_chk(file) {
		return user_prompt("Overwrite file " + file + "?", yes)
	}
	return true
}

// write trimmed string to checked file
func file_write_str(file, data string, yes bool) (bool) {
	if file_write_chk(file, yes) {
		file_make_dir(file)
		err := os.WriteFile(file, []byte(strings.TrimSpace(data)), 0666); 
		if err != nil { log.Fatal(err) }
		return true
	}
	return false
}

// read trimmed string from file
func file_read_str(file string) (string) {
	file_bytes, err := os.ReadFile(file); if err != nil { log.Fatal(err) }
	return strings.TrimSpace(string(file_bytes))
}

// get name and contents of all *<ext> files in <dir>
func get_dir_strs(dir, ext string) ([]string, []string) {
	var name_strs []string
	var data_strs []string
	dir = filepath.Clean(dir)
	files, err := os.ReadDir(dir); if err != nil { log.Fatal(err) }
	for _, file := range files {
		if (filepath.Ext(file.Name()) == ext) && (file.IsDir() == false) {
			file_str := file_read_str(filepath.Join(dir, file.Name()))
			name_strs = append(name_strs, strings.TrimSuffix(file.Name(), ext))
			data_strs = append(data_strs, file_str)
		}
    }
    return name_strs, data_strs
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
			file_str := file_read_str(filepath.Join(dir, file.Name()))
			f_map[file_str] = strings.TrimSuffix(file.Name(), ext)
		}
    }
    return f_map
}


////////////
// CONFIG //
////////////

// set key value in config file
// create file if it doesn't exist
func cfg_set(key, value string) {
	if !path_exists_chk(CFG_FILE) {	 // create file
		file_write_str(CFG_FILE, "", true) 
	}
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	file_str := file_read_str(CFG_FILE)
	lines := strings.Split(file_str, "\n")
	str := ""
	found := false
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == key {
			str += key + " " + value + "\n"
			found = true
		} else { str += line + "\n"	}
	}
	if !found { str += key + " " + value + "\n" }
	file_write_str(CFG_FILE, str, true)
}

// get first matching key value from config file
// return "" if file doesn't exist or no key match
func cfg_get(key string) (string) {
	if !path_exists_chk(CFG_FILE) { return "" }
	key = strings.TrimSpace(key)
	file_str := file_read_str(CFG_FILE)
	lines := strings.Split(file_str, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == key {
			return fields[1]
		}
	}
	return ""
}

