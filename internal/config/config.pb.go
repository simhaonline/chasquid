// Code generated by protoc-gen-go.
// source: config.proto
// DO NOT EDIT!

/*
Package config is a generated protocol buffer package.

It is generated from these files:
	config.proto

It has these top-level messages:
	Config
*/
package config

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Config struct {
	// Hostname to use when we say hello.
	// For aesthetic purposes, but may help if our ip address resolves to it.
	// Default: machine hostname.
	Hostname string `protobuf:"bytes,1,opt,name=hostname" json:"hostname,omitempty"`
	// Maximum email size, in megabytes.
	// Default: 50.
	MaxDataSizeMb int64 `protobuf:"varint,2,opt,name=max_data_size_mb,json=maxDataSizeMb" json:"max_data_size_mb,omitempty"`
	// Addresses to listen on for SMTP (usually port 25).
	// Default: "systemd", which means systemd passes sockets to us.
	// systemd sockets must be named with "FileDescriptorName=smtp".
	SmtpAddress []string `protobuf:"bytes,3,rep,name=smtp_address,json=smtpAddress" json:"smtp_address,omitempty"`
	// Addresses to listen on for submission (usually port 587).
	// Default: "systemd", which means systemd passes sockets to us.
	// systemd sockets must be named with "FileDescriptorName=submission".
	SubmissionAddress []string `protobuf:"bytes,4,rep,name=submission_address,json=submissionAddress" json:"submission_address,omitempty"`
	// Address for the monitoring http server.
	// Default: no monitoring http server.
	MonitoringAddress string `protobuf:"bytes,5,opt,name=monitoring_address,json=monitoringAddress" json:"monitoring_address,omitempty"`
	// Mail delivery agent (MDA, also known as LDA) to use.
	// This should point to the binary to use to deliver email to local users.
	// The content of the email will be passed via stdin.
	// If it exits unsuccessfully, we assume the mail was not delivered.
	// Default: "procmail".
	MailDeliveryAgentBin string `protobuf:"bytes,6,opt,name=mail_delivery_agent_bin,json=mailDeliveryAgentBin" json:"mail_delivery_agent_bin,omitempty"`
	// Command line arguments for the mail delivery agent. One per argument.
	// Some replacements will be done:
	//  - "%user%"   -> local user (anything before the @)
	//  - "%domain%" -> domain (anything after the @)
	// Default: "-d", "%user"  (adequate for procmail)
	MailDeliveryAgentArgs []string `protobuf:"bytes,7,rep,name=mail_delivery_agent_args,json=mailDeliveryAgentArgs" json:"mail_delivery_agent_args,omitempty"`
	// Directory where we store our persistent data.
	// Default: "/var/lib/chasquid"
	DataDir string `protobuf:"bytes,8,opt,name=data_dir,json=dataDir" json:"data_dir,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*Config)(nil), "Config")
}

func init() { proto.RegisterFile("config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 251 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x90, 0x41, 0x2f, 0x04, 0x31,
	0x14, 0xc7, 0xb3, 0x3b, 0xcc, 0x8e, 0x67, 0x25, 0xb6, 0x21, 0xca, 0x69, 0xb9, 0x70, 0xe1, 0x22,
	0xe2, 0x3c, 0xcc, 0xd5, 0x65, 0x7d, 0x80, 0xa6, 0xd5, 0x1a, 0x2f, 0xd9, 0xb6, 0x9b, 0xbe, 0x12,
	0x7c, 0x53, 0xdf, 0x46, 0x5b, 0xcc, 0x4a, 0x38, 0xbe, 0xff, 0xef, 0xf7, 0x6f, 0x5f, 0x1e, 0x4c,
	0x1f, 0xbc, 0x7b, 0xc4, 0xfe, 0x62, 0x15, 0x7c, 0xf4, 0x27, 0x1f, 0x63, 0xa8, 0x6f, 0x4b, 0xc0,
	0x8e, 0xa0, 0x79, 0xf2, 0x14, 0x9d, 0xb4, 0x86, 0x8f, 0xe6, 0xa3, 0xb3, 0xad, 0xc5, 0x30, 0xb3,
	0x53, 0xd8, 0xb5, 0xf2, 0x55, 0x68, 0x19, 0xa5, 0x20, 0x7c, 0x37, 0xc2, 0x2a, 0x3e, 0x4e, 0x4e,
	0xb5, 0xd8, 0x49, 0x79, 0x97, 0xe2, 0xfb, 0x94, 0xde, 0x29, 0x76, 0x0c, 0x53, 0xb2, 0x71, 0x25,
	0xa4, 0xd6, 0xc1, 0x10, 0xf1, 0x6a, 0x5e, 0xa5, 0x87, 0xb6, 0x73, 0xd6, 0x7e, 0x45, 0xec, 0x1c,
	0x18, 0x3d, 0x2b, 0x8b, 0x44, 0xe8, 0xdd, 0x20, 0x6e, 0x14, 0x71, 0xb6, 0x26, 0xbf, 0x74, 0xeb,
	0x1d, 0x46, 0x1f, 0xd0, 0xf5, 0x83, 0xbe, 0x59, 0x16, 0x9c, 0xad, 0xc9, 0x8f, 0x7e, 0x05, 0x07,
	0x56, 0xe2, 0x52, 0x68, 0xb3, 0xc4, 0x17, 0x13, 0xde, 0x84, 0xec, 0x8d, 0x8b, 0x42, 0xa1, 0xe3,
	0x75, 0xe9, 0xec, 0x65, 0xdc, 0x7d, 0xd3, 0x36, 0xc3, 0x1b, 0x74, 0xec, 0x1a, 0xf8, 0x7f, 0x35,
	0x19, 0x7a, 0xe2, 0x93, 0xb2, 0xda, 0xfe, 0x9f, 0x5e, 0x9b, 0x20, 0x3b, 0x84, 0xa6, 0x5c, 0x45,
	0x63, 0xe0, 0x4d, 0xf9, 0x60, 0x92, 0xe7, 0x0e, 0x83, 0xaa, 0xcb, 0x89, 0x2f, 0x3f, 0x03, 0x00,
	0x00, 0xff, 0xff, 0xc5, 0x2e, 0x79, 0xef, 0x72, 0x01, 0x00, 0x00,
}
