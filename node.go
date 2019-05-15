package paxos

type Node interface {
	Send(receiver *Processor, message Message)
	Receive(sender *Processor, message Message)
}

func (processor *Processor) Send(receiver *Processor, message Message) {
	network.Send(processor, receiver, message)
}

func (processor *Processor) Receive(sender *Processor, message Message) {
	switch body := message.body.(type) {
	// Acceptor
	case BallotNumber:
		maxBallotNumber := processor.MaxBallotNumber()
		if maxBallotNumber.Less(body) {
			processor.SetMaxBallotNumber(body)
			// Promise not to accept proposals with less ballot number
			processor.Promise(sender, body, processor.AcceptedProposal())
		}
	case Proposal:
		if processor.MaxBallotNumber().LessOrEqual(body.ballotNumber) {
			processor.SetAcceptedProposal(body)
			// Accept value with ballot number at least promised
			processor.Accept(sender)
		}
	// Proposer
	case Promise:
		// Ignore old messages
		if body.ballotNumber == processor.BallotNumber() {
			processor.AddPromise(body)
		}
	case Ack:
		processor.AddAccept()
	}
}
