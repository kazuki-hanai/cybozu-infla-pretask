package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

type Recver struct {
	num int
	raw string
}

type Sender struct {
	num    int
	hexsha [32]byte
}

func Worker(recver chan Recver, sender chan Sender) {
	for v := range recver {
		sender <- Sender{v.num, sha256.Sum256([]byte(v.raw))}
	}
}

func Buffer(sender chan Sender) {
	i := 0
	buffer := make(map[int][32]byte)
	for v := range sender {
		buffer[v.num] = v.hexsha
		hexsha, ok := buffer[i]
		if ok {
			fmt.Println(hexsha)
			i += 1
		}
	}
}

func processConcurrent() {
	file, err := os.Open("./test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// setup worker
	const WORKER_NUM = 10
	recvers := make([]chan Recver, WORKER_NUM)
	sender := make(chan Sender, 10000)
	defer close(sender)
	for i := 0; i < WORKER_NUM; i++ {
		recvers[i] = make(chan Recver)
		go Worker(recvers[i], sender)
	}
	go Printer(sender)

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		recvers[i%WORKER_NUM] <- Recver{i, scanner.Text()}
		i += 1
	}
	for i := 0; i < WORKER_NUM; i++ {
		close(recvers[i])
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
		// sha256.Sum256([]byte(scanner.Text()))
		r := sha256.Sum256([]byte(scanner.Text()))
		fmt.Println(hex.EncodeToString(r[:]))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// processSingle()
	processConcurrent()
}
