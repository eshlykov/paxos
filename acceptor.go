package paxos

type Acceptor interface {
	Promise(sender *Processor)
	Accept(sender *Processor)
}

type AcceptorImpl struct {
	maxBallotNumber  BallotNumber
	acceptedProposal Proposal
}

func (processor *Processor) Promise(sender *Processor, maxBallotNumber BallotNumber, acceptedProposal Proposal) {
	processor.Send(sender, Message{Promise{processor, processor.maxBallotNumber, processor.acceptedProposal}})
}

func (processor *Processor) Accept(sender *Processor) {
	processor.Send(sender, Message{Ack{}})
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Promise struct {
	acceptor         *Processor
	ballotNumber     BallotNumber
	acceptedProposal Proposal
}

type Ack struct{}
