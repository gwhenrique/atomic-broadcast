package HierarchicalConsensus

import (
	"strconv"
	"strings"
)

import PPLink "../PPLink"
import Members "../Members"

var Channel chan Message
var link *PPLink.PPLink

func Send(message Message, to *Members.Member) {

	convertedMessage := convertTo(message)
	convertedTo := to.Address

	packet := PPLink.PP2PLink_Req_Message{To: convertedTo, Message: convertedMessage}
	link.Req <- packet

	// target := strconv.Itoa(to.Name)
	// from := strconv.Itoa(message.Member.Name)
	// msgType := strconv.Itoa(int(message.Type))

}
func Init(member *Members.Member) {

	Channel = make(chan Message)
	link = PPLink.Init(member.Address)

	go listen()

}

func convertFrom(message string) Message {

	var elements = strings.Split(message, "/")
	var result = Message{}

	fromn, _ := strconv.Atoi(elements[0])
	result.Member = Members.Find(fromn)

	result.Instance, _ = strconv.Atoi(elements[1])
	i, _ := strconv.Atoi(elements[2])
	result.Type = MessageType(i)

	for j := 3; j < len(elements); j++ {
		var k, _ = strconv.Atoi(elements[j])
		result.Values = append(result.Values, k)
	}

	return result

}
func convertTo(message Message) string {

	var result string

	result += strconv.Itoa(message.Member.Name) + "/"
	result += strconv.Itoa(message.Instance) + "/"
	result += strconv.Itoa(int(message.Type)) + "/"
	for _, val := range message.Values {
		result += strconv.Itoa(val) + "/"
	}

	return result

}
func listen() {

	for {

		y := <-link.Ind
		Channel <- convertFrom(y.Message)

	}

}
