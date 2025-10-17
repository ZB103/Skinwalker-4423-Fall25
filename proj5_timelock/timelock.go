package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
)

//function that takes in user-provided epoch and sys time
//and performs timelock, then returns derived code
func runTimelock(epoch string, sysTime string) string{
	return "1234"
}

func main(){
	//getting the user-provided epoch
	reader := bufio.NewReader(os.Stdin)
	epoch, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("ERR: Issue getting epoch - ", err)
		panic(0)
	}
	
	//get system time
	sysTime := time.Now().String()
	fmt.Println(epoch)
	fmt.Println(sysTime)
	
	//calculate 4-character code
	code := runTimelock(epoch, sysTime)
	fmt.Println("Your code is: ", code)
}