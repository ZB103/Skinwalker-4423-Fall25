package main

import (
	"fmt"
	"os"
)

const sentinel = "\x00\xFF\x00\x00\xFF\x00"

// command format:
// storing a hidden message: ./Steg -s -(b or B) -o<offset> -i<interval> -w"<wrapper.file>" -h"<hidden.file>"
// retrieving a hidden message: ./Steg -r -(b or B) -o<offset> -i<interval> -w"<wrapper.file>" > <new.file>

// BYTE METHOD ----------------------------------------------

// store hidden data in byte mode
func storeByteMode(wrapper []byte, hidden []byte, offset int, interval int) []byte {
	// i := 0
	// for i < len(hidden) {
		// if offset >= len(wrapper) {
			// fmt.Fprintf(os.Stderr, "Error: Wrapper file too small for hidden data.\n")
			// os.Exit(1)
		// }
		// wrapper[offset] = hidden[i]
		// offset += interval
		// i++
	// }
	
	// i = 0
	// for i < len(sentinel) {
		// wrapper[offset] = sentinel[i]
		// offset += interval
		// i += 1
	// }
	
	return wrapper
}

// retrieve hidden data in byte mode
func retrieveByteMode(wrapper []byte, offset int, interval int) []byte {
	hidden := []byte{}
	match := 0

	for offset + (7*interval) < len(wrapper) {
		b := wrapper[offset]

		if b == sentinel[match] {
			match++
			// full sentinel matched, stop reading
			if match == len(sentinel) {
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

// BIT MODE ----------------------------------------------------

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
			wrapper[offset] |= ((b & 0x80) >> 7)
			// shift hidden byte left by 1 for next bit
			hidden[i] <<= 1	//could result in values > 1 byte
			offset += interval
		}
	}
	
	sentinelBytes := []byte(sentinel)
	for i := 0; i < len(sentinelBytes); i++ {
		//clear LSB of wrapper byte
		wrapper[offset] &= 0xFE
		wrapper[offset] |= ((sentinelBytes[i] & 0x80) >> 7)
		sentinelBytes[i] <<= 1
		offset += interval
	}
	//mySentinel = string(sentinelBytes)
	
	return wrapper
}

// retrieve hidden data in bit mode
func retrieveBitMode(wrapper []byte, offset int, interval int) []byte {
	hidden := []byte{}
	match := 0

	for offset + (7*interval) < len(wrapper) {
		var b byte = 0
		for j := 0; j < 8; j++ {
			b <<= 1
			b |= (wrapper[offset] & 0x01)
			offset += interval
		}

		// check for sentinel match
		if b == sentinel[match] {
			match++
			if match == len(sentinel) {
				break
			}
		} else {
			if match > 0 {
				hidden = append(hidden, sentinel[:match]...)
				match = 0
			}
			hidden = append(hidden, b)
		}
	}
	return hidden
}

//main - getting info from user
func main() {
	args := os.Args[1:]
	offset := 0
	interval := 1
	wrapperFile := ""
	hiddenFile := ""
	
	// get offset - -o required
	if args[2][:2] == "-o" {
		fmt.Sscanf(args[2], "-o%d", &offset)
	} else{
		panic(0)
	}
	//get interval
	if args[3][:2] == "-i" {
		fmt.Sscanf(args[3], "-i%d", &interval)
		//get -w, -i present
		wrapperFile = args[4][2:]
		//get -h if present
		if len(args) > 5 && args[5][:2] == "-h" {
			hiddenFile = args[5][2:]
		}	
	//get -w, no -i
	} else if args[3][:2] == "-w" {
		wrapperFile = args[3][2:]
		//get -h if present
		if len(args) > 4 && args[4][:2] == "-h" {
			hiddenFile = args[4][2:]
		}
	}else {
		fmt.Println("Error in reading inputs 2-5")
		panic(0)
	}
	
	// read the wrapper file
	wrapperBytes, err := os.ReadFile(wrapperFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading wrapper file: %v\n", err)
		os.Exit(1)
	}
	var hiddenBytes []byte
	if args[0] == "-s" {
		hiddenBytes, err = os.ReadFile(hiddenFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading hidden file: %v\n", err)
			os.Exit(1)
		}
		hiddenBytes = append(hiddenBytes, []byte(sentinel)...)
		fmt.Println("Wrapper size: ", len(wrapperBytes))
		fmt.Println("Hidden + sentinel size: ", len(hiddenBytes))
	}

	// byte check
	if args[1] == "-B" {
		// store
		if args[0] == "-s" {
			wrapperBytes = storeByteMode(wrapperBytes, hiddenBytes, offset, interval)
			os.Stdout.Write(wrapperBytes)
			// get
		} else if args[0] == "-r" {
			hidden := retrieveByteMode(wrapperBytes, offset, interval)
			os.Stdout.Write(hidden)
		} else {
			panic(0)
		}
	// bit check
	} else if args[1] == "-b" {
		// store
		if args[0] == "-s" {
			wrapperBytes = storeBitMode(wrapperBytes, hiddenBytes, offset, interval)
			os.Stdout.Write(wrapperBytes)
			// get
		} else if args[0] == "-r" {
			hidden := retrieveBitMode(wrapperBytes, offset, interval)
			os.Stdout.Write(hidden)
		} else {
			panic(0)
		}
	} else {
		panic(0)
	}

}
