package paxos

type Processor struct {
	ProposerImpl
	AcceptorImpl
	network Network
}

func NewProcessor(network Network) *Processor {
	processor := new(Processor)
	processor.ID = generateNewID()
	processor.ballotNumber = BallotNumber{0, processor.ID}
	processor.promises = make([]Promise, 0)
	processor.maxBallotNumber = BallotNumber{0, 0}
	processor.acceptedProposal = Proposal{BallotNumber{0, 0}, nil}
	processor.network = network
	return processor
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var proposerID uint = 0

func generateNewID() uint {
	currentID := proposerID
	proposerID++
	return currentID
}
