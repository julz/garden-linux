// Code generated by protoc-gen-gogo.
// source: net_in.proto
// DO NOT EDIT!

package garden

import proto "code.google.com/p/gogoprotobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type NetInRequest struct {
	Handle           *string `protobuf:"bytes,1,req,name=handle" json:"handle,omitempty"`
	HostPort         *uint32 `protobuf:"varint,3,opt,name=host_port" json:"host_port,omitempty"`
	ContainerPort    *uint32 `protobuf:"varint,2,opt,name=container_port" json:"container_port,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *NetInRequest) Reset()         { *m = NetInRequest{} }
func (m *NetInRequest) String() string { return proto.CompactTextString(m) }
func (*NetInRequest) ProtoMessage()    {}

func (m *NetInRequest) GetHandle() string {
	if m != nil && m.Handle != nil {
		return *m.Handle
	}
	return ""
}

func (m *NetInRequest) GetHostPort() uint32 {
	if m != nil && m.HostPort != nil {
		return *m.HostPort
	}
	return 0
}

func (m *NetInRequest) GetContainerPort() uint32 {
	if m != nil && m.ContainerPort != nil {
		return *m.ContainerPort
	}
	return 0
}

type NetInResponse struct {
	HostPort         *uint32 `protobuf:"varint,1,req,name=host_port" json:"host_port,omitempty"`
	ContainerPort    *uint32 `protobuf:"varint,2,req,name=container_port" json:"container_port,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *NetInResponse) Reset()         { *m = NetInResponse{} }
func (m *NetInResponse) String() string { return proto.CompactTextString(m) }
func (*NetInResponse) ProtoMessage()    {}

func (m *NetInResponse) GetHostPort() uint32 {
	if m != nil && m.HostPort != nil {
		return *m.HostPort
	}
	return 0
}

func (m *NetInResponse) GetContainerPort() uint32 {
	if m != nil && m.ContainerPort != nil {
		return *m.ContainerPort
	}
	return 0
}

func init() {
}
