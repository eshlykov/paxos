package paxos

import (
	"math/rand"
	"sync"
	"time"
)

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

	processor.Wait()
}

func (processor *Processor) Propose() {
	promises, ballotNumber, value := processor.TryUpdateValue()

	for _, promise := range promises {
		processor.Send(promise.acceptor, Message{Proposal{ballotNumber, value}})
	}

	processor.Wait()
}

func (processor *Processor) OnPromise(promise Promise) {
	processor.lockProposer.Lock()
	defer processor.lockProposer.Unlock()

	if promise.ballotNumber == processor.ballotNumber {
		processor.promises = append(processor.promises, promise)
	}
}

func (processor *Processor) OnAccept() {
	processor.lockProposer.Lock()
	defer processor.lockProposer.Unlock()

	processor.acceptCount++
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
		if promise.acceptedProposal.value == nil {
			continue
		}

		// If proposer have promise with non-empty value,  then some proposer was faster,
		// and we should vote for its value. Also, chosen value should have max ballot number.
		// Note: this ballot numbers cannot be less than current ballot number,
		// otherwise acceptor would not send promise.
		if promiseIndex == -1 || proposer.promises[promiseIndex].acceptedProposal.ballotNumber.Less(promise.ballotNumber) {
			promiseIndex = i
		}
	}

	if promiseIndex != -1 {
		proposer.value = proposer.promises[promiseIndex].acceptedProposal.value
	}

	// Otherwise, proposer will vote for its initial value.
	// Note: ballot number is not changed.

	return proposer.promises, proposer.ballotNumber, proposer.value
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (proposer *ProposerImpl) PromiseCount() int {
	proposer.lockProposer.Lock()
	defer proposer.lockProposer.Unlock()

	return len(proposer.promises)
}

func (proposer *ProposerImpl) AcceptCount() int {
	proposer.lockProposer.Lock()
	defer proposer.lockProposer.Unlock()

	return proposer.acceptCount
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (proposer *ProposerImpl) Wait() {
	time.Sleep(time.Duration(RandomInRange(50, 100)) * time.Millisecond)
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var source = rand.NewSource(time.Now().UnixNano())
var random = rand.New(source)
var lock sync.Mutex

func RandomInRange(low, high float64) float64 {
	lock.Lock()
	defer lock.Unlock()

	return low + (high-low)*random.Float64()
}
