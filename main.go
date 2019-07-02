package main

import (
	"fmt"
	"time"
	"os"
	"strconv"

	CBTOB "./CBTOBroadcast"
)

/*
func main() {

	if (len(os.Args) < 2) {
		fmt.Println("Please specify at least one address:port!")
		return
	}

	addresses := os.Args[1:]
	fmt.Println(addresses)

	mod := BestEffortBroadcast_Module{
		Req: make(chan BestEffortBroadcast_Req_Message),
		Ind: make(chan BestEffortBroadcast_Ind_Message) }
	mod.Init(addresses[0])

	msg := BestEffortBroadcast_Req_Message{
		Addresses: addresses,
		Message: "BATATA!" }

	yy := make(chan string)
	mod.Req <- msg
	<- yy
}
*/

func main() {

	if (len(os.Args) < 2) {
		fmt.Println("Please specify at least one address:port!")
		return
	}

	myRank, _ := strconv.Atoi(os.Args[1])
	message := os.Args[2]



	var array [2]string

	bla := make(chan int)

	array[0] = "127.0.0.1:6000"
	array[1] = "127.0.0.1:7000"
	cbto := CBTOB.Init(array[myRank], array[:], myRank)
	ReqMessage := CBTOB.Req_CBTOB_Message{
		ProcSender: array[myRank],
		To:         array[:],
		MessageId:  1,
		Message:    message,
	}

	time.Sleep(2 * time.Second)


	go func() {
		cbto.Req <- ReqMessage
		for{
			msg := <-cbto.Ind
			fmt.Println("APP: CBTO: got a message! " + msg.Message)
		}
	}()

	<-bla

	// cbto2 := CBTOB.Init("127.0.0.1:7000", array[:], 1)
	// ReqMessage2 := CBTOB.Req_CBTOB_Message{
	// 	ProcSender: "127.0.0.1:7000",
	// 	To:         array[:],
	// 	MessageId:  1,
	// 	Message:    "drones",
	// }


	// go func() {
	// 	cbto2.Req <- ReqMessage2
	// 	for{
	// 		msg := <-cbto2.Ind
	// 		fmt.Println("APP: CBTO2: got a message! " + msg.Message)
	// 	}
	// }()

	// <- bla

	// ReqMessage.MessageId = 2
	// ReqMessage.Message = "abecedário"

	// cbto.Req <- ReqMessage

	// msg = <-cbto.Ind
	// fmt.Println("APP: got a message! " + msg.Message)

}
