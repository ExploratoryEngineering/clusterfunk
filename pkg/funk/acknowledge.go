package funk

import (
	"sync/atomic"
	"time"

	"github.com/stalehd/clusterfunk/pkg/toolbox"
)

// ackCollection handles acknowledgement and timeouts on acknowledgements
// The collection can be acked one time only and
type ackCollection interface {
	// StartAck starts the acking
	StartAck(nodes []string, shardIndex uint64, timeout time.Duration)

	// Ack adds another node to the acknowledged list. Returns true if that node
	// is in the list of nodes that should ack
	Ack(nodeID string, shardIndex uint64) bool

	// ShardIndex returns the shard index that is acked. This returns 0 when
	// the ack is completed.
	ShardIndex() uint64

	// MissingAck returns a channel that sends a list of nodes that haven't acknowledged within the timeout.
	// If something is sent on the MissingAck channel the Completed channel won't trigger.
	MissingAck() <-chan []string

	// Completed returns a channel that is triggered when all nodes have
	// acknowledged
	Completed() <-chan struct{}

	// Done clears up the resources and closes all open channels
	Done()
}

func newAckCollection() ackCollection {
	return &ackColl{
		nodes:         toolbox.NewStringSet(),
		completedChan: make(chan struct{}),
		missingChan:   make(chan []string),
		shardIndex:    new(uint64),
	}
}

type ackColl struct {
	nodes         toolbox.StringSet
	completedChan chan struct{}
	missingChan   chan []string
	shardIndex    *uint64
}

func (a *ackColl) StartAck(nodes []string, shardIndex uint64, timeout time.Duration) {
	atomic.StoreUint64(a.shardIndex, shardIndex)
	a.nodes.Sync(nodes...)
	go func() {
		time.Sleep(timeout)
		if a.nodes.Size() > 0 {
			a.missingChan <- a.nodes.List()
		}
	}()
}

func (a *ackColl) Ack(nodeID string, shardIndex uint64) bool {
	if atomic.LoadUint64(a.shardIndex) != shardIndex {
		return false
	}
	if a.nodes.Remove(nodeID) {
		if a.nodes.Size() == 0 {
			go func() { a.completedChan <- struct{}{} }()
		}
		return true
	}
	return false
}

func (a *ackColl) MissingAck() <-chan []string {
	return a.missingChan
}

func (a *ackColl) Completed() <-chan struct{} {
	return a.completedChan
}

func (a *ackColl) Done() {
	close(a.missingChan)
	close(a.completedChan)
	atomic.StoreUint64(a.shardIndex, 0)
}

func (a *ackColl) ShardIndex() uint64 {
	return atomic.LoadUint64(a.shardIndex)
}
