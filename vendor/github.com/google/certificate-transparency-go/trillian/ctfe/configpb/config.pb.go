// Code generated by protoc-gen-go. DO NOT EDIT.
// source: config.proto

/*
Package configpb is a generated protocol buffer package.

It is generated from these files:
	config.proto

It has these top-level messages:
	LogBackend
	LogBackendSet
	LogConfigSet
	LogConfig
*/
package configpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import keyspb "github.com/google/trillian/crypto/keyspb"
import google_protobuf "github.com/golang/protobuf/ptypes/any"
import google_protobuf1 "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type LogBackend struct {
	// name defines the name of the log backend for use in LogConfig messages and must be unique.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// backend_spec defines the RPC endpoint that clients should use to send requests
	// to this log backend. These should be in the same format as rpcBackendFlag in the
	// CTFE main and must not be an empty string.
	BackendSpec string `protobuf:"bytes,2,opt,name=backend_spec,json=backendSpec" json:"backend_spec,omitempty"`
}

func (m *LogBackend) Reset()                    { *m = LogBackend{} }
func (m *LogBackend) String() string            { return proto.CompactTextString(m) }
func (*LogBackend) ProtoMessage()               {}
func (*LogBackend) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *LogBackend) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *LogBackend) GetBackendSpec() string {
	if m != nil {
		return m.BackendSpec
	}
	return ""
}

// LogBackendSet supports a configuration where a single set of frontends handle
// requests for multiple backends. For example this could be used to run different
// backends in different geographic regions.
type LogBackendSet struct {
	Backend []*LogBackend `protobuf:"bytes,1,rep,name=backend" json:"backend,omitempty"`
}

func (m *LogBackendSet) Reset()                    { *m = LogBackendSet{} }
func (m *LogBackendSet) String() string            { return proto.CompactTextString(m) }
func (*LogBackendSet) ProtoMessage()               {}
func (*LogBackendSet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *LogBackendSet) GetBackend() []*LogBackend {
	if m != nil {
		return m.Backend
	}
	return nil
}

// LogConfigSet is a set of LogConfig messages.
type LogConfigSet struct {
	Config []*LogConfig `protobuf:"bytes,1,rep,name=config" json:"config,omitempty"`
}

func (m *LogConfigSet) Reset()                    { *m = LogConfigSet{} }
func (m *LogConfigSet) String() string            { return proto.CompactTextString(m) }
func (*LogConfigSet) ProtoMessage()               {}
func (*LogConfigSet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *LogConfigSet) GetConfig() []*LogConfig {
	if m != nil {
		return m.Config
	}
	return nil
}

// LogConfig describes the configuration options for a log instance.
type LogConfig struct {
	LogId        int64                `protobuf:"varint,1,opt,name=log_id,json=logId" json:"log_id,omitempty"`
	Prefix       string               `protobuf:"bytes,2,opt,name=prefix" json:"prefix,omitempty"`
	RootsPemFile []string             `protobuf:"bytes,3,rep,name=roots_pem_file,json=rootsPemFile" json:"roots_pem_file,omitempty"`
	PrivateKey   *google_protobuf.Any `protobuf:"bytes,4,opt,name=private_key,json=privateKey" json:"private_key,omitempty"`
	// The public key is included for the convenience of test tools (and obviously
	// should match the private key above); it is not used by the CT personality.
	PublicKey     *keyspb.PublicKey `protobuf:"bytes,5,opt,name=public_key,json=publicKey" json:"public_key,omitempty"`
	RejectExpired bool              `protobuf:"varint,6,opt,name=reject_expired,json=rejectExpired" json:"reject_expired,omitempty"`
	ExtKeyUsages  []string          `protobuf:"bytes,7,rep,name=ext_key_usages,json=extKeyUsages" json:"ext_key_usages,omitempty"`
	// not_after_start defines the start of the range of acceptable NotAfter
	// values, inclusive.
	// Leaving this unset implies no lower bound to the range.
	NotAfterStart *google_protobuf1.Timestamp `protobuf:"bytes,8,opt,name=not_after_start,json=notAfterStart" json:"not_after_start,omitempty"`
	// not_after_limit defines the end of the range of acceptable NotAfter values,
	// exclusive.
	// Leaving this unset implies no upper bound to the range.
	NotAfterLimit *google_protobuf1.Timestamp `protobuf:"bytes,9,opt,name=not_after_limit,json=notAfterLimit" json:"not_after_limit,omitempty"`
	// accept_only_ca controls whether or not *only* certificates with the CA bit
	// set will be accepted.
	AcceptOnlyCa bool `protobuf:"varint,10,opt,name=accept_only_ca,json=acceptOnlyCa" json:"accept_only_ca,omitempty"`
	// backend_name if set indicates which backend serves this log. The name must be
	// one of those defined in the LogBackendSet.
	LogBackendName string `protobuf:"bytes,11,opt,name=log_backend_name,json=logBackendName" json:"log_backend_name,omitempty"`
}

