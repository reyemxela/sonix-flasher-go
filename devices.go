package main

// DeviceNames maps vid:pid to device name strings
var DeviceNames = map[int]map[int]string{
	0x0c45: {
		0x7010: "SN32F268F (bootloader)",
		0x7040: "SN32F248B (bootloader)",
		0x7900: "SN32F248 (bootloader)",

		0x652f: "Glorious GMMK / Tecware Phantom",
		0x766b: "Kemove",
		0x7698: "Womier",
		0x5004: "Redragon",
		0x5104: "Redragon",
		0x8513: "Sharkoon",
		0x8508: "SPCGear",
		0x7903: "Ajazz",
	},

	0x05ac: {
		0x024f: "Keychron / Flashquark Horizon Z",
	},

	0x320F: {
		0x5013: "Akko",
	},
}

// BootloaderConfig holds max firmware sizes and default offsets for a device
type BootloaderConfig struct {
	maxFirmwareSize int
	offset          int
}

const (
	F260 = 30 * 1024 // 30K
	F240 = 64 * 1024 // 64K
)

// BootloaderConfigs maps vid:pid to a config with
// max firmware sizes and default offsets
var BootloaderConfigs = map[int]map[int]BootloaderConfig{
	0x0c45: {
		0x7010: { // SN32F268F
			maxFirmwareSize: F260,
			offset:          0x200,
		},
		0x7040: { // SN32F248B
			maxFirmwareSize: F240,
			offset:          0x00,
		},
		0x7900: { // SN32F248
			maxFirmwareSize: F240,
			offset:          0x00,
		},
	},
}
