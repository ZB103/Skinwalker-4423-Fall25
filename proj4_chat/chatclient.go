package main

import (
	"fmt"
	"github.com/secsy/goftp"
)

const (
	address  = "138.47.99.228" // FTP server IP
	port     = "21"            // FTP port
	username = "anonymous"     // FTP username
	password = ""              // FTP password
	path     = "/"           // directory path on FTP server

)

func main(){
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
	
	fmt.Print(entries)
}