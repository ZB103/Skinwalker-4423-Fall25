package main

import (
	"fmt"
	"os"
	"bufio"
)

func main(){
	//getting the user-provided epoch
	reader := bufio.NewReader(os.Stdin)
	epoch, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("ERR: Issue getting epoch - ", err)
		panic(0)
	}
	
	fmt.Print(epoch)
}