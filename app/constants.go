package main

// Epochs
const MessageConfidence = uint64(1)

// BuildVersion is the local build version, set by build system
const BuildVersion = "1.1.1"

func UserVersion() string {
	return BuildVersion
}
