package main

/*
 * d-lev support functions
*/

import (
	"math"
	"fmt"
	"strings"
	"strconv"
	"log"
	"os"
	"path/filepath"
	"math/rand"
)

type param_t struct {
	ptype int
	plabel string
	pname string
}

var pro_params = []param_t {
	{0x01, "50Hz",    "s_p0_ds"},	// 0
	{0x07, "Dith",    "s_p1_ds"},	// 1
	{0x03, "P<>V",    "s_p2_ds"},	// 2
	{0x01, "Erev",    "s_p3_ds"},	// 3
	{0x07, "Dith",    "s_p4_ds"},	// 4
	{0x24, "LCD ",    "s_p5_ds"},	// 5
	{0xcb, "Pcal",    "p_p0_ds"},	// 6
	{0xc0, "Lin ",    "p_p1_ds"},	// 7
	{0x26, "Ofs-",    "p_p2_ds"},	// 8
	{0xcb, "Sens",    "p_p3_ds"},	// 9
	{0x26, "Ofs+",    "p_p4_ds"},	// 10
	{0xfd, "Cent",    "p_p5_ds"},	// 11
	{0xcb, "Vcal",    "v_p0_ds"},	// 12
	{0xc0, "Lin ",    "v_p1_ds"},	// 13
	{0x26, "Ofs-",    "v_p2_ds"},	// 14
	{0xcb, "Sens",    "v_p3_ds"},	// 15
	{0x24, "Drop",    "v_p4_ds"},	// 16
	{0x26, "Ofs+",    "v_p5_ds"},	// 17
	{0x35, "Mon ",    "v_p6_ds"},	// 18
	{0x01, "Out ",    "v_p7_ds"},	// 19
	{0x20, "LED ",    "t_p0_ds"},	// 20
	{0x04, "Qant",    "t_p1_ds"},	// 21
	{0x03, "Post",    "t_p2_ds"},	// 22
	{0xab, "Note",    "t_p3_ds"},	// 23
	{0xaf, "Oct ",    "t_p4_ds"},	// 24
	{0xc5, "Bass",   "eq_p0_ds"},	// 25
	{0xc5, "Treb",   "eq_p1_ds"},	// 26
	{0x35, "Line",    "v_p8_ds"},	// 27
	{0x7e, "Wait",    "s_p6_ds"},	// 28
	{0x24, "Lift",    "p_p6_ds"},	// 29
	{0x7e, "Auto",    "p_p7_ds"},	// 30
}

var not_params = []param_t {
	{0x7f, "    ", "menu_pg_ds"},	// 31 - NOT stored in *.dlp !  MENU_PG_IDX!
	{0xfe, "load",   "ps_p0_ds"},	// 32 - NOT stored in *.dlp !
	{0x7e, "stor",   "ps_p1_ds"},	// 33 - NOT stored in *.dlp !
	{0xa7, "Load",   "ps_p2_ds"},	// 34 - NOT stored in *.dlp !
	{0x07, "Stor",   "ps_p3_ds"},	// 35 - NOT stored in *.dlp !
}

