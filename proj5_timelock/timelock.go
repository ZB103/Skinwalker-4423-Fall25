package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
	"strconv"
	"strings"
	"crypto/md5"
	"encoding/hex"
)

//function that takes in user-provided epoch and sys time
//and performs timelock, then returns derived code
func runTimelock(epoch time.Time, sysTime time.Time) string{
	var code int
	
	//get time elapsed since epoch
	timeDifference := sysTime.Sub(epoch)
	
	//encode using md5
	hash := md5.Sum([]byte(timeDifference.String()))
	hashString := hex.EncodeToString(hash[:])
	fmt.Println(hashString)
	
	//extract and concat first two letters of hash LtoR
	
	//extract and concat first two single-digit ints of hash LtoR
	
	//concat two extracted values
	
	//return code as a string
	return strconv.Itoa(code)
}

func main(){
	//getting the user-provided epoch
	reader := bufio.NewReader(os.Stdin)
	epochStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("ERR: Issue getting epoch - ", err)
		panic(0)
	}
	//convert user input to valid Date
	//YYYY MM DD HH mm SS
	epochStr = strings.TrimRight(epochStr, "\n")
	epoch, err := time.Parse("2006 01 02 15 04 05", epochStr)
	if err != nil{
		fmt.Println("ERR: Issue converting input to Time - ", err)
		panic(0)
	}
	
	//get system time
	sysTime := time.Now().Round(time.Second)
	fmt.Println(epoch)
	fmt.Println(sysTime)
	
	//calculate 4-character code
	code := runTimelock(epoch, sysTime)
	fmt.Println("Your code is: ", code)
}