# sonix-flasher-go

This is a firmware flashing utility for Sonix SN32 chips, built for use with the [SonixQMK](https://github.com/SonixQMK/qmk_firmware) keyboard firmware.

## Installing

Just download the correct binary for your OS from Releases and run it.

## Usage
```
./sonix-flasher-go -f <firmware.bin> [-o 0x200] [-v 0c45] [-p 7010] [-no-confirm]

  -f string
        Firmware file to flash
  -no-confirm
        Doesn't ask before flashing
  -o string
        Offset position to start flash (in hex format)
  -p string
        USB device PID (in hex format) (default "0")
  -v string
        USB device VID (in hex format) (default "0")
```

