// Test to let Go ping my virutal machine

package main

import (
	//"fmt"
	//"bufio"
	"fmt"
	"os/exec"
)

func main() {
	ip := "192.168.1.187"
	cmd := exec.Command("ping", ip)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println(string(output))
}
