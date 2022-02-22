// Construido como parte da disciplina de Computação Concorrente
// PUCRS - Escola Politecnica
// Fillipe Almeida da Silva

/*
LANCAR N PROCESSOS EM SHELL's DIFERENTES COM NUMEROS DE IDENTIFICAÇÃO UNICOS
go run chat.go 1
go run chat.go 2
go run chat.go ...
*/

/*
	USO DOS COMANDOS
	ADIÇÃO:		add|<conteúdo da mensagem>
	EDIÇÃO:		edt|<idUsuario,idAviso>|<novo conteudo do aviso>
	DELEÇÃO		del|<idUsuario,idAviso>
*/

package main

import (
	//"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	. "./BEB"
)

func main() {

	//Define constantes de comandos
	const CMD_ADICIONA = "add|"
	const CMD_EDITA = "edt|"
	const CMD_DELETA = "del|"

	//Exemplo de utilização
	if len(os.Args) < 2 {
		//fmt.Println("Selecione pelo menos um endereço:porta. Ex:")
		//fmt.Println("go run chat.go 1 127.0.0.1:5001  127.0.0.1:6001   127.0.0.1:7001")
		fmt.Println("Obrigatório selecionar um ID para o processo.")
		fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
		return
	}

	//Inicializa o quadro de avisos como um mapa de idGlobal ("idUsuario,idAviso") em string (mensagem)
	quadAviso := make(map[string]string)

	autoEndereco := [][]string{{"127.0.0.1:5001", "127.0.0.1:6001", "127.0.0.1:7001"},
		{"127.0.0.1:6001", "127.0.0.1:5001", "127.0.0.1:7001"},
		{"127.0.0.1:7001", "127.0.0.1:5001", "127.0.0.1:6001"}}

	autoMsg := []string{"add|Teste 1", "add|Helloooo", "add|Wassup", "add|Ablablabla", "add|aaaaaaha", "add|yosss bepis", "add|do a barrel roll", "add|hewwo owo",
		"edt|1,1|Replace 1", "edt|2,1|Stolennn", "edt|3,1|get outta here", "edt|2,2|asuehuasheau", "edt|1,2|[REDACTED]", "edt|3,2|dance dance", "edt|1,3|JOJO",
		"del|1,2", "del|1,4", "del|2,3", "del|2,5", "del|3,1", "del|3,4"}

	//Guarda os endereços de comunicação com os outros usuários
	idAviso := 1
	idUsuario := os.Args[1]
	idAux, _ := strconv.Atoi(idUsuario)
	idAux -= 1

	var registro []string
	addresses := autoEndereco[idAux] //os.Args[2:]

	fmt.Println("Endereços selecionados:")
	fmt.Println(addresses)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")

	beb := BestEffortBroadcast_Module{
		Req: make(chan BestEffortBroadcast_Req_Message),
		Ind: make(chan BestEffortBroadcast_Ind_Message)}

	beb.Init(addresses[0])

	// enviador de broadcasts, irá enviar apenas comandos válidos
	//ALTERAR PARA MANDAR MENSAGENS AUTOMÁTICAS
	go func(autoMsg []string) {
		//scanner := bufio.NewScanner(os.Stdin)
		var msg string

		for {
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second) //dormir por tempo aleatório
			//if scanner.Scan() {
			//Selecionar mensagem aleatória
			msg = autoMsg[rand.Intn(len(autoMsg))] //scanner.Text()
			//fmt.Println("Comando enviado: ", msg)

			//Verifica se a mensagem é um comando
			if msg[:4] == CMD_ADICIONA || msg[:4] == CMD_EDITA || msg[:4] == CMD_DELETA {
				//Verifica primeiros 4 caracteres da mensagem para verificar se é adição
				if msg[:4] == CMD_ADICIONA {
					//Se adicionando, colocar idGlobal de aviso na mensagem
					msg += "|" + idUsuario + "," + strconv.Itoa(idAviso)
					idAviso += 1
				}

				msg += "§" + addresses[0] //Adiciona o endereço de origem do usuário na mensagem, o idUsuario e o idAviso
				//Ex: add|dgfsdifhsufh|1,4§127.0.0.1:5001

				req := BestEffortBroadcast_Req_Message{
					Addresses: addresses[0:], //Envia para todos os usuários listados (incluindo si mesmo)
					Message:   msg}
				beb.Req <- req // ENVIA PARA TODOS PROCESSOS ENDERECADOS NO INICIO
			} else {
				fmt.Println("Comando inválido recebido: " + msg[:4])
			}
			//}
		}
	}(autoMsg)

	// receptor de broadcasts
	go func(quadAviso map[string]string) {
		for {
			in := <-beb.Ind // RECEBE MENSAGEM DE QUALQUER PROCESSO NO FORMATO "<codOP>|<mensagem com idGlobal>§<endereço origem>"
			message := strings.Split(in.Message, "§")
			//message = ["<codOP>|<mensagem com idGlobal>"", "<endereço origem>"]
			in.From = message[1]
			registro = append(registro, in.Message)
			//in.Message = message[0]
			fmt.Println("Mensagem recebida: ", message[0])
			msgBody := strings.Split(message[0], "|")

			/*Realiza o tratamento de códigos
			ESTRUTURA DOS COMANDOS
						0	1						2
			ADIÇÃO:		add|<conteúdo da mensagem>|<idUsuario,idAviso>
			EDIÇÃO:		edt|<idUsuario,idAviso>|	<novo conteudo do aviso>
			DELEÇÃO		del|<idUsuario,idAviso>*/
			if msgBody[0]+"|" == CMD_ADICIONA { //Se é uma adição
				quadAviso[msgBody[2]] = msgBody[1]

			} else if msgBody[0]+"|" == CMD_EDITA { //Se é uma edição
				_, ok := quadAviso[msgBody[1]] //Se existe um aviso com essa identificação, então edita, se não, ignora
				if ok {
					quadAviso[msgBody[1]] = msgBody[2]
				}

			} else if msgBody[0]+"|" == CMD_DELETA { //Se é uma deleção
				_, ok := quadAviso[msgBody[1]] //Se existe um aviso com essa identificação, então edita, se não, ignora
				if ok {
					delete(quadAviso, msgBody[1])
				}

			} // o envio permite apenas mensagens de comando, então não é necessário um else

			//Reimprime o quadro de avisos
			for cha, val := range quadAviso {
				chaveDiv := strings.Split(cha, ",")
				fmt.Println("Usuário ", chaveDiv[0], ", Aviso ", chaveDiv[1], ": ", val)
			}
			fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
			//fmt.Printf("Message from %v: %v\n", in.From, in.Message)
		}
	}(quadAviso)

	blq := make(chan int)
	<-blq
}
