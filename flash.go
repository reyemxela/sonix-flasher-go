package main

import (
	"bytes"
	"fmt"

	"github.com/reyemxela/hid"
)

// Flash does the actual initialize>prepare>flash process.
//
// It is assumed that the firmware and offset have been sanity-checked already.
func Flash(dev *hid.Device, offset int, firmware []byte) error {
	// initialize
	err := HIDSendReport(dev, ToBytes(cmdInit), true)
	if err != nil {
		return fmt.Errorf("initialization failed: %v", err)
	}

	// prepare
	err = HIDSendReport(dev, ToBytes(cmdPrepare, offset, len(firmware)/64), true)
	if err != nil {
		return fmt.Errorf("initialization failed: %v", err)
	}

	// flash
	for i := 0; i < len(firmware); i += 64 {
		fmt.Printf("\rFlashing chunk %d/%d", i+64, len(firmware))
		err = HIDSendReport(dev, firmware[i:i+64], false)
		if err != nil {
			return fmt.Errorf("flashing failed: %v", err)
		}
	}
	fmt.Printf("\nDone!\n")

	// reboot
	err = HIDSendReport(dev, ToBytes(cmdReboot), false)
	if err != nil {
		return fmt.Errorf("rebooting failed: %v", err)
	}

	return nil
}

// HIDSendReport is a wrapper for the hid SendFeatureReport,
// and (opionally) checking the status reported back.
func HIDSendReport(dev *hid.Device, report []byte, checkStatus bool) error {
	if len(report) > 64 {
		return fmt.Errorf("report size must be less than 64 bytes")
	}

	pad := make([]byte, 64-len(report))
	_, err := dev.SendFeatureReport(append([]byte{0x00}, append(report, pad...)...))
	if err != nil {
		return err
	}

	if checkStatus {
		data := make([]byte, 65)
		_, err := dev.GetFeatureReport(data)
		// data, err := HIDGetReport(d)
		if err != nil {
			return err
		}

		recvc := data[1:5]
		recvs := data[5:9]

		if !bytes.Equal(recvc, report[0:4]) {
			return fmt.Errorf("response (0x%x) does not match expected (0x%x)", recvc, report[0:4])
		}

		if es := ToBytes(expectedStatus); !bytes.Equal(recvs, es) {
			return fmt.Errorf("response status (0x%x) does not match expected status (0x%x)", recvs, es)
		}
	}

	return nil
}
