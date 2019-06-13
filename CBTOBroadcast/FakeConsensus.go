package CBTOBroadcast

type FakeConsensus struct {
	Propose chan *list.List
	Decide chan *list.List
	Proposals []*list.List
}

FakeConsensus