package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/reyemxela/hid"
)

const (
	version        = "v1.0.1"
	cmdBase        = 0x55AA00
	cmdInit        = cmdBase + 1
	cmdPrepare     = cmdBase + 5
	cmdReboot      = cmdBase + 7
	expectedStatus = 0xFAFAFAFA
)

func main() {
	fmt.Printf("sonix-flasher-go - %s\n\n", version)

	var (
		filename  = flag.String("f", "", "Firmware file to flash")
		offsetStr = flag.String("o", "", "Offset position to start flash (in hex format)")
		vidStr    = flag.String("v", "0", "USB device VID (in hex format)")
		pidStr    = flag.String("p", "0", "USB device PID (in hex format)")
		noConfirm = flag.Bool("no-confirm", false, "Doesn't ask before flashing")
	)

	flag.Parse()

	// at least filename is mandatory, so check that
	if *filename == "" {
		flag.PrintDefaults()
		return
	}

	// enumerate devices
	vid := HexInt(*vidStr)
	pid := HexInt(*pidStr)
	devices := []hid.DeviceInfo{}
	for _, d := range hid.Enumerate(uint16(vid), uint16(pid)) {
		if d.Interface <= 0 {
			devices = append(devices, d)
			fmt.Printf("found device %04x:%04x - %s\n\n", d.VendorID, d.ProductID, DeviceNames[int(d.VendorID)][int(d.ProductID)])
		}
	}

	if len(devices) < 1 {
		ErrorExit("ERROR: no devices")
	} else if len(devices) > 1 {
		ErrorExit("ERROR: too many devices found. please specify vid and pid")
	}

	dev, err := devices[0].Open()
	if err != nil {
		ErrorExit("ERROR: opening device: %v", err.Error())
	}
	vid = int(dev.VendorID)
	pid = int(dev.ProductID)
	conf := BootloaderConfigs[vid][pid]

	// if we don't have a device entry, set a default
	if conf.maxFirmwareSize == 0 {
		conf.maxFirmwareSize = F240
	}

	// get/check offset
	offset := 0
	if *offsetStr == "" {
		offset = conf.offset
		fmt.Printf("offset not specified, using default for device: 0x%x\n", offset)
	} else {
		offset = HexInt(*offsetStr)
	}

	if offset%64 != 0 {
		ErrorExit("ERROR: offset must be divisible by 64")
	}

	// read firmware file
	firmware, err := ReadFirmware(*filename)
	if err != nil {
		ErrorExit("ERROR: reading firmware: %v", err.Error())
	}
	if len(firmware)+offset > conf.maxFirmwareSize {
		ErrorExit("ERROR: firmware + offset too large\n"+
			"input:    %v bytes\n"+
			"chip max: %v bytes", len(firmware)+offset, conf.maxFirmwareSize)
	}

	fmt.Printf("------------------------------------------\n"+
		"file:    %v\n"+
		"size:    %v bytes\n"+
		"offset:  0x%x\n"+
		"vid:     %04x\n"+
		"pid:     %04x\n"+
		"------------------------------------------\n",
		filepath.Base(*filename), len(firmware), offset, vid, pid)

	// if noConfirm is false, ask for confirmation before flashing
	if !*noConfirm {
		fmt.Printf("Proceed with flash? [y/N] ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if strings.ToUpper(scanner.Text()) != "Y" {
			fmt.Printf("cancelled\n\n")
			return
		}
	}
	Flash(dev, offset, firmware)
}
