package main

import (
	"fmt"
	"io"
	"net"
	"time"
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
	times := []string{}

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
		delay := time.Since(readTime)
		times = append(times, delay.String())
		// Increment Counter

		// Update read time
		readTime = time.Now()
		if Display_Time {
			fmt.Printf("Received: '%s' | Delay since last: %s\n", string(buffer), delay)
		}
	}
	fmt.Println("----------------------------------------")
}

func map(){
	
}
