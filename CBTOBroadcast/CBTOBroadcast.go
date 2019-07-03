package CBTOBroadcast

//no modelo atual, PP2PLink e BEBroadcast vão ser recriados diversas vezes
//revisar isso para que os módulos superiores compartilhem BEB?

//processo, processos, id, mensagem

//trocar address e nome pra Member
//arrumar nomes de variáveis - procsender = from?
//aaaaaaaa

// 127.0.0.1:7000 --- BEB: Deliver: Received '' from 127.0.0.1:53486
// panic: runtime error: slice bounds out of range

// goroutine 10 [running]:
// _/home/15110525/Desktop/atomic-broadcast/CBTOBroadcast.RBIndToCBTOB(0xc000124090, 0xf, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
// 	/home/15110525/Desktop/atomic-broadcast/CBTOBroadcast/CBTOBroadcast.go:69 +0x2a5
// _/home/15110525/Desktop/atomic-broadcast/CBTOBroadcast.(*CBTOBroadcast).AddUndeliveredMessage(0xc000084680, 0xc000124090, 0xf, 0x0, 0x0, 0x0, 0x0)
// 	/home/15110525/Desktop/atomic-broadcast/CBTOBroadcast/CBTOBroadcast.go:161 +0x5a
// _/home/15110525/Desktop/atomic-broadcast/CBTOBroadcast.(*CBTOBroadcast).CheckMessage(0xc000084680, 0xc000124090, 0xf, 0x0, 0x0, 0x0, 0x0)
// 	/home/15110525/Desktop/atomic-broadcast/CBTOBroadcast/CBTOBroadcast.go:148 +0x217
// _/home/15110525/Desktop/atomic-broadcast/CBTOBroadcast.(*CBTOBroadcast).Start.func1(0xc000084680)
// 	/home/15110525/Desktop/atomic-broadcast/CBTOBroadcast/CBTOBroadcast.go:122 +0x177
// created by _/home/15110525/Desktop/atomic-broadcast/CBTOBroadcast.(*CBTOBroadcast).Start
// 	/home/15110525/Desktop/atomic-broadcast/CBTOBroadcast/CBTOBroadcast.go:116 +0x3f
// exit status 2

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"

	RB "../ReliableBroadcast"
	// "../HierarchicalConsensus"
)

type Ind_CBTOB_Message struct {
	ProcSender  string
	MessageId   int
	Message     string
	MessageType string
}

type Req_CBTOB_Message struct {
	ProcSender string
	To         []string
	MessageId  int
	Message    string
}

type CBTOBroadcast struct {
	// Ind                   chan HierarchicalConsensus.Message
	// Red                   chan HierarchicalConsensus.Message
	Ind chan Ind_CBTOB_Message
	Req chan Req_CBTOB_Message
	// hierarchicalConsensus *FakeConsensus
	unordered *list.List
	delivered *list.List
	rb        RB.ReliableBroadcast_Module
	Addresses []string
	Address   string
	round     int
	rank      int
	wait      bool
}

//CBTOBReqToRB: Conversão de pedido da aplicação para envio por reliable broadcast
func CBTOBReqToRB(message Req_CBTOB_Message) RB.ReliableBroadcast_Req_Message {

	RBMessage := RB.ReliableBroadcast_Req_Message{
		Sender:    message.ProcSender,
		Addresses: message.To,
		Message:   "nothing;" + strconv.Itoa(message.MessageId) + ";" + message.Message, //adicionar messageID aqui
	}

	fmt.Println("Mensagem criada para o RB: " + RBMessage.Sender)

	return RBMessage
}

//RBIndToCBTOB: Indicação do reliable broadcast para mensagem para a aplicação
func RBIndToCBTOB(message RB.ReliableBroadcast_Ind_Message) Ind_CBTOB_Message {
	parts := strings.Split(message.Message, ";")
	content := strings.Join(parts[2:], ";")
	messageType := parts[0] //Decide ou whatever
	messageID, _ := strconv.Atoi(parts[1])
	fmt.Println("CBTOB: RBIndToCBTOB: content is: " + content)

	CBTOBMessage := Ind_CBTOB_Message{
		ProcSender:  message.Sender,
		Message:     content,
		MessageId:   messageID,
		MessageType: messageType,
	}

	return CBTOBMessage
}

func Init(address string, addresses []string, processRank int) *CBTOBroadcast {

	var module *CBTOBroadcast
	module = &CBTOBroadcast{}
	module.Ind = make(chan Ind_CBTOB_Message)
	module.Req = make(chan Req_CBTOB_Message)
	module.rb = RB.ReliableBroadcast_Module{
		Self:      address,
		Addresses: addresses,
		Req:       make(chan RB.ReliableBroadcast_Req_Message),
		Ind:       make(chan RB.ReliableBroadcast_Ind_Message),
		Delivered: make(map[string]bool),
	}
	module.rb.Init()
	module.Address = address
	module.Addresses = addresses
	module.unordered = list.New()
	module.delivered = list.New()

	module.rank = processRank
	module.round = 0
	module.wait = false

	// module.hierarchicalConsensus = HierarchicalConsensus.NewMultiHierarchicalConsensus() //importar corretamente deopis

	module.Start()
	fmt.Println("CBTOB Init!")
	return module

}

