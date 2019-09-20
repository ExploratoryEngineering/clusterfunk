// Code generated by protoc-gen-go. DO NOT EDIT.
// source: management.proto

package clustermgmt

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GetStateResponse_ClusterState int32

const (
	GetStateResponse_OK     GetStateResponse_ClusterState = 0
	GetStateResponse_VOTING GetStateResponse_ClusterState = 1
)

var GetStateResponse_ClusterState_name = map[int32]string{
	0: "OK",
	1: "VOTING",
}

var GetStateResponse_ClusterState_value = map[string]int32{
	"OK":     0,
	"VOTING": 1,
}

func (x GetStateResponse_ClusterState) String() string {
	return proto.EnumName(GetStateResponse_ClusterState_name, int32(x))
}

func (GetStateResponse_ClusterState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_edc174f991dc0a25, []int{1, 0}
}

type NodeInfo_NodeKind int32

const (
	NodeInfo_LEADER    NodeInfo_NodeKind = 0
	NodeInfo_VOTER     NodeInfo_NodeKind = 1
	NodeInfo_NONVOTER  NodeInfo_NodeKind = 2
	NodeInfo_NONMEMBER NodeInfo_NodeKind = 3
)

var NodeInfo_NodeKind_name = map[int32]string{
	0: "LEADER",
	1: "VOTER",
	2: "NONVOTER",
	3: "NONMEMBER",
}

var NodeInfo_NodeKind_value = map[string]int32{
	"LEADER":    0,
	"VOTER":     1,
	"NONVOTER":  2,
	"NONMEMBER": 3,
}

func (x NodeInfo_NodeKind) String() string {
	return proto.EnumName(NodeInfo_NodeKind_name, int32(x))
}

func (NodeInfo_NodeKind) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_edc174f991dc0a25, []int{3, 0}
}

// NodeState is the state of nodes. See godoc comments for a description
type NodeInfo_NodeState int32

const (
	NodeInfo_INITIALIZED  NodeInfo_NodeState = 0
	NodeInfo_READY        NodeInfo_NodeState = 1
	NodeInfo_SERVING      NodeInfo_NodeState = 2
	NodeInfo_REORGANIZING NodeInfo_NodeState = 3
	NodeInfo_DRAINING     NodeInfo_NodeState = 4
	NodeInfo_TERMINATING  NodeInfo_NodeState = 5
)

var NodeInfo_NodeState_name = map[int32]string{
	0: "INITIALIZED",
	1: "READY",
	2: "SERVING",
	3: "REORGANIZING",
	4: "DRAINING",
	5: "TERMINATING",
}

var NodeInfo_NodeState_value = map[string]int32{
	"INITIALIZED":  0,
	"READY":        1,
	"SERVING":      2,
	"REORGANIZING": 3,
	"DRAINING":     4,
	"TERMINATING":  5,
}

func (x NodeInfo_NodeState) String() string {
	return proto.EnumName(NodeInfo_NodeState_name, int32(x))
}

func (NodeInfo_NodeState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_edc174f991dc0a25, []int{3, 1}
}

type GetStateRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetStateRequest) Reset()         { *m = GetStateRequest{} }
func (m *GetStateRequest) String() string { return proto.CompactTextString(m) }
func (*GetStateRequest) ProtoMessage()    {}
func (*GetStateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_edc174f991dc0a25, []int{0}
}