var pre_params = []param_t {
	{0x35, "osc ",  "o_p0_ds"},	// 0
	{0x24, "odd ",  "o_p1_ds"},	// 1
	{0x24, "harm",  "o_p2_ds"},	// 2
	{0xca, "pmod",  "o_p3_ds"},	// 3
	{0xca, "vmod",  "o_p4_ds"},	// 4
	{0xa7, "oct ",  "o_p5_ds"},	// 5
	{0xf0, "offs",  "o_p6_ds"},	// 6
	{0x24, "xmix",  "o_p7_ds"},	// 7
	{0x24, "fm  ",  "o_p8_ds"},	// 8
	{0x70, "freq",  "o_p9_ds"},	// 9
	{0x76, "reso", "o_p10_ds"},	// 10
	{0xa4, "mode", "o_p11_ds"},	// 11
	{0xca, "pmod", "o_p12_ds"},	// 12
	{0xca, "vmod", "o_p13_ds"},	// 13
	{0xc5, "bass", "o_p14_ds"},	// 14
	{0xc5, "treb", "o_p15_ds"},	// 15
	{0xf0, "hmul", "o_p16_ds"},	// 16
	{0xf0, "hmul", "o_p17_ds"},	// 17
	{0xf0, "offs", "o_p18_ds"},	// 18
	{0x35, "sprd", "o_p19_ds"},	// 19
	{0x24, "xmix", "o_p20_ds"},	// 20
	{0x74, "nois",  "n_p0_ds"},	// 21
	{0x70, "freq",  "n_p3_ds"},	// 22
	{0x76, "reso",  "n_p4_ds"},	// 23
	{0xa4, "mode",  "n_p5_ds"},	// 24
	{0xca, "pmod",  "n_p6_ds"},	// 25
	{0xca, "vmod",  "n_p7_ds"},	// 26
	{0xca, "pmod",  "n_p8_ds"},	// 27
	{0xca, "vmod",  "n_p9_ds"},	// 28
	{0x34, "puls", "n_p10_ds"},	// 29
	{0xc5, "bass", "n_p11_ds"},	// 30
	{0x24, "xmix", "n_p12_ds"},	// 31
	{0xc5, "treb", "n_p13_ds"},	// 32
	{0x24, "duty", "n_p14_ds"},	// 33
	{0xf1, "reso",  "r_p0_ds"},	// 34
	{0xc6, "harm",  "r_p1_ds"},	// 35
	{0x72, "freq",  "r_p2_ds"},	// 36
	{0xc6, "tap ",  "r_p3_ds"},	// 37
	{0x71, "hpf ",  "r_p4_ds"},	// 38
	{0xc5, "xmix",  "r_p5_ds"},	// 39
	{0xa2, "mode",  "r_p6_ds"},	// 40
	{0x70, "freq",  "f_p0_ds"},	// 41
	{0x35, "levl",  "f_p1_ds"},	// 42
	{0x70, "freq",  "f_p2_ds"},	// 43
	{0x35, "levl",  "f_p3_ds"},	// 44
	{0x70, "freq",  "f_p4_ds"},	// 45
	{0x35, "levl",  "f_p5_ds"},	// 46
	{0xca, "pmod",  "f_p6_ds"},	// 47
	{0xca, "vmod",  "f_p7_ds"},	// 48
	{0xca, "pmod",  "f_p8_ds"},	// 49
	{0xca, "vmod",  "f_p9_ds"},	// 50
	{0xca, "pmod", "f_p10_ds"},	// 51
	{0xca, "vmod", "f_p11_ds"},	// 52
	{0x76, "reso", "f_p12_ds"},	// 53
	{0x76, "reso", "f_p13_ds"},	// 54
	{0x70, "freq", "f_p14_ds"},	// 55
	{0x35, "levl", "f_p15_ds"},	// 56
	{0x70, "freq", "f_p16_ds"},	// 57
	{0x35, "levl", "f_p17_ds"},	// 58
	{0x70, "freq", "f_p18_ds"},	// 59
	{0x35, "levl", "f_p19_ds"},	// 60
	{0x76, "reso", "f_p20_ds"},	// 61
	{0x70, "freq", "f_p22_ds"},	// 62
	{0x35, "levl", "f_p23_ds"},	// 63
	{0x70, "freq", "f_p24_ds"},	// 64
	{0x35, "levl", "f_p25_ds"},	// 65
	{0xca, "pmod", "f_p26_ds"},	// 66
	{0xca, "vmod", "f_p27_ds"},	// 67
	{0x76, "reso", "f_p28_ds"},	// 68
	{0x24, "cntr", "pc_p0_ds"},	// 69
	{0x24, "rate", "pc_p1_ds"},	// 70
	{0x44, "span", "pc_p2_ds"},	// 71
	{0x24, "corr", "pc_p3_ds"},	// 72
	{0xc9, "vmod", "pc_p4_ds"},	// 73
	{0x25, "kloc",  "e_p0_ds"},	// 74
	{0x44, "knee",  "e_p1_ds"},	// 75
	{0x76, "fall",  "e_p2_ds"},	// 76
	{0x75, "rise",  "e_p3_ds"},	// 77
	{0x34, "velo",  "e_p4_ds"},	// 78
	{0x73, "damp",  "e_p5_ds"},	// 79
	{0x25, "dloc",  "e_p6_ds"},	// 80
	{0x74, "prev", "pp_p0_ds"},	// 81
	{0xc5, "harm", "pp_p1_ds"},	// 82
	{0xa7, "oct ", "pp_p2_ds"},	// 83
	{0xca, "pmod", "pp_p3_ds"},	// 84
	{0xca, "vmod", "pp_p4_ds"},	// 85
	{0xc5, "treb", "pp_p5_ds"},	// 86
	{0xc5, "bass", "pp_p6_ds"},	// 87
	{0xb0, "chan",  "m_p0_ds"},	// 88
	{0x25, "vloc",  "m_p1_ds"},	// 89
	{0x42, "bend",  "m_p2_ds"},	// 90
	{0xa7, "oct ",  "m_p3_ds"},	// 91
	{0x34, "velo",  "m_p4_ds"},	// 92
	{0xf2, "cc  ",  "m_p5_ds"},	// 93
	{0x45, "cloc",  "m_p6_ds"},	// 94
	{0x8b, "prvw", "pp_p7_ds"},	// 95
	{0xa3, "bank",  "b_p0_ds"},	// 96
}

