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

var pro_params = []param_t {  // these are in preset / profile / slot order
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
	{0x31, "Mon ",    "v_p6_ds"},	// 18
	{0x01, "Out ",    "v_p7_ds"},	// 19
	{0x20, "LED ",    "t_p0_ds"},	// 20
	{0x04, "Qant",    "t_p1_ds"},	// 21
	{0x03, "Post",    "t_p2_ds"},	// 22
	{0xab, "Note",    "t_p3_ds"},	// 23
	{0xaf, "Oct ",    "t_p4_ds"},	// 24
	{0xc5, "Bass",   "eq_p0_ds"},	// 25
	{0xc5, "Treb",   "eq_p1_ds"},	// 26
	{0x31, "Line",    "v_p8_ds"},	// 27
	{0x7d, "Wait",    "s_p6_ds"},	// 28
	{0x24, "Lift",    "p_p6_ds"},	// 29
	{0x7e, "Auto",    "p_p7_ds"},	// 30
}

var not_params = []param_t {  // these are in sequence
	{0x7f, "    ", "menu_pg_ds"},	// 31 - NOT stored in *.dlp !  MENU_PG_IDX!
	{0x7e, "load",   "ps_p0_ds"},	// 32 - NOT stored in *.dlp !
	{0x7e, "stor",   "ps_p1_ds"},	// 33 - NOT stored in *.dlp !
	{0x05, "Load",   "ps_p2_ds"},	// 34 - NOT stored in *.dlp !
	{0x05, "Stor",   "ps_p3_ds"},	// 35 - NOT stored in *.dlp !
}

