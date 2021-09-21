package main

import (
	"fmt"
	"encoding/base64"
	"os"
)

func main() {
	// Input parameters
	var inputData = string(os.Args[1])
    
    decode, err := base64.StdEncoding.DecodeString(inputData)
    if err != nil {
        panic("malformed input")
    }
    fmt.Printf("          ")
    fmt.Println(string(decode))
}
