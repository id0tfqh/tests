package main

import (
	"fmt"
	"encoding/base64"
	"os"
)

func main() {
	// Input parameters
	var inputData = string(os.Args[1])
    
    encode := base64.StdEncoding.EncodeToString([]byte(inputData))
    fmt.Printf("          ")
    fmt.Printf("%v\n", encode)

}
