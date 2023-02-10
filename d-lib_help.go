package main

var help_str = `
Usage: d-* [ command ] [ -flags <options>, ... ]

Where: d-* = d-lin (Linux); d-win (Win); d-mac (Mac Intel); d-mm1 (Mac M1)

Commands & Flags:
  <command> -h                              help with command flags
  help                                      extended help
  ports <port>                              list ports / set port
  slots <directory>                         list slots w/ DLP file names
  view  -k                                  view knobs
  view  -s=<slot> <-pro>                    view a slot
  view  -f=<file> <-pro>                    view a DLP file
  ktof  -f=<file> <-pro>                    download knobs to DLP file
  ftok  -f=<file> <-pro>                    upload DLP file to knobs
  stof  -f=<file> -s=<slot> <-pro>          download slot to DLP file
  ftos  -f=<file> -s=<slot> <-pro>          upload DLP file to slot
  btos  -f=<file> -s=<slot> <-pro>          upload BNK DLP files
  stob  -f=<file>                           slot DLP names to BNK file
  dump  -f=<file> <-eeprom|-spi|-pre|-pro>  bulk download to file
  pump  -f=<file> <-eeprom|-spi|-pre|-pro>  bulk upload from file
  hcl   <ver|crc|acal|0 9 rr|22 rk|...>     execute HCL command**
`

var help_verbose_str = `
Notes:
- Flags may be entered in any order.
- Flag prefix either "-" or "--" (e.g. -s=5; --s=5).
- Flags / values separator either "=" or space (e.g. -s=5; -s 5).
- If not provided, the file extension is added automatically.
- If provided, an incorrect file extension flags an error.
- A user prompt precedes any file overwrite.
- The command "pump -pre" now skips over writing to profile slots.
- The "btos" and "stob" commands uses the -f path to locate all files.
- If "btos" -s=n is negative, then the next slot written is n--, else n++.
- The "btos" command skips over lines in *.bnk files that begin with "//".
- Individual preset and profile files share the same *.dlp file extension.
- The serial port number is stored in the config file "d-lev.cfg".
- The "ports" command updates the config file if a port number is given.
- If missing, a config file will be automatically generated.
- If "view" output doesn't fit in the command window, change the font/layout.
- Linux requires an executable file be prefaced with: "./"
- Linux users may need to join the "dialout" group for serial port access.
- ** See files HCL.txt, REGS.txt, and KNOBS.txt for details.

Usage Examples: (e.g. Windows build)
- Show version & short help:
    d-win
- Show version & extended help:
    d-win help
    d-win -help
    d-win -h
- List all available serial ports:
    d-win ports
- List ports and set port to 5:
    d-win ports 5
- List slots compared to DLP files in "_ALL_" directory:
    d-win slots _ALL_
- View all current knob values:
    d-win view -k
- View preset in slot 20:
    d-win view -s 20
- View profile in slot -2:
    d-win view -pro -s -2
- View preset in file "bassoon.dlp":
    d-win view -f bassoon
- View profile in file "some_pro.dlp":
    d-win view -pro -f somepro
- Download preset knobs to preset file "him_her.dlp":
    d-win ktof -f him_her
- Upload preset file "flute.dlp" to preset knobs:
    d-win ftok -f flute
- Download profile knobs to profile file "my_prof_4.dlp":
    d-win ktof -pro -f my_prof_4
- Upload profile file "some_prof.dlp" to profile knobs:
    d-win ftok -pro -f some_prof
- Download preset slot -5 to preset file "female7.dlp":
    d-win stof -s -5 -f female7
- Upload preset file "cello8.dlp" to preset slot 9:
    d-win ftos -f cello8 -s 9
- Download profile slot 0 to profile file "my_sys.dlp":
    d-win stof -pro -s 0 -f my_sys
- Upload profile file "_sys_9.dlp" to profile slot 3
    d-win ftos -pro -f _sys_9 -s 3
- Upload bank of presets in bank file "mybank.bnk" to preset slots 10, 11, 12, etc.:
    d-win btos -f mybank -s 10
- Upload bank of profiles in bank file "oldprofs.bnk" to profile slots -1, -2, etc.:
    d-win btos -pro -f oldprofs -s -1
- Create bank file "mybnk.bnk" from preset & profile files in current directory:
    d-win stob -f mybnk
- Download software & all presets & profiles to file "2022-01-23.eeprom":
    d-win dump -eeprom -f 2022-01-23
- Upload software & all presets & profiles from file "factory.eeprom":
    d-win pump -eeprom -f factory
- Download software to file "sw_backup.spi":
    d-win dump -spi -f sw_backup
- Upload software from file "f9e1c5c7.spi":
    d-win pump -spi -f f9e1c5c7
- Download all presets to file "old_presets.pre":
    d-win dump -pre -f old_presets
- Upload all presets from file "my_dlev.pre":
    d-win pump -pre -f my_dlev
- Download all profiles to file "my_setup.pro":
    d-win dump -pro -f my_setup
- Upload all profiles from file "your_setup.pro":
    d-win pump -pro -f your_setup
- Read the software version:
    d-win hcl ver
- Calculate the EEPROM software CRC (s/b debb20e3):
    d-win hcl crc
- Perform an ACAL:
    d-win hcl acal
- Read processor registers 0 thru 9:
    d-win hcl 0 9 rr
- Read knob 42:
    d-win hcl 42 rk
`
