package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/semaphore"
)

type Lightswitch struct {
	counter int
	mutex   *semaphore.Weighted
}

func lock(ls *Lightswitch, sem *semaphore.Weighted) {
	ls.mutex.Acquire(ctx, 1)
	ls.counter += 1
	if ls.counter == 1 {
		sem.Acquire(ctx, 1)
	}

	ls.mutex.Release(1)
}

func unlock(ls *Lightswitch, sem *semaphore.Weighted) {
	ls.mutex.Acquire(ctx, 1)
	ls.counter -= 1
	if ls.counter == 0 {
		sem.Release(1)
	}
	ls.mutex.Release(1)
}

// a context is required for the weighted semaphore pkg.
var ctx = context.Background()

var mutex = semaphore.NewWeighted(int64(10))    //
var noReader = semaphore.NewWeighted(int64(10)) // Impede outros leitores de acessar
var noWriter = semaphore.NewWeighted(int64(10)) //Impede outros escritores de acessar

var data = "a"

func reader(id int, readSwitch *Lightswitch) {
	for true {

		noReader.Acquire(ctx, 1)
		lock(readSwitch, noWriter)
		noReader.Release(1)

		//Seção Crítica
		fmt.Println("Data: ", data, ". ID: ", id)

		unlock(readSwitch, noWriter)

		n := rand.Intn(4) // n will be between 0 and 10
		time.Sleep(time.Duration(n) * time.Second)
	}
}

func writer(id string, writeSwitch *Lightswitch) {

	for true {

		lock(writeSwitch, noReader)
		noWriter.Acquire(ctx, 1)
		data += "|" + id
		fmt.Println("Escrita por: ", id)
		noWriter.Release(1)
		unlock(writeSwitch, noReader)

		n := rand.Intn(4) // n will be between 0 and 10
		time.Sleep(time.Duration(n) * time.Second)
	}
}

func main() {

	readSwitch := Lightswitch{
		counter: 0,
		mutex:   semaphore.NewWeighted(int64(10)),
	}
	writeSwitch := Lightswitch{
		counter: 0,
		mutex:   semaphore.NewWeighted(int64(10)),
	}

	fmt.Println("Iniciando...")
	index := []string{"1", "2", "3", "4", "5"}

	//Dispara leitores
	for i := 0; i < 5; i++ {
		go reader(i+1, &readSwitch)
	}
	//Dispara escritores
	for i := 0; i < 5; i++ {
		go writer(index[i], &writeSwitch)
	}

	fmt.Println("Pronto.")
	for true {
	}

}
