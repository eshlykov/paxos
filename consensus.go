package paxos

type Consensus interface {
	Decide(value interface{})
}

func (processor *Processor) Decide(value interface{}) interface{} {
	processor.value = value

	quorumNumber := (len(processor.acceptors) + 1) / 2

	for {
		processor.Prepare()
		if len(processor.promises) < quorumNumber {
			continue
		}

		processor.Propose()
		if processor.acceptCount >= quorumNumber {
			return processor.value
		}
	}
}
