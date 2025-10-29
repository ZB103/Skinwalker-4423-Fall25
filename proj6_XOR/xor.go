package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// Take in plain text
	m, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		os.Exit(1)
	}

	// Read from key file
	k, err2 := os.ReadFile("key")
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Error reading key file: %v\n", err2)
		os.Exit(1)
	}

	// Make lengths of key file and plain text
	keyLen := len(k)
	msgLen := len(m)

	// If key is empty throw error
	if keyLen == 0 {
		fmt.Fprintf(os.Stderr, "Error: Key file is empty.\n")
		os.Exit(1)
	}

	// Create slice that will hold final result
	result := make([]byte, msgLen)

	// Loop that repeatidly xor's bit from key to message
	for i := 0; i < msgLen; i++ {
		result[i] = m[i] ^ k[i%keyLen]
	}

	// Write out the result to stdout
	_, err = os.Stdout.Write(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to stdout: %v\n", err)
		os.Exit(1)
	}
}
