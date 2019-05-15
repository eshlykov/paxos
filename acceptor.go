package paxos

import "sync"

type Acceptor interface {
	Promise(receiver *Processor, axBallotNumber BallotNumber, acceptedProposal Proposal)
	Accept(receiver *Processor)
}

type AcceptorImpl struct {
	maxBallotNumber  BallotNumber
	acceptedProposal Proposal
	lockAcceptor     sync.Mutex
}

func (processor *Processor) Promise(receiver *Processor, maxBallotNumber BallotNumber, acceptedProposal Proposal) {
	processor.Send(receiver, Message{Promise{
		processor,
		maxBallotNumber,
		acceptedProposal,
	}})
}

func (processor *Processor) Accept(receiver *Processor) {
	processor.Send(receiver, Message{Ack{}})
}

func (processor *Processor) OnPrepare(sender *Processor, ballotNumber BallotNumber) {
	var acceptedProposal Proposal

	{
		processor.lockAcceptor.Lock()
		defer processor.lockAcceptor.Unlock()

		if !processor.maxBallotNumber.Less(ballotNumber) {
			// Ignore old ballot numbers.
			return
		}

		processor.maxBallotNumber = ballotNumber
		acceptedProposal = processor.acceptedProposal
	}

	// Promise not to accept proposals with less ballot number.
	processor.Promise(sender, ballotNumber, acceptedProposal)
}

func (processor *Processor) OnPropose(sender *Processor, proposal Proposal) {
	{
		processor.lockAcceptor.Lock()
		defer processor.lockAcceptor.Unlock()

		if !processor.maxBallotNumber.LessOrEqual(proposal.ballotNumber) {
			// Ignore old proposals.
			return
		}

		processor.acceptedProposal = proposal
	}

	// Accept value with ballot number at least promised.
	processor.Accept(sender)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Promise struct {
	acceptor         *Processor
	ballotNumber     BallotNumber
	acceptedProposal Proposal
}

type Ack struct{}
