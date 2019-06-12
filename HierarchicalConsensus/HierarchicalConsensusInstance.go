package HierarchicalConsensus

import (
	"fmt"
	"strconv"
	"../Members"
	"../BestEffortBroadcast"
)

type HierarchicalConsensusInstance struct {
	myRank int // my id in nProcs
	nProcs int // num of procs   ????
    instNumber int // number of consensus instance
	// Channels
	bebReq chan ...
	
	pfdCrash chan ...
	
	Req chan ...      // Proposal
	Ind chan ...      // Decide

    round int
	detectedranks       // set of numbers
	proposal            // type of proposed value, list of numbers?
	proposer int 
	delivered           // set of messages
	broadcast bool          
}


func newHierarchicalConsensusInstance(instanceNumber int, 
	// myId int, 
	// nProcs int  from Members package  getSelf and get...
	bebR chan ..., bebI chan ..., 
	pfdI chan ...
	) *HierarchicalConsensusInstance {

	result := HierarchicalConsensusInstance{}

	result.myRank = Members.GetSelf()   
	result.nProcs = Members.Count()
 
	result.bebReq = bebR // best effort request channel (outgoing)
	result.bebInd = bebI // best effort indication channel (incoming) 

	result.pfdCrash = pfdI // perf fail detc indication channel (incoming)
	
	result.Req = 
	result.Ind = 
	
	result.MsgReq = make(chan PaxosMessages.Message)
	result.LeadReq = make(chan PaxosMessages.Message)
	
	result.instNumber = instanceNumber
	
    result.round = 1
	result.detectedranks       // = empty
	result.proposal            // = nil
	result.proposer = 0 
	result.delivered           // nProcs falses
	result.broadcast = false   
	
	go result.Start()
	return &result
}


func (hci *HierarchicalConsensusInstance) start() {
	for {

		select {
		case y := <-hci.myReq:
			 ...			
             checkGuards() // twice ?
			 break
		case y := <-hci.bebReq:
			...
			checkGuards() // twice ?
			break
		case y := <-hci.pfdCrash:
			....;
			checkGuards() // twice ?
	        break

		}
	}
}

func (hci *HierarchicalConsensusInstance) checkGuards() {
	hci.checkNewRound()
	hci.checkDecide() 
}

func (hci *HierarchicalConsensusInstance) checkNewRound() {
	if ...
	round ++
}

func (hci *HierarchicalConsensusInstance) checkDecide() {
	if ...
	broadcast = true
	trigger ...
	trigger ...
}






