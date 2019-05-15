package paxos

import "sync"

type Proposer interface {
	Prepare()
	Propose()
}

type ProposerImpl struct {
	ID           uint
	acceptors    []*Processor
	ballotNumber BallotNumber
	value        interface{}
	promises     []Promise
	acceptCount  int
	lockProposer sync.Mutex
}

func (processor *Processor) Prepare() {
	ballotNumber := processor.CountUpBallotNumber()

	for _, acceptor := range processor.acceptors {
		processor.Send(acceptor, Message{ballotNumber})
	}
}

func (processor *Processor) Propose() {
	promises, ballotNumber, value := processor.TryUpdateValue()

	for _, promise := range promises {
		processor.Send(promise.acceptor, Message{Proposal{ballotNumber, value}})
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (proposer *ProposerImpl) AddAcceptors(acceptors []*Processor) {
	proposer.acceptors = acceptors
}

func (proposer *ProposerImpl) CountUpBallotNumber() BallotNumber {
	proposer.lockProposer.Lock()
	defer proposer.lockProposer.Unlock()

	proposer.ballotNumber = proposer.ballotNumber.Next()
	proposer.promises = make([]Promise, 0, 0)
	proposer.acceptCount = 0

	return proposer.ballotNumber
}

func (proposer *ProposerImpl) TryUpdateValue() ([]Promise, BallotNumber, interface{}) {
	proposer.lockProposer.Lock()
	defer proposer.lockProposer.Unlock()

	promiseIndex := -1
	for i, promise := range proposer.promises {
		// If proposer have promise with non-empty value,
		// then some proposer was faster, and we should vote for its value
		if promise.acceptedProposal.value != nil {
			// Also, chosen value should have max ballot number
			// Note: this ballot numbers are not less than current ballot number,
			// otherwise acceptor would not send promise
			if promiseIndex == -1 ||
				proposer.promises[promiseIndex].acceptedProposal.ballotNumber.Less(promise.ballotNumber) {
				promiseIndex = i
			}
		}
	}

	if promiseIndex != -1 {
		proposer.value = proposer.promises[promiseIndex].acceptedProposal.value
	}

	// Otherwise, proposer will vote for its initial value
	// Note: ballot number is not changed

	return proposer.promises, proposer.ballotNumber, proposer.value
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (proposer *ProposerImpl) BallotNumber() BallotNumber {
	proposer.lockProposer.Lock()
	defer proposer.lockProposer.Unlock()

	return proposer.ballotNumber
}

func (proposer *ProposerImpl) AddPromise(promise Promise) {
	proposer.lockProposer.Lock()
	defer proposer.lockProposer.Unlock()

	proposer.promises = append(proposer.promises, promise)
}

func (proposer *ProposerImpl) AddAccept() {
	proposer.lockProposer.Lock()
	defer proposer.lockProposer.Unlock()

	proposer.acceptCount++
}

func (proposer *ProposerImpl) Promises() []Promise {
	proposer.lockProposer.Lock()
	defer proposer.lockProposer.Unlock()

	return proposer.promises
}

func (proposer *ProposerImpl) AcceptCount() int {
	proposer.lockProposer.Lock()
	defer proposer.lockProposer.Unlock()

	return proposer.acceptCount
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type BallotNumber struct {
	timestamp  uint
	proposerID uint
}

func (ballotNumber BallotNumber) Equal(other BallotNumber) bool {
	return ballotNumber.timestamp == other.timestamp && ballotNumber.proposerID == other.proposerID
}

func (ballotNumber BallotNumber) Less(other BallotNumber) bool {
	return ballotNumber.timestamp < other.timestamp ||
		ballotNumber.timestamp == other.timestamp && ballotNumber.proposerID < other.proposerID
}

func (ballotNumber BallotNumber) LessOrEqual(other BallotNumber) bool {
	return ballotNumber.Less(other) || ballotNumber.Equal(other)
}

func (ballotNumber BallotNumber) Next() BallotNumber {
	return BallotNumber{ballotNumber.timestamp + 1, ballotNumber.proposerID}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Proposal struct {
	ballotNumber BallotNumber
	value        interface{}
}

func (first Proposal) Less(second Proposal) bool {
	return first.ballotNumber.Less(second.ballotNumber)
}
