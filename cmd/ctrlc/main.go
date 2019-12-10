package main
//
//Copyright 2019 Telenor Digital AS
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/ExploratoryEngineering/params"
	"github.com/ExploratoryEngineering/clusterfunk/pkg/clientfunk"
	"github.com/ExploratoryEngineering/clusterfunk/pkg/funk/clustermgmt"
	"github.com/ExploratoryEngineering/clusterfunk/pkg/toolbox"
	"google.golang.org/grpc"
)

type parameters struct {
	ClusterName string `param:"desc=Cluster name;default=clusterfunk"`
	Zeroconf    bool   `param:"desc=Use zeroconf discovery for Serf;default=true"`
	Endpoint    string `param:"desc=Management endpoint to use"`
}

const (
	cmdStatus     = "status"
	cmdNodes      = "nodes"
	cmdEndpoints  = "endpoints"
	cmdAddNode    = "add-node"
	cmdRemoveNode = "remove-node"
	cmdShards     = "shards"
	cmdStepDown   = "step-down"
	cmdHelp       = "help"
)

func main() {
	rootArgs := []string{}
	remainingArgs := []string{}
	for i, v := range os.Args[1:] {
		if v[0] == '-' {
			rootArgs = append(rootArgs, v)
		} else {
			remainingArgs = append(remainingArgs, os.Args[i+1:]...)
			break
		}

	}
	var config parameters
	if err := params.NewEnvFlag(&config, rootArgs); err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(remainingArgs) == 0 {
		fmt.Fprintf(os.Stderr, "No command specfied\n")
		return
	}
	if remainingArgs[0] == cmdHelp {
		fmt.Println(`
ctrlc [--cluster-name] [--zeroconf] [--endpoint] [cmd]

	--cluster-name: Name of cluster
	--zeroconf     Enable or disable zeroconf
		--endpoint     Management endpoint to use

Available commands:

	status           Show node and cluster status
	nodes            List nodes in cluster
	endpoints [name] List endpoints
	add-node [id]    Add node to cluster
	remove-node [id] Remove node from cluster
	shards           List shard distribution
	step-down        Leader step down
	`)
		return
	}
	if config.Endpoint == "" && config.Zeroconf {
		if config.Zeroconf && config.ClusterName == "" {
			fmt.Fprintf(os.Stderr, "Needs a cluster name if zeroconf is to be used for discovery")
			return
		}
		ep, err := clientfunk.ZeroconfManagementLookup(config.ClusterName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to do zeroconf lookup: %v\n", err)
			return
		}
		config.Endpoint = ep
	}

	if config.Endpoint == "" {
		fmt.Fprintf(os.Stderr, "Need an endpoint for one of the cluster nodes")
		return
	}

	grpcParams := toolbox.GRPCClientParam{
		ServerEndpoint: config.Endpoint,
	}
	opts, err := toolbox.GetGRPCDialOpts(grpcParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create GRPC dial options: %v\n", err)
		return
	}
	conn, err := grpc.Dial(config.Endpoint, opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not dial management endpoint: %v\n", err)
		return
	}
	client := clustermgmt.NewClusterManagementClient(conn)

	switch remainingArgs[0] {
	case cmdStatus:
		status(client)

	case cmdAddNode:
		if len(remainingArgs) != 2 {
			fmt.Fprintf(os.Stderr, "Need ID for %s\n", cmdAddNode)
			return
		}
		addNode(remainingArgs[1], client)

	case cmdRemoveNode:
		if len(remainingArgs) != 2 {
			fmt.Fprintf(os.Stderr, "Need ID for %s\n", cmdRemoveNode)
			return
		}
		removeNode(remainingArgs[1], client)

	case cmdEndpoints:
		epFilter := ""
		if len(remainingArgs) > 1 {
			epFilter = remainingArgs[1]
		}
		listEndpoints(epFilter, client)

	case cmdNodes:
		listNodes(client)

	case cmdShards:
		listShards(client)

	case cmdStepDown:
		stepDown(client)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", remainingArgs[0])
		return
	}
}

const (
	gRPCTimeout = 1 * time.Second
)