var knob_pnames = []string {  // these are in UI page order (hcl rk & wk knob order)
	"v_p6_ds",  "v_p7_ds",  "v_p0_ds",  "p_p0_ds",  "ps_p1_ds", "b_p0_ds",  "ps_p0_ds", "menu_pg_ds",  // [0:7] D-LEV
	"v_p6_ds",  "v_p8_ds",  "pp_p0_ds", "eq_p1_ds", "o_p0_ds",  "eq_p0_ds", "n_p0_ds",  "menu_pg_ds",  // [8:15] LEVELS
	"pp_p4_ds", "pp_p3_ds", "pp_p0_ds", "pp_p5_ds", "pp_p1_ds", "pp_p6_ds", "pp_p2_ds", "menu_pg_ds",  // [16:23] PREVIEW
	"m_p1_ds",  "m_p4_ds",  "m_p6_ds",  "m_p5_ds",  "m_p2_ds",  "m_p0_ds",  "m_p3_ds",  "menu_pg_ds",  // [24:31] MIDI
	"e_p0_ds",  "e_p3_ds",  "e_p1_ds",  "e_p2_ds",  "e_p4_ds",  "e_p5_ds",  "e_p6_ds",  "menu_pg_ds",  // [32:39] VOLUME
	"pc_p4_ds", "pc_p3_ds", "pc_p1_ds", "pc_p2_ds", "pc_p0_ds", "pp_p7_ds", "t_p2_ds",  "menu_pg_ds",  // [40:47] PITCH
	"n_p9_ds",  "n_p8_ds",  "n_p0_ds",  "n_p13_ds", "n_p10_ds", "n_p11_ds", "n_p14_ds", "menu_pg_ds",  // [48:55] NOISE
	"n_p7_ds",  "n_p6_ds",  "n_p3_ds",  "n_p0_ds",  "n_p5_ds",  "n_p12_ds", "n_p4_ds",  "menu_pg_ds",  // [56:63] FLT_NOISE
	"o_p4_ds",  "o_p3_ds",  "o_p2_ds",  "o_p15_ds", "o_p1_ds",  "o_p14_ds", "o_p5_ds",  "menu_pg_ds",  // [64:71] 0_OSC
	"o_p6_ds",  "o_p18_ds", "o_p16_ds", "o_p17_ds", "o_p8_ds",  "o_p19_ds", "o_p7_ds",  "menu_pg_ds",  // [72:79] 1_OSC
	"o_p13_ds", "o_p12_ds", "o_p9_ds",  "o_p0_ds",  "o_p11_ds", "o_p20_ds", "o_p10_ds", "menu_pg_ds",  // [80:87] FLT_OSC
	"r_p2_ds",  "r_p3_ds",  "r_p4_ds",  "r_p1_ds",  "r_p6_ds",  "r_p5_ds",  "r_p0_ds",  "menu_pg_ds",  // [95:95] RESON
	"f_p7_ds",  "f_p6_ds",  "f_p0_ds",  "f_p1_ds",  "f_p14_ds", "f_p15_ds", "f_p12_ds", "menu_pg_ds",  // [96:103] 0_FORM
	"f_p9_ds",  "f_p8_ds",  "f_p2_ds",  "f_p3_ds",  "f_p16_ds", "f_p17_ds", "f_p13_ds", "menu_pg_ds",  // [104:111] 1_FORM
	"f_p11_ds", "f_p10_ds", "f_p4_ds",  "f_p5_ds",  "f_p22_ds", "f_p23_ds", "f_p20_ds", "menu_pg_ds",  // [112:119] 2_FORM
	"f_p27_ds", "f_p26_ds", "f_p18_ds", "f_p19_ds", "f_p24_ds", "f_p25_ds", "f_p28_ds", "menu_pg_ds",  // [120:127] 3_FORM
	"v_p0_ds",  "v_p4_ds",  "v_p1_ds",  "s_p4_ds",  "v_p2_ds",  "v_p5_ds",  "v_p3_ds",  "menu_pg_ds",  // [128:135] V_FIELD
	"p_p0_ds",  "p_p6_ds",  "p_p1_ds",  "s_p1_ds",  "p_p2_ds",  "p_p4_ds",  "p_p3_ds",  "menu_pg_ds",  // [136:143] P_FIELD
	"t_p0_ds",  "p_p5_ds",  "s_p5_ds",  "t_p3_ds",  "t_p1_ds",  "t_p4_ds",  "t_p2_ds",  "menu_pg_ds",  // [144:151] DISPLAY
	"s_p6_ds",  "p_p7_ds",  "s_p2_ds",  "s_p0_ds",  "ps_p3_ds",  "s_p3_ds", "ps_p2_ds", "menu_pg_ds",  // [152:159] SYSTEM
}

