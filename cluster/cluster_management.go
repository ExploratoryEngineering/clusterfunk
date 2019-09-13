package cluster

import (
	"time"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/stalehd/clusterfunk/cluster/clustermgmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// This is the cluster management service implementations

// getGRPCOpts returns gRPC server options for the configuration
func (cf *clusterfunkCluster) getGRPCOpts(config GRPCServerParameters) ([]grpc.ServerOption, error) {
	if !config.TLS {
		return []grpc.ServerOption{}, nil
	}
	if config.CertFile == "" || config.KeyFile == "" {
		return nil, errors.New("missing cert file and key file parameters for GRPC server")
	}
	creds, err := credentials.NewServerTLSFromFile(config.CertFile, config.KeyFile)
	if err != nil {
		return nil, err
	}
	return []grpc.ServerOption{grpc.Creds(creds)}, nil
}

func (cf *clusterfunkCluster) startManagementServices() error {
	opts, err := cf.getGRPCOpts(cf.config.Management)
	if err != nil {
		return err
	}
	cf.mgmtServer = grpc.NewServer(opts...)

	clustermgmt.RegisterClusterManagementServer(cf.mgmtServer, cf)

	listener, err := net.Listen("tcp", cf.config.Management.Endpoint)
	if err != nil {
		return err
	}

	fail := make(chan error)
	go func(ch chan error) {
		if err := cf.mgmtServer.Serve(listener); err != nil {
			log.Printf("Unable to launch node management gRPC server: %v", err)
			ch <- err
		}
	}(fail)

	select {
	case err:=<-fail:
		return err
	case <-time.After(250 * time.Millisecond):
		// ok
	}
	cf.setTag(ManagementEndpoint, listener.Addr().String())
	return nil
}

// Node management implementation
// -----------------------------------------------------------------------------

func (cf *clusterfunkCluster) GetState(context.Context, *clustermgmt.GetStateRequest) (*clustermgmt.GetStateResponse, error) {
	ret := &clustermgmt.GetStateResponse{
		NodeId:    cf.config.NodeID,
		RaftState: cf.ra.State().String(),
	}
	cfg := cf.ra.GetConfiguration()
	if cfg.Error() != nil {
		return nil, cfg.Error()
	}
	ret.RaftNodeCount = int32(len(cfg.Configuration().Servers))
	ret.SerfNodeCount = int32(len(cf.se.Members()))
	return ret, nil
}

func (cf *clusterfunkCluster) ListSerfNodes(context.Context, *clustermgmt.ListSerfNodesRequest) (*clustermgmt.ListSerfNodesResponse, error) {
	ret := &clustermgmt.ListSerfNodesResponse{
		NodeId: cf.config.NodeID,
	}
	members := cf.se.Members()
	ret.Swarm = make([]*clustermgmt.SerfNodeInfo, len(members))

	for i, v := range members {
		ret.Swarm[i] = &clustermgmt.SerfNodeInfo{
			Id:       v.Name,
			Endpoint: fmt.Sprintf("%s:%d", v.Addr.String(), v.Port),
			Status:   v.Status.String(),
		}
		for k, v := range v.Tags {
			if strings.HasPrefix(k, clusterEndpointPrefix) {
				ret.Swarm[i].ServiceEndpoints = append(ret.Swarm[i].ServiceEndpoints, &clustermgmt.SerfEndpoint{Name: k, HostPort: v})
			}
		}
	}
	return ret, nil
}

// Leader management implementation, ie all Raft-related functions not covered by the node management implementation
// -----------------------------------------------------------------------------
func (cf *clusterfunkCluster) ListRaftNodes(context.Context, *clustermgmt.ListRaftNodesRequest) (*clustermgmt.ListRaftNodesResponse, error) {
	ret := &clustermgmt.ListRaftNodesResponse{
		NodeId: cf.config.NodeID,
	}
	config := cf.ra.GetConfiguration()
	if err := config.Error(); err != nil {
		return nil, err
	}
	leader := cf.ra.Leader()

	members := config.Configuration().Servers
	ret.Members = make([]*clustermgmt.RaftNodeInfo, len(members))
	for i, v := range members {
		ret.Members[i] = &clustermgmt.RaftNodeInfo{
			Id:        string(v.ID),
			RaftState: v.Suffrage.String(),
			IsLeader:  (v.Address == leader),
		}
	}
	return ret, nil
}