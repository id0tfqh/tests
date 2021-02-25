package main
/*
* ToDo
* $ go build crypto-hash.go
* $ ./ crypto-hash <arg1> <arg2> <arg3>
* $ ./ crypto-hash <type hash> <file from dictionary> <hash>
* Example:
* $ ./ crypto-hash sha512 dictionary kt4tihriheiuthirughi4uhe44qiuhq
 */

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)
import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
//	"golang.org/x/crypto/sha3"
)

func main() {
	if len(os.Args) < 4 {
		get_error("args < 4")
	}

	check_mode(os.Args[1])

	// Передаются три аргумента:
	// os.Args[1] - тип хэш функции
	// os.Args[2] - словарь с паролями
	// os.Args[3] - хэш-функция, которую надо найти
	hack_password(os.Args[1], os.Args[2], os.Args[3])
}

func check_mode(mode string) {
	if mode != "sha256" && mode != "sha512" && mode != "md5" && mode != "sha1" {

		get_error("crypto-hash function is not found")
	}
}

// mode -- <type hash> тип данных -- string;
// dictionary -- <file from dictionary> файл словаря;
// stat_hash -- <hash> тип данных -- string.
func hack_password(mode string, dictionary, stat_hash string) {
	file, err := os.Open(dictionary)
	check_error(err)
	defer file.Close()

	var reader *bufio.Reader = bufio.NewReader(file)
	var passw, dynamic_hash string

	// Украшательство
	fmt.Println("---------------------------------------------------------------")
	for {
		passw, _ = reader.ReadString('\n')
		passw = strings.Replace(passw, "\n", "", -1)
		dynamic_hash = encrypt(mode, passw)

		if stat_hash == dynamic_hash {
			fmt.Println("-------------------------------------------------------------")
			fmt.Println("[Succes]: ", passw)
			fmt.Println("-------------------------------------------------------------")
			os.Exit(0)
		} else {
			fmt.Println("[Failure]: ", passw)
		}
	}
}

func encrypt(crypt string, text string) string {
	if crypt == "md5" {
		hash := md5.New()
		hash.Write([]byte(text))
		return hex.EncodeToString(hash.Sum(nil))

	} else if crypt == "sha256" {
		hash := sha256.New()
		hash.Write([]byte(text))
		return hex.EncodeToString(hash.Sum(nil))

	} else if crypt == "sha512" {
		hash := sha512.New()
		hash.Write([]byte(text))
		return hex.EncodeToString(hash.Sum(nil))

	} else {
		hash := sha1.New()
		hash.Write([]byte(text))
		return hex.EncodeToString(hash.Sum(nil))
	}
}

func check_error(err error) {
	if err != nil {
		fmt.Println("Error is :", err)
		os.Exit(1)
	}
}

func get_error(err string) {
	fmt.Println("Error is: ", err)
	os.Exit(1)
}
