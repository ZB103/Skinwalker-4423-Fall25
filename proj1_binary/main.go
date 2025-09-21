package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, err := io.ReadAll(os.Stdin)

	if err != nil {
		panic(err)
	}
	// Convert bytes into a string
	s := strings.TrimSpace(string(data))
	// Calls function to see if bits are 7 or 8 bit
	n := bit_type(s)
	// Sperates bits based on bit type
	str_bit := seperate_bits(s, n)
	// Convert byte into ASCII character
	result := binary_convert(str_bit)
	fmt.Println(result)

}

// check if 7 or 8 bit
func bit_type(bit_string string) int {
	// Checks length of string to see if divisible by 7 or 8
	byte_len := len(bit_string)
	if (byte_len % 7) == 0 {
		return 7
	} else {
		return 8
	}

}

// Converting bits into 7 or 8 bit chunks
func seperate_bits(bit_string string, byte_len int) string {
	// Converting the string of bits into a rune (array)
	in := []rune(bit_string)
	out := make([]rune, 0, len(in)+len(in)/byte_len)

	// i is the index, r is the acutal rune
	// Iterating through in and appending to out
	for i, r := range in {
		out = append(out, r)
		if (i+1)%byte_len == 0 && i != len(in)-1 {
			out = append(out, ' ')
		}
	}
	// Turn out into one string.
	joined := string(out)
	return joined
}

func binary_convert(str_bits string) string {
	// Splitting string into substrings
	bits := strings.Split(str_bits, " ")
	var result string

	for _, b := range bits {
		// Coverting binary to base 2 Integer
		val, err := strconv.ParseInt(b, 2, 64)
		if err != nil {
			panic(err)
		}
		// Convert to character
		result += string(rune(val))
	}
	return result
}
