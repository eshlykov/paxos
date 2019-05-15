package paxos

import (
	"math/rand"
	"sync"
	"time"
)

type Message struct {
	body interface{}
}

type Network interface {
	Send(sender *Processor, receiver *Processor, message Message)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SyncNetwork struct {
}

func NewSyncNetwork() *SyncNetwork {
	return new(SyncNetwork)
}

func (network *SyncNetwork) Send(sender *Processor, receiver *Processor, message Message) {
	receiver.Receive(sender, message)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type AsyncNetwork struct {
}

func NewAsyncNetwork() *AsyncNetwork {
	return new(AsyncNetwork)
}

func (network *AsyncNetwork) Send(sender *Processor, receiver *Processor, message Message) {
	go receiver.Receive(sender, message)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SlowNetwork struct {
	random *rand.Rand
	lock   sync.Mutex
}

func NewSlowNetwork() *SlowNetwork {
	network := new(SlowNetwork)
	source := rand.NewSource(time.Now().UnixNano())
	network.random = rand.New(source)
	return network
}

func (network *SlowNetwork) Delay() {
	delay := 0.

	{
		network.lock.Lock()
		defer network.lock.Unlock()

		delay = 5 * network.random.Float64()
	}

	time.Sleep(time.Duration(delay) * time.Millisecond)
}

func (network *SlowNetwork) Send(sender *Processor, receiver *Processor, message Message) {
	go func() {
		network.Delay()
		receiver.Receive(sender, message)
	}()
}
