package main

/*
 * Librarian for the D-Lev Theremin
 * See file "CHANGE_LOG.txt" for details
*/

import (
	"flag"
	"log"
	"os"
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
		help_cmd(false)  // short help
    } else {
		// parse subcommands
		switch os.Args[1] {
			case "update": update(port)
			case "help": help_cmd(true)
			case "ports": port = ports(port)
			case "view": view(port)
			case "match": match(port)
			case "diff": diff(port)
			case "ktof": ktof(port)
			case "ftok": ftok(port)
			case "stof": stof(port)
			case "ftos": ftos(port)
			case "btos": btos(port)
			case "dump": dump(port)
			case "pump": pump(port)
			case "split": split()
			case "join": join()
			case "morph": morph(port)
			case "batch": batch()
			case "knob": knob(port)
			case "hcl": hcl_cmd(port)
			case "loop":loop_cmd(port)
			case "ver": ver_cmd(port)
			case "acal": acal_cmd(port)
			case "reset": reset_cmd(port)
			case "dev": dev()  // dev stuff
			default: log.Fatalln("> Unknown command:", os.Args[1])
		}
	}
}  // end of main()


////////////////////
// main functions //
////////////////////

// list free serial ports / set port
func ports(port int) (int) {
	sub := flag.NewFlagSet("ports", flag.ExitOnError)
	port_str := sub.String("p", "", "`port` number")
	sub.Parse(os.Args[2:])
	//
	return ports_cmd(port, *port_str)
}

