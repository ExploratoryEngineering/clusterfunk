package cluster

// TODO: get rid of this (user supplied)
const numberOfShards = 10000

// NodeState is the enumeration of different states a node can be in.
type NodeState byte

// These are the (local) states the cluster node can be in
const (
	Invalid     NodeState = iota // Invalid or unknown state
	Joining                      // Joining the cluster
	Operational                  // Operational, normal operation
	Voting                       // Leader election in progress
	Resharding                   // Leader is elected, resharding in progress
	Starting                     // Starting the node
	Stopping                     // Stopping the node
)

// NodeRole is the roles the node can have in the cluster
type NodeRole byte

// These are the roles the node might have in the cluster
const (
	Unknown  NodeRole = iota // Uknown state
	Follower                 // A follower in a cluster
	Leader                   // The current leader node
	NonVoter                 // Non voting role in cluster
)

// Event is the interface for cluster events that are triggered
type Event struct {
	// State is the current cluster state
	LocalState NodeState
}

// Cluster is a wrapper for the Serf and Raft libraries. It will handle typical
// cluster operations.
type Cluster interface {

	// Name returns the cluster's name
	Name() string

	// Start launches the cluster, ie joins a Serf cluster and announces its
	// presence
	Start() error

	// Stop stops the cluster
	Stop()

	// Role is the current role of the node
	Role() NodeRole

	// State is the current cluster state
	LocalState() NodeState

	// Nodes return a list of the active nodes in the cluster
	Nodes() []Node

	// LeaderNode returns the leader of the cluster
	LeaderNode() Node

	// LocalNode returns the local node
	LocalNode() Node

	// Events returns an event channel for the cluster. The channel will
	// be closed when the cluster is stopped. Events are for information only
	Events() <-chan Event

	// AddLocalEndpoint registers an endpoint on the local node
	// TOOD(stalehd) add LocalNode() + AddEndpoint method to node
	AddLocalEndpoint(name, endpoint string)
}

// The following are internal tags and values for nodes
const (
	clusterEndpointPrefix = "ep."
	RaftNodeID            = "raft.nodeid"
	NodeType              = "kind"
	VoterKind             = "member"
	NonvoterKind          = "nonvoter"
	NodeRaftState         = "raft.state"
)

// The following is a list of well-known endpoints on nodes
const (
	//These are
	//MetricsEndpoint    = "ep.metrics"    // MetricsEndpoint is the metrics endpoint
	//HTTPEndpoint       = "ep.http"       // HTTPEndpoint is the HTTP endpoint
	SerfEndpoint       = "ep.serf"
	RaftEndpoint       = "ep.raft"
	LeaderEndpoint     = "ep.leader"
	ManagementEndpoint = "ep.management" //  gRPC endpoint for management
)

const (
	// StateLeader is the state reported in the Serf cluster tags
	StateLeader = "leader"
	// StateFollower is the state reported when the node is in the follower state
	StateFollower = "follower"
	// StateNone is the state reported when the node is in an unknown (raft) state
	StateNone = "none"
)
