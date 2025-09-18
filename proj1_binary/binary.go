package main

import (
	"fmt"
	"bufio"
	"os"
	//"encoding/binary"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		//get string input from file
		binary := scanner.Text()
		
		//decode binary
		for i := 0; i < len(binary); i++ {
			fmt.Printf("")
		}
		
		//output final decoded string
		fmt.Println(binary)
	}
}
