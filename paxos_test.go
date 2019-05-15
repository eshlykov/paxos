package paxos

import (
	"errors"
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

func validate(proposed, decieded chan interface{}) error {
	p := make(map[interface{}]bool)
	for value := range proposed {
		p[value] = true
	}

	d := make(map[interface{}]bool)
	for value := range decieded {
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

func testPaxos(processorCount int) error {
	processors := initProcessors(processorCount)

	proposed := make(chan interface{}, processorCount)
	decided := make(chan interface{}, processorCount)

	runProcessors(processors, proposed, decided)

	return validate(proposed, decided)
}

func TestPaxosUnitTest(testing *testing.T) {
	if ok := testPaxos(1); ok != nil {
		testing.Errorf("UnitTest(): %s", ok.Error())
	}
}

func TestPaxosStressTest(testing *testing.T) {
	for _, processorCount := range []int{1, 2, 5, 10, 25} {
		for i := 0; i < 100*processorCount; i++ {
			if ok := testPaxos(processorCount); ok != nil {
				testing.Errorf("StressTest(processorCount: %d): %s", processorCount, ok.Error())
			}
		}
	}
}