func (module *CBTOBroadcast) Start() {
	go func() {
		for {
			select {
			case y := <-module.Req:
				module.Broadcast(y)
			case y := <-module.rb.Ind:
				module.CheckMessage(y)
				// module.AddUndeliveredMessage(y)
				module.TryNewConsensus()
				// case y := <- module.hierarchicalConsensus.Decide:
				// 	module.SortMessages(y) //mudar posteriormente,
				// 	//no momento não existirá mensagem do consenso
				// 	module.TryNewConsensus()
			}
		}
	}()

}

func (module *CBTOBroadcast) CheckMessage(message RB.ReliableBroadcast_Ind_Message) {
	// fmt.Println("CBTOB: CheckMessage: got message: " + message.Message)
	parts := strings.Split(message.Message, ";")
	content := strings.Join(parts[1:], ";")
	messageType := parts[0]

	//adicionar cabeçalho para diferenciar entre consenso e atomic broadcast

	if messageType == "Decide" {
		// fmt.Println("CBTOB: Got decision message: " + message.Message)
		module.SortMessages(content)
	} else {
		// fmt.Println("CBTOB: got some message, adding to undelivered: " + message.Message)
		module.AddUndeliveredMessage(message)
	}

}

func (module *CBTOBroadcast) Broadcast(message Req_CBTOB_Message) {
	// fmt.Println("CBTOB: Broadcasting " + message.Message)
	rbMessage := CBTOBReqToRB(message)
	module.rb.Req <- rbMessage
	// fmt.Println("CBTOB: enviei pro RB")
}

func (module *CBTOBroadcast) AddUndeliveredMessage(message RB.ReliableBroadcast_Ind_Message) {
	CBTOBMessage := RBIndToCBTOB(message)
	module.unordered.PushBack(CBTOBMessage)
	// fmt.Println("CBTOB: undordered message added")
}

//rever tudo isso aqui, tá horroroso
//método não remove mensagens do conjunto de não ordenadas
func (module *CBTOBroadcast) SortMessages(message string) {
	// if module.r == message.round {
	// decidedMessages := ConsensusToDecidedOrder(message)
	// fmt.Println("CBTOB: sorting: " + message)
	parts := strings.Split(message, ";")
	parts = parts[:len(parts)-1]
	tmpList := list.New()
	// fmt.Println("CBTOB: parts length: " + strconv.Itoa(len(parts)))
	for i := 0; i < len(parts); i += 2 {
		for e := module.unordered.Front(); e != nil; e = e.Next() {
			currMsg := e.Value.(Ind_CBTOB_Message)
			msgID, _ := strconv.Atoi(parts[i+1])
			if (currMsg.ProcSender == parts[i]) && (currMsg.MessageId == msgID) {
				fmt.Println("CBTOB: SortMessages: encontrei mensagem igual: " + parts[i])
				module.Ind <- currMsg
				fmt.Println("CBTOB: Enviado para o APP")
				module.delivered.PushBack(currMsg)
				tmpList.PushBack(currMsg)
			}

		}
	}
	// module.delivered.PushBackList(decidedMessages)
	for e := tmpList.Front(); e != nil; e = e.Next() {
		msg := e.Value.(Ind_CBTOB_Message)
		for m := module.unordered.Front(); m != nil; m = m.Next() {
			b := m.Value.(Ind_CBTOB_Message)
			if (msg.ProcSender == b.ProcSender) && (msg.MessageId == b.MessageId) {
				module.unordered.Remove(m)

			}
		}
	}

	fmt.Println("CBTOB: SortMessages: size of unordered after sorting: " + strconv.Itoa(module.unordered.Len()))

	module.round++
	module.wait = false
	// }
}

func (module *CBTOBroadcast) TryNewConsensus() {
	// fmt.Println("CBTOB: Trying to create consensus")
	if module.unordered.Front() != nil && !module.wait {
		module.wait = true
		// consensusInstance := module.hierarchicalConsensus.CreateInstance()
		// consensusInstance.Proposal <- MessageOrder()
		if module.rank == 0 {
			// fmt.Println("CBTOB: creating consensus")
			// fmt.Println("CBTOB: sent decision to RB")
			decision := module.DecisionToRB()
			fmt.Println("CBTOB: decision message: " + decision.Message)
			module.rb.Req <- module.DecisionToRB()
		}

		// module.hierarchicalConsensus.Propose <- module.unordered
		//envia mensagem de propose para consenso:
		//  mensagem contendo a lista unordered

	}
}

func (module *CBTOBroadcast) MessageOrder() string {

	MessageOrder := ""

	for e := module.unordered.Front(); e != nil; e = e.Next() {
		currMsg := e.Value.(Ind_CBTOB_Message)
		MessageOrder += currMsg.ProcSender + ";" + strconv.Itoa(currMsg.MessageId) + ";"
	}

	return MessageOrder

}

//pensar melhor sobre esse método pirado aqui
func (module *CBTOBroadcast) DecisionToRB() RB.ReliableBroadcast_Req_Message {

	decidedMessageOrder := module.MessageOrder()

	RBMessage := RB.ReliableBroadcast_Req_Message{
		Addresses: module.Addresses,
		Message:   "Decide;" + decidedMessageOrder,
		Sender:    module.Address,
	}
	// fmt.Println("CBTOB: created decision message :::: " + RBMessage.Message)
	// a := ""
	// for i := 0; i < len(module.Addresses); i++ {
	// 	a += module.Addresses[i]
	// }
	// fmt.Println(a)
	return RBMessage
}
