package paxos

type Message struct {
	body interface{}
}

type Network interface {
	Send(sender *Processor, receiver *Processor, message Message)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SyncNetwork struct {
}

func (network *SyncNetwork) Send(sender *Processor, receiver *Processor, message Message) {
	receiver.Receive(sender, message)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type AsyncNetwork struct {
}

func (network *AsyncNetwork) Send(sender *Processor, receiver *Processor, message Message) {
	go receiver.Receive(sender, message)
}