var page_names = []string {
	"    D-LEV",
	"   LEVELS",
	"  PREVIEW",
	"     MIDI",
	"   VOLUME",
	"    PITCH",
	"    NOISE",
	"FLT_NOISE",
	"    0_OSC",
	"    1_OSC",
	"  FLT_OSC",
	"    RESON",
	"   0_FORM",
	"   1_FORM",
	"   2_FORM",
	"   3_FORM",
	"  V_FIELD",
	"  P_FIELD",
	"  DISPLAY",
	"   SYSTEM",
}

// return filter freq value (type 0x70, 0x71)
// 7041 * EXP2((ENC * (2^27) / 24) + 3/4)
// input: [0:192]
// output: 27 to 7040 (Hz)
func filt_freq(enc int) (int) {
	enc_mo := float64(int64(enc) * ((1 << 27) / 24) + 0xc0000000)
	return int(7041 * (math.Pow(2, enc_mo / math.Pow(2, 27)) / math.Pow(2, 32)))
}

// return reson freq value (type 0x72)
// 48001 / ((((((~(ENC<<25))^4)*0.871)+((~(ENC<<25))>>3))>>22)+4)
// input: [0:127]
// output: 46 to 9600 (Hz)
func reson_freq(enc int) (int) {
	fs_rev := uint64(^(uint32(enc) << 25))
	sq := (fs_rev * fs_rev) >> 32
	qd := (sq * sq) >> 32
	return int(48001 / ((((uint64(float64(qd) * 0.871) + (fs_rev >> 3)) >> 22) + 4)))
}

// given encoder value and type, return display string[5]
func enc_disp(enc int, ptype int) (string) {
	switch ptype {
		case 0x70, 0x71 : enc = filt_freq(enc)
		case 0x72 : enc = reson_freq(enc)
		default : if ptype >= 0x80 { enc = int(int8(uint8(enc))) }  // signed
	}
	return fmt.Sprintf("%5v", enc)
}

