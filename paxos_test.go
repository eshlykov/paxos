package paxos

import (
	"sync"
	"testing"
)

func initProcessors(processorCount int) []*Processor {
	processors := make([]*Processor, 0)
	for i := 0; i < processorCount; i++ {
		processors = append(processors, NewProcessor())
	}
	for _, processor := range processors {
		processor.AddAcceptors(processors)
	}
	return processors
}

func routine(wg *sync.WaitGroup, channel chan interface{}, processor *Processor, value interface{}) {
	channel <- processor.Decide(value)
	wg.Done()
}

func runProcessors(processors []*Processor, channel chan interface{}) {
	wg := new(sync.WaitGroup)

	for i, processor := range processors {
		wg.Add(1)
		go routine(wg, channel, processor, i)
	}

	wg.Wait()
	close(channel)
}

func validate(channel chan interface{}) bool {
	decided := <-channel
	for {
		value, ok := <-channel
		if !ok {
			return true
		}
		if value != decided {
			return false
		}
	}
}

func TestPaxos(testing *testing.T) {
	processorCount := 100
	processors := initProcessors(processorCount)

	channel := make(chan interface{}, processorCount)
	runProcessors(processors, channel)

	if !validate(channel) {
		testing.Errorf("Validity is broken!")
	}
}
