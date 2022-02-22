package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Lightswitch struct {
	counter int
	mutex   chan struct{}
}

func lock(ls *Lightswitch, ch chan struct{}) {
	<-ls.mutex
	ls.counter += 1
	if ls.counter == 1 {
		<-ch
	}

	ls.mutex <- struct{}{}
}

func unlock(ls *Lightswitch, ch chan struct{}) {
	<-ls.mutex
	ls.counter -= 1
	if ls.counter == 0 {
		ch <- struct{}{}
	}
	ls.mutex <- struct{}{}
}

var data = "a"

func reader(id int, readSwitch *Lightswitch, noReader chan struct{}, noWriter chan struct{}) {
	for true {

		<-noReader
		lock(readSwitch, noWriter)
		noReader <- struct{}{}

		//Seção Crítica
		fmt.Println("Data: ", data, ". ID: ", id)

		unlock(readSwitch, noWriter)

		n := rand.Intn(4) // n will be between 0 and 10
		time.Sleep(time.Duration(n) * time.Second)
	}
}

func writer(id string, writeSwitch *Lightswitch, noReader chan struct{}, noWriter chan struct{}) {

	for true {

		lock(writeSwitch, noReader)
		<-noWriter
		data += "|" + id
		fmt.Println("Escrita por: ", id)
		noWriter <- struct{}{}
		unlock(writeSwitch, noReader)

		n := rand.Intn(4) // n will be between 0 and 10
		time.Sleep(time.Duration(n) * time.Second)
	}
}

func main() {

	fmt.Println("Iniciando...")
	index := []string{"1", "2", "3", "4", "5"}

	mutex := make(chan struct{}, 10)
	noReader := make(chan struct{}, 10)
	noWriter := make(chan struct{}, 10)

	readSwitch := Lightswitch{
		counter: 0,
		mutex:   make(chan struct{}, 10),
	}
	writeSwitch := Lightswitch{
		counter: 0,
		mutex:   make(chan struct{}, 10),
	}

	for i := 0; i < 10; i++ {
		mutex <- struct{}{}
		noReader <- struct{}{}
		noWriter <- struct{}{}
		readSwitch.mutex <- struct{}{}
		writeSwitch.mutex <- struct{}{}
	}

	//Dispara leitores
	for i := 0; i < 5; i++ {
		go reader(i+1, &readSwitch, noReader, noWriter)
		go writer(index[i], &writeSwitch, noReader, noWriter)
	}

	fmt.Println("Pronto.")
	for true {
	}

}
