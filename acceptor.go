package paxos

type Acceptor interface {
	Promise(sender *Processor)
	Accept(sender *Processor)
}

type AcceptorImpl struct {
	maxBallotNumber  BallotNumber
	acceptedProposal Proposal
}

func (processor *Processor) Promise(receiver *Processor, maxBallotNumber BallotNumber, acceptedProposal Proposal) {
	processor.Send(receiver, Message{Promise{processor, processor.maxBallotNumber, processor.acceptedProposal}})
}

func (processor *Processor) Accept(receiver *Processor) {
	processor.Send(receiver, Message{Ack{}})
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Promise struct {
	acceptor         *Processor
	ballotNumber     BallotNumber
	acceptedProposal Proposal
}

type Ack struct{}
