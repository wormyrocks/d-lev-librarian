package main

var help_str = `
Usage: d-* [ command ] [ -flag <option> -flag <option> ... ]
Where: d-* = d-win (Windows); d-mac (Mac Intel); d-mm1 (Mac M1);
             d-lin (Linux Intel); d-arm (Linux ARM64); d-a32 (Linux ARM32)

Commands & Flags:
  update                                               interactive menu to do update
  <command> -h                                         help with individual command flags
  help                                                 extended librarian help with examples
  ports <-p port>                                      list ports / set port
  view  <-k|-s slot|-f file> <-pro>                    view knobs|slot|DLP file
  match <-d dir> <-s|-d2 dir> <-g> <-pro>              match DLP files with slots|DLP files
  diff  <-f file> <-k|-s slot|-f2 file> <-pro>         compare DLP file to knobs|slot|DLP file2
  diff  <-s slot> <-k|-s2 slot> <-pro>                 compare slot to knobs|slot2
  ktof  <-f file> <-pro> <-y>                          download knobs to DLP file
  ftok  <-f file> <-pro>                               upload DLP file to knobs
  stof  <-f file> <-s slot> <-pro> <-y>                download slot to DLP file
  ftos  <-f file> <-s slot> <-pro>                     upload DLP file to slot
  btos  <-f file> <-s slot> <-pro>                     upload BNK DLP files
  dump  <-f file.ext>                                  bulk download to PRE|PRO|EEPROM file
  pump  <-f file.ext>                                  bulk upload from PRE|PRO|EEPROM file
  split <-f file.ext> <-y>                             split PRE|PRO|EEPROM container file into sub files
  join  <-f file.ext> <-y>                             join DLP|PRE|PRO|SPI sub files into container file
  morph <-k|-s slot|-f file> <-mo|n|e|f|r|> <-i seed>  morph knobs|slot|DLP file to knobs
  batch <-d dir> <-d2 dir> <-m|u|r> <pro> <-y>         batch convert DLP files
  knob  <-k page:knob> <-o offset|-v val>              read/set/offset knob value
  hcl   <ver|crc|acal|22 rk|...>                       issue HCL command**
  loop  <"some text to loop back">                     serial port loop back test***
  ver                                                  get software version
  acal                                                 issue acal
  reset                                                issue processor reset
`

