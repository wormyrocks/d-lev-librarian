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

// check for blank file name
func file_blank_chk(file string) {
	if strings.TrimSpace(file) == "" { 
		log.Fatal("> Missing file name!") 
	}
}

// check file ext, add if missing
func file_ext_chk(file string, ext string) (string) {
    f_ext := filepath.Ext(file)
    if len(f_ext) == 0 { 
		file += ext
	} else if f_ext != ext { 
		log.Fatal("> Wrong file extension: ", f_ext, " (expecting: ", ext, " or none)") 
	}
	return file
}

// check if file exists, prompt to overwrite
func file_exists_chk(file string) {
    _, err := os.Stat(file)
    if !errors.Is(err, os.ErrNotExist) {
		fmt.Print("> Overwrite file ", file, " ?  <y|n> ")
		var input string
		fmt.Scanln(&input)
		if input != "y" { log.Fatalln("> Abort, exiting program...") }
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

// get dlp file contents, check and update filename
func get_dlp(file *string) ([]byte) {
	file_blank_chk(*file)
	*file = file_ext_chk(*file, ".dlp")
	file_bytes, err := os.ReadFile(*file); if err != nil { log.Fatal(err) }
	return file_bytes
}
