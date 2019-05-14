package paxos

type Message struct {
	body interface{}
}

func Send(sender *Processor, receiver *Processor, message Message) {
	receiver.Receive(sender, message)
}