func (m *GetStateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetStateRequest.Unmarshal(m, b)
}
func (m *GetStateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetStateRequest.Marshal(b, m, deterministic)
}
func (m *GetStateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetStateRequest.Merge(m, src)
}
func (m *GetStateRequest) XXX_Size() int {
	return xxx_messageInfo_GetStateRequest.Size(m)
}
func (m *GetStateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetStateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetStateRequest proto.InternalMessageInfo

type GetStateResponse struct {
	NodeId               string                        `protobuf:"bytes,1,opt,name=NodeId,proto3" json:"NodeId,omitempty"`
	State                GetStateResponse_ClusterState `protobuf:"varint,2,opt,name=State,proto3,enum=clustermgmt.GetStateResponse_ClusterState" json:"State,omitempty"`
	NodeCount            int32                         `protobuf:"varint,3,opt,name=NodeCount,proto3" json:"NodeCount,omitempty"`
	VoterCount           int32                         `protobuf:"varint,4,opt,name=VoterCount,proto3" json:"VoterCount,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *GetStateResponse) Reset()         { *m = GetStateResponse{} }
func (m *GetStateResponse) String() string { return proto.CompactTextString(m) }
func (*GetStateResponse) ProtoMessage()    {}
func (*GetStateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_edc174f991dc0a25, []int{1}
}

func (m *GetStateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetStateResponse.Unmarshal(m, b)
}
func (m *GetStateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetStateResponse.Marshal(b, m, deterministic)
}
func (m *GetStateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetStateResponse.Merge(m, src)
}
func (m *GetStateResponse) XXX_Size() int {
	return xxx_messageInfo_GetStateResponse.Size(m)
}
func (m *GetStateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetStateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetStateResponse proto.InternalMessageInfo

func (m *GetStateResponse) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *GetStateResponse) GetState() GetStateResponse_ClusterState {
	if m != nil {
		return m.State
	}
	return GetStateResponse_OK
}

func (m *GetStateResponse) GetNodeCount() int32 {
	if m != nil {
		return m.NodeCount
	}
	return 0
}

func (m *GetStateResponse) GetVoterCount() int32 {
	if m != nil {
		return m.VoterCount
	}
	return 0
}

type EndpointInfo struct {
	Name                 string   `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	HostPort             string   `protobuf:"bytes,2,opt,name=HostPort,proto3" json:"HostPort,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EndpointInfo) Reset()         { *m = EndpointInfo{} }
func (m *EndpointInfo) String() string { return proto.CompactTextString(m) }
func (*EndpointInfo) ProtoMessage()    {}
func (*EndpointInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_edc174f991dc0a25, []int{2}
}

func (m *EndpointInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EndpointInfo.Unmarshal(m, b)
}
func (m *EndpointInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EndpointInfo.Marshal(b, m, deterministic)
}
func (m *EndpointInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EndpointInfo.Merge(m, src)
}
func (m *EndpointInfo) XXX_Size() int {
	return xxx_messageInfo_EndpointInfo.Size(m)
}
func (m *EndpointInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_EndpointInfo.DiscardUnknown(m)
}

var xxx_messageInfo_EndpointInfo proto.InternalMessageInfo

func (m *EndpointInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EndpointInfo) GetHostPort() string {
	if m != nil {
		return m.HostPort
	}
	return ""
}

// NodeInfo contains information on a single node in the cluster/swarm
type NodeInfo struct {
	// NodeId is the node's identifier.
	NodeId               string             `protobuf:"bytes,1,opt,name=NodeId,proto3" json:"NodeId,omitempty"`
	Kind                 NodeInfo_NodeKind  `protobuf:"varint,2,opt,name=Kind,proto3,enum=clustermgmt.NodeInfo_NodeKind" json:"Kind,omitempty"`
	State                NodeInfo_NodeState `protobuf:"varint,3,opt,name=State,proto3,enum=clustermgmt.NodeInfo_NodeState" json:"State,omitempty"`
	Endpoints            []*EndpointInfo    `protobuf:"bytes,4,rep,name=Endpoints,proto3" json:"Endpoints,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *NodeInfo) Reset()         { *m = NodeInfo{} }
func (m *NodeInfo) String() string { return proto.CompactTextString(m) }
func (*NodeInfo) ProtoMessage()    {}
func (*NodeInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_edc174f991dc0a25, []int{3}
}

func (m *NodeInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeInfo.Unmarshal(m, b)
}
func (m *NodeInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeInfo.Marshal(b, m, deterministic)
}
func (m *NodeInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeInfo.Merge(m, src)
}
func (m *NodeInfo) XXX_Size() int {
	return xxx_messageInfo_NodeInfo.Size(m)
}
func (m *NodeInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeInfo.DiscardUnknown(m)
}

var xxx_messageInfo_NodeInfo proto.InternalMessageInfo

func (m *NodeInfo) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *NodeInfo) GetKind() NodeInfo_NodeKind {
	if m != nil {
		return m.Kind
	}
	return NodeInfo_LEADER
}

func (m *NodeInfo) GetState() NodeInfo_NodeState {
	if m != nil {
		return m.State
	}
	return NodeInfo_INITIALIZED
}

func (m *NodeInfo) GetEndpoints() []*EndpointInfo {
	if m != nil {
		return m.Endpoints
	}
	return nil
}

// ListSerfNodesResponse is the response to a ListSerfNodes call. The number of
// Serf nodes will initially be quite small so we won't need a stream response
// here.
type ListNodesResponse struct {
	NodeId               string      `protobuf:"bytes,1,opt,name=NodeId,proto3" json:"NodeId,omitempty"`
	Nodes                []*NodeInfo `protobuf:"bytes,2,rep,name=Nodes,proto3" json:"Nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ListNodesResponse) Reset()         { *m = ListNodesResponse{} }
func (m *ListNodesResponse) String() string { return proto.CompactTextString(m) }
func (*ListNodesResponse) ProtoMessage()    {}
func (*ListNodesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_edc174f991dc0a25, []int{4}
}

func (m *ListNodesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListNodesResponse.Unmarshal(m, b)
}
func (m *ListNodesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListNodesResponse.Marshal(b, m, deterministic)
}
func (m *ListNodesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListNodesResponse.Merge(m, src)
}
func (m *ListNodesResponse) XXX_Size() int {
	return xxx_messageInfo_ListNodesResponse.Size(m)
}
func (m *ListNodesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListNodesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListNodesResponse proto.InternalMessageInfo

func (m *ListNodesResponse) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *ListNodesResponse) GetNodes() []*NodeInfo {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type ListNodesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListNodesRequest) Reset()         { *m = ListNodesRequest{} }
func (m *ListNodesRequest) String() string { return proto.CompactTextString(m) }
func (*ListNodesRequest) ProtoMessage()    {}
func (*ListNodesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_edc174f991dc0a25, []int{5}
}

func (m *ListNodesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListNodesRequest.Unmarshal(m, b)
}
func (m *ListNodesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListNodesRequest.Marshal(b, m, deterministic)
}
func (m *ListNodesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListNodesRequest.Merge(m, src)
}
func (m *ListNodesRequest) XXX_Size() int {
	return xxx_messageInfo_ListNodesRequest.Size(m)
}
func (m *ListNodesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListNodesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListNodesRequest proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("clustermgmt.GetStateResponse_ClusterState", GetStateResponse_ClusterState_name, GetStateResponse_ClusterState_value)
	proto.RegisterEnum("clustermgmt.NodeInfo_NodeKind", NodeInfo_NodeKind_name, NodeInfo_NodeKind_value)
	proto.RegisterEnum("clustermgmt.NodeInfo_NodeState", NodeInfo_NodeState_name, NodeInfo_NodeState_value)
	proto.RegisterType((*GetStateRequest)(nil), "clustermgmt.GetStateRequest")
	proto.RegisterType((*GetStateResponse)(nil), "clustermgmt.GetStateResponse")
	proto.RegisterType((*EndpointInfo)(nil), "clustermgmt.EndpointInfo")
	proto.RegisterType((*NodeInfo)(nil), "clustermgmt.NodeInfo")
	proto.RegisterType((*ListNodesResponse)(nil), "clustermgmt.ListNodesResponse")
	proto.RegisterType((*ListNodesRequest)(nil), "clustermgmt.ListNodesRequest")
}

func init() { proto.RegisterFile("management.proto", fileDescriptor_edc174f991dc0a25) }

var fileDescriptor_edc174f991dc0a25 = []byte{
	// 492 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0xdd, 0x8e, 0x12, 0x31,
	0x14, 0xde, 0x19, 0x60, 0x64, 0x0e, 0xe8, 0x96, 0x93, 0x68, 0x46, 0xb2, 0x22, 0xe9, 0x15, 0xd1,
	0x84, 0x0b, 0x8c, 0xf1, 0x6e, 0xe3, 0xb8, 0x34, 0x58, 0x17, 0x3a, 0xa6, 0x4b, 0x88, 0xee, 0x1d,
	0x4a, 0xdd, 0x90, 0x38, 0x53, 0x64, 0xca, 0x3b, 0xf9, 0x34, 0xbe, 0x89, 0xef, 0x60, 0xa6, 0xe5,
	0x67, 0xd8, 0x88, 0xde, 0xf5, 0x9c, 0xf3, 0x7d, 0xdf, 0xe9, 0xf9, 0x7a, 0x0a, 0x24, 0x9d, 0x67,
	0xf3, 0x3b, 0x95, 0xaa, 0xcc, 0xf4, 0x57, 0x6b, 0x6d, 0x34, 0x36, 0xbe, 0x7e, 0xdf, 0xe4, 0x46,
	0xad, 0xd3, 0xbb, 0xd4, 0xd0, 0x16, 0x9c, 0x8f, 0x94, 0xb9, 0x31, 0x73, 0xa3, 0xa4, 0xfa, 0xb1,
	0x51, 0xb9, 0xa1, 0xbf, 0x3c, 0x20, 0x87, 0x5c, 0xbe, 0xd2, 0x59, 0xae, 0xf0, 0x09, 0x04, 0x42,
	0x2f, 0x14, 0x5f, 0x44, 0x5e, 0xd7, 0xeb, 0x85, 0x72, 0x1b, 0xe1, 0x5b, 0xa8, 0x59, 0x60, 0xe4,
	0x77, 0xbd, 0xde, 0xa3, 0xc1, 0x8b, 0x7e, 0x49, 0xbc, 0x7f, 0x5f, 0xa5, 0x7f, 0xe5, 0x8a, 0x2e,
	0xe9, 0x88, 0x78, 0x01, 0x61, 0xa1, 0x75, 0xa5, 0x37, 0x99, 0x89, 0x2a, 0x5d, 0xaf, 0x57, 0x93,
	0x87, 0x04, 0x76, 0x00, 0x66, 0xda, 0xa8, 0xb5, 0x2b, 0x57, 0x6d, 0xb9, 0x94, 0xa1, 0x14, 0x9a,
	0x65, 0x51, 0x0c, 0xc0, 0x4f, 0xae, 0xc9, 0x19, 0x02, 0x04, 0xb3, 0x64, 0xca, 0xc5, 0x88, 0x78,
	0xf4, 0x12, 0x9a, 0x2c, 0x5b, 0xac, 0xf4, 0x32, 0x33, 0x3c, 0xfb, 0xa6, 0x11, 0xa1, 0x2a, 0xe6,
	0xa9, 0xda, 0x4e, 0x62, 0xcf, 0xd8, 0x86, 0xfa, 0x7b, 0x9d, 0x9b, 0x8f, 0x7a, 0x6d, 0xec, 0x28,
	0xa1, 0xdc, 0xc7, 0xf4, 0xb7, 0x0f, 0x75, 0x3b, 0x6e, 0x41, 0x3e, 0x65, 0xc4, 0x00, 0xaa, 0xd7,
	0xcb, 0x6c, 0xb1, 0xf5, 0xa1, 0x73, 0xe4, 0xc3, 0x8e, 0x6c, 0x0f, 0x05, 0x4a, 0x5a, 0x2c, 0xbe,
	0xde, 0x99, 0x57, 0xb1, 0xa4, 0xe7, 0xa7, 0x49, 0x47, 0x8e, 0xbd, 0x81, 0x70, 0x37, 0x4f, 0x1e,
	0x55, 0xbb, 0x95, 0x5e, 0x63, 0xf0, 0xf4, 0x88, 0x5a, 0x9e, 0x56, 0x1e, 0xb0, 0xf4, 0xd2, 0xcd,
	0x61, 0x7b, 0x03, 0x04, 0x63, 0x16, 0x0f, 0x99, 0x24, 0x67, 0x18, 0x42, 0x6d, 0x96, 0x4c, 0x99,
	0x24, 0x1e, 0x36, 0xa1, 0x2e, 0x12, 0xe1, 0x22, 0x1f, 0x1f, 0x42, 0x28, 0x12, 0x31, 0x61, 0x93,
	0x77, 0x4c, 0x92, 0x0a, 0x55, 0xee, 0xa9, 0xdc, 0x2d, 0xce, 0xa1, 0xc1, 0x05, 0x9f, 0xf2, 0x78,
	0xcc, 0x6f, 0xd9, 0xd0, 0xa9, 0x48, 0x16, 0x0f, 0x3f, 0x13, 0x0f, 0x1b, 0xf0, 0xe0, 0x86, 0xc9,
	0x59, 0x61, 0xbf, 0x8f, 0x04, 0x9a, 0x92, 0x25, 0x72, 0x14, 0x0b, 0x7e, 0x5b, 0x64, 0x2a, 0x45,
	0x93, 0xa1, 0x8c, 0xb9, 0x28, 0xa2, 0x6a, 0x21, 0x34, 0x65, 0x72, 0xc2, 0x45, 0x6c, 0xdf, 0xab,
	0x46, 0x3f, 0x41, 0x6b, 0xbc, 0xcc, 0x4d, 0xd1, 0x2a, 0xff, 0xef, 0x02, 0xbe, 0x84, 0x9a, 0x05,
	0x46, 0xbe, 0x35, 0xe2, 0xf1, 0x5f, 0x3d, 0x94, 0x0e, 0x43, 0x11, 0x48, 0x49, 0xd9, 0xae, 0xfb,
	0xe0, 0xa7, 0x07, 0xad, 0xed, 0x0a, 0x4d, 0xf6, 0x5f, 0x05, 0x47, 0x50, 0xdf, 0x6d, 0x2f, 0x5e,
	0x9c, 0x58, 0x6a, 0xcb, 0x6f, 0x3f, 0xfb, 0xe7, 0xca, 0xe3, 0x07, 0x08, 0xf7, 0x2d, 0xf1, 0x18,
	0x7b, 0xff, 0x2a, 0xed, 0xce, 0xa9, 0xb2, 0xd3, 0xfa, 0x12, 0xd8, 0x0f, 0xfc, 0xea, 0x4f, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xa3, 0x0a, 0xf8, 0xbf, 0xd4, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ClusterManagementClient is the client API for ClusterManagement service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ClusterManagementClient interface {
	// GetState returns the node's internal state. Actual content is WIP.
	GetState(ctx context.Context, in *GetStateRequest, opts ...grpc.CallOption) (*GetStateResponse, error)
	// ListNodes lists the known nodes in the cluster, both voters, non-voters and
	// bystanders (ie those that are member of the Serf swarm but not the Raft cluster)
	ListNodes(ctx context.Context, in *ListNodesRequest, opts ...grpc.CallOption) (*ListNodesResponse, error)
}

type clusterManagementClient struct {
	cc *grpc.ClientConn
}

func NewClusterManagementClient(cc *grpc.ClientConn) ClusterManagementClient {
	return &clusterManagementClient{cc}
}

func (c *clusterManagementClient) GetState(ctx context.Context, in *GetStateRequest, opts ...grpc.CallOption) (*GetStateResponse, error) {
	out := new(GetStateResponse)
	err := c.cc.Invoke(ctx, "/clustermgmt.ClusterManagement/GetState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterManagementClient) ListNodes(ctx context.Context, in *ListNodesRequest, opts ...grpc.CallOption) (*ListNodesResponse, error) {
	out := new(ListNodesResponse)
	err := c.cc.Invoke(ctx, "/clustermgmt.ClusterManagement/ListNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterManagementServer is the server API for ClusterManagement service.
type ClusterManagementServer interface {
	// GetState returns the node's internal state. Actual content is WIP.
	GetState(context.Context, *GetStateRequest) (*GetStateResponse, error)
	// ListNodes lists the known nodes in the cluster, both voters, non-voters and
	// bystanders (ie those that are member of the Serf swarm but not the Raft cluster)
	ListNodes(context.Context, *ListNodesRequest) (*ListNodesResponse, error)
}

func RegisterClusterManagementServer(s *grpc.Server, srv ClusterManagementServer) {
	s.RegisterService(&_ClusterManagement_serviceDesc, srv)
}

func _ClusterManagement_GetState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterManagementServer).GetState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clustermgmt.ClusterManagement/GetState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterManagementServer).GetState(ctx, req.(*GetStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterManagement_ListNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterManagementServer).ListNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clustermgmt.ClusterManagement/ListNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterManagementServer).ListNodes(ctx, req.(*ListNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ClusterManagement_serviceDesc = grpc.ServiceDesc{
	ServiceName: "clustermgmt.ClusterManagement",
	HandlerType: (*ClusterManagementServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetState",
			Handler:    _ClusterManagement_GetState_Handler,
		},
		{
			MethodName: "ListNodes",
			Handler:    _ClusterManagement_ListNodes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "management.proto",
}