var help_verbose_str = `
Notes:
- Flags may be entered in any order.
- Flag prefix either "-" or "--" (e.g. -s=5; --s=5).
- Flags / values separator either space or "=" (e.g. -s 5; -s=5).
- If not provided, the *.dlp file extension is added automatically.
- If provided, an incorrect file extension flags an error.
- If the specified target directory doesn't exist it will be created.
- The dump, pump, split, and join commands require a file extension to know what to do.
- A <y|n> user prompt precedes most file overwrites (-y flag overrides prompt).
- The "knob" command page name matches first chars, and is case agnostic.
- The "knob" command knob number [0:6]: 0 @ upper left, 1 @ upper right, etc.
- The "btos" command uses the file path to locate all files.
- The "btos" command skips over lines in *.bnk files that begin with "//".
- The "batch" command options are mono, stereo, update, and rob.
- Preset and profile files share the same *.dlp file extension.
- The serial port number is stored in the config file "d-lib.cfg".
- The "ports" command updates the config file if a port number is given.
- If missing, the config file will be automatically generated.
- If "view" output doesn't fit in the window, resize it or change the font/layout.
- Linux & Mac require executable files to be prefaced with: "./" e.g. "./d-mac".
- Windows powershell requires executable files to be prefaced with: ".\" e.g. ".\d-win".
- Linux users may need to join the "dialout" group for serial port access.
- If the librarian hangs, CTRL-C will usually kill a terminal program.
- ** See files HCL.txt, REGS.txt, and KNOBS.txt for details.
- *** Requires USB dongle RX and TX wires to be connected together.

Usage Examples: (e.g. Windows build)
- Interactive update menu:
    d-win update
- Show librarian version & compact help:
    d-win
- Show librarian version & extended help:
    d-win -h
- List all serial ports & current port:
    d-win ports
- List all serial ports & set port to 5:
    d-win ports -p 5
- View all current knob values:
    d-win view -k
- View preset in slot 20:
    d-win view -s 20
- View profile in slot 2:
    d-win view -s 2 -pro
- View preset file "bassoon.dlp":
    d-win view -f bassoon
- View profile file "some_pro.dlp":
    d-win view -f some_pro -pro 
- Match slots with DLP files in "_ALL_" directory:
    d-win match -s -d _ALL_
- Match slots with DLP files in "_ALL_" directory with best guess:
    d-win match -s -d _ALL_ -g
- Match DLP files in "_OLD_" directory with DLP files in "_ALL_" directory:
    d-win match -d2 _OLD_ -d _ALL_
- Compare current knob values to preset file "mimi.dlp":
    d-win diff -f mimi -k
- Compare preset in slot 7 to file "saw.dlp":
    d-win diff -f saw -s 7
- Compare preset file "trixie.dlp" to file "patsy.dlp":
    d-win diff -f patsy -f2 trixie
- Compare profile file "_sys_3.dlp" to file "_sys_0.dlp":
    d-win diff -f _sys_0 -f2 _sys_3 -pro
- Compare preset in slot 20 to preset in slot 45:
    d-win diff -s 45 -s2 20
- Compare current knob values to preset in slot 3:
    d-win diff -s 3 -k
- Compare profile in slot 3 to profile in slot 0:
    d-win diff -s 0 -s2 3 -pro
- Download preset knobs to preset file "him_her.dlp":
    d-win ktof -f him_her
- Upload preset file "flute.dlp" to preset knobs:
    d-win ftok -f flute
- Download profile knobs to profile file "my_prof_4.dlp":
    d-win ktof -f my_prof_4 -pro
- Upload profile file "some_prof.dlp" to profile knobs:
    d-win ftok -f some_prof -pro
- Download preset slot 5 to preset file "female7.dlp":
    d-win stof -s 5 -f female7
- Upload preset file "cello8.dlp" to preset slot 9:
    d-win ftos -f cello8 -s 9
- Download profile slot 0 to profile file "my_sys.dlp":
    d-win stof -s 0 -f my_sys -pro
- Upload profile file "_sys_9.dlp" to profile slot 3
    d-win ftos -f _sys_9 -s 3 -pro
- Upload bank of presets in bank file "mybank.bnk" to preset slots 10, 11, 12, etc.:
    d-win btos -f mybank -s 10
- Download software & all presets & profiles to file "2022-01-23.eeprom":
    d-win dump -f 2022-01-23.eeprom
- Upload software & all presets & profiles from file "factory.eeprom":
    d-win pump -f factory.eeprom
- Download software to file "sw_backup.spi":
    d-win dump -f sw_backup.spi
- Upload software from file "f9e1c5c7.spi":
    d-win pump -f f9e1c5c7.spi
- Download all presets to file "old_presets.pre":
    d-win dump -f old_presets.pre
- Upload all presets from file "my_dlev.pre":
    d-win pump -f my_dlev.pre
- Download all profiles to file "my_setup.pro":
    d-win dump -f my_setup.pro
- Upload all profiles from file "your_setup.pro":
    d-win pump -f your_setup.pro
- Split file "some.eeprom" into "some.pre", "some.pro", some.spi":
    d-win split -f some.eeprom
- Split file "my_setup.pro" into "pro_000.dlp" thru "pro_005.dlp":
    d-win split -f my_setup.pro
- Split file "my_new.pre" into "000.dlp" thru "249.dlp":
    d-win split -f my_new.pre
- Join files "some.pre", "some.pro" and "some.spi" into "some.eeprom":
    d-win join -f some.eeprom
- Join files "pro_000.dlp" thru "pro_005.dlp" to "stuff.pro":
    d-win join -f stuff.pro
- Join files "000.dlp" thru "249.dlp" to "some.pre":
    d-win join -f some.pre
- Morph knobs (osc):
    d-win morph -mo 12
- Morph slot 23 (filters, resonator, seed):
    d-win morph -s 23 -mf 5 -mr 20 -i 9
- Morph file "cello_8" (osc, filters, resonator):
    d-win morph -f cello_8 -mo 10 -mf 10 -mr 10
- Batch convert all presets in the _ALL_ directory to mono in the _MONO_ directory:
    d-win batch -d _ALL_ -d2 _MONO_ -m
- Batch update all presets in the _ALL_ directory and overwrite them:
    d-win batch -d _ALL_ -d2 _ALL_ -u
- Read knob RESON:mode:
    d-win knob -k re:4
- Set knob RESON:mode to 10:
    d-win knob -k re:4 -v 10
- Offset knob 1_FORM:reso by 4:
    d-win knob -k 1_f:6 -o 4
- Offset knob FLT_OSC:xmix by -2:
    d-win knob -k flt_o:5 -o -2
- Calculate the EEPROM software CRC (s/b debb20e3):
    d-win hcl crc
- Read processor registers 0 thru 9:
    d-win hcl 0 9 rr
- Loop back serial port text "testing 123":
    d-win loop "testing 123"
- Read the software version:
    d-win ver
- Perform an ACAL:
    d-win acal
- Reset the processor:
    d-win reset
`
