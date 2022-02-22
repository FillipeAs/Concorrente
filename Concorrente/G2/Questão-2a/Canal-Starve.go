package main

import (
	"fmt"
	"math/rand"
	"time"
)

var readers = 0 //Num leitores acessando a memória
var data = "a"

func reader(id int, mutex chan struct{}, roomEmpty chan struct{}) {
	for true {
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

func writer(id string, roomEmpty chan struct{}) {

	for true {
		<-roomEmpty // mesmo que sem.Wait() visto na disciplina
		data += "|" + id
		fmt.Println("Escrita por: ", id)
		roomEmpty <- struct{}{} // mesmo que sem.Signal() visto na disciplina

		n := rand.Intn(4) // n will be between 0 and 10
		time.Sleep(time.Duration(n) * time.Second)
	}
}

func main() {

	fmt.Println("Iniciando...")
	index := []string{"1", "2", "3", "4", "5"}

	mutex := make(chan struct{}, 10)
	roomEmpty := make(chan struct{}, 10)
	for i := 0; i < 10; i++ {
		mutex <- struct{}{}
		roomEmpty <- struct{}{}
	}

	for i := 0; i < 5; i++ {

		go reader(i+1, mutex, roomEmpty)
		go writer(index[i], roomEmpty)
	}

	fmt.Println("Pronto.")
	for true {
	}

}
