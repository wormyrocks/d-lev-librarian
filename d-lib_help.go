package main

var help_str = `
Usage: d-* [ command ] [ -flags <options>, ... ]
Where: d-* = d-lin (Linux); d-win (Win); d-mac (Mac Intel); d-mm1 (Mac M1)

Commands & Flags:
  <command> -h                        help with command flags
  help                                extended help
  reset                               issue processor reset
  ports -p <port>                     list ports / set port
  slots -d <directory>                list slots w/ DLP file names
  view  -k                            view knobs
  view  -s=<slot> <-pro>              view a slot
  view  -f=<file> <-pro>              view a DLP file
  diff  -f=<file> -k <-pro>           compare DLP file to knobs
  diff  -f=<file> -s=<slot> <-pro>    compare DLP file to slot
  diff  -f=<file> -f2=<file> <-pro>   compare DLP file to file
  ktof  -f=<file> <-pro>              download knobs to DLP file
  ftok  -f=<file> <-pro>              upload DLP file to knobs
  stof  -f=<file> -s=<slot> <-pro>    download slot to DLP file
  ftos  -f=<file> -s=<slot> <-pro>    upload DLP file to slot
  btos  -f=<file> -s=<slot> <-pro>    upload BNK DLP files
  stob  -f=<file> <-hdr>              slot DLP names to BNK file
  dump  -f=<file.ext>                 bulk download to file
  pump  -f=<file.ext>                 bulk upload from file
  split -f=<file.ext>                 split container file into sub files
  hcl   <ver|crc|acal|22 rk|...>      issue HCL command**
`

var help_verbose_str = `
Notes:
- Flags may be entered in any order.
- Flag prefix either "-" or "--" (e.g. -s=5; --s=5).
- Flags / values separator either "=" or space (e.g. -s=5; -s 5).
- If not provided, the file extension is added automatically.
- If provided, an incorrect file extension flags an error.
- The dump, pump, and split commands require a file extension to know what to do.
- A <y|n> user prompt precedes most file overwrites.
- The "btos" and "stob" commands uses the file path to locate all files.
- The "btos" command skips over lines in *.bnk files that begin with "//".
- Individual preset and profile files share the same *.dlp file extension.
- The serial port number is stored in the config file "d-lev.cfg".
- The "ports" command updates the config file if a port number is given.
- If missing, a config file will be automatically generated.
- If "view" output doesn't fit in the window, resize it or change the font/layout.
- Linux & Mac require executable files to be prefaced with: "./" e.g. "./d-mac".
- Linux users may need to join the "dialout" group for serial port access.
- If the librarian hangs, CTRL-C will usually kill a terminal program.
- ** See files HCL.txt, REGS.txt, and KNOBS.txt for details.

Usage Examples: (e.g. Windows build)
- Show version & compact help:
    d-win
- Show version & extended help:
    d-win -h
- Reset the D-Lev processor:
    d-win reset
- List all serial ports & current port:
    d-win ports
- List all serial ports & set port to 5:
    d-win ports -p 5
- List slots compared to DLP files in "_ALL_" directory:
    d-win slots -d _ALL_
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
- Compare current knob values to preset file "mimi.dlp":
    d-win diff -f mimi -k
- Compare preset in slot 7 to file "saw.dlp":
    d-win diff -f saw -s 7
- Compare preset file "trixie.dlp" to file "patsy.dlp":
    d-win diff -f patsy -f2 trixie
- Compare profile file "_sys_3.dlp" to file "_sys_0.dlp":
    d-win diff -f _sys_0 -f2 _sys_3 -pro
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
- Create bank file "mybnk.bnk" from preset & profile files in current directory:
    d-win stob -f mybnk
- Create bank file "0junk.bnk" from preset & profile files in "_ALL_" directory:
    d-win stob -f _ALL_/0junk
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
- Read the software version:
    d-win hcl ver
- Calculate the EEPROM software CRC (s/b debb20e3):
    d-win hcl crc
- Perform an ACAL:
    d-win hcl acal
- Read knob 41 (PITCH:corr):
    d-win hcl 41 rk
- Write 0 to knob 1 (D-LEV:Out):
    d-win hcl 1 0 wk
- Read processor registers 0 thru 9:
    d-win hcl 0 9 rr
`
