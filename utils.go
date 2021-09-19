package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// HexInt converts strings into hex ints, stripping off any "0x" that might be present.
func HexInt(s string) int {
	// get rid of any "0x" prefixes
	sr := strings.Replace(strings.ToUpper(s), "0X", "", -1)
	i, err := strconv.ParseInt(sr, 16, strconv.IntSize)
	if err != nil {
		ErrorExit("ERROR: unable to parse input: " + s)
	}
	return int(i)
}

// ToBytes converts int(s) to a little-endian byte array.
func ToBytes(i ...int) []byte {
	buf := new(bytes.Buffer)
	for _, b := range i {
		binary.Write(buf, binary.LittleEndian, uint32(b))
	}
	return buf.Bytes()
}

// ErrorExit prints out the specified format string, then exits with code 1.
func ErrorExit(s string, f ...interface{}) {
	fmt.Printf(s+"\n\nExiting\n", f...)
	os.Exit(1)
}