var pre_params = []param_t {  // these are in preset / profile / slot order
	// oscillators:
	{0x31, "osc ",  "o_p0_ds"},	// 0
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
	{0x31, "sprd", "o_p19_ds"},	// 19
	{0x24, "xmix", "o_p20_ds"},	// 20
	// noise gen:
	{0x74, "nois",  "n_p0_ds"},	// 21
	{0x70, "freq",  "n_p3_ds"},	// 22
	{0x76, "reso",  "n_p4_ds"},	// 23
	{0xa4, "mode",  "n_p5_ds"},	// 24
	{0xca, "pmod",  "n_p6_ds"},	// 25
	{0xca, "vmod",  "n_p7_ds"},	// 26
	{0xca, "pmod",  "n_p8_ds"},	// 27
	{0xca, "vmod",  "n_p9_ds"},	// 28
	{0x30, "puls", "n_p10_ds"},	// 29
	{0xc5, "bass", "n_p11_ds"},	// 30
	{0x24, "xmix", "n_p12_ds"},	// 31
	{0xc5, "treb", "n_p13_ds"},	// 32
	{0x24, "duty", "n_p14_ds"},	// 33
	// resonator:
	{0xf1, "reso",  "r_p0_ds"},	// 34
	{0xc6, "harm",  "r_p1_ds"},	// 35
	{0x72, "freq",  "r_p2_ds"},	// 36
	{0xc6, "tap ",  "r_p3_ds"},	// 37
	{0x71, "hpf ",  "r_p4_ds"},	// 38
	{0xc5, "xmix",  "r_p5_ds"},	// 39
	{0xa2, "mode",  "r_p6_ds"},	// 40
	// formants:
	{0x70, "freq",  "f_p0_ds"},	// 41
	{0xf2, "levl",  "f_p1_ds"},	// 42
	{0x70, "freq",  "f_p2_ds"},	// 43
	{0xf2, "levl",  "f_p3_ds"},	// 44
	{0x70, "freq",  "f_p4_ds"},	// 45
	{0xf2, "levl",  "f_p5_ds"},	// 46
	{0xca, "pmod",  "f_p6_ds"},	// 47
	{0xca, "vmod",  "f_p7_ds"},	// 48
	{0xca, "pmod",  "f_p8_ds"},	// 49
	{0xca, "vmod",  "f_p9_ds"},	// 50
	{0xca, "pmod", "f_p10_ds"},	// 51
	{0xca, "vmod", "f_p11_ds"},	// 52
	{0x76, "reso", "f_p12_ds"},	// 53
	{0x76, "reso", "f_p13_ds"},	// 54
	{0x70, "freq", "f_p14_ds"},	// 55
	{0xf2, "levl", "f_p15_ds"},	// 56
	{0x70, "freq", "f_p16_ds"},	// 57
	{0xf2, "levl", "f_p17_ds"},	// 58
	{0x70, "freq", "f_p18_ds"},	// 59
	{0xf2, "levl", "f_p19_ds"},	// 60
	{0x76, "reso", "f_p20_ds"},	// 61
	{0x70, "freq", "f_p22_ds"},	// 62
	{0xf2, "levl", "f_p23_ds"},	// 63
	{0x70, "freq", "f_p24_ds"},	// 64
	{0xf2, "levl", "f_p25_ds"},	// 65
	{0xca, "pmod", "f_p26_ds"},	// 66
	{0xca, "vmod", "f_p27_ds"},	// 67
	{0x76, "reso", "f_p28_ds"},	// 68
	// pitch correction:
	{0x24, "cntr", "pc_p0_ds"},	// 69
	{0x24, "rate", "pc_p1_ds"},	// 70
	{0x44, "span", "pc_p2_ds"},	// 71
	{0x24, "corr", "pc_p3_ds"},	// 72
	{0xc9, "vmod", "pc_p4_ds"},	// 73
	// envelope gen:
	{0x25, "kloc",  "e_p0_ds"},	// 74
	{0x44, "knee",  "e_p1_ds"},	// 75
	{0x76, "fall",  "e_p2_ds"},	// 76
	{0x75, "rise",  "e_p3_ds"},	// 77
	{0x30, "velo",  "e_p4_ds"},	// 78
	{0x73, "damp",  "e_p5_ds"},	// 79
	{0x25, "dloc",  "e_p6_ds"},	// 80
	// pitch preview:
	{0x74, "prev", "pp_p0_ds"},	// 81
	{0xc5, "harm", "pp_p1_ds"},	// 82
	{0xa7, "oct ", "pp_p2_ds"},	// 83
	{0xca, "pmod", "pp_p3_ds"},	// 84
	{0xca, "vmod", "pp_p4_ds"},	// 85
	{0xc5, "treb", "pp_p5_ds"},	// 86
	{0xc5, "bass", "pp_p6_ds"},	// 87
	// midi:
	{0xb0, "chan",  "m_p0_ds"},	// 88
	{0x25, "vloc",  "m_p1_ds"},	// 89
	{0x42, "bend",  "m_p2_ds"},	// 90
	{0xa7, "oct ",  "m_p3_ds"},	// 91
	{0x30, "velo",  "m_p4_ds"},	// 92
	{0xfc, "cc  ",  "m_p5_ds"},	// 93
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

// given pidx & pro, return kidx, kflg
func knob_lookup(pidx int, pro bool) (int, bool) {
	if pro {
		if pidx >= len(pro_params) { return 0, false }
		for kidx, kname := range knob_pnames {
			if kname == pro_params[pidx].pname { return kidx, true }
		}
	} else {
		if pidx >= len(pre_params) { return 0, false }
		for kidx, kname := range knob_pnames {
			if kname == pre_params[pidx].pname { return kidx, true }
		}
	}
	return 0, false
}

// put knob hex str in preset / slot order, return hex string
func knob_pre_order(knob_str string, pro bool) (string) {
	str_split := (strings.Split(strings.TrimSpace(knob_str), "\n"))
	if len(str_split) < KNOBS { log.Fatalln("> Bad knob info!") }
	hex_str := ""
	line_str := ""
	for pidx:=0; pidx<SLOT_BYTES; pidx++ {
		kidx, kflg := knob_lookup(pidx, pro)
		if kflg { line_str = fmt.Sprintf("%02s", str_split[kidx]) + line_str } else { line_str = "00" + line_str }
		if pidx % 4 == 3 { 
			hex_str += line_str + "\n" 
			line_str = ""
		}
	}
	return hex_str
}

// generate knob ui display strings
func knob_ui_strs(hex_str string) ([]string) {
	ints := hexs_to_ints(hex_str, 1)
	if len(ints) < KNOBS { log.Fatalln("> Bad knob info!") }
	var strs []string
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
	ints := hexs_to_ints(hex_str, 4)
	if len(ints) < SLOT_BYTES { log.Fatalln("> Bad file / slot info!") }
	var strs []string
	for idx, pname := range knob_pnames {
		ptype, plabel, pidx, pgroup := pname_lookup(pname)
		if idx % UI_KNOBS == UI_PAGE_KNOB { 
			strs = append(strs, page_names[idx / UI_KNOBS])
		} else { 
			if pro == (pgroup == "pro") && pgroup != "not" {
				strs = append(strs, plabel + enc_disp(ints[pidx], ptype)) 
			} else {
				strs = append(strs, plabel + "     ") 
			}
		}
	}
	return strs
}

// render ui display strings to printable string
func ui_prn_str(strs []string) (string) {
	if len(strs) < len(knob_pnames) { log.Fatalln("> Bad input info!") }
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

// generate diff display strings
func diff_strs(hex_str0, hex_str1 string, pro bool) ([]string, []string, []bool) {
	ints0 := hexs_to_ints(hex_str0, 4)
	ints1 := hexs_to_ints(hex_str1, 4)
	if (len(ints0) < SLOT_BYTES) || (len(ints1) < SLOT_BYTES) { log.Fatalln("> Bad file info!") }
	var strs0 []string
	var strs1 []string
	var diffs []bool
	for idx, pname := range knob_pnames {
		ptype, plabel, pidx, pgroup := pname_lookup(pname)
		if idx % UI_KNOBS == UI_PAGE_KNOB { 
			strs0 = append(strs0, page_names[idx / UI_KNOBS])
			strs1 = append(strs1, page_names[idx / UI_KNOBS])
			diffs = append(diffs, false)
		} else { 
			if pro == (pgroup == "pro") && pgroup != "not" {
				strs0 = append(strs0, plabel + enc_disp(ints0[pidx], ptype)) 
				if ints0[pidx] != ints1[pidx] {
					strs1 = append(strs1, plabel + enc_disp(ints1[pidx], ptype)) 
					diffs = append(diffs, true)
				} else {
					strs1 = append(strs1, plabel + "     ") 
					diffs = append(diffs, false)
				}
			} else {
				strs0 = append(strs0, plabel + "     ") 
				strs1 = append(strs1, plabel + "     ") 
				diffs = append(diffs, false)
			}
		}
	}
	return strs0, strs1, diffs
}

// render ui display strings to printable string
func diff_prn_str(strs0, strs1 []string, diffs []bool) (string) {
	if (len(strs0) < len(knob_pnames)) || (len(strs1) < len(knob_pnames)) || (len(diffs) < len(knob_pnames)) { 
		log.Fatalln("> Bad input info!") 
	}
	h_line_sub := "+" + strings.Repeat("-", 22);
	h_line := strings.Repeat(h_line_sub, 2) + "+\n";
	prn_str := ""
	chgs := 0
	for uipg:=0; uipg<UI_PAGES; uipg++ {
		pg_str := ""
		chg_f := false
		for uirow:=0; uirow<UI_ROWS; uirow++ {
			idx := uipg*UI_ROWS*UI_COLS + uirow*UI_COLS
			pg_str += "| " + strs0[idx] + "  " + strs0[idx+1] + " "
			pg_str += "| " + strs1[idx] + "  " + strs1[idx+1] + " "
			pg_str += "|\n"
			if diffs[idx] { chg_f = true; chgs++ }
			if diffs[idx+1] { chg_f = true; chgs++ }
		}
		pg_str += h_line
		if chg_f { prn_str += pg_str }
	}
	if chgs != 0 { prn_str = h_line + prn_str }  // top line
	prn_str += fmt.Sprintln("> differences", chgs)
	return strings.TrimSpace(prn_str)
}

// map slots in preset string, return slice
func map_slots(pre string, file_map map[string]string) ([]string) {
	split_strs := strings.Split(pre, "\n")
	if len(split_strs) < SLOTS * SLOT_BYTES/4 { log.Fatalln("> Bad slots info!") }
	var strs []string
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
	if len(strs) < SLOTS { log.Fatalln("> Bad slots info!") }
	var prn_str string
	for row:=0; row<PRE_SLOTS/5; row++ {
		for col:=0; col<5; col++ {
			idx := col*50 + row
			prn_str += fmt.Sprintf("[%2v] %-15s", idx, strings.TrimSpace(strs[idx]))
		}
		if row < PRO_SLOTS {  // pros
			prn_str += fmt.Sprintf("[%1v] %s", row, strings.TrimSpace(strs[row+PRE_SLOTS]))
		}
		prn_str += "\n"
	}
	return prn_str
}

// render slots list to BNK file string
func slots_bnk_str(strs []string, headers bool) (string) {
	var bnk_str string
	for i, _ := range strs {
		if headers && (i % 10 == 0) {
			ofs := 0
			inc := 9
			bnk_str += "// "
			if i >= PRE_SLOTS {
				ofs = -PRE_SLOTS
				inc = PRO_SLOTS-1
				bnk_str += "PRO"
			}
			bnk_str += fmt.Sprint("[", i+ofs, ":", i+ofs+inc, "]\n")
		}
		bnk_str += strs[i] + "\n"
	}
	return bnk_str
}


////////////
// update //
////////////

// read, update, write all *.dlp files in dir
func update_dlp(dir string, pro bool, dry bool) {
	dir = filepath.Clean(dir)
	// prompt user
	if !dry { 
		fmt.Print("> Update all DLP files in directory ", dir, " ?  <y|n> ")
		var input string
		fmt.Scanln(&input)
		if input != "y" { log.Fatalln("> Abort, exiting program..." ) }
	}
	files, err := os.ReadDir(dir); if err != nil { log.Fatal(err) }
	upd_cnt := 0
	dlp_cnt := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".dlp" && file.IsDir() == false {
			dlp_cnt++
			file_name := file.Name()
			file_path := filepath.Join(dir, file_name)
			fmt.Println(dlp_cnt, "-", file_name)
			// read in
			file_bytes, err := os.ReadFile(file_path); if err != nil { log.Fatal(err) }
			ints := hexs_to_ints(string(file_bytes), 4)
			upd_f := false
			if pro {  // PROFILES
				// zero out high fluff
				nz_f := false
				for idx, param := range ints {
					if idx >= len(pro_params) { 
						if param != 0 { 
							ints[idx] = 0
							nz_f = true 
						}
					}
				}
				if nz_f { 
					fmt.Println("- fluff zeroed")
					upd_f = true 
				}
			} else {  // PRESETS
				// zero out high fluff
				nz_f := false
				for idx, param := range ints {
					if idx >= len(pre_params) { 
						if param != 0 { 
							ints[idx] = 0
							nz_f = true 
						}
					}
				}
				if nz_f { 
					fmt.Println("- fluff zeroed")
					upd_f = true 
				}
				// update all formant levels (if necessary)
				levls := []int{  // signed
					int(int8(ints[42])), 
					int(int8(ints[44])), 
					int(int8(ints[46])), 
					int(int8(ints[56])), 
					int(int8(ints[58])), 
					int(int8(ints[60])), 
					int(int8(ints[63])), 
					int(int8(ints[65])) }
				min := 0
				max := 0
				for _, levl := range levls {
					if levl > max { max = levl }
					if levl < min { min = levl }
				}
				abs_max := max
				if -min > max { abs_max = -min }
				// only update if formants in use
				if abs_max != 0 {
					delta := 8  // +6dB
					if abs_max + delta > 63 { delta = 63 - abs_max }
					if delta < 8 { 
						fmt.Println("- FORM:delta", delta) 
					}
					if delta != 0 {
						for idx, _ := range levls {
							if levls[idx] < 0 { 
								levls[idx] -= delta
							} else if levls[idx] > 0 { 
								levls[idx] += delta 
							}
						}
						fmt.Println("- FORM:levl", levls)
						ints[42] = levls[0]
						ints[44] = levls[1]
						ints[46] = levls[2]
						ints[56] = levls[3]
						ints[58] = levls[4]
						ints[60] = levls[5]
						ints[63] = levls[6]
						ints[65] = levls[7]
						upd_f = true
					}
					// update reson xmix level (if necessary)
					xmix := int(int8(ints[39]))  // signed
					mode := int(int8(ints[40]))  // signed
					if (mode <= 0) && (xmix != 0) && (delta < 8) {
						db_ratio := math.Pow(10, ((float64(delta)-8)*3/4)/20)
						xm_norm := float64(xmix) / 64
						xm_sign := math.Copysign(1, xm_norm)
						xm_abs := xm_sign * xm_norm
						xm_ratio := xm_abs / (1 + xm_abs)
						xm_ratio_new := xm_ratio * db_ratio
						xmix_new := int(math.Round(64 * xm_sign * xm_ratio_new / (1 - xm_ratio_new)))
						ints[39] = xmix_new
						fmt.Println("- RESON:xmix", xmix, "=>", xmix_new)
						upd_f = true
					}
				}
				// update reson mode (if necessary)
				reson_mode := int(int8(ints[40]))  // signed
				if reson_mode != 0 {
					mode_new := -reson_mode
					ints[40] = mode_new
					fmt.Println("- RESON:mode", reson_mode, "=>", mode_new)
					upd_f = true
				}
				// update noise knobs
				if ints[21] == 0 {  // if nois[0] kill everything
					nz_f := false
					for idx, levl := range ints {
						if (idx >= 22) && (idx <= 33) { 
							if levl != 0 { 
								nz_f = true
								ints[idx] = 0
							}
						}
					}
					if nz_f { 
						fmt.Println("- NOISE:all: => 0")
						upd_f = true 
					}
				} else { // adjust nois
					nois := float64(ints[21])  // unsigned
					vmod := float64(int8(ints[28]))  // signed
					nois_new := int(math.Round(nois - 8 + 0.025*vmod*vmod))
					ints[21] = nois_new
					fmt.Println("- NOISE:nois", nois, "=>", nois_new)
					upd_f = true
					if ints[28] != 0 {  // if vmod != 0 : adjust vmod
						vmod := float64(int8(ints[28]))  // signed
						vmod_new := int(math.Round(-1.2 * (1 + vmod)))
						ints[28] = vmod_new
						fmt.Println("- NOISE:vmod", vmod, "=>", vmod_new)
					}
					if ints[27] != 0 {  // if pmod != 0 : 1/2 strength
						pmod := int(int8(ints[27]))  // signed
						pmod_norm := float64(pmod) / 64
						pmod_sign := math.Copysign(1, pmod_norm)
						pmod_new := int(math.Round(pmod_sign * math.Sqrt(pmod_norm*pmod_norm / 2) * 64))
						ints[27] = pmod_new
						fmt.Println("- NOISE:pmod", pmod, "=>", pmod_new)
					}
				}
				// normalize pitch correction
				cntr := ints[69]
				rate := ints[70]
				span := ints[71]
				corr := ints[72]
				vmod := int(int8(ints[73]))
				vmod_old := vmod
				corr = 0  // default
				vmod = -15  // default
				if span == 0 { rate = 12; cntr = 12; span = 31 }  // defaults
				if rate < 15 { rate = 12; cntr = 12; span = 31 }  // defaults
				if cntr != ints[69] { fmt.Println("- PITCH:cntr", ints[69], "=>", cntr); upd_f = true }
				if rate != ints[70] { fmt.Println("- PITCH:rate", ints[70], "=>", rate); upd_f = true }
				if span != ints[71] { fmt.Println("- PITCH:span", ints[71], "=>", span); upd_f = true }
				if corr != ints[72] { fmt.Println("- PITCH:corr", ints[72], "=>", corr); upd_f = true }
				if vmod != vmod_old { fmt.Println("- PITCH:vmod", vmod_old, "=>", vmod); upd_f = true }
 				ints[69] = cntr
				ints[70] = rate
				ints[71] = span
				ints[72] = corr
				ints[73] = vmod
				// normaize stereo
				reson_xmix := int(int8(ints[39]))
				reson_mode = int(int8(ints[40]))
				if (reson_xmix == 0) && (reson_mode == 0) {  // do stereo defaults
					reso := int(int8(ints[34]))
					harm := int(int8(ints[35]))
					freq := ints[36]
					tap  := int(int8(ints[37]))
					hpf  := ints[38]
					//
					reso_new := 0
					harm_new := 17
					freq_new := 0   // 46Hz
					tap_new  := 42
					hpf_new  := 99  // 479Hz
					mode_new := -2
					if reso != reso_new { fmt.Println("- RESO:reso", reso, "=>", reso_new); upd_f = true }
					if harm != harm_new { fmt.Println("- RESO:harm", harm, "=>", harm_new); upd_f = true }
					if freq != freq_new { fmt.Println("- RESO:freq", freq, "=>", freq_new); upd_f = true }
					if tap  != tap_new { fmt.Println("- RESO:tap",  tap, "=>", tap_new); upd_f = true }
					if hpf  != hpf_new { fmt.Println("- RESO:hpf",  hpf, "=>", hpf_new); upd_f = true }
					if reson_mode != mode_new { fmt.Println("- RESO:mode", reson_mode, "=>", mode_new); upd_f = true }
					ints[34] = reso_new
					ints[35] = harm_new
					ints[36] = freq_new
					ints[37] = tap_new
					ints[38] = hpf_new
					ints[40] = mode_new
				}
			}
			// write back
			if upd_f {
				hexs := ints_to_hexs(ints, 4)
				if !dry { err = os.WriteFile(file_path, []byte(hexs), 0666); if err != nil { log.Fatal(err) } }
				fmt.Println("")
				upd_cnt++
			} else {
				fmt.Println("- no changes -\n")
			}
		}
    }
	fmt.Println("> updated", upd_cnt, "of", dlp_cnt, "DLP files in", dir, "directory")
	if dry { fmt.Println("\n- DRY RUN, NO FILES UPDATED - (use -dry=false to update)\n") }
}


/////////
// dev //
/////////

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

// find DLP files with various values
func find_dlp(dir string) {
	dir = filepath.Clean(dir)
	files, err := os.ReadDir(dir); if err != nil { log.Fatal(err) }
	dlp_cnt := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".dlp" && !file.IsDir() {
			dlp_cnt++
			file_name := file.Name()
			file_path := filepath.Join(dir, file_name)
			// read in
			file_bytes, err := os.ReadFile(file_path); if err != nil { log.Fatal(err) }
			ints := hexs_to_ints(string(file_bytes), 4)

			if ints[39] == 0 {
				fmt.Println("xmix", int8(ints[39]))
				fmt.Println("mode", int8(ints[40]))
				fmt.Println("tap ", int8(ints[37]))
				fmt.Println("harm", int8(ints[35]))
				fmt.Println("reso", int8(ints[34]))
				fmt.Println("freq", ints[36], reson_freq(ints[36]))
				fmt.Println("hpf ", ints[38], filt_freq(ints[38]))
				fmt.Println("-", file.Name(), "\n")
			}


/*
			// formant levels
			levls := []int{  // signed
				int(int8(ints[42])), 
				int(int8(ints[44])), 
				int(int8(ints[46])), 
				int(int8(ints[56])), 
				int(int8(ints[58])), 
				int(int8(ints[60])), 
				int(int8(ints[63])), 
				int(int8(ints[65])) }
			// formant resonances
			resos := []int{  // unsigned
				ints[53],
				ints[54],
				ints[61],
				ints[68] }
			fmt.Println("osc ", ints[0])
			fmt.Println("levl", levls)
			fmt.Println("reso", resos)
			fmt.Println("nois", ints[21])
			fmt.Println("-", file.Name(), "\n")
*/
/*
			// PITCH:corr[72], PITCH:vmod[73]
			if ints[72] > 0 { 
				fmt.Println(">", file, "corr", ints[72], "vmod", ints[73])
			}
*/
/*
			cntr := ints[69]
			rate := ints[70]
			span := ints[71]
			corr := ints[72]
			vmod := int(int8(ints[73]))
			pitch := []int{ corr, rate, cntr, vmod, span }
			flag := "***"
			if corr != 0 { corr = 0 }
			if span == 0 { rate = 12; cntr = 12; vmod = -15; span = 31 }
			if rate < 15 { rate = 12; cntr = 12; vmod = -15; span = 31 }
			new := []int{ corr, rate, cntr, vmod, span }
			for i, old := range pitch {
				if old != new[i] { flag = "" }
			}
			fmt.Println(dlp_cnt, pitch, new, file_name, flag)
*/
/*
			cntr := ints[69]
			rate := ints[70]
			span := ints[71]
			corr := ints[72]
			vmod := int(int8(ints[73]))
			pitch := []int{ corr, rate, cntr, vmod, span }
			fmt.Println(dlp_cnt, pitch, file_name)
*/
/*
			// NOISE:nois[21], vmod[28], pmod[27]
			if ints[28] != 0 { 
				fmt.Println("nois", ints[21], "vmod", int(int8(ints[28])), "pmod", int(int8(ints[27])), file.Name())
			}
*/			
/*
			// NOISE:nois[21], the rest [22:33]
			if ints[21] == 0 { 
				for idx, levl := range ints {
					if (idx >= 22) && (idx <= 33) { 
						if levl != 0 { 
							fmt.Println("idx:", idx, "levl:", levl, ":", file.Name())
						}
					}
				}
			}
*/
		}
    }
}

