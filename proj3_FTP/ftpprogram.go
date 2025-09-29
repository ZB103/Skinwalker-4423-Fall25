/*
* author:    anky (with help from ChatGPT)
* date:      09/23/2025
* description:
*   Decodes hidden messages from FTP file permissions.
*   METHOD selects the scheme
 */

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/secsy/goftp"
)

// ====== CONFIG ======
const (
	METHOD   = 10              // choose decoding method: 7 or 10
	address  = "138.47.99.228" // FTP server IP
	port     = "21"            // FTP port
	username = "anonymous"     // FTP username
	password = ""              // FTP password
	path     = "/10"           // directory path on FTP server

)

// Convert permission characters into bits.
// '-'/'S'/'T' -> '0'
// Any other character -> '1'
func permCharsToBits(seg string) string {
	var b strings.Builder
	b.Grow(len(seg))
	for _, c := range seg {
		switch c {
		case '-', 'S', 'T':
			b.WriteByte('0')
		default:
			b.WriteByte('1')
		}
	}
	return b.String()
}

// Build the bit string for a given entry depending on METHOD.
func entryBits(perm string) string {
	if METHOD == 10 {
		// METHOD 10 -> take ALL 10 permission characters
		if len(perm) < 10 {
			return ""
		}
		return permCharsToBits(perm[:10])
	}
	// METHOD 7 -> take the last 7 characters only
	if len(perm) < 7 {
		return ""
	}
	return permCharsToBits(perm[len(perm)-7:])
}

// convert a 7 bit binary string into a rune.
// only returns true if the rune is printable.
func bin7ToRune(bin7 string) (rune, bool) {
	if len(bin7) != 7 {
		return 0, false
	}
	v, err := strconv.ParseInt(bin7, 2, 64)
	if err != nil || v < 0 || v > 127 {
		return 0, false
	}
	r := rune(v)
	// only accept printable runes and common whitespace
	if !unicode.IsPrint(r) && r != '\n' && r != '\t' && r != ' ' {
		return 0, false
	}
	return r, true
}

// split a bitstream into groups of 7 bits,
// convert each group into ASCII, and join them into a string.
func decodeBitstreamToASCII(bitstream string) string {
	var out strings.Builder
	for i := 0; i+7 <= len(bitstream); i += 7 {
		chunk := bitstream[i : i+7]         // take 7 bits
		if r, ok := bin7ToRune(chunk); ok { // try convert to rune
			out.WriteRune(r) // append if valid
		}
	}
	return out.String()
}

func main() {
	// create FTP client configuration
	config := goftp.Config{
		User:            username,
		Password:        password,
		ActiveTransfers: true,
		Timeout:         20 * time.Second,
	}

	// connect to the FTP server
	client, err := goftp.DialConfig(config, address+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close() // make sure to close connection later

	// get directory listing from the server
	entries, err := client.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	// builder to collect the full hidden bitstream
	var bitstream strings.Builder

	// process each entry in the directory
	for _, e := range entries {
		perm := e.Mode().String() // permission string like "-rw-r--r--"

		// METHOD 7 -> only use entries where the first 3 characters are "---"
		if METHOD == 7 && len(perm) >= 3 && perm[:3] != "---" {
			continue
		}

		// METHOD 10 -> include ALL entries
		seg := entryBits(perm) // get this entry’s bit string
		if seg == "" {
			continue
		}
		bitstream.WriteString(seg) // append bits to global stream
	}

	// Get the final bitstream
	bits := bitstream.String()

	// Decode the bitstream into ASCII
	decoded := decodeBitstreamToASCII(bits)

	// Print results
	fmt.Println(decoded)
}