// given pname, return ptype, plabel, pidx, pgroup
func pname_lookup(pname string) (int, string, int, string) {
	for pidx, param := range pre_params {
		if pname == param.pname { return param.ptype, param.plabel, pidx, "pre" }
	}
	for pidx, param := range pro_params {
		if pname == param.pname { return param.ptype, param.plabel, pidx, "pro" }
	}
	for pidx, param := range not_params {
		if pname == param.pname { return param.ptype, param.plabel, pidx, "not" }
	}
	return 0, "", 0, ""
}

// generate knob ui display strings
func knob_ui_strs(hex_str string) ([]string) {
	var strs []string
	ints := hexs_to_ints(hex_str, 1)
	for kidx, kname := range knob_pnames {
		ptype, plabel, _, _ := pname_lookup(kname)
		if kidx % UI_KNOBS == UI_PAGE_KNOB { 
			strs = append(strs, page_names[kidx / UI_KNOBS])
		} else { 
			strs = append(strs, plabel + enc_disp(ints[kidx], ptype)) 
		}
	}
	return strs
}

// generate pre / pro / slot ui display strings
func pre_ui_strs(hex_str string, pro bool) ([]string) {
	var strs []string
	ints := hexs_to_ints(hex_str, 4)
	for i, pname := range knob_pnames {
		ptype, plabel, pidx, pgroup := pname_lookup(pname)
		if i % UI_KNOBS == UI_PAGE_KNOB { 
			strs = append(strs, page_names[i / UI_KNOBS])
		} else { 
			if pro == (pgroup == "pro") && pgroup != "not" {
				strs = append(strs, plabel + enc_disp(ints[pidx], ptype)) 
			} else {
				strs = append(strs, plabel + enc_disp(0, ptype)) 
			}
		}
	}
	return strs
}

// render ui display strings to printable string
func ui_prn_str(strs []string) (string) {
	h_line_sub := "+" + strings.Repeat("-", 22);
	h_line := strings.Repeat(h_line_sub, UI_PRN_PG_COLS) + "+\n";
	prn_str := h_line
	for prow:=0; prow<UI_PRN_PG_ROWS; prow++ {
		for uirow:=0; uirow<UI_ROWS; uirow++ {
			for pcol:=0; pcol<UI_PRN_PG_COLS; pcol++ {
				idx := (prow * UI_COLS * UI_ROWS * UI_PRN_PG_COLS) + (uirow * UI_COLS) + (pcol * UI_COLS * UI_ROWS)
				prn_str += "| " + strs[idx] + "  " + strs[idx+1] + " "
			}
			prn_str += "|\n"
		}
		prn_str += h_line
	}
	return strings.TrimSpace(prn_str)
}

// map slots in preset string, return slice
func map_slots(pre string, file_map map[string]string) ([]string) {
	var strs []string
	split_strs := strings.Split(pre, "\n")
	for s:=0; s<SLOTS; s++ {
		pre_str := ""
		for i:=s*SLOTS/4; i<(s+1)*SLOTS/4; i++ {
			pre_str += split_strs[i] + "\n"
		}
		file, exists := file_map[strings.TrimSpace(pre_str)]
		if !exists { file = "_??_" }
		strs = append(strs, file)
	}
	return strs
}

// render slots display strings to printable string
func slots_prn_str(strs []string) (string) {
	var prn_str string
	for row:=0; row<=120; row++ {
		if row == 0 { // pre 0 filler
			prn_str += strings.Repeat(" ", 25)
		} else { // pres < 0
			prn_str += fmt.Sprintf("[%4v] %-18s", -row, strings.TrimSpace(strs[256-row]))
		}  // pres >= 0
		prn_str += fmt.Sprintf("[%3v] %-18s", row, strings.TrimSpace(strs[row]))
		//
		if row < 8 {  // pros >= 0
			prn_str += fmt.Sprintf("[%1v] %-18s", row, strings.TrimSpace(strs[row+128]))
			if row > 0 {  // pros < 0
				prn_str += fmt.Sprintf("[%2v] %-18s", -row, strings.TrimSpace(strs[128-row]))
			}
		}
		prn_str += "\n"
	}
	return prn_str
}


