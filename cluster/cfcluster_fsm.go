package cluster

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/stalehd/clusterfunk/cluster/clusterproto"

	"github.com/stalehd/clusterfunk/utils"
	"google.golang.org/grpc"
)

// Cluster mutations -- this is best expressed as a state machine. The state
// machine runs independently of the rest of the code. Mostly.

type internalFSMState int

// The FSM even is used to set the state with an optional node parameter
type fsmEvent struct {
	State  internalFSMState
	NodeID string
}

// These are the internal states. See implementation for a description. Most of these
// operations can be interrupted if there's a leader change midway in the process.
const (
	clusterSizeChanged internalFSMState = iota
	assumeLeadership
	assumeFollower
	ackReceived
	ackCompleted
	reshardCluster
	updateRaftNodeList
	newShardMapReceived
	commitLogReceived
	leaderLost
)

func (t internalFSMState) String() string {
	switch t {
	case clusterSizeChanged:
		return "clusterSizeChanged"
	case assumeLeadership:
		return "assumeLeadership"
	case assumeFollower:
		return "assumeFollower"
	case ackReceived:
		return "ackReceived"
	case ackCompleted:
		return "ackCompleted"
	case reshardCluster:
		return "reshardCluster"
	case updateRaftNodeList:
		return "updateRaftNodeList"
	case newShardMapReceived:
		return "newShardMapReceived"
	case commitLogReceived:
		return "commitLogReceived"
	case leaderLost:
		return "leaderLost"
	default:
		panic(fmt.Sprintf("Unknown state: %d", t))
	}
}

// Sets the new state unless a different state is waiting.
func (c *clusterfunkCluster) setFSMState(newState internalFSMState, nodeID string) {
	select {
	case c.stateChannel <- fsmEvent{State: newState, NodeID: nodeID}:
	case <-time.After(1 * time.Second):
		log.Printf("Unable to set cluster FSM state (%d) after 1 second", newState)
		// channel is already full - skip
	}
}

