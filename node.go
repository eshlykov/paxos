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
	case BallotNumber:
		processor.OnPrepare(sender, body)
	case Proposal:
		processor.OnPropose(sender, body)
	case Promise:
		processor.OnPromise(body)
	case Ack:
		processor.OnAccept()
	}
}
