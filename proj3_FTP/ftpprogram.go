package main // every Go file starts with a package name; main makes it a program

import ( // list of packages we use
	"bufio"   // buffered reading/writing on connections
	"fmt"     // for printing and formatting strings
	"net"     // network connections
	"os"      // command line args and printing
	"strconv" // convert strings to numbers
	"strings" // string functions
)

var (
	FTP_HOST = "138.47.156.109"      // server address
	FTP_PORT = 21                    // ftp port
	FTP_DIR  = "/home/anonymous/ftp" // directory to list
	// METHOD picks how many bits per character to use
	METHOD = 7 // default is 7-bit
)

func connectAndList(host string, port int, dir string) []string { // connect and get permission strings
	addr := fmt.Sprintf("%s:%d", host, port) // make host:port
	conn, _ := net.Dial("tcp", addr)         // open TCP connection
	defer conn.Close()                       // close it when done

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)) // read and write lines

	// Read the welcome line
	rw.ReadString('\n')

	// Send our username: anonymous
	fmt.Fprintf(rw, "USER anonymous\r\n") // write command
	rw.Flush()                            // actually send it
	// Server replies
	line, _ := rw.ReadString('\n') // read reply line
	if strings.HasPrefix(line, "331") {
		// If it asked for a password, send one
		fmt.Fprintf(rw, "PASS anonymous\r\n") // password also "anonymous'
		rw.Flush()                            // send it
		rw.ReadString('\n')                   // read the login success reply
	}

	// Change to the directory that has the files
	fmt.Fprintf(rw, "CWD %s\r\n", dir) // CWD
	rw.Flush()                         // send it
	rw.ReadString('\n')                // read the OK line

	// Enter passive mode so we can open a data connection
	fmt.Fprintf(rw, "PASV\r\n")               // ask for PASV
	rw.Flush()                                // send it
	pasvLine, _ := rw.ReadString('\n')        // read the PASV reply
	start := strings.IndexByte(pasvLine, '(') // find where the numbers start
	end := strings.IndexByte(pasvLine, ')')   // and end
	inside := pasvLine[start+1 : end]         // grab the inside text
	parts := strings.Split(inside, ",")       // split 6 numbers
	ip1 := parts[0]
	ip2 := parts[1]
	ip3 := parts[2]
	ip4 := parts[3]
	hostIP := strings.Join([]string{ip1, ip2, ip3, ip4}, ".") // rebuild IP
	// parse the last two numbers
	p1 := stringToInt(strings.TrimSpace(parts[4])) // high byte of port
	p2 := stringToInt(strings.TrimSpace(parts[5])) // low byte of port
	dataPort := p1*256 + p2                        // combine to full port

	dataConn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", hostIP, dataPort)) // open data connection
	defer dataConn.Close()                                                 // close it later

	// Ask the server to LIST the directory
	fmt.Fprintf(rw, "LIST\r\n") // send LIST
	rw.Flush()                  // send it
	rw.ReadString('\n')         // read the "150 opening" line

	dataReader := bufio.NewReader(dataConn) // read from data connection
	var lines []string                      // will hold each listing line
	for {                                   // read until it stops
		line, err := dataReader.ReadString('\n') // read one line
		if err != nil {                          // if no more, stop
			break
		}
		lines = append(lines, strings.TrimRight(line, "\r\n")) // keep line 
	}
	dataConn.Close()    // close data connection
	rw.ReadString('\n') // read the final "226 transfer complete" line

	items := make([]string, 0, len(lines)) // final permissions
	for _, ln := range lines {             // look at each LIST line
		f := strings.Fields(ln) // split by spaces
		if len(f) < 9 {         //
			continue
		}
		perm := f[0]                           // first thing is the permission string
		if len(perm) >= 10 && perm[0] != 'd' { // only files, not directories
			items = append(items, perm[:10]) // keep first 10 chars like -rw-r--r--
		}
	}
	return items // send back permissions only
}

// Turn permission strings into a message using the chosen method
// METHOD 7: use only the right-most 7 bits
// METHOD 10: use all 10 bits (type + 9 perms)
func decodePermissions(perms []string, method int) string { // translate perms to text
	var bits []byte           // where we collect '0' and '1'
	for _, p := range perms { // go through each permission string
		if len(p) < 10 { // must look like -rwxr-xr--
			continue
		}

		if method == 7 { // 7-bit path
			permBits := p[1:10]                              // drop the type, keep 9 perm chars
			if strings.ContainsAny(permBits[0:3], "rwxst") { // if user bits set, skip
				continue
			}
			bitVector := permissionToBits(p) // make a 10-bit vector
			if len(bitVector) != 10 {        // maek suyre
				continue
			}
			last7 := bitVector[3:10]      // take the last 7 bits
			bits = append(bits, last7...) // add them to our list
			continue
		}

		// 10-bit path: take all 10 bits as-is
		bitVector := permissionToBits(p) // make 10 bits
		if len(bitVector) == 10 {        // if good, append
			bits = append(bits, bitVector...)
		}
	}

	// Now turn the bits into characters
	msg := ""
	for i := 0; i+7 < len(bits); i += 7 { // 7-bit ASCII groups
		if method == 7 {
			val := bitsToInt(bits[i : i+7]) // convert 7 bits to number
			msg = msg + string(byte(val))   // append to the string
		}
	}
	if method == 10 { // 10-bit uses 8-bit ASCII for output
		for i := 0; i+7 < len(bits); i += 8 {
			val := bitsToInt(bits[i : i+8]) // convert 8 bits to number
			msg = msg + string(byte(val))   // append character
		}
	}
	finalMsg := msg
	return finalMsg // final message
}

// permissionToBits converts a 10-char permission string to 10 bits
func permissionToBits(p string) []byte { // make '0'/'1' bytes from a permission string
	if len(p) < 10 { // must be at least 10 chars
		return nil
	}
	result := make([]byte, 0, 10) // will hold 10 bits
	if p[0] == '-' {
		result = append(result, '0')
	} else {
		result = append(result, '1')
	}
	// next nine bits from rwx letters: if a letter is present, put 1 else 0
	for i := 1; i < 10; i++ {
		c := p[i]
		if c == 'r' || c == 'w' || c == 'x' || c == 's' || c == 't' {
			result = append(result, '1')
		} else {
			result = append(result, '0')
		}
	}
	return result // ten bits total
}

func bitsToInt(bits []byte) int { // turn a slice of '0'/'1' into a number
	v := 0                   // start at zero
	for _, b := range bits { // for each bit
		v <<= 1       // shift left by one
		if b == '1' { // if the bit is one
			v |= 1 // set the lowest bit
		}
	}
	return v // final value
}

func main() { //
	// We can choose the method from the command line: "7" or "10"
	method := METHOD      // start with default
	if len(os.Args) > 1 { // if user passed an argument
		if os.Args[1] == "10" { // choose 10-bit
			method = 10
		} else if os.Args[1] == "7" { // choose 7-bit
			method = 7
		}
	}

	perms := connectAndList(FTP_HOST, FTP_PORT, FTP_DIR) // get permission strings
	message := decodePermissions(perms, method)          // turn them into text
	text := message                                      // extra temp variable
	fmt.Print(text)                                      // print message only
}

// tiny helper that turns a string into an int 
func stringToInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
