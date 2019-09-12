package cluster

import (
	"errors"
	"fmt"
	"sync"
)

type nodeData struct {
	NodeID       string
	TotalWeights int
	Shards       []Shard
}

func newNodeData(nodeID string) *nodeData {
	return &nodeData{NodeID: nodeID, TotalWeights: 0, Shards: make([]Shard, 0)}
}

func (nd *nodeData) AddShard(shard Shard) {
	shard.SetNodeID(nd.NodeID)
	nd.TotalWeights += shard.Weight()
	nd.Shards = append(nd.Shards, shard)
}

func (nd *nodeData) RemoveShard(preferredWeight int) Shard {
	for i, v := range nd.Shards {
		if v.Weight() <= preferredWeight {
			nd.Shards = append(nd.Shards[:i], nd.Shards[i+1:]...)
			nd.TotalWeights -= v.Weight()
			return v
		}
	}
	if len(nd.Shards) == 0 {
		panic("no shards remaining")
	}
	// This will cause a panic if there's no shards left. That's OK.
	// Future me might disagree.
	ret := nd.Shards[0]
	if len(nd.Shards) > 1 {
		nd.Shards = nd.Shards[1:]
	}
	nd.TotalWeights -= ret.Weight()
	return ret
}

type weightedShardManager struct {
	shards      []Shard
	mutex       *sync.Mutex
	totalWeight int
	nodes       map[string]*nodeData
}

// NewShardManager creates a new ShardManager instance.
func NewShardManager() ShardManager {
	return &weightedShardManager{
		shards:      make([]Shard, 0),
		mutex:       &sync.Mutex{},
		totalWeight: 0,
		nodes:       make(map[string]*nodeData),
	}
}

func (sm *weightedShardManager) Init(maxShards int, weights []int) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if maxShards < 1 {
		return errors.New("maxShards must be > 0")
	}
	if len(sm.shards) != 0 {
		return errors.New("shards already set")
	}

	if len(weights) != maxShards {
		return errors.New("maxShards and len(weights) must be the same")
	}

	sm.shards = make([]Shard, maxShards)
	for i := range sm.shards {
		weight := 1
		if weights != nil {
			weight = weights[i]
		}
		sm.totalWeight += weight
		sm.shards[i] = NewShard(i, weight)
	}
	return nil
}

func (sm *weightedShardManager) AddNode(nodeID string) []ShardTransfer {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	newNode := newNodeData(nodeID)
	ret := make([]ShardTransfer, 0)

	// Invariant: First node
	if len(sm.nodes) == 0 {
		for i := range sm.shards {
			newNode.AddShard(sm.shards[i])
			ret = append(ret, ShardTransfer{
				Shard:             sm.shards[i],
				SourceNodeID:      "",
				DestinationNodeID: newNode.NodeID,
			})
		}
		sm.nodes[nodeID] = newNode
		return ret
	}

	//Invariant: Node # 2 or later
	targetCount := sm.totalWeight / (len(sm.nodes) + 1)

	for k, v := range sm.nodes {
		for v.TotalWeights > targetCount && v.TotalWeights > 0 {
			shardToMove := v.RemoveShard(targetCount - v.TotalWeights)
			newNode.AddShard(shardToMove)
			ret = append(ret, ShardTransfer{
				Shard:             shardToMove,
				SourceNodeID:      v.NodeID,
				DestinationNodeID: newNode.NodeID,
			})
		}
		sm.nodes[k] = v
	}
	sm.nodes[nodeID] = newNode
	return ret
}

func (sm *weightedShardManager) RemoveNode(nodeID string) []ShardTransfer {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	nodeToRemove, exists := sm.nodes[nodeID]
	if !exists {
		panic(fmt.Sprintf("Unknown node ID: %s", nodeID))
	}
	delete(sm.nodes, nodeID)
	// Invariant: This is the last node in the cluster. No point in
	// generating transfers
	if len(sm.nodes) == 0 {
		return []ShardTransfer{}
	}

	ret := make([]ShardTransfer, 0)
	targetCount := sm.totalWeight / len(sm.nodes)
	for k, v := range sm.nodes {
		//		fmt.Printf("Removing node %s: Node %s w=%d target=%d\n", nodeID, k, v.TotalWeights, targetCount)
		for v.TotalWeights < targetCount && nodeToRemove.TotalWeights > 0 {
			shardToMove := nodeToRemove.RemoveShard(targetCount - v.TotalWeights)
			v.AddShard(shardToMove)
			ret = append(ret, ShardTransfer{
				Shard:             shardToMove,
				SourceNodeID:      nodeToRemove.NodeID,
				DestinationNodeID: v.NodeID,
			})
		}
		sm.nodes[k] = v
	}
	return ret
}

func (sm *weightedShardManager) MapToNode(shardID int) Shard {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if shardID > len(sm.shards) || shardID < 0 {
		// This might be too extreme but useful for debugging.
		// another alternative is to return a catch-all node allowing
		// the proxying to fix it but if the shard ID is invalid it is
		// probably an error with the shard function itself and warrants
		// a panic() from the library.
		panic(fmt.Sprintf("shard ID is outside range [0-%d]: %d", len(sm.shards), shardID))
	}
	return sm.shards[shardID]
}

func (sm *weightedShardManager) Shards() []Shard {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	return sm.shards[:]
}

func (sm *weightedShardManager) TotalWeight() int {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	return sm.totalWeight
}

func (sm *weightedShardManager) MarshalBinary() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (sm *weightedShardManager) UnmarshalBinary(data []byte) error {
	return errors.New("not implemented")
}