package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	//"time"
)

//função auxiliar para imprimir matriz
func toString(mat [N][N]float64) {
	var s string = ""

	for i := 0; i < N; i++ {
		s = "["
		for j := 0; j < N; j++ {
			s += strconv.FormatFloat(mat[i][j], 'f', 4, 64)
			s += ","
		}

		s = s[:len(s)-1] + "]"
		fmt.Println(s)
	}
}

func mediaMat(i int, j int, wg1 *sync.WaitGroup, wg2 *sync.WaitGroup) {
	var med float64 = 0
	//Irá determinar qual a posição do processo na matriz para tratar de acordo
	/*	i
		[00 01 ... 0N-1] j
		[10 11 ... 1N-1]
		[.. ... ... ...]
		[N0 N1 ... N-1N-1]
	*/

	//Realiza várias iterações
	for i := 0; i < I; i++ {
		//Inicia Fase 1
		//Se for um canto, tem apenas dois vizinhos
		if i == 0 && j == 0 {
			med = (matriz[i][j] + matriz[i][j+1] + matriz[i+1][j]) / 3
		} else if i == 0 && j == (N-1) {
			med = (matriz[i][j] + matriz[i][j-1] + matriz[i+1][j]) / 3
		} else if i == (N-1) && j == 0 {
			med = (matriz[i][j] + matriz[i][j+1] + matriz[i-1][j]) / 3
		} else if i == (N-1) && j == (N-1) {
			med = (matriz[i][j] + matriz[i][j-1] + matriz[i-1][j]) / 3
		} else if i == 0 { //Se for uma aresta, tem apenas 3 vizinhos
			med = (matriz[i][j] + matriz[i][j-1] + matriz[i][j+1] + matriz[i+1][j]) / 4
		} else if i == N-1 {
			med = (matriz[i][j] + matriz[i][j-1] + matriz[i][j+1] + matriz[i-1][j]) / 4
		} else if j == 0 {
			med = (matriz[i][j] + matriz[i+1][j] + matriz[i][j+1] + matriz[i-1][j]) / 4
		} else if j == N-1 {
			med = (matriz[i][j] + matriz[i+1][j] + matriz[i][j-1] + matriz[i-1][j]) / 4
		} else { //Se for no meio, trata normalmente com 4 vizinhos
			med = (matriz[i][j] + matriz[i-1][j] + matriz[i+1][j] + matriz[i][j+1] + matriz[i][j-1]) / 5
		}
		wg1.Done()
		wg1.Wait()
		//time.Sleep(100 * time.Millisecond)

		//Inicia Fase 2
		matriz[i][j] = med
		wg2.Done()
		wg2.Wait()
	}
}

//Declara constante N do tamanho da matriz e a constante I do número de iterações
const N int = 5
const I int = 4

//Declara matriz
var matriz [N][N]float64

func main() {
	//Declara Wait Groups para realizar a sincronização
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup

	//Inicializa a matriz de valores
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			//Atribui à cada posição da matriz um valor real no intervalo [-10,10)
			matriz[i][j] = (rand.Float64() * 20) - 10
		}
	}

	//Informa o WaitGroup que serão lançados processos no número de posições da matriz
	wg1.Add(N * N)
	//wg2.Add(N * N)

	//Dispara processos para cada indice da matriz
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			go mediaMat(i, j, &wg1, &wg2)
		}
	}

	//Irá começar as iterações. Main fica responsável por monitorar o fim de cada fase, imprimir a matriz no fim da segunda e reciclar o WaitGroup conforme as fases acabam
	for i := 0; i < I; i++ {
		//Espera a fase 1 terminar
		println("Iteração", i)
		println("\tFase 1 Iniciando...")
		wg1.Wait()
		println("\tFase 1 Concluída.")
		//Adiciona os processos novamente para iniciar a fase 2
		wg2.Add(N * N)
		println("\tFase 2 Iniciando...")
		wg2.Wait()
		println("\tFase 2 Concluída. Matriz:")
		toString(matriz)
		//Adiciona os processos novamente para iniciar a próxima fase 1
		wg1.Add(N * N)
	}

	println("Fim da execução.")
}
