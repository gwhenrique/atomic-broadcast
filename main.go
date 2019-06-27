package main

import (
	"fmt"
	"strings"
	CBTOB "./CBTOBroadcast"
)

func main() {
	cbto := CBTOB.init("127.0.0.1:6000", ["127.0.0.1:6000"], 1)
	ReqMessage := CBTOB.Req_CBTOB_Message{
		procSender: "127.0.0.1:6000",
		To: ["127.0.0.1:6000"],
		MessageId: 1,
		Message: "memes"
	}
	cbto.Req <- ReqMessage

	for {
		msg := <- cbto.Ind
		fmt.Println(msg.Message)
	}
}