// clusterStateMachine is the FSM for
func (c *clusterfunkCluster) clusterStateMachine() {
	log.Printf("STATE: Launching")
	var unacknowledgedNodes []string
	oldState := leaderLost
	shardMapLogIndex := uint64(0)
	for newState := range c.stateChannel {

		timeCall(func() {
			log.Printf("STATE: %s", newState.State.String())
			switch newState.State {

			case assumeLeadership:
				c.setRole(Leader)
				// start with updating the node list
				c.setFSMState(updateRaftNodeList, "")

			case clusterSizeChanged:
				c.setLocalState(Resharding)
				c.setFSMState(updateRaftNodeList, "")

			case updateRaftNodeList:
				// list is updated. Start resharding.
				c.updateNodeList()
				c.setFSMState(reshardCluster, "")

			case reshardCluster:
				// reshard cluster, distribute via replicated log.

				// Reset the list of acked nodes.
				list := c.getNodes()
				// TODO: ShardManager needs a rewrite
				c.shardManager.UpdateNodes(list...)
				proposedShardMap, err := c.shardManager.MarshalBinary()
				if err != nil {
					panic(fmt.Sprintf("Can't marshal the shard map: %v", err))
				}
				mapMessage := NewLogMessage(ProposedShardMap, c.config.LeaderEndpoint, proposedShardMap)
				// Build list of unacked nodes
				// Note that this might include the local node as well, which
				// is OK. The client part will behave like all other parts.
				unacknowledgedNodes = append([]string{}, list...)

				// Replicate proposed shard map via log
				buf, err := mapMessage.MarshalBinary()
				if err != nil {
					panic(fmt.Sprintf("Unable to marshal the log message containing shard map: %v", err))
				}
				timeCall(func() {
					index, err := c.raftNode.AppendLogEntry(buf)
					if err != nil {
						// We might have lost the leadership here. Log and continue.
						if err := c.raftNode.ra.VerifyLeader().Error(); err == nil {
							panic("I'm the leader but I could not write the log")
						}
						// otherwise -- just log it and continue
						log.Printf("Could not write log entry for new shard map")
					}
					atomic.StoreUint64(c.reshardingLogIndex, index)
				}, "Appending shard map log entry")
				// This is the index we want commits for.
				shardMapLogIndex = c.raftNode.LastIndex()
				log.Printf("Shard map index = %d", shardMapLogIndex)
				//c.updateNodes(c.shardManager.NodeList())
				// Next messages will be ackReceived when the changes has replicated
				// out to the other nodes.
				// No new state here - wait for a series of ackReceived states
				// from the nodes.

			case ackReceived:
				// when a new ack message is received the ack is noted for the node and
				// until all nodes have acked the state will stay the same.
				for i, v := range unacknowledgedNodes {
					if v == newState.NodeID {
						unacknowledgedNodes = append(unacknowledgedNodes[:i], unacknowledgedNodes[i+1:]...)
					}
				}
				log.Printf("STATE: ack received from %s, %d nodes remaining", newState.NodeID, len(unacknowledgedNodes))
				// Timeouts are handled when calling the other nodes via gRPC
				allNodesHaveAcked := false
				if len(unacknowledgedNodes) == 0 {
					allNodesHaveAcked = true
				}

				if allNodesHaveAcked {
					c.setFSMState(ackCompleted, "")
				}
				return // continue

			case ackCompleted:
				// TODO: Log final commit message, establishing the new state in the cluster
				// ack is completed. Enable the new shard map for the cluster by
				// sending a commit log message. No further processing is required
				// here.
				commitMessage := NewLogMessage(ShardMapCommitted, c.raftNode.LocalNodeID(), []byte{})
				buf, err := commitMessage.MarshalBinary()
				if err != nil {
					panic(fmt.Sprintf("Unable to marshal commit message: %v", err))
				}
				if _, err := c.raftNode.AppendLogEntry(buf); err != nil {
					// We might have lost the leadership here. Panic if we're still
					// the leader
					if err := c.raftNode.ra.VerifyLeader().Error(); err == nil {
						panic("I'm the leader but I could not write the log")
					}
					// otherwise -- just log it and continue
					log.Printf("Could not write log entry for new shard map: %v", err)
				}
				c.setLocalState(Operational)
				return // continue

			case assumeFollower:
				c.setRole(Follower)
				c.setLocalState(Resharding)
				// Not much happens here but the next state should be - if all
				// goes well - a shard map log message from the leader.

			case newShardMapReceived:
				// update internal map and ack map via gRPC
				c.setLocalState(Resharding)

				// No new state - the next is commitLogReceived which is set
				// via the replicated log events

			case commitLogReceived:
				// commit log received, set state operational and resume normal
				// operations. Signal to the rest of the library (channel)
				c.setLocalState(Operational)

			case leaderLost:
				c.setLocalState(Voting)
				// leader is lost - stop processing until a leader is elected and
				// the commit log is received
			}
		}, fmt.Sprintf("STATE: handle %s -> %s state transition", oldState.String(), newState.State.String()))
		oldState = newState.State

	}
}

func (c *clusterfunkCluster) ackShardMap(endpoint string) {

	// Step 1 Leader ID
	// Confirm the shard map
	clientParam := utils.GRPCClientParam{
		ServerEndpoint: endpoint,
		TLS:            false,
		CAFile:         "",
	}
	opts, err := utils.GetGRPCDialOpts(clientParam)
	if err != nil {
		//panic(fmt.Sprintf("Unable to acknowledge gRPC client parameters: %v", err))
		log.Printf("Unable to acknowledge gRPC client parameters: %v", err)
		return
	}
	conn, err := grpc.Dial(clientParam.ServerEndpoint, opts...)
	if err != nil {
		//panic(fmt.Sprintf("Unable to dial server when acking shard map: %v", err))
		log.Printf("Unable to dial server when acking shard map: %v", err)
		return
	}
	defer conn.Close()
	logIndex := atomic.LoadUint64(c.reshardingLogIndex)

	client := clusterproto.NewClusterLeaderServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	resp, err := client.ConfirmShardMap(ctx, &clusterproto.ConfirmShardMapRequest{
		NodeID:   c.raftNode.LocalNodeID(),
		LogIndex: int64(logIndex),
	})
	if err != nil {
		//panic(fmt.Sprintf("Unable to confirm shard map: %v", err))
		log.Printf("Unable to confirm shard map: %v", err)
		return
	}
	if !resp.Success {
		log.Printf("Leader rejected ack. I got index %d and leader wants %d", logIndex, resp.CurrentIndex)
		atomic.StoreUint64(c.wantedShardLogIndex, uint64(resp.CurrentIndex))
		return
	}
	log.Printf("Shard map ack successfully sent to leader (leader=%s, index=%d)", endpoint, logIndex)
	atomic.StoreUint64(c.wantedShardLogIndex, 0)
}