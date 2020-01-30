package ctrlc

import (
	"context"
	"fmt"
	"os"

	"github.com/ExploratoryEngineering/clusterfunk/pkg/funk/clustermgmt"
)

// NodeCommand is the subcommand to add and remove nodes
type NodeCommand struct {
	Add addNodeCommand    `kong:"cmd,help='Add node to cluster'"`
	Rm  removeNodeCommand `kong:"cmd,help='Remove node from cluster'"`
	ID  string            `kong:"required,help='Node ID',short='N'"`
}

type addNodeCommand struct {
}

func (c *addNodeCommand) Run(param *Command) error {
	client := connectToManagement(param)
	if client == nil {
		return errStd
	}
	ctx, done := context.WithTimeout(context.Background(), gRPCTimeout)
	defer done()
	res, err := client.AddNode(ctx, &clustermgmt.AddNodeRequest{
		NodeId: param.Node.ID,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding node: %v\n", err)
		return errStd
	}
	if res.Error != nil {
		fmt.Fprintf(os.Stderr, "Leader could not add node: %v\n", res.Error.Message)
		return errStd
	}
	fmt.Printf("Node %s added to cluster\n", param.Node.ID)
	return nil
}

type removeNodeCommand struct {
}

func (c *removeNodeCommand) Run(param *Command) error {
	client := connectToManagement(param)
	if client == nil {
		return errStd
	}

	ctx, done := context.WithTimeout(context.Background(), gRPCTimeout)
	defer done()
	res, err := client.RemoveNode(ctx, &clustermgmt.RemoveNodeRequest{
		NodeId: param.Node.ID,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error removing node: %v\n", err)
		return errStd
	}
	if res.Error != nil {
		fmt.Fprintf(os.Stderr, "Leader could not remove node: %v\n", res.Error.Message)
		return errStd
	}
	fmt.Printf("Node %s removed from cluster\n", param.Node.ID)

	return nil
}
