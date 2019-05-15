package paxos

import "sync"

type Acceptor interface {
	Promise(sender *Processor)
	Accept(sender *Processor)
}

type AcceptorImpl struct {
	maxBallotNumber  BallotNumber
	acceptedProposal Proposal
	lockAcceptor     sync.Mutex
}

func (processor *Processor) Promise(receiver *Processor, maxBallotNumber BallotNumber, acceptedProposal Proposal) {
	processor.Send(receiver, Message{Promise{
		processor,
		processor.MaxBallotNumber(),
		processor.AcceptedProposal(),
	}})
}

func (processor *Processor) Accept(receiver *Processor) {
	processor.Send(receiver, Message{Ack{}})
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (acceptor *AcceptorImpl) MaxBallotNumber() BallotNumber {
	acceptor.lockAcceptor.Lock()
	defer acceptor.lockAcceptor.Unlock()

	return acceptor.maxBallotNumber
}

func (acceptor *AcceptorImpl) SetMaxBallotNumber(ballotNumber BallotNumber) {
	acceptor.lockAcceptor.Lock()
	defer acceptor.lockAcceptor.Unlock()

	acceptor.maxBallotNumber = ballotNumber
}

func (acceptor *AcceptorImpl) AcceptedProposal() Proposal {
	acceptor.lockAcceptor.Lock()
	defer acceptor.lockAcceptor.Unlock()

	return acceptor.acceptedProposal
}

func (acceptor *AcceptorImpl) SetAcceptedProposal(proposal Proposal) {
	acceptor.lockAcceptor.Lock()
	defer acceptor.lockAcceptor.Unlock()

	acceptor.acceptedProposal = proposal
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Promise struct {
	acceptor         *Processor
	ballotNumber     BallotNumber
	acceptedProposal Proposal
}

type Ack struct{}
