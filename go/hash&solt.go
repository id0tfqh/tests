package main

/*
* ToDo
* $ go build hash&salt.go
* $ ./hash&salt <type hash> <volume salt in Byte> <text for encrypt>
* Example:
* $ ./hash&salt sha512 72 МирTrueМай
* - - - - - - - - - - - - - - - - - -
* For reference:
* $ ./hash&salt -i
* *  or
* $ ./hash&salt --info
 */

import (
	"crypto/rand"
	"fmt"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

//	"golang.org/x/crypto/sha3"

func main() {
	if len(os.Args) < 4 {
		get_error("args < 4")
	}

	check_mode(os.Args[1])

	// Передаются три аргумента:
	// os.Args[1] - тип хэш функции
	// os.Args[2] - размер соли в байтах
	// os.Args[3] - строковое зрачение, которое надо преобразовать
	hashing(os.Args[1], os.Args[2], os.Args[3])

	// Если первый аргумент функции (параметр, передаваемый в командной строке)
	// равен "-i" или равен "--info" (нестрогая дизъюнкция),
	// тогда вывести информацию о программе - вызов функции get_info.
	if os.Args[1] == "-i" || os.Args[1] == "--info" {
		// Вызов функции по программе.
		getInfo()
		// В противном случае закончить работу функции
		// и вернуться обратно в функцию main.
	} else {
		return
	}

}

// Проверка корректности вводимых данных
func check_mode(mode string) {
	if mode != "sha224" && mode != "sha256" && mode != "sha384" && mode != "sha512" && mode != "sha1" {

		get_error("crypto-hash function is not found")
	}
}

// mode -- <type hash> тип данных -- string;
// dictionary -- <file from dictionary> файл словаря;
// stat_hash -- <hash> тип данных -- string.

// https://www.programmersought.com/article/34797549181/
func hashing(mode string, volume int, text string) {
	var text string

	salt := make([]byte, volume)
	rand.Read(salt)

	exitHash := pbkdf2.Key(text, salt, 65536, 32, hash)

	// Украшательство
	fmt.Println("-------------------------------------------------------------")
	fmt.Printf("hash=%X\nsalt=%X\n", exitHash, salt)
	fmt.Println("-------------------------------------------------------------")
	os.Exit(0)
}

// Функция вывода информации по работе программы
func getInfo() {
	// Пример: ./hash&salt -i
	// или: ./hash&salt --info
	fmt.Println(
		// Текст выводимого сообщения
		`ToDo:
    ./hash&salt <type hash> <volume salt in Byte> <text for encrypt>

Example:
    $ ./hash&salt sha256 16 Text`)
	// Закрыть программу без генерации ошибки
	os.Exit(0)
}

// Func check error
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
