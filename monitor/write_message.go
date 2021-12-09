package monitor

import "time"

type writeMessage struct {
	// name file name
	name string
	// count trigger count
	count int
	// last the last update time, unix nano
	last int64
}

func newWriteMessage(name string, count int, last int64) *writeMessage {
	return &writeMessage{
		name:  name,
		count: count,
		last:  last,
	}
}

func newDefaultWriteMessage(name string) *writeMessage {
	return newWriteMessage(name, 1, time.Now().UnixNano())
}

type writeMessageList []*writeMessage

func (list writeMessageList) Len() int {
	return len(list)
}

func (list writeMessageList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list writeMessageList) Less(i, j int) bool {
	if list[i].last < list[j].last {
		return true
	} else if list[i].last == list[j].last {
		return list[i].count > list[j].count
	} else {
		return false
	}
}
