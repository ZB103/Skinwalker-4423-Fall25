package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"unicode"
)

//function that takes in user-provided epoch and sys time
//and performs timelock, then returns derived code
func runTimelock(epoch time.Time, sysTime time.Time) string{
	//get time elapsed since epoch
	timeDifference := sysTime.Sub(epoch)
	//deal with potential daylight savings time
	if sysTime.IsDST() {
		fmt.Println("sysTime is in DST")
	}
	
	//double encode using md5
	hash := md5.Sum([]byte(timeDifference.String()))
	hashString := hex.EncodeToString(hash[:])
	dblHash := md5.Sum([]byte(hashString))
	dblHashString := hex.EncodeToString(dblHash[:])
	
	//extract and concat first two letters of hash LtoR
	var firstLetters string
	for i := range len(dblHashString){
		//is it a letter?
		if unicode.IsLetter(rune(dblHashString[i])){
			firstLetters = firstLetters + string(dblHashString[i])
		}
		//did we get two letters?
		if len(firstLetters) == 2{
			break
		}
	}
	
	//extract and concat first two single-digit ints of hash LtoR
	var firstNumbers string
	for i := len(dblHashString)-1; i >= 0; i--{
		//is it a letter?
		if unicode.IsNumber(rune(dblHashString[i])){
			firstNumbers = firstNumbers + string(dblHashString[i])
		}
		//did we get two letters?
		if len(firstNumbers) == 2{
			break
		}
	}
	
	//concat two extracted values
	code := firstLetters + firstNumbers
	
	//return code as a string
	return code
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
	sysTime := time.Now().Round(time.Second).UTC()
	fmt.Println(epoch)
	fmt.Println(sysTime)
	
	//calculate 4-character code
	code := runTimelock(epoch, sysTime)
	fmt.Println(code)
}