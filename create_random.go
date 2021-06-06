package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

func create_random_str() string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 64
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String() + "\n"
}

func main() {
	file, err := os.Create("./test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for i := 0; i < 27145828; i++ {
		if i%1000000 == 0 {
			fmt.Println(i)
		}
		file.WriteString(create_random_str())
	}
}
