package HierarchicalConsensus

// DEFINIR TIPO DA CONTEUDO (values)

import "../Members"

type MessageType int
type Message struct {
	Instance int
	Type     MessageType
	Member   *Members.Member
	Values   []int
}

const Decided MessageType = 0
