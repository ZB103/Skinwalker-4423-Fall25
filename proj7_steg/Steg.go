package main

import (
	"fmt"
	"os"
)

const sentinel = "\x00\xFF\x00\x00\xFF\x00"

// I've already tried retrieving what anky gave us but its not working

// command format:
// storing a hidden message: ./Steg -s -(b or B) -o<offset> -i<interval> -w"<wrapper.file>" -h"<hidden.file>"
// retrieving a hidden message: ./Steg -r -(b or B) -o<offset> -i<interval> -w"<wrapper.file>" > <new.file>

// BYTE METHOD

// store hidden data in byte mode
func storeByteMode(wrapper []byte, hidden []byte, offset int, interval int) []byte {
	i := 0
	for i < len(hidden) {
		if offset >= len(wrapper) {
			fmt.Fprintf(os.Stderr, "Error: Wrapper file too small for hidden data.\n")
			os.Exit(1)
		}
		wrapper[offset] = hidden[i]
		offset += interval
		i++
	}
	return wrapper
}

// retrieve hidden data in byte mode
func retrieveByteMode(wrapper []byte, offset int, interval int) []byte {
	hidden := []byte{}
	sentinelLen := len(sentinel)
	match := 0

	for offset < len(wrapper) {
		b := wrapper[offset]

		if b == sentinel[match] {
			match++
			// full sentinel matched, stop reading
			if match == sentinelLen {
				break
			}
		} else {
			// if partial match failed, append previously matched bytes
			if match > 0 {
				hidden = append(hidden, sentinel[:match]...)
				match = 0
			}
			hidden = append(hidden, b)
		}
		offset += interval
	}
	return hidden
}

// BIT MODE

// store hidden data in bit mode
func storeBitMode(wrapper []byte, hidden []byte, offset int, interval int) []byte {
	for i := 0; i < len(hidden); i++ {
		b := hidden[i]
		for j := 0; j < 8; j++ {
			if offset >= len(wrapper) {
				fmt.Fprintf(os.Stderr, "Error: Wrapper file too small for hidden data.\n")
				os.Exit(1)
			}
			// clear LSB of wrapper byte
			wrapper[offset] &= 0xFE
			// take MSB of hidden byte and store in LSB of wrapper
			wrapper[offset] |= (b & 0x80) >> 7
			// shift hidden byte left by 1 for next bit
			b <<= 1
			offset += interval
		}
	}
	return wrapper
}

// retrieve hidden data in bit mode
func retrieveBitMode(wrapper []byte, offset int, interval int) []byte {
	hidden := []byte{}
	const SENTINEL = "\x00\xFF\x00\x00\xFF\x00"
	sentinelLen := len(SENTINEL)
	match := 0

	for offset < len(wrapper) {
		var b byte = 0
		for j := 0; j < 8; j++ {
			b <<= 1
			b |= wrapper[offset] & 0x01
			offset += interval
		}

		// check for sentinel match
		if b == SENTINEL[match] {
			match++
			if match == sentinelLen {
				break
			}
		} else {
			if match > 0 {
				hidden = append(hidden, SENTINEL[:match]...)
				match = 0
			}
			hidden = append(hidden, b)
		}
	}
	return hidden
}

func main() {
	args := os.Args[1:]

	mode := ""
	method := ""
	offset := 0
	interval := 1
	wrapperFile := ""
	hiddenFile := ""

	for _, arg := range args {
		switch {
		case arg == "-s":
			mode = "s"
		case arg == "-r":
			mode = "r"
		case arg == "-b":
			method = "b"
		case arg == "-B":
			method = "B"
		case len(arg) > 2 && arg[:2] == "-o":
			fmt.Sscanf(arg, "-o%d", &offset)
		case len(arg) > 2 && arg[:2] == "-i":
			fmt.Sscanf(arg, "-i%d", &interval)
		case len(arg) > 2 && arg[:2] == "-w":
			wrapperFile = arg[2:]
		case len(arg) > 2 && arg[:2] == "-h":
			hiddenFile = arg[2:]
		}
	}

	// read the wrapper file
	wrapperBytes, err := os.ReadFile(wrapperFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading wrapper file: %v\n", err)
		os.Exit(1)
	}
	var hiddenBytes []byte
	if mode == "s" {
		hiddenBytes, err = os.ReadFile(hiddenFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading hidden file: %v\n", err)
			os.Exit(1)
		}
	}

	if mode == "s" {
		hiddenBytes = append(hiddenBytes, []byte(sentinel)...)
	}

	fmt.Println("Wrapper size: ", len(wrapperBytes))
	if mode == "s" {
		fmt.Println("Hidden + sentinel size: ", len(hiddenBytes))
	}

	// byte check
	if method == "B" {
		// store
		if mode == "s" {
			wrapperBytes = storeByteMode(wrapperBytes, hiddenBytes, offset, interval)
			os.Stdout.Write(wrapperBytes)
			// get
		} else {
			hidden := retrieveByteMode(wrapperBytes, offset, interval)
			os.Stdout.Write(hidden)
		}
	}

	// bit check
	if method == "b" {
		// store
		if mode == "s" {
			wrapperBytes = storeBitMode(wrapperBytes, hiddenBytes, offset, interval)
			os.Stdout.Write(wrapperBytes)
			// get
		} else {
			hidden := retrieveBitMode(wrapperBytes, offset, interval)
			os.Stdout.Write(hidden)
		}
	}

}
