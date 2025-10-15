package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
	"strings"
)

const (
	Display_Time = true
)

func main() {
	host := "138.47.99.228"
	port := "31337"

	conn, err := net.Dial("tcp", host+":"+port)

	if err != nil {
		fmt.Printf("Error Connecting: %v\n", err)
	}

	fmt.Printf("Connected to %s. Waiting for data...\n", host)

	buffer := make([]byte, 1)
	// Make slice to store all times
	times := []float64{}

	//Record first time
	readTime := time.Now()

	for {
		_, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nEnd of File. Connection Closed")
			} else {
				fmt.Printf("\n Error Reading: %v", err)
			}
			break
		}

		// Calculate times
		delay := time.Since(readTime).Seconds()
		times = append(times, delay)
		// Increment Counter

		// Update read time
		readTime = time.Now()
		if Display_Time {
			fmt.Printf("Received: '%s' | Delay since last: %.3f seconds\n", string(buffer), delay)
		}
	}
	fmt.Println("----------------------------------------")

	var total float64
	for _, t := range times {
		if t > 0 {
			total += t
		}
	}

	avg := total / float64(len(times))

	fmt.Printf("The average time is %.4f\n", avg)

	bin := []int{}

	for _, x := range times {
		if x >= avg {
			bin = append(bin, 1)
		} else {
			bin = append(bin, 0)
		}
	}

	bin_str := ""

	for _, x := range bin {
		bin_str += strconv.Itoa(x)
	}

	fmt.Println(bin_str)
	fmt.Println(len(bin_str))

	decode(bin_str)

}

func decode(data string) {

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