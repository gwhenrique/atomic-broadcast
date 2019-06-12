
//                  |                                                    ^                                   
//                  |                                                    |                                   
//                  |                                                    |                                   
//                  |                                                    |                                   
//                  |MyReq                                               |     MyInd                         
//                  |                                                    |                                   
// +----------------|----------------------------+                       |                                   
// |Multi           |                            |                       |                                   
// |Hierarchical    v                            |                       |                                   
// |Consensus         C                          |                       |                                   
// |                              +--+           |                       |                                   
// |                              |  |        C  |        +----------------------------+                     
// |                              |--|          ------------->           | Hierarchical|                     
// |                              |  |           |        | MyReq             sensus   |                     
// |                              |--|           |        |                Instance    |                     
// |                              |  |           |        |                            |                     
// |                              |--|          ------------->                         |                     
// |                              |  |        B  |        | PFDCrash                   |                     
// |                              |--|           |        |                            |                     
// |                              |  |           |        |                            |                     
// |                              |--|           |        |                            |                     
// |                              |  |          ------------->                         |                     
// |                              |--|        A  |        | BEBInd                     |                     
// |                              |  |           |        |                            |                     
// |                              +--+           |        |                            |                     
// |                                             |        |                            |                     
// |         B                          A        |        |                            |                     
// |       ^                          ^          |        +----------------------------+                     
// |       |PFDInd              BEBInd|          |               |                                           
// |       |                          |          |               | BEBReq                                    
// +-------|--------------------------|----------+               |                                           
//         |                          |                          |                                           
//         |                          |                          |                                           
//         |                          |                          v                                           
//  +------------+ +-------------------------------------------------------------------+                     
//  |            | |                                                                   |                     
//  |            | |                                                                   |                     
//  |  PerfectFD | |     BEBroadcast                                                   |                     
//  |            | |                                                                   |                     
//  |            | |                                                                   |                     
//  +------------+ +-------------------------------------------------------------------+                     
                                                                                                          
                                                                                                          
//  RESTRICTIONS
//    1)  we will not use a failure detector in our first scenarios
//        i.e we assume no failures


package HierarchicalConsensus

import (
	"fmt"
	"sort"
	"sync"
 	"strconv"
  "../Members"
	"../BestEffortBroadcast"
)
// import "../PFD" - not for now

type Proposal struct {

}
type MultiHierarchicalConsensus struct {
	// Channels
	Ind chan HierarchicalConsensusMessages.Message
	Req chan HierarchicalConsensusMessages.Message          // propose, c.round, value

	BestEffortBroadcast BEB.BestEffortBroadcast_Module
	
	// PFD module -  we are not using a perfect failure detector now

	currentInstanceNbr int  
	valueToPropose Proposal  
	instances      map[int]*HierarchicalConsensusInstance
}

func NewMultiHierarchicalConsensus() *MultiHierarchicalConsensus {

	result := MultiHierarchicalConsensus{}

	result.Ind = make(chan HierarchicalConsensusMessages.Message)
	result.Req = make(chan HierarchicalConsensusMessages.Message)

	result.BestEffortBroadcast = BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message)}
	result.BestEffortBroadcast.Init(Members.GetSelf().Adress)

	result.currentInstanceNbr = 1  
	result.valueToPropose = nil
	result.instances = make(map[int]*HierarchicalConsensusInstance)
 	
	// PFD init -  we are not using a perfect failure detector now

	// api with upper level
	go result.handleRequests()

	// multiplex - messages arrive and are sent to hirarchical consensus instances
	go result.CheckInMailbox()

	return &result
}

//
// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------
// The interaction with the application layer sending proposal and indicating decision

func (mhc *MultiHierarchicalConsensus) handleRequests() {
	for {
		r:= <- mhc.Req                       // new proposal
		instance := mhc.CreateInstance()
		instance.Req <- r                    // forward proposal to instance
	}
}


// Receive Messages from Messages Module (network), identify hc instance (or create it)
// and forward to the specific hc instance
func (mhc *MultiHierarchicalConsensus) CheckInMailbox() {
	for {
		msg :=  <- mhc.BestEfforBroadcast.Ind 
		// find instance
		// if doesnt exist create new hcinstance
		//       instance := mhc.CreateInstance()
		instance.Req <- r   
	}
}


// If crash detected  -  not for now
// func (mhc *MultiHierarchicalConsensus) CheckCrash() {
// 	for {
// 		crash := <-PFDInd
// 		forall hc instance hci, hci.PDFCrash <- crash
// 	}
// }


func (mhc *MultiHierarchicalConsensus) CreateInstance() *HierarchicalConsensusInstance {
	// assume mp.instances is locked when called

	instance := newHierarchicalConsensusInstance(
								 mhc.currentInstanceNbr
								 mhc.MyInd
								 mhc.BestEfforBroadcast.Req
								)

	mhc.instances[currentInstanceNbr] = instance
	mhc.currentInstanceNbr++

	return instance
}