// view knobs, DLP file, slot
func view(port int) {
	sub := flag.NewFlagSet("view", flag.ExitOnError)
	file := sub.String("f", "", "view `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	knobs := sub.Bool("k", false, "view knobs")
	slot := sub.String("s", "", "view `slot` number")
	sub.Parse(os.Args[2:])
	//
	view_cmd(port, *file, *pro, *knobs, *slot)
}

// twiddle knob
func knob(port int) {
	sub := flag.NewFlagSet("knobs", flag.ExitOnError)
	knob := sub.String("k", "", "page:knob[0:6]")
	offset := sub.String("o", "", "knob offset value")
	val := sub.String("s", "", "knob set value")
	sub.Parse(os.Args[2:])
	//
	knob_cmd(port, *knob, *offset, *val)
}

// diff DLP file(s) / slot(s) / knobs
func diff(port int) {
	sub := flag.NewFlagSet("diff", flag.ExitOnError)
	file := sub.String("f", "", "compare `file` name")
	file2 := sub.String("f2", "", "compare `file2` name")
	pro := sub.Bool("pro", false, "profile mode")
	knobs := sub.Bool("k", false, "compare knobs")
	slot := sub.String("s", "", "compare `slot` number")
	slot2 := sub.String("s2", "", "compare `slot2` number")
	sub.Parse(os.Args[2:])
	//
	diff_cmd(port, *file, *file2, *pro, *knobs, *slot, *slot2)
}

// match slots / DLP files w/ DLP files & list
func match(port int) {
	sub := flag.NewFlagSet("match", flag.ExitOnError)
	dir := sub.String("d", ".", "`directory` name")
	dir2 := sub.String("d2", ".", "`directory` name")
	pro := sub.Bool("pro", false, "profile mode")
	hdr := sub.Bool("h", false, "header format")
	guess := sub.Bool("g", false, "guess")
	slots := sub.Bool("s", false, "slots")
	sub.Parse(os.Args[2:])
	//
	match_cmd(port, *dir, *dir2, *pro, *hdr, *guess, *slots)
}


//////////////////////
// download to file //
//////////////////////

// dump to file
func dump(port int) {
	sub := flag.NewFlagSet("dump", flag.ExitOnError)
	file := sub.String("f", "", "target `file` name")
	yes := sub.Bool("y", false, "overwrite files")
	sub.Parse(os.Args[2:])
	//
	dump_cmd(port, *file, *yes)
}

// knobs => *.dlp
func ktof(port int) {
	sub := flag.NewFlagSet("ktof", flag.ExitOnError)
	file := sub.String("f", "", "target `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	yes := sub.Bool("y", false, "overwrite files")
	sub.Parse(os.Args[2:])
	//
	ktof_cmd(port, *file, *pro, *yes)
}

// slot => *.dlp
func stof(port int) {
	sub := flag.NewFlagSet("stof", flag.ExitOnError)
	slot := sub.String("s", "", "source `slot` number")
	file := sub.String("f", "", "target `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	yes := sub.Bool("y", false, "overwrite files")
	sub.Parse(os.Args[2:])
	//
	stof_cmd(port, *slot, *file, *pro, *yes)
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
	pump_cmd(port, *file)
}

// *.dlp => knobs
func ftok(port int) {
	sub := flag.NewFlagSet("ftok", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	ftok_cmd(port, *file, *pro)
}

// *.dlp => slot
func ftos(port int) {
	sub := flag.NewFlagSet("ftos", flag.ExitOnError)
	slot := sub.String("s", "", "target `slot` number")
	file := sub.String("f", "", "source `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	ftos_cmd(port, *slot, *file, *pro)
}

// *.bnk => *.dlps => slots
func btos(port int) {
	sub := flag.NewFlagSet("btos", flag.ExitOnError)
	slot := sub.String("s", "", "starting `slot` number")
	file := sub.String("f", "", "bank `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	btos_cmd(port, *slot, *file, *pro)
}

// split bulk files
func split() {
	sub := flag.NewFlagSet("split", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	yes := sub.Bool("y", false, "overwrite files")
	sub.Parse(os.Args[2:])
	//
	split_cmd(*file, *yes)
}

// join bulk files
func join() {
	sub := flag.NewFlagSet("join", flag.ExitOnError)
	file := sub.String("f", "", "target `file` name")
	yes := sub.Bool("y", false, "overwrite files")
	sub.Parse(os.Args[2:])
	//
	join_cmd(*file, *yes)
}


///////////
// morph //
///////////

func morph(port int) {
	sub := flag.NewFlagSet("morph", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	knobs := sub.Bool("k", false, "source knobs")
	slot := sub.String("s", "", "source `slot` number")
	seed := sub.Int("i", timeseed(), "initial seed")
	mo := sub.Int("mo", 0, "oscillator `mult`iplier")
	mn := sub.Int("mn", 0, "noise `mult`iplier")
	me := sub.Int("me", 0, "eq (bass & treble) `mult`iplier")
	mf := sub.Int("mf", 0, "filter `mult`iplier")
	mr := sub.Int("mr", 0, "resonator `mult`iplier")
	sub.Parse(os.Args[2:])
	morph_cmd(port, *file, *knobs, *slot, *seed, *mo, *mn, *me, *mf, *mr)
}


//////////////
// updating //
//////////////

// read, processs, write all *.dlp in dir => dir2
func batch() {
	sub := flag.NewFlagSet("batch", flag.ExitOnError)
	dir := sub.String("d", ".", "source `directory` name")
	dir2 := sub.String("d2", ".", "target `directory` name")
	pro := sub.Bool("pro", false, "profile mode")
	mono := sub.Bool("m", false, "to mono")
	update := sub.Bool("u", false, "update")
	robs := sub.Bool("r", false, "Rob S. PP")
	yes := sub.Bool("y", false, "overwrite files")
	sub.Parse(os.Args[2:])
	//
	process_dlps(*dir, *dir2, *pro, *mono, *update, *robs, *yes)
}

// do a bunch of update stuff via interactive menu
func update(port int) {
	sub := flag.NewFlagSet("update", flag.ExitOnError)
	dir_all := sub.String("d", "_ALL_", "source `directory` name")
	dir_work := sub.String("d2", "_WORK_", "source `directory` name")
	sub.Parse(os.Args[2:])
	//
	update_cmd(port, *dir_all, *dir_work)
}


/////////
// dev //
/////////

func dev() {
	/*
	// find DLP files with certain characteristics
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	dir := sub.String("d", ".", "`directory` name")
	sub.Parse(os.Args[2:])
	//
	find_dlp(*dir)
	*/
	dev_cmd()
}
