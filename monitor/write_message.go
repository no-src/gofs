package monitor

import "time"

type writeMessage struct {
	// name the file name
	name string
	// size the file size
	size int64
	// count the trigger count
	count int
	// last the last update time, unix nano
	last int64
	// cancel the current writeMessage is canceled or not
	cancel bool
}

func newWriteMessage(name string, size int64, count int, last int64) *writeMessage {
	return &writeMessage{
		name:   name,
		size:   size,
		count:  count,
		last:   last,
		cancel: false,
	}
}

func newDefaultWriteMessage(name string, size int64) *writeMessage {
	return newWriteMessage(name, size, 1, time.Now().UnixNano())
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
