package CBTOBroadcast


//no modelo atual, PP2PLink e BEBroadcast vão ser recriados diversas vezes
//revisar isso para que os módulos superiores compartilhem BEB?


//processo, processos, id, mensagem


import (
	"container/list"
	"strconv"

	// "../HierarchicalConsensus"
)

type Ind_CBTOB_Message struct {
	ProcSender string
	MessageId int
	Message string,
	MessageType string
}

type Req_CBTOB_Message struct {
	procSender string
	To []string
	MessageId int
	Message string
}


//list: id e proc de cada mensagem
//gabriel_g_santos@hotmail.com git do gabriel

func CBTOBReqToRB(message Req_CBTOB_Message) ReliableBroadcast_Req_Message {
	RBMessage := ReliableBroadcast_Req_Message {
		From: message.ProcSender,
		Addresses: message.To,
		Message: "nothing;" + strconv.Itoa(MessageId) + ";" + message.Message //adicionar messageID aqui
	}
	return RBMessage
}

func RBIndToCBTOB(message ReliableBroadcast_Ind_Message) Ind_CBTOB_Message {
	parts := strings.Split(message.Message, ";")
	content := strings.Join(parts[1:], "")
	messageType := parts[0] //Decide ou whatever
	messageid := parts[1]



	CBTOBMessage := Ind_CBTOB_Message {
		ProcSender: message.Sender,
		Message: content,
		MessageId: messageId,
		MessageType: messageType
	}

	return CBTOBMessage
}

type CBTOBroadcast struct {
	Ind                   chan HierarchicalConsensus.Message
	Red                   chan HierarchicalConsensus.Message
	hierarchicalConsensus *FakeConsensus
	unordered             *list.List
	delivered             *list.List
	rb                    *RBroadcast
	Addresses             []string
	round                 int
	rank                  int
	wait                  bool
}

func Init(address string, processRank int, pp2p *PP2PLink) *CBTOB {

	var module *CBTOB

	module.Ind = make(chan Ind_CBTOB_Message)
	module.Req = make(chan Req_CBTOB_Message)
	module.rb = ReliableBroadcast.ReliableBroadcast_Module{
		
	}

	module.unordered = list.New()
	module.delivered = list.New()

	module.rank = processRank
	module.round = 0
	module.wait = false

	// module.hierarchicalConsensus = HierarchicalConsensus.NewMultiHierarchicalConsensus() //importar corretamente deopis

	module.Start()

	return module

}

func (module *CBTOB) Start() {
	go func() {
		for {
			select {
			case y := <-module.Req:
				module.Broadcast(y)
			case y := <-module.rb.Ind:
				module.AddUndeliveredMessage(y)
				module.TryNewConsensus()
			case y := <- module.hierarchicalConsensus.Decide:
				module.SortMessages(y) //mudar posteriormente,
				//no momento não existirá mensagem do consenso
				module.TryNewConsensus()
			}
		}
	}()

}

func (module *CBTOBroadcast) Broadcast(message Req_CBTOB_Message) {
	rbMessage := CBTOBReqToRB(message)
	module.rb.Req <- rbMessage
}

func (module *CBTOBroadcast) AddUndeliveredMessage(message ReliableBroadcast_Ind_Message) {
	CBTOBMessage := RBIndToCBTOB(message)
	module.unordered(message)
}

func (module *CBTOBroadcast) SortMessages(message Consensus_Message) {
	if module.r == round {
		decidedMessages := ConsensusToDecidedOrder(message)
		for e := decidedMessages.Front(); e != nil; e = e.Next() {
			module.Ind <- e
		}
		module.delivered.PushBackList(decidedMessages)
		for e := decidedMessages.Front(); e != nil; e = e.Next() {
			module.unordered.Remove(e)
		}
		module.round++
		module.wait = false
	}
}

func (module *CBTOBroadcast) TryNewConsensus() {
	if module.unordered.front() != nil && !module.wait {
		module.wait = true

		consensusInstance := module.hierarchicalConsensus.CreateInstance()

		consensusInstance.Proposal <- MessageOrder()

		// module.hierarchicalConsensus.Propose <- module.unordered
		//envia mensagem de propose para consenso:
		//  mensagem contendo a lista unordered

	}
}

func (module *CBTOBMessage) MessageOrder() string {
	
	MessageOrder := ""

	for e := unordered.Front(); e != nil; e = e.Next() {
		MessageOrder += strconv.Itoa(e.ProcSender) + ";" strconv.Itoa(e.MessageId) + ";"
	}
	
	retun MessageOrder

}

//pensar melhor sobre esse método pirado aqui
func (module *CBTOBroadcast) DecisionToRB() ReliableBroadcast_Req_Message {
	
	decidedMessageOrder := MessageOrder()

	RBMessage := ReliableBroadcast_Req_Message {
		Addresses: module.Addresses,
		Message: "Decide;" + decidedMessageOrder,
		Sender: module.rank
	}
}
