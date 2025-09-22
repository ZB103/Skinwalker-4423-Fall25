package main

import (
	"fmt"
	"os"
	"bufio"
	"unicode"
)

func encrypt(key string, charArray []rune){
	reader := bufio.NewReader(os.Stdin)
	
	//continuously await input
	for{
		//get input string
		data, _ := reader.ReadString('\n')
		
		//encrypt string
		var encryptedData string
		j := 0	//index of key
		for i := 0; i < len(data); i++{
			//skip non-letter chars
			if !unicode.IsLetter(rune(data[i])) {
				encryptedData = encryptedData + string(data[i])
			}else{
			
				//get shift amnt from key
				var shift int
				for k := 0; k < len(charArray); k++{
					if charArray[k] == rune(key[j]){shift = k}
				}
				
				//increment place in key
				j++
				if j >= len(key){j = 0}
				
				//shift letter
				var oldPos int
				for k := 0; k < len(charArray); k++{
					if charArray[k] == rune(unicode.ToLower(rune(data[i]))){oldPos = k}
				}
				var newPos int = oldPos + shift
				if newPos > 25 {newPos -= 26}
				
				//add shifted letter to encrypted string
				encryptedChar := string(charArray[newPos])
				
				//check for capitalization
				if(unicode.IsUpper(rune(data[i]))){
					encryptedChar = string(unicode.ToUpper(rune(encryptedChar[0])))
				}
				encryptedData = encryptedData + encryptedChar
			}
		}
		
		//print final encrypted string
		fmt.Printf(encryptedData)
	}
}

func decrypt(key string, charArray []rune){
	reader := bufio.NewReader(os.Stdin)
	
	//continuously await input
	for{
		//get input string
		data, _ := reader.ReadString('\n')
		
		//decrypt string
		var decryptedData string
		j := 0	//index of key
		for i := 0; i < len(data); i++{
			//skip non-letter chars
			if !unicode.IsLetter(rune(data[i])) {
				decryptedData = decryptedData + string(data[i])
			}else{
			
				//get shift amnt from key
				var shift int
				for k := 0; k < len(charArray); k++{
					if charArray[k] == rune(key[j]){shift = k}
				}
				
				//increment place in key
				j++
				if j >= len(key){j = 0}
				
				//shift letter
				var oldPos int
				for k := 0; k < len(charArray); k++{
					if charArray[k] == rune(unicode.ToLower(rune(data[i]))){oldPos = k}
				}
				var newPos int = oldPos - shift
				if newPos < 0 {newPos += 26}
				
				//add shifted letter to encrypted string
				decryptedChar := string(charArray[newPos])
				
				//check for capitalization
				if(unicode.IsUpper(rune(data[i]))){
					decryptedChar = string(unicode.ToUpper(rune(decryptedChar[0])))
				}
				decryptedData = decryptedData + decryptedChar
			}
		}
		
		//print final decrypted string
		fmt.Printf(decryptedData)
	}
}

func main() {
	//take in argument of -e or -d
	args := os.Args[1:]
	
	//set array
	var charArray = []rune {'a','b','c','d','e','f','g','h','i','j','k','l','m','n','o','p','q','r','s','t','u','v','w','x','y','z'}
	
	//get the key
	var rawKey string
	if args[1] != ""{
		rawKey = args[1]
	} else {
		panic(0)
	}
	//convert key to lowercase
	var key string
	for i := 0; i < len(rawKey); i++{
		if unicode.IsLetter(rune(rawKey[i])) {
			key = key + string(unicode.ToLower(rune(rawKey[i])))
		}
	}
	
	//figure out whether we are encrypting (-e) or decrypting (-d)
	if args[0] == "-e"{
		encrypt(key, charArray)
	}else if args[0] == "-d"{
		decrypt(key, charArray)
	}else{
		panic(0)
	}
}