func status(client clustermgmt.ClusterManagementClient) {
	ctx, done := context.WithTimeout(context.Background(), gRPCTimeout)
	defer done()
	res, err := client.GetStatus(ctx, &clustermgmt.GetStatusRequest{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving status: %v\n", err)
		return
	}
	fmt.Printf("Cluster name: %s\n", res.ClusterName)
	fmt.Printf("Node ID:      %s\n", res.LocalNodeId)
	fmt.Printf("State:        %s\n", res.LocalState)
	fmt.Printf("Role:         %s\n", res.LocalRole)
	fmt.Printf("Leader ID:    %s\n", res.LeaderNodeId)
	fmt.Printf("Nodes:        %d Raft, %d Serf\n", res.RaftNodeCount, res.SerfNodeCount)
	fmt.Printf("Shards:       %d (total weight: %d)\n\n", res.ShardCount, res.ShardWeight)
}

func addNode(id string, client clustermgmt.ClusterManagementClient) {
	ctx, done := context.WithTimeout(context.Background(), gRPCTimeout)
	defer done()
	res, err := client.AddNode(ctx, &clustermgmt.AddNodeRequest{
		NodeId: id,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding node: %v\n", err)
		return
	}
	if res.Error != nil {
		fmt.Fprintf(os.Stderr, "Leader could not add node: %v\n", res.Error.Message)
		return
	}
	fmt.Printf("Node %s added to cluster\n", id)
}

func removeNode(id string, client clustermgmt.ClusterManagementClient) {
	ctx, done := context.WithTimeout(context.Background(), gRPCTimeout)
	defer done()
	res, err := client.RemoveNode(ctx, &clustermgmt.RemoveNodeRequest{
		NodeId: id,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error removing node: %v\n", err)
		return
	}
	if res.Error != nil {
		fmt.Fprintf(os.Stderr, "Leader could not remove node: %v\n", res.Error.Message)
		return
	}
	fmt.Printf("Node %s removed from cluster\n", id)
}

func listEndpoints(filter string, client clustermgmt.ClusterManagementClient) {
	ctx, done := context.WithTimeout(context.Background(), gRPCTimeout)
	defer done()

	res, err := client.FindEndpoint(ctx, &clustermgmt.EndpointRequest{EndpointName: filter})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching for endpoint: %v\n", err)
		return
	}
	if res.Error != nil {
		fmt.Fprintf(os.Stderr, "Unable to search for endpoint: %v\n", res.Error.Message)
		return
	}

	fmt.Printf("Node ID              Name                 Endpoint\n")

	sort.Slice(res.Endpoints, func(i, j int) bool {
		return res.Endpoints[i].Name < res.Endpoints[j].Name
	})
	for _, v := range res.Endpoints {
		fmt.Printf("%-20s %-20s %s\n", v.NodeId, v.Name, v.HostPort)
	}
	fmt.Printf("\nReporting node: %s\n", res.NodeId)
}

func listNodes(client clustermgmt.ClusterManagementClient) {
	ctx, done := context.WithTimeout(context.Background(), gRPCTimeout)
	defer done()

	res, err := client.ListNodes(ctx, &clustermgmt.ListNodesRequest{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing nodes: %v\n", err)
		return
	}
	if res.Error != nil {
		fmt.Fprintf(os.Stderr, "Unable to list nodes: %v\n", res.Error.Message)
		return
	}

	fmt.Printf("  Node ID              Raft       Serf\n")
	for _, v := range res.Nodes {
		leader := ""
		if v.Leader {
			leader += "*"
		}
		fmt.Printf("%-2s%-20s %-10s %s\n", leader, v.NodeId, v.RaftState, v.SerfState)
	}
	fmt.Printf("\nReporting node: %s   Leader node: %s\n", res.NodeId, res.LeaderId)
}

func listShards(client clustermgmt.ClusterManagementClient) {
	ctx, done := context.WithTimeout(context.Background(), gRPCTimeout)
	defer done()

	res, err := client.ListShards(ctx, &clustermgmt.ListShardsRequest{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing shards: %v\n", err)
		return
	}
	if res.Error != nil {
		fmt.Fprintf(os.Stderr, "Unable to list shards: %v\n", res.Error.Message)
		return
	}

	fmt.Println("Node ID              Shards             Weight")
	for _, v := range res.Shards {
		shardPct := float32(v.ShardCount) / float32(res.TotalShards) * 100.0
		weightPct := float32(v.ShardWeight) / float32(res.TotalWeight) * 100.0
		fmt.Printf("%-20s %10d (%3.1f%%) %10d (%3.1f%%)\n", v.NodeId, v.ShardCount, shardPct, v.ShardWeight, weightPct)
	}
	fmt.Printf("\nReporting node: %s    Total shards: %d    Total weight: %d\n", res.NodeId, res.TotalShards, res.TotalWeight)
}

func stepDown(client clustermgmt.ClusterManagementClient) {
	ctx, done := context.WithTimeout(context.Background(), gRPCTimeout)
	defer done()

	res, err := client.StepDown(ctx, &clustermgmt.StepDownRequest{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error asking leader to step down: %v\n", err)
		return
	}
	if res.Error != nil {
		fmt.Fprintf(os.Stderr, "Leader is unable to step down: %v\n", res.Error.Message)
		return
	}
	fmt.Printf("Leader node %s has stepped down\n", res.NodeId)
}
