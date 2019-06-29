package CBTOBroadcast

//no modelo atual, PP2PLink e BEBroadcast vão ser recriados diversas vezes
//revisar isso para que os módulos superiores compartilhem BEB?

//processo, processos, id, mensagem

//trocar address e nome pra Member
//arrumar nomes de variáveis - procsender = from?
//aaaaaaaa

import (
	"container/list"
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
	procSender string
	To         []string
	MessageId  int
	Message    string
}

//list: id e proc de cada mensagem

func CBTOBReqToRB(message Req_CBTOB_Message) RB.ReliableBroadcast_Req_Message {

	RBMessage := RB.ReliableBroadcast_Req_Message{
		Sender:    message.procSender,
		Addresses: message.To,
		Message:   "nothing;" + strconv.Itoa(message.MessageId) + ";" + message.Message, //adicionar messageID aqui
	}
	return RBMessage
}

func RBIndToCBTOB(message RB.ReliableBroadcast_Ind_Message) Ind_CBTOB_Message {
	parts := strings.Split(message.Message, ";")
	content := strings.Join(parts[1:], "")
	messageType := parts[0] //Decide ou whatever
	messageId, _ := strconv.Atoi(parts[1])

	CBTOBMessage := Ind_CBTOB_Message{
		ProcSender:  message.Sender,
		Message:     content,
		MessageId:   messageId,
		MessageType: messageType,
	}

	return CBTOBMessage
}

type CBTOBroadcast struct {
	// Ind                   chan HierarchicalConsensus.Message
	// Red                   chan HierarchicalConsensus.Message
	Ind chan Ind_CBTOB_Message
	Req chan Req_CBTOB_Message
	// hierarchicalConsensus *FakeConsensus
	unordered *list.List
	delivered *list.List
	rb        *RB.ReliableBroadcast_Module
	Addresses []string
	Address   string
	round     int
	rank      int
	wait      boolmessage
}

func Init(address string, addresses []string, processRank int) *CBTOBroadcast {

	var module *CBTOBroadcast

	module.Ind = make(chan Ind_CBTOB_Message)
	module.Req = make(chan Req_CBTOB_Message)
	module.rb = &RB.ReliableBroadcast_Module{
		Self:      address,
		Addresses: addresses,
	}

	module.rb.Init()
	module.Address = address
	module.unordered = list.New()
	module.delivered = list.New()

	module.rank = processRank
	module.round = 0
	module.wait = false

	// module.hierarchicalConsensus = HierarchicalConsensus.NewMultiHierarchicalConsensus() //importar corretamente deopis

	module.Start()

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
	parts := strings.Split(message.Message, ";")
	content := strings.Join(parts[1:], ";")
	messageType := parts[0]

	if messageType == "Decide" {
		module.SortMessages(content)
	} else {
		module.AddUndeliveredMessage(message)
	}

}

func (module *CBTOBroadcast) Broadcast(message Req_CBTOB_Message) {
	rbMessage := CBTOBReqToRB(message)
	module.rb.Req <- rbMessage
}

func (module *CBTOBroadcast) AddUndeliveredMessage(message RB.ReliableBroadcast_Ind_Message) {
	CBTOBMessage := RBIndToCBTOB(message)
	module.unordered.PushBack(CBTOBMessage)
}

//rever tudo isso aqui, tá horroroso
func (module *CBTOBroadcast) SortMessages(message string) {
	// if module.r == message.round {
	// decidedMessages := ConsensusToDecidedOrder(message)

	parts := strings.Split(message, ";")
	tmpList := list.New()
	for i := 0; i < len(parts); i += 2 {
		for e := module.unordered.Front(); e != nil; e = e.Next() {
			currMsg := e.Value.(Ind_CBTOB_Message)
			msgID, _ := strconv.Atoi(parts[i+1])
			if (currMsg.ProcSender == parts[i]) && (currMsg.MessageId == msgID) {
				module.Ind <- currMsg
				module.delivered.PushBack(currMsg)
				tmpList.PushBack(currMsg)
			}

		}
	}
	// module.delivered.PushBackList(decidedMessages)
	for e := tmpList.Front(); e != nil; e = e.Next() {
		module.unordered.Remove(e)
	}

	module.round++
	module.wait = false
	// }
}

func (module *CBTOBroadcast) TryNewConsensus() {
	if module.unordered.Front() != nil && !module.wait {
		module.wait = true

		// consensusInstance := module.hierarchicalConsensus.CreateInstance()
		// consensusInstance.Proposal <- MessageOrder()
		if module.rank == 0 {
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

	return RBMessage
}
