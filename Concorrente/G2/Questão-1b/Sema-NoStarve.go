package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/semaphore"
)

// a context is required for the weighted semaphore pkg.
var ctx = context.Background()

var readers = 0                                  //Num leitores acessando a memória
var mutex = semaphore.NewWeighted(int64(10))     //
var roomEmpty = semaphore.NewWeighted(int64(10)) //Semáforo que indica a quantidade de processos na sala
var catraca = semaphore.NewWeighted(int64(10))   //Semáforo que indica a quantidade de processos na sala

var data = "a"

func reader(id int) {
	for true {

		//Tranca todos os leitores até que o escritor na fila libere o acesso. Escritores trancam a memória e quando liberam, escritores na fila trancam leitores novamente.
		catraca.Acquire(ctx, 1)
		catraca.Release(1)

		mutex.Acquire(ctx, 1)

		readers += 1
		if readers == 1 {
			roomEmpty.Acquire(ctx, 1)
		}

		mutex.Release(1)
		//Seção Crítica
		fmt.Println("Data: ", data, ". ID: ", id)

		mutex.Acquire(ctx, 1)

		readers -= 1
		if readers == 0 {
			roomEmpty.Release(1)
		}

		mutex.Release(1)
		n := rand.Intn(4) // n will be between 0 and 10
		time.Sleep(time.Duration(n) * time.Second)
	}
}

func writer(id string) {

	for true {
		catraca.Acquire(ctx, 1)
		//{ Ativa a catraca para barrar leitores novos de entrarem na fila
		roomEmpty.Acquire(ctx, 1) // mesmo que sem.Wait() visto na disciplina
		//{
		data += "|" + id
		fmt.Println("Escrita por: ", id)
		//} Libera a catraca
		catraca.Release(1)
		//}
		roomEmpty.Release(1) // mesmo que sem.Signal() visto na disciplina

		n := rand.Intn(4) // n will be between 0 and 10
		time.Sleep(time.Duration(n) * time.Second)
	}
}

func main() {

	fmt.Println("Iniciando...")
	index := []string{"1", "2", "3", "4", "5"}

	//Dispara leitores
	for i := 0; i < 5; i++ {
		go reader(i + 1)
	}
	//Dispara escritores
	for i := 0; i < 5; i++ {
		go writer(index[i])
	}

	fmt.Println("Pronto.")
	for true {
	}

}