// render slots list to BNK file string
func slots_bnk_str(strs []string) (string) {
	var bnk_str string
	for i, _ := range strs {
		if i <= 120 {
			if i == 0 { 
				bnk_str += "////////////////\n"
				bnk_str += "// pre[0:120] //\n"
				bnk_str += "////////////////\n"
			}
			bnk_str += strs[i] + "\n"
		} else if i >= 121 && i < 128 {
			if i == 121 { 
				bnk_str += "////////////////\n"
				bnk_str += "// pro[-1:-7] //\n"
				bnk_str += "////////////////\n"
			}
			bnk_str += strs[248-i] + "\n"
		} else if i >= 128 && i < 136 {
			if i == 128 {
				bnk_str += "//////////////\n"
				bnk_str += "// pro[0:7] //\n"
				bnk_str += "//////////////\n"
			}
			bnk_str += strs[i] + "\n"
		} else if i >= 136 && i < 256 {
			if i == 136 {
				bnk_str += "//////////////////\n"
				bnk_str += "// pre[-1:-120] //\n"
				bnk_str += "//////////////////\n"
			}
			bnk_str += strs[256-i] + "\n"
		}
	}
	return bnk_str
}


/////////
// dev //
/////////

// find DLP files with various values
func find_dlp(dir string) {
	dir = filepath.Clean(dir)
	files, err := os.ReadDir(dir); if err != nil { log.Fatal(err) }
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".dlp" && !file.IsDir() {
			file := filepath.Join(dir, file.Name())
			// read in
			file_bytes, err := os.ReadFile(file); if err != nil { log.Fatal(err) }
			ints := hexs_to_ints(string(file_bytes), 4)
			// PITCH:corr idx=72, PITCH:vmod idx=73
			if ints[72] > 0 { 
				fmt.Println(">", file, "corr", ints[72], "vmod", ints[73])
			}
		}
    }
}

// read, update, write all *.dlp files in dir
func update_dlp(dir string, pro bool) {
	dir = filepath.Clean(dir)
	// prompt user
	fmt.Print("> Update all DLP files in directory ", dir, " ?  <y|n> ")
	var input string
	fmt.Scanln(&input)
	if input != "y" { log.Fatalln("> Abort, exiting program..." ) }
	files, err := os.ReadDir(dir); if err != nil { log.Fatal(err) }
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".dlp" && file.IsDir() == false {
			file := filepath.Join(dir, file.Name())
			// read in
			file_bytes, err := os.ReadFile(file); if err != nil { log.Fatal(err) }
			ints := hexs_to_ints(string(file_bytes), 4)
			if pro {  // PROFILES
				// zero out high fluff
				for i, _ := range ints {
					if i >= len(pro_params) { ints[i] = 0 }
				}
				// V_FIELD:Drop strength 1/2
				ints[16] = ints[16] / 2
			} else {  // PRESETS
				// zero out high fluff
				for i, _ := range ints {
					if i >= len(pre_params) { ints[i] = 0 }
				}
				// PITCH:corr = 0
				ints[72] = 0
				// PITCH:vmod = -15
				ints[73] = -15
			}
			// write back
			hexs := ints_to_hexs(ints, 4)
			err = os.WriteFile(file, []byte(hexs), 0666); if err != nil { log.Fatal(err) }
			fmt.Println("> updated file", file)
		}
    }
}

// make some dlp files for testing
func gen_test_dlps(dir string) {
	dir = filepath.Clean(dir)
	// prompt user
	fmt.Print("> Generate test DLP files in directory ", dir, " ?  <y|n> ")
	var input string
	fmt.Scanln(&input)
	if input != "y" { log.Fatalln("> Abort, exiting program..." ) }
	for f:=0; f<256; f++ {  // generate 256 files
		name := strconv.Itoa(f) + ".dlp"
		file := filepath.Join(dir, name)
		var str string
		for ln:=0; ln<64; ln++ {  // 64 lines
			str += strconv.FormatInt(int64(rand.Intn(0x100000000)), 16) +"\n"
		}
		// write file
		err := os.WriteFile(file, []byte(str), 0666); if err != nil { log.Fatal(err) }
		fmt.Println("> created file", file)
	}
}
