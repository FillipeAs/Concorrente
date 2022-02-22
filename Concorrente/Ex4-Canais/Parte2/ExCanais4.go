package main

import (
	"fmt"
	"math/rand"
	"time"
)

var tabela = make(chan [][]bool, 8)

var algumaCoisa = []string{"Hello world, Im right here!", "Waddup bruh", "Eyyyyyy", "Hello there", "GENERAL KENOBI", "Hewwo", "Greetings"}

type Mensagem struct {
	orig []int
	dest int
	mens string
	ttl  int
}

//TODO:	Criar canais para comunocação de mensagens no formato mensagem = [id_orig, id_dest, mensagem, timeToLive] com buffer de tamanho 16 caso todos os nodos resolvam todos mandar mensagens para um destino
//		Criar o comportamento paralelo à topologia de troca de mensagens, imprimindo sempre que ocorrer evento de mensagem
//			Nodo irá tentar ler uma mensagem sem bloquear
//			Ele então irá ver se é para si, senão irá repassar utilizando um Time to Live nas mensagens
//			Ele então irá decidir se vai mandar alguma mensagem própria e então enviar para um nodo qualquer
//				Se não houver este nodo na topologia conhecida, ele irá declarar que tentou enviar mensagem impossível de ser entregue
//				Se houver, enviar para todos os vizinhos
func main() {

	//Inicializa a semente de números aleatórios
	rand.Seed(time.Now().UnixNano())

	//Canal de saída geral de cada nodo com tamanho de buffer apropriado para o número de ligações. Ordenados de nodo 1-8 sendo ChansSaida[i] onde i == id do nodo -1
	var ChansSaida = []chan [][]bool{make(chan [][]bool, 1), make(chan [][]bool, 2), make(chan [][]bool, 3), make(chan [][]bool, 3),
		make(chan [][]bool, 5), make(chan [][]bool, 2), make(chan [][]bool, 2), make(chan [][]bool, 2)}

	//Canais de ENTRADA de mensagens para cada nodo, com buffer de tamanho N-1 (7) para que suporte afunilamento de mensagens para o mesmo nodo.
	var ChansINMen [8]chan Mensagem
	for i := range ChansINMen {
		ChansINMen[i] = make(chan Mensagem, 16)
	}

	//Array de canais de entrada de cada nodo, quais nodos mandam mensagem para este em particular
	var ChanTo1 = []chan [][]bool{ChansSaida[3]}
	var ChanTo2 = []chan [][]bool{ChansSaida[2], ChansSaida[3]}
	var ChanTo3 = []chan [][]bool{ChansSaida[1], ChansSaida[4], ChansSaida[5]}
	var ChanTo4 = []chan [][]bool{ChansSaida[0], ChansSaida[1], ChansSaida[4]}
	var ChanTo5 = []chan [][]bool{ChansSaida[2], ChansSaida[3], ChansSaida[5], ChansSaida[6], ChansSaida[7]}
	var ChanTo6 = []chan [][]bool{ChansSaida[2], ChansSaida[4]}
	var ChanTo7 = []chan [][]bool{ChansSaida[4], ChansSaida[7]}
	var ChanTo8 = []chan [][]bool{ChansSaida[4], ChansSaida[6]}

	//Arrays de canais de SAIDA de cada nodo para carregar mensagens, nodos só podem mandar para vizinhos
	var ChanOUTMen1 = []chan Mensagem{ChansINMen[3]}
	var ChanOUTMen2 = []chan Mensagem{ChansINMen[2], ChansINMen[3]}
	var ChanOUTMen3 = []chan Mensagem{ChansINMen[1], ChansINMen[4], ChansINMen[5]}
	var ChanOUTMen4 = []chan Mensagem{ChansINMen[0], ChansINMen[1], ChansINMen[4]}
	var ChanOUTMen5 = []chan Mensagem{ChansINMen[2], ChansINMen[3], ChansINMen[5], ChansINMen[6], ChansINMen[7]}
	var ChanOUTMen6 = []chan Mensagem{ChansINMen[2], ChansINMen[4]}
	var ChanOUTMen7 = []chan Mensagem{ChansINMen[4], ChansINMen[7]}
	var ChanOUTMen8 = []chan Mensagem{ChansINMen[4], ChansINMen[6]}

	/*Main to Nodes*/
	var ChansToNode [8]chan bool
	for i := range ChansToNode {
		ChansToNode[i] = make(chan bool)
	}

	var ChansFromNode [8]chan bool
	for i := range ChansFromNode {
		ChansFromNode[i] = make(chan bool)
	}

	/*Dispara processos*/
	go nodo(1, []int{4}, ChansToNode[0], ChansFromNode[0], ChansSaida[0], ChanTo1, ChansINMen[0], ChanOUTMen1)
	go nodo(2, []int{3, 4}, ChansToNode[1], ChansFromNode[1], ChansSaida[1], ChanTo2, ChansINMen[1], ChanOUTMen2)
	go nodo(3, []int{2, 5, 6}, ChansToNode[2], ChansFromNode[2], ChansSaida[2], ChanTo3, ChansINMen[2], ChanOUTMen3)
	go nodo(4, []int{1, 2, 5}, ChansToNode[3], ChansFromNode[3], ChansSaida[3], ChanTo4, ChansINMen[3], ChanOUTMen4)
	go nodo(5, []int{3, 4, 6, 7, 8}, ChansToNode[4], ChansFromNode[4], ChansSaida[4], ChanTo5, ChansINMen[4], ChanOUTMen5)
	go nodo(6, []int{3, 5}, ChansToNode[5], ChansFromNode[5], ChansSaida[5], ChanTo6, ChansINMen[5], ChanOUTMen6)
	go nodo(7, []int{5, 8}, ChansToNode[6], ChansFromNode[6], ChansSaida[6], ChanTo7, ChansINMen[6], ChanOUTMen7)
	go nodo(8, []int{5, 7}, ChansToNode[7], ChansFromNode[7], ChansSaida[7], ChanTo8, ChansINMen[7], ChanOUTMen8)

	/*Contador de rodadas*/
	for i := 0; i < 8; i++ {
		fmt.Printf("Rodada %d...\n", i+1)

		for j := range ChansToNode {
			fmt.Printf("\tMandando sinal %d\n", j+1)
			ChansToNode[j] <- true
		}

		for j := range ChansFromNode {
			fmt.Printf("\tEsperando sinal %d\n", j+1)
			<-ChansFromNode[j]
		}
		fmt.Printf("Fim Rodada %d...\n", i+1)
	}

	/*Para a atualização das matrizes dos nodos*/
	for i := range ChansToNode {
		fmt.Printf("\tEncerrando sinal %d\n", i+1)
		ChansToNode[i] <- false
	}

	/*Imprime tabela de vizinhos*/
	fmt.Println("Tabela de vizinhos:")
	toString(<-tabela)
}

