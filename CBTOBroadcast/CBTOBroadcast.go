package CBTOBroadcast


//no modelo atual, PP2PLink e BEBroadcast vão ser recriados diversas vezes
//revisar isso para que os módulos superiores compartilhem BEB?



import (
	"container/list"

	"../HierarchicalConsensus"
)

type Ind_CBTOB_Message struct {
}

type Req_CBTOB_Message struct {
}

type CBTOBroadcast struct {
	Ind                   chan HierarchicalConsensus.Message
	Red                   chan HierarchicalConsensus.Message
	Decide                chan Decide_Consensus_Message
	hierarchicalConsensus *MultiHierarchicalConsensus
	unordered             *list.List
	delivered             *list.List
	rb                    *RBroadcast
	round                 int
	rank                  int
	wait                  bool
}

func Init(address string, processRank int, pp2p *PP2PLink) *CBTOB {

	var module *CBTOB

	module.Ind = make(chan Ind_CBTOB_Message)
	module.Req = make(chan Req_CBTOB_Message)
	module.Decide = make(chan Decide_Consensus_Message)
	module.rb = ReliableBroadcast.ReliableBroadcast_Module{
		
	}

	module.unordered = list.New()
	module.delivered = list.New()

	module.rank = processRank
	module.round = 0
	module.wait = false

	module.hierarchicalConsensus = HierarchicalConsensus.NewMultiHierarchicalConsensus() //importar corretamente deopis

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
				//conversão da mensagem para tipo entendível
			case y := <- module.hierarchicalConsensus.Ind:
				module.SortMessages(y)
				module.TryNewConsensus()
			}
		}
	}()

}

func (module *CBTOBroadcast) Broadcast(message Req_CBTOB_Message) {
	rbMessage := ReqCBTOBToRB(message)
	module.rb.Req <- rbMessage
}

func (module *CBTOBroadcast) AddUndeliveredMessage(message ...) {
	//talvez converter para outro tipo?
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

		//envia mensagem de propose para consenso:
		//  mensagem contendo a lista unordered

	}
}
