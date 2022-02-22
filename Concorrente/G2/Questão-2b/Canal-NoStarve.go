package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// a context is required for the weighted semaphore pkg.
var ctx = context.Background()

var readers = 0 //Num leitores acessando a memória

var data = "a"

func reader(id int, mutex chan struct{}, roomEmpty chan struct{}, catraca chan struct{}) {
	for true {

		//Tranca todos os leitores até que o escritor na fila libere o acesso. Escritores trancam a memória e quando liberam, escritores na fila trancam leitores novamente.
		<-catraca
		catraca <- struct{}{}

		<-mutex

		readers += 1
		if readers == 1 {
			<-roomEmpty
		}

		mutex <- struct{}{}
		//Seção Crítica
		fmt.Println("Data: ", data, ". ID: ", id)

		<-mutex

		readers -= 1
		if readers == 0 {
			roomEmpty <- struct{}{}
		}

		mutex <- struct{}{}
		n := rand.Intn(4) // n will be between 0 and 10
		time.Sleep(time.Duration(n) * time.Second)
	}
}

func writer(id string, roomEmpty chan struct{}, catraca chan struct{}) {

	for true {
		<-catraca
		//{ Ativa a catraca para barrar leitores novos de entrarem na fila
		<-roomEmpty
		//{
		data += "|" + id
		fmt.Println("Escrita por: ", id)
		//} Libera a catraca
		catraca <- struct{}{}
		//}
		roomEmpty <- struct{}{}

		n := rand.Intn(4) // n will be between 0 and 10
		time.Sleep(time.Duration(n) * time.Second)
	}
}

func main() {

	fmt.Println("Iniciando...")
	index := []string{"1", "2", "3", "4", "5"}

	mutex := make(chan struct{}, 10)
	roomEmpty := make(chan struct{}, 10)
	catraca := make(chan struct{}, 10)
	for i := 0; i < 10; i++ {
		mutex <- struct{}{}
		roomEmpty <- struct{}{}
		catraca <- struct{}{}
	}

	for i := 0; i < 5; i++ {

		go reader(i+1, mutex, roomEmpty, catraca)
		go writer(index[i], roomEmpty, catraca)
	}

	fmt.Println("Pronto.")
	for true {
	}

}
