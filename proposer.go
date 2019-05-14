package paxos

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
}

func (processor *Processor) Prepare() {
	processor.CountUpBallotNumber()

	for _, acceptor := range processor.acceptors {
		processor.Send(acceptor, Message{processor.ballotNumber})
	}

	//time.Sleep(1000 * time.Millisecond)
}

func (processor *Processor) Propose() {
	processor.TryUpdateValue()

	for _, promise := range processor.promises {
		processor.Send(promise.acceptor, Message{Proposal{processor.ballotNumber, processor.value}})
	}

	//time.Sleep(1000 * time.Millisecond)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (proposer *ProposerImpl) AddAcceptors(acceptors []*Processor) {
	proposer.acceptors = acceptors
}

func (proposer *ProposerImpl) CountUpBallotNumber() {
	proposer.ballotNumber = proposer.ballotNumber.Next()
	proposer.promises = make([]Promise, 0)
	proposer.acceptCount = 0
}

func (proposer *ProposerImpl) TryUpdateValue() {
	promiseIndex := -1
	for i, promise := range proposer.promises {
		// If proposer have promise with non-empty value, then some proposer was faster, and we should vote for its value
		if promise.acceptedProposal.value != nil {
			// Also, chosen value should have max ballot number
			// Note: this ballot numbers are not less than current ballot number, otherwise acceptor would not send promise
			if promiseIndex == -1 || proposer.promises[promiseIndex].acceptedProposal.ballotNumber.Less(promise.ballotNumber) {
				promiseIndex = i
			}
		}
	}

	if promiseIndex != -1 {
		proposer.value = proposer.promises[promiseIndex].acceptedProposal.value
	}
	// Otherwise, proposer will vote for its initial value
	// Note: ballot number is not changed
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type BallotNumber struct {
	timestamp  uint // Now()
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
