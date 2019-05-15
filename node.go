package paxos

type Node interface {
	Send(receiver *Processor, message Message)
	Receive(sender *Processor, message Message)
}

func (processor *Processor) Send(receiver *Processor, message Message) {
	network.Send(processor, receiver, message)
}

func (processor *Processor) Receive(sender *Processor, message Message) {
	switch data := message.body.(type) {
	// Acceptor
	case BallotNumber:
		if processor.maxBallotNumber.Less(data) {
			processor.maxBallotNumber = data
			// Promise not to accept proposals with less ballot number
			processor.Promise(sender, processor.maxBallotNumber, processor.acceptedProposal)
		}
	case Proposal:
		if processor.maxBallotNumber.LessOrEqual(data.ballotNumber) {
			processor.acceptedProposal = data
			// Accept value with ballot number at least promised
			processor.Accept(sender)
		}
	// Proposer
	case Promise:
		// Ignore old messages
		if data.ballotNumber == processor.ballotNumber {
			processor.promises = append(processor.promises, data)
		}
	case Ack:
		processor.acceptCount++
	}
}
