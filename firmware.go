package main

import "os"

// ReadFirmware opens the specified file, pads it if needed,
// and returns the firmware as a byte array.
func ReadFirmware(f string) ([]byte, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	s := stat.Size()
	b := make([]byte, s)
	_, err = file.Read(b)
	if err != nil {
		return nil, err
	}

	// pad jumploader, if this is one
	if len(b) < 512 {
		b = append(b, make([]byte, 512-len(b))...)
	}

	if len(b)%64 != 0 {
		ErrorExit("ERROR: firmware size must be divisible by 64")
	}

	return b, nil
}