/*FUNÇÃO DOS NODOS*/
/*Parâmetros: id-> numero do nodo*/
/*enviaPrimeiro-> comportamento do nodo, determina se recebe ou envia primeiro*/
/*vi-> vizinhos adjacentes*/
/*fromMain-> canal sinc vindo do main*/
/*toMain-> canal sinc para o main*/
/*viOUT-> Canal de saída das arestas adjacentes*/
/*viIN-> Array de canais de entrada das arestas adjacentes*/
func nodo(id int, vi []int, fromMain chan bool, toMain chan bool, viOUT chan [][]bool, viIN []chan [][]bool, chanINMen chan Mensagem, chansOUTMen []chan Mensagem) {

	var vizinhos = [][]bool{{false, false, false, false, false, false, false, false}, {false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false}, {false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false}, {false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false}, {false, false, false, false, false, false, false, false}}

	var vizinhos_temp = [][]bool{{false, false, false, false, false, false, false, false}, {false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false}, {false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false}, {false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false}, {false, false, false, false, false, false, false, false}}

	var vi_l = vi

	/*Inicialização da matriz de vizinhos com os adjacentes*/
	for i := 0; i < 8; i++ {
		for j := i; j < 8; j++ {
			if (i+1) == id && acha(vi_l, j+1) {
				vizinhos[i][j] = true
				vizinhos[j][i] = true

			}
		}
	}

	/*Irá executar enquanto o main estiver mandando true*/
	rodada := true
	for rodada {

		rodada = <-fromMain
		if !rodada {
			break
		}

		fmt.Printf("\t\tProcesso %d is go\n", id)

		//COMPORTAMENTO DE TOPOLOGIA
		/*Irá mandar sua matriz de vizinhos para os vizinhos*/
		for i := 0; i < cap(viOUT); i++ {
			fmt.Printf("\t\t\tProcesso %d enviando matriz\n", id)
			viOUT <- vizinhos
		}

		/*Irá receber a matriz de vizinhos dos vizinhos e atualizar a sua*/
		k := 0
		for k < len(viIN) {
			vizinhos_temp = <-viIN[k]
			fmt.Printf("\t\t\tProcesso %d recebeu matriz\n", id)
			k++

			/*Atualiza matriz via OU lógico de cada posição*/
			for i := 0; i < 8; i++ {
				for j := i; j < 8; j++ {
					vizinhos[i][j] = vizinhos[i][j] || vizinhos_temp[i][j]
					vizinhos[j][i] = vizinhos[j][i] || vizinhos_temp[j][i]
				}
			}
		}

		//COMPORTAMENTO DE MENSAGEM
		//Verifica se existe entrada de mensagens no buffer de entrada
		for e := 0; e < len(chanINMen); e++ {

			i := <-chanINMen
			if i.dest == id && i.ttl > -1 {
				fmt.Printf("\t\t\t\tProcesso %d recebeu mensagem de %d: "+i.mens+"\tCaminho de origem: %v\n", id, i.orig[0], i.orig)
			} else if !acha(i.orig, id) && i.ttl > -1 { //Se a mensagem não é para mim E eu já não enviei E ainda é válida, repassa para todos os vizinhos com o ttl reduzido e meu id na lista de origem
				j := 0
				for j < len(chansOUTMen) {

					m := Mensagem{append(i.orig, id), i.dest, i.mens, i.ttl - 1}

					select {
					case chansOUTMen[j] <- m:
						fmt.Printf("\t\t\tProcesso %d enviou mensagem para o vizinho %d de %d: "+m.mens+"\n", id, j+1, len(chansOUTMen)+1)
					default:
						fmt.Printf("\t\t\tProcesso %d não pode enviar mensagem para o vizinho %d de %d: "+m.mens+"\n", id, j+1, len(chansOUTMen)+1)
					}
					j++
				}
			} else { //Se a mensagem não é mais válida ou fui eu que mandei, descarta
				fmt.Printf("\t\t\t\tProcesso %d descartou a mensagem de %d para %d com ttl %d: "+i.mens+"\tCaminho de origem: %v\n", id, i.orig[0], i.dest, i.ttl, i.orig)
			}
		}

		//Decide se vai enviar uma mensagem e para quem. 50% de chance de gerar uma mensagem
		if rand.Intn(100) < 50 {

			nodoDest := id
			for nodoDest == id {
				nodoDest = rand.Intn(8) + 1
			}

			//Verifica se o nodo destino tem algum vizinho
			chegaLa := false
			for i := 0; i < 8; i++ {
				chegaLa = chegaLa || vizinhos[nodoDest-1][i]
			}

			//Se o nodo não tem vizinhos conhecidos, não manda nada
			if !chegaLa {
				fmt.Printf("\t\t\tProcesso %d selecionou o destino inválido para %d.\n", id, nodoDest)
			} else {

				//Gera uma mensagem para o nodo destino escolhido
				m := Mensagem{[]int{id}, nodoDest, algumaCoisa[rand.Intn(7)], 4}
				fmt.Printf("\t\t\tProcesso %d enviando mensagem para o %d: "+m.mens+"\n", id, nodoDest)

				j := 0
				for j < len(chansOUTMen) {

					select {
					case chansOUTMen[j] <- m:
						fmt.Printf("\t\t\tProcesso %d enviou mensagem para o vizinho %d de %d: "+m.mens+"\n", id, j+1, len(chansOUTMen)+1)
					default:
						fmt.Printf("\t\t\tProcesso %d não pode enviar mensagem para o vizinho %d de %d: "+m.mens+"\n", id, j+1, len(chansOUTMen)+1)
					}

					j++
				}
			}
		}

		/*Diz que está pronto e espera Main mandar a próxima rodada*/
		fmt.Printf("\t\tProcesso %d pronto pra prox\n", id)
		toMain <- true
	}

	/*Manda a matriz completa para o main*/
	//fmt.Printf("\t\tProcesso %d recebeu rodada = %t\n", id, rodada)
	tabela <- vizinhos
}

/*Função auxiliar para dizer se elemento está presente no array*/
func acha(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

/*Função auxiliar para escrever matriz em forma de string*/
func toString(mat [][]bool) {

	var s string = ""

	for i := 0; i < 8; i++ {
		s = "["

		for j := 0; j < 8; j++ {
			if mat[i][j] {
				s += "V,"
			} else {
				s += "F,"
			}
		}

		s = s[:len(s)-1] + "]"
		fmt.Println(s)
	}
}
