package main

import (
	"fmt"
	"io"
	"os"
)

// Test File doing random Stuff

func main() {

	data, err := io.ReadAll(os.Stdin)

	if err != nil {
		panic(err)
	}
	length := len(data)

	fmt.Printf("The length of the cipherText is %d\n", length)

	content, err := os.ReadFile("key")

	if err != nil {
		panic(err)
	}

	length2 := len(content)
	fmt.Printf("The length of the key is %d", length2)
}
