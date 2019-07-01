package main

import (
	"fmt"

	CBTOB "./CBTOBroadcast"
)

func main() {
	var array [1]string

	array[0] = "127.0.0.1:6000"
	cbto := CBTOB.Init("127.0.0.1:6000", array[:], 0)
	ReqMessage := CBTOB.Req_CBTOB_Message{
		ProcSender: "127.0.0.1:6000",
		To:         array[:],
		MessageId:  1,
		Message:    "memes",
	}
	cbto.Req <- ReqMessage

	msg := <-cbto.Ind
	fmt.Println("APP: got a message! " + msg.Message)

	ReqMessage.MessageId = 2
	ReqMessage.Message = "abecedário"

	cbto.Req <- ReqMessage

	msg = <-cbto.Ind
	fmt.Println("APP: got a message! " + msg.Message)

}
