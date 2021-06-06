package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

type Mapper struct {
	num int
	val string
}

func Worker(recver, sender chan Mapper, done chan bool) {
	for v := range recver {
		r := sha256.Sum256([]byte(v.val))
		sender <- Mapper{v.num, hex.EncodeToString(r[:])}
	}
	done <- true
}

func Printer(sender chan Mapper) {
	i := 0
	buffer := make(map[int]string)
	for v := range sender {
		buffer[v.num] = v.val
		for hexsha, ok := buffer[i]; ok; hexsha, ok = buffer[i] {
			fmt.Println(i, hexsha)
			delete(buffer, i)
			i += 1
		}
	}
}

func processConcurrent2() {
	file, err := os.Open("./test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// setup worker
	const WORKER_NUM = 1
	recvers := make([]chan Mapper, WORKER_NUM)
	sender := make(chan Mapper, 10000)
	done := make([]chan bool, WORKER_NUM)
	defer close(sender)
	for i := 0; i < WORKER_NUM; i++ {
		recvers[i] = make(chan Mapper)
		done[i] = make(chan bool)
		defer close(done[i])
		go Worker(recvers[i], sender, done[i])
	}
	go Printer(sender)

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		recvers[i%WORKER_NUM] <- Mapper{i, scanner.Text()}
		i += 1
	}
	for i := 0; i < WORKER_NUM; i++ {
		close(recvers[i])
		<-done[i]
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func processConcurrent1() {
	file, err := os.Open("./test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		printChecksum := func(s string) {
			r := sha256.Sum256([]byte(s))
			fmt.Println(hex.EncodeToString(r[:]))
		}
		go printChecksum(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func processSingle() {
	file, err := os.Open("./test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r := sha256.Sum256([]byte(scanner.Text()))
		fmt.Println(hex.EncodeToString(r[:]))
		// _ = hex.EncodeToString(r[:])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	processSingle()
	// processConcurrent1()
	// processConcurrent2()
}
