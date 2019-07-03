package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

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

func ParseAddresses(address_file string) []string {
	var addresses []string

	file, _ := os.Open(address_file)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		addresses = append(addresses, scanner.Text())
	}

	return addresses

}

func main() {

	myRank, _ := strconv.Atoi(os.Args[1])

	// var array [2]string

	addresses := ParseAddresses(os.Args[2])

	bla := make(chan int)

	// array[0] = "127.0.0.1:6000"
	// array[1] = "127.0.0.1:7000"
	cbto := CBTOB.Init(addresses[myRank], addresses, myRank)

	time.Sleep(2 * time.Second)
	currID := 1
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			message, _ := reader.ReadString('\n')
			message = message[:len(message)-1]
			ReqMessage := CBTOB.Req_CBTOB_Message{
				ProcSender: addresses[myRank],
				To:         addresses[:],
				MessageId:  currID,
				Message:    message,
			}
			currID++
			cbto.Req <- ReqMessage
		}
	}()

	go func() {
		for {
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
