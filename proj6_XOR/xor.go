package main

import (
	"fmt"
	"io"
	"os"
)

func main() {

	key_file := "key"
	m, err := io.ReadAll(os.Stdin)
	k, err2 := os.ReadFile(key_file)
	b := len(m)
	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err)
	}

	result := make([]byte, b)

	for i := range m {
		result[i] = (m[i] ^ k[i%len(k)])
	}

	fmt.Println("XOR Result: ", result)
	fmt.Println("As String: ", string(result))

}
