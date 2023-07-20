package main

/*
 * Librarian for the D-Lev Theremin
 * See file "CHANGE_LOG.txt" for details
*/

import (
	"flag"
	"log"
	"os"
)

func main() {

	// do update if no args
	if len(os.Args) < 2 {
		update_cmd(WORK_DIR)
    } else {
		// parse subcommands
		switch os.Args[1] {
		case "update" : update()
		case "help", "-help", "-h", "/h" : help()
		case "ports" : ports()
		case "view" : view()
		case "match" : match()
		case "diff" : diff()
		case "ktof" : ktof()
		case "ftok" : ftok()
		case "stof" : stof()
		case "ftos" : ftos()
		case "btos" : btos()
		case "dump" : dump()
		case "pump" : pump()
		case "split" : split()
		case "join" : join()
		case "morph" : morph()
		case "batch" : batch()
		case "knob" : knob()
		case "hcl" : hcl_cmd()
		case "loop" :loop_cmd()
		case "ver" : ver_cmd()
		case "acal" : acal_cmd()
		case "reset" : reset_cmd()
		case "dev" : dev()  // dev stuff
		default : log.Fatalln("> Unknown command:", os.Args[1])
		}
	}
}  // end of main()


////////////////////
// main functions //
////////////////////

// show help
func help() {
	sub := flag.NewFlagSet("help", flag.ExitOnError)
	verbose := sub.Bool("v", false, "verbose mode")
	sub.Parse(os.Args[2:])
	//
	help_cmd(*verbose)
}

// list free serial ports / set port
func ports() {
	sub := flag.NewFlagSet("ports", flag.ExitOnError)
	port_str := sub.String("p", "", "`port` number")
	sub.Parse(os.Args[2:])
	//
	ports_cmd(*port_str)
}

// view knobs, DLP file, slot
func view() {
	sub := flag.NewFlagSet("view", flag.ExitOnError)
	file := sub.String("f", "", "view `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	knobs := sub.Bool("k", false, "view knobs")
	slot := sub.String("s", "", "view `slot` number")
	sub.Parse(os.Args[2:])
	//
	view_cmd(*file, *pro, *knobs, *slot)
}

// twiddle knob
func knob() {
	sub := flag.NewFlagSet("knob", flag.ExitOnError)
	knob := sub.String("k", "", "page:knob[0:6]")
	offset := sub.String("o", "", "knob offset value")
	val := sub.String("s", "", "knob set value")
	sub.Parse(os.Args[2:])
	//
	knob_cmd(*knob, *offset, *val)
}

// diff DLP file(s) / slot(s) / knobs
func diff() {
	sub := flag.NewFlagSet("diff", flag.ExitOnError)
	file := sub.String("f", "", "compare `file` name")
	file2 := sub.String("f2", "", "compare `file2` name")
	pro := sub.Bool("pro", false, "profile mode")
	knobs := sub.Bool("k", false, "compare knobs")
	slot := sub.String("s", "", "compare `slot` number")
	slot2 := sub.String("s2", "", "compare `slot2` number")
	sub.Parse(os.Args[2:])
	//
	diff_cmd(*file, *file2, *pro, *knobs, *slot, *slot2)
}

// match slots / DLP files w/ DLP files & list
func match() {
	sub := flag.NewFlagSet("match", flag.ExitOnError)
	dir := sub.String("d", ".", "`directory` name")
	dir2 := sub.String("d2", ".", "`directory` name")
	pro := sub.Bool("pro", false, "profile mode")
	hdr := sub.Bool("hdr", false, "header format")
	guess := sub.Bool("g", false, "guess")
	slots := sub.Bool("s", false, "slots")
	sub.Parse(os.Args[2:])
	//
	match_cmd(*dir, *dir2, *pro, *hdr, *guess, *slots)
}


//////////////////////
// download to file //
//////////////////////

// dump to file
func dump() {
	sub := flag.NewFlagSet("dump", flag.ExitOnError)
	file := sub.String("f", "", "target `file` name")
	yes := sub.Bool("y", false, "overwrite files")
	sub.Parse(os.Args[2:])
	//
	dump_cmd(*file, *yes)
}

// knobs => *.dlp
func ktof() {
	sub := flag.NewFlagSet("ktof", flag.ExitOnError)
	file := sub.String("f", "", "target `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	yes := sub.Bool("y", false, "overwrite files")
	sub.Parse(os.Args[2:])
	//
	ktof_cmd(*file, *pro, *yes)
}

// slot => *.dlp
func stof() {
	sub := flag.NewFlagSet("stof", flag.ExitOnError)
	slot := sub.String("s", "", "source `slot` number")
	file := sub.String("f", "", "target `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	yes := sub.Bool("y", false, "overwrite files")
	sub.Parse(os.Args[2:])
	//
	stof_cmd(*slot, *file, *pro, *yes)
}


//////////////////////
// upload from file //
//////////////////////

// pump from file
func pump() {
	sub := flag.NewFlagSet("pump", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	sub.Parse(os.Args[2:])
	//
	pump_cmd(*file)
}

// *.dlp => knobs
func ftok() {
	sub := flag.NewFlagSet("ftok", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	ftok_cmd(*file, *pro)
}

// *.dlp => slot
func ftos() {
	sub := flag.NewFlagSet("ftos", flag.ExitOnError)
	slot := sub.String("s", "", "target `slot` number")
	file := sub.String("f", "", "source `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	ftos_cmd(*slot, *file, *pro)
}

// *.bnk => *.dlps => slots
func btos() {
	sub := flag.NewFlagSet("btos", flag.ExitOnError)
	slot := sub.String("s", "", "starting `slot` number")
	file := sub.String("f", "", "bank `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	btos_cmd(*slot, *file, *pro)
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

func morph() {
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
	morph_cmd(*file, *knobs, *slot, *seed, *mo, *mn, *me, *mf, *mr)
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
func update() {
	sub := flag.NewFlagSet("update", flag.ExitOnError)
	dir := sub.String("d", WORK_DIR, "work `directory` name")
	sub.Parse(os.Args[2:])
	//
	update_cmd(*dir)	
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