func (m *LogConfig) Reset()                    { *m = LogConfig{} }
func (m *LogConfig) String() string            { return proto.CompactTextString(m) }
func (*LogConfig) ProtoMessage()               {}
func (*LogConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *LogConfig) GetLogId() int64 {
	if m != nil {
		return m.LogId
	}
	return 0
}

func (m *LogConfig) GetPrefix() string {
	if m != nil {
		return m.Prefix
	}
	return ""
}

func (m *LogConfig) GetRootsPemFile() []string {
	if m != nil {
		return m.RootsPemFile
	}
	return nil
}

func (m *LogConfig) GetPrivateKey() *google_protobuf.Any {
	if m != nil {
		return m.PrivateKey
	}
	return nil
}

func (m *LogConfig) GetPublicKey() *keyspb.PublicKey {
	if m != nil {
		return m.PublicKey
	}
	return nil
}

func (m *LogConfig) GetRejectExpired() bool {
	if m != nil {
		return m.RejectExpired
	}
	return false
}

func (m *LogConfig) GetExtKeyUsages() []string {
	if m != nil {
		return m.ExtKeyUsages
	}
	return nil
}

func (m *LogConfig) GetNotAfterStart() *google_protobuf1.Timestamp {
	if m != nil {
		return m.NotAfterStart
	}
	return nil
}

func (m *LogConfig) GetNotAfterLimit() *google_protobuf1.Timestamp {
	if m != nil {
		return m.NotAfterLimit
	}
	return nil
}

func (m *LogConfig) GetAcceptOnlyCa() bool {
	if m != nil {
		return m.AcceptOnlyCa
	}
	return false
}

func (m *LogConfig) GetLogBackendName() string {
	if m != nil {
		return m.LogBackendName
	}
	return ""
}

func init() {
	proto.RegisterType((*LogBackend)(nil), "configpb.LogBackend")
	proto.RegisterType((*LogBackendSet)(nil), "configpb.LogBackendSet")
	proto.RegisterType((*LogConfigSet)(nil), "configpb.LogConfigSet")
	proto.RegisterType((*LogConfig)(nil), "configpb.LogConfig")
}

