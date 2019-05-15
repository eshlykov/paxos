package paxos

import (
	"errors"
	"sync"
	"testing"
)

func initProcessors(network Network, processorCount int) []*Processor {
	processors := make([]*Processor, 0)
	for i := 0; i < processorCount; i++ {
		processors = append(processors, NewProcessor(network))
	}
	for _, processor := range processors {
		processor.AddAcceptors(processors)
	}
	return processors
}

func routine(wg *sync.WaitGroup, proposed, decided chan interface{}, processor *Processor, value interface{}) {
	proposed <- value
	decided <- processor.Decide(value)
	wg.Done()
}

func runProcessors(processors []*Processor, proposed, decided chan interface{}) {
	wg := new(sync.WaitGroup)

	for i, processor := range processors {
		wg.Add(1)
		go routine(wg, proposed, decided, processor, i)
	}

	wg.Wait()
	close(proposed)
	close(decided)
}

func validate(proposed, decided chan interface{}) error {
	p := make(map[interface{}]bool)
	for value := range proposed {
		p[value] = true
	}

	d := make(map[interface{}]bool)
	for value := range decided {
		if _, ok := p[value]; !ok {
			return errors.New("validity is violated")
		}
		d[value] = true
	}

	if len(d) > 1 {
		return errors.New("agreement is violated")
	}

	return nil
}

func testPaxos(network Network, processorCount int) error {
	processors := initProcessors(network, processorCount)

	proposed := make(chan interface{}, processorCount)
	decided := make(chan interface{}, processorCount)

	runProcessors(processors, proposed, decided)

	return validate(proposed, decided)
}

func TestPaxosUnitTest(testing *testing.T) {
	if ok := testPaxos(NewSyncNetwork(), 1); ok != nil {
		testing.Errorf("UnitTest(): %s", ok.Error())
	}
}

func TestPaxosStressTest(testing *testing.T) {
	for _, network := range []Network{NewSyncNetwork(), NewAsyncNetwork(), NewSlowNetwork()} {
		for _, processorCount := range []int{1, 3, 5} {
			for i := 0; i < 8; i++ {
				if ok := testPaxos(network, processorCount); ok != nil {
					testing.Errorf("StressTest(network: %v, processorCount: %d): %s",
						network, processorCount, ok.Error())
				}
			}
		}
	}
}