func init() { proto.RegisterFile("config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 487 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x4b, 0x8f, 0xd3, 0x30,
	0x10, 0x56, 0xe9, 0x6e, 0x77, 0x3b, 0x7d, 0x00, 0xe6, 0xa1, 0xd0, 0x0b, 0xa5, 0x02, 0xa9, 0x12,
	0x52, 0x8a, 0x40, 0x7b, 0xe2, 0x80, 0x76, 0x2b, 0x90, 0x50, 0x2b, 0x58, 0xa5, 0x70, 0xb6, 0x1c,
	0x77, 0x1a, 0x4c, 0x9d, 0xd8, 0x72, 0x5c, 0x54, 0xff, 0x69, 0x7e, 0x03, 0x8a, 0xed, 0xec, 0x6a,
	0xe1, 0xc2, 0x29, 0xf6, 0xf7, 0x98, 0x7c, 0xe3, 0x19, 0x18, 0x72, 0x55, 0xed, 0x44, 0x91, 0x6a,
	0xa3, 0xac, 0x22, 0xe7, 0xe1, 0xa6, 0xf3, 0xc9, 0x45, 0x21, 0xec, 0x8f, 0x43, 0x9e, 0x72, 0x55,
	0x2e, 0x0a, 0xa5, 0x0a, 0x89, 0x0b, 0x6b, 0x84, 0x94, 0x82, 0x55, 0x0b, 0x6e, 0x9c, 0xb6, 0x6a,
	0xb1, 0x47, 0x57, 0xeb, 0x3c, 0x7e, 0x42, 0x81, 0xc9, 0xb3, 0xa8, 0xf5, 0xb7, 0xfc, 0xb0, 0x5b,
	0xb0, 0xca, 0x45, 0xea, 0xf9, 0xdf, 0x94, 0x15, 0x25, 0xd6, 0x96, 0x95, 0x3a, 0x08, 0x66, 0x4b,
	0x80, 0xb5, 0x2a, 0xae, 0x18, 0xdf, 0x63, 0xb5, 0x25, 0x04, 0x4e, 0x2a, 0x56, 0x62, 0xd2, 0x99,
	0x76, 0xe6, 0xfd, 0xcc, 0x9f, 0xc9, 0x0b, 0x18, 0xe6, 0x81, 0xa6, 0xb5, 0x46, 0x9e, 0xdc, 0xf3,
	0xdc, 0x20, 0x62, 0x1b, 0x8d, 0x7c, 0xf6, 0x01, 0x46, 0xb7, 0x45, 0x36, 0x68, 0x49, 0x0a, 0x67,
	0x91, 0x4f, 0x3a, 0xd3, 0xee, 0x7c, 0xf0, 0xf6, 0x71, 0xda, 0x36, 0x99, 0xde, 0x2a, 0xb3, 0x56,
	0x34, 0x7b, 0x0f, 0xc3, 0xb5, 0x2a, 0x96, 0x5e, 0xd2, 0xf8, 0x5f, 0x43, 0x2f, 0xe8, 0xa3, 0xfd,
	0xd1, 0x1d, 0x7b, 0xd0, 0x65, 0x51, 0x32, 0xfb, 0xdd, 0x85, 0xfe, 0x0d, 0x4a, 0x9e, 0x40, 0x4f,
	0xaa, 0x82, 0x8a, 0xad, 0x6f, 0xa2, 0x9b, 0x9d, 0x4a, 0x55, 0x7c, 0xde, 0x92, 0xa7, 0xd0, 0xd3,
	0x06, 0x77, 0xe2, 0x18, 0xf3, 0xc7, 0x1b, 0x79, 0x09, 0x63, 0xa3, 0x94, 0xad, 0xa9, 0xc6, 0x92,
	0xee, 0x84, 0xc4, 0xa4, 0x3b, 0xed, 0xce, 0xfb, 0xd9, 0xd0, 0xa3, 0xd7, 0x58, 0x7e, 0x12, 0x12,
	0xc9, 0x05, 0x0c, 0xb4, 0x11, 0xbf, 0x98, 0x45, 0xba, 0x47, 0x97, 0x9c, 0x4c, 0x3b, 0xbe, 0xa7,
	0xf0, 0xb8, 0x69, 0xfb, 0xb8, 0xe9, 0x65, 0xe5, 0x32, 0x88, 0xc2, 0x15, 0x3a, 0xf2, 0x06, 0x40,
	0x1f, 0x72, 0x29, 0xb8, 0x77, 0x9d, 0x7a, 0xd7, 0xc3, 0x34, 0xce, 0xee, 0xda, 0x33, 0x2b, 0x74,
	0x59, 0x5f, 0xb7, 0x47, 0xf2, 0x0a, 0xc6, 0x06, 0x7f, 0x22, 0xb7, 0x14, 0x8f, 0x5a, 0x18, 0xdc,
	0x26, 0xbd, 0x69, 0x67, 0x7e, 0x9e, 0x8d, 0x02, 0xfa, 0x31, 0x80, 0x4d, 0x6a, 0x3c, 0xda, 0xa6,
	0x2a, 0x3d, 0xd4, 0xac, 0xc0, 0x3a, 0x39, 0x0b, 0xa9, 0xf1, 0x68, 0x57, 0xe8, 0xbe, 0x7b, 0x8c,
	0x5c, 0xc1, 0xfd, 0x4a, 0x59, 0xca, 0x76, 0x16, 0x0d, 0xad, 0x2d, 0x33, 0x36, 0x39, 0xf7, 0x19,
	0x26, 0xff, 0x24, 0xff, 0xd6, 0xae, 0x45, 0x36, 0xaa, 0x94, 0xbd, 0x6c, 0x1c, 0x9b, 0xc6, 0x70,
	0xb7, 0x86, 0x14, 0xa5, 0xb0, 0x49, 0xff, 0xff, 0x6b, 0xac, 0x1b, 0x43, 0x93, 0x96, 0x71, 0x8e,
	0xda, 0x52, 0x55, 0x49, 0x47, 0x39, 0x4b, 0xc0, 0x37, 0x35, 0x0c, 0xe8, 0xd7, 0x4a, 0xba, 0x25,
	0x23, 0x73, 0x78, 0xd0, 0x0c, 0xae, 0xdd, 0x35, 0xbf, 0x87, 0x03, 0x3f, 0xab, 0xb1, 0xbc, 0x59,
	0x99, 0x2f, 0xac, 0xc4, 0xbc, 0xe7, 0x7f, 0xf9, 0xee, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xfd,
	0x31, 0xa1, 0xca, 0x47, 0x03, 0x00, 0x00,
}