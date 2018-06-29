// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/config/ratelimit/v2/rls.proto

/*
	Package v2 is a generated protocol buffer package.

	It is generated from these files:
		envoy/config/ratelimit/v2/rls.proto

	It has these top-level messages:
		RateLimitServiceConfig
*/
package v2

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import envoy_api_v2_core1 "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
import _ "github.com/lyft/protoc-gen-validate/validate"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Rate limit :ref:`configuration overview <config_rate_limit_service>`.
type RateLimitServiceConfig struct {
	// Types that are valid to be assigned to ServiceSpecifier:
	//	*RateLimitServiceConfig_ClusterName
	//	*RateLimitServiceConfig_GrpcService
	ServiceSpecifier isRateLimitServiceConfig_ServiceSpecifier `protobuf_oneof:"service_specifier"`
	// Specifies if Envoy should use the data-plane-api client
	// :repo:`api/envoy/service/ratelimit/v2/rls.proto` or the legacy
	// client :repo:`source/common/ratelimit/ratelimit.proto` when
	// making requests to the rate limit service.
	//
	// .. note::
	//
	//   The legacy client will be used by
	//   default until the start of the 1.9.0 release cycle. At the start of the
	//   1.9.0 release cycle this field will be removed and only the data-plane-api
	//   proto will be supported. This means that your rate limit service needs to
	//   have support for the data-plane-api proto by the start of the 1.9.0 release cycle.
	//   Lyft's `reference implementation <https://github.com/lyft/ratelimit>`_
	//   supports the data-plane-api version as of v1.1.0.
	UseDataPlaneProto bool `protobuf:"varint,3,opt,name=use_data_plane_proto,json=useDataPlaneProto,proto3" json:"use_data_plane_proto,omitempty"`
}

func (m *RateLimitServiceConfig) Reset()                    { *m = RateLimitServiceConfig{} }
func (m *RateLimitServiceConfig) String() string            { return proto.CompactTextString(m) }
func (*RateLimitServiceConfig) ProtoMessage()               {}
func (*RateLimitServiceConfig) Descriptor() ([]byte, []int) { return fileDescriptorRls, []int{0} }

type isRateLimitServiceConfig_ServiceSpecifier interface {
	isRateLimitServiceConfig_ServiceSpecifier()
	MarshalTo([]byte) (int, error)
	Size() int
}

type RateLimitServiceConfig_ClusterName struct {
	ClusterName string `protobuf:"bytes,1,opt,name=cluster_name,json=clusterName,proto3,oneof"`
}
type RateLimitServiceConfig_GrpcService struct {
	GrpcService *envoy_api_v2_core1.GrpcService `protobuf:"bytes,2,opt,name=grpc_service,json=grpcService,oneof"`
}

func (*RateLimitServiceConfig_ClusterName) isRateLimitServiceConfig_ServiceSpecifier() {}
func (*RateLimitServiceConfig_GrpcService) isRateLimitServiceConfig_ServiceSpecifier() {}

func (m *RateLimitServiceConfig) GetServiceSpecifier() isRateLimitServiceConfig_ServiceSpecifier {
	if m != nil {
		return m.ServiceSpecifier
	}
	return nil
}

func (m *RateLimitServiceConfig) GetClusterName() string {
	if x, ok := m.GetServiceSpecifier().(*RateLimitServiceConfig_ClusterName); ok {
		return x.ClusterName
	}
	return ""
}

func (m *RateLimitServiceConfig) GetGrpcService() *envoy_api_v2_core1.GrpcService {
	if x, ok := m.GetServiceSpecifier().(*RateLimitServiceConfig_GrpcService); ok {
		return x.GrpcService
	}
	return nil
}

func (m *RateLimitServiceConfig) GetUseDataPlaneProto() bool {
	if m != nil {
		return m.UseDataPlaneProto
	}
	return false
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*RateLimitServiceConfig) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _RateLimitServiceConfig_OneofMarshaler, _RateLimitServiceConfig_OneofUnmarshaler, _RateLimitServiceConfig_OneofSizer, []interface{}{
		(*RateLimitServiceConfig_ClusterName)(nil),
		(*RateLimitServiceConfig_GrpcService)(nil),
	}
}

func _RateLimitServiceConfig_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*RateLimitServiceConfig)
	// service_specifier
	switch x := m.ServiceSpecifier.(type) {
	case *RateLimitServiceConfig_ClusterName:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		_ = b.EncodeStringBytes(x.ClusterName)
	case *RateLimitServiceConfig_GrpcService:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.GrpcService); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("RateLimitServiceConfig.ServiceSpecifier has unexpected type %T", x)
	}
	return nil
}

func _RateLimitServiceConfig_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*RateLimitServiceConfig)
	switch tag {
	case 1: // service_specifier.cluster_name
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.ServiceSpecifier = &RateLimitServiceConfig_ClusterName{x}
		return true, err
	case 2: // service_specifier.grpc_service
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(envoy_api_v2_core1.GrpcService)
		err := b.DecodeMessage(msg)
		m.ServiceSpecifier = &RateLimitServiceConfig_GrpcService{msg}
		return true, err
	default:
		return false, nil
	}
}

func _RateLimitServiceConfig_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*RateLimitServiceConfig)
	// service_specifier
	switch x := m.ServiceSpecifier.(type) {
	case *RateLimitServiceConfig_ClusterName:
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.ClusterName)))
		n += len(x.ClusterName)
	case *RateLimitServiceConfig_GrpcService:
		s := proto.Size(x.GrpcService)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*RateLimitServiceConfig)(nil), "envoy.config.ratelimit.v2.RateLimitServiceConfig")
}
func (m *RateLimitServiceConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RateLimitServiceConfig) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.ServiceSpecifier != nil {
		nn1, err := m.ServiceSpecifier.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn1
	}
	if m.UseDataPlaneProto {
		dAtA[i] = 0x18
		i++
		if m.UseDataPlaneProto {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	return i, nil
}

func (m *RateLimitServiceConfig_ClusterName) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	dAtA[i] = 0xa
	i++
	i = encodeVarintRls(dAtA, i, uint64(len(m.ClusterName)))
	i += copy(dAtA[i:], m.ClusterName)
	return i, nil
}
func (m *RateLimitServiceConfig_GrpcService) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.GrpcService != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintRls(dAtA, i, uint64(m.GrpcService.Size()))
		n2, err := m.GrpcService.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}
func encodeVarintRls(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *RateLimitServiceConfig) Size() (n int) {
	var l int
	_ = l
	if m.ServiceSpecifier != nil {
		n += m.ServiceSpecifier.Size()
	}
	if m.UseDataPlaneProto {
		n += 2
	}
	return n
}

func (m *RateLimitServiceConfig_ClusterName) Size() (n int) {
	var l int
	_ = l
	l = len(m.ClusterName)
	n += 1 + l + sovRls(uint64(l))
	return n
}
func (m *RateLimitServiceConfig_GrpcService) Size() (n int) {
	var l int
	_ = l
	if m.GrpcService != nil {
		l = m.GrpcService.Size()
		n += 1 + l + sovRls(uint64(l))
	}
	return n
}

func sovRls(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozRls(x uint64) (n int) {
	return sovRls(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RateLimitServiceConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRls
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RateLimitServiceConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RateLimitServiceConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRls
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRls
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ServiceSpecifier = &RateLimitServiceConfig_ClusterName{string(dAtA[iNdEx:postIndex])}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GrpcService", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRls
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRls
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &envoy_api_v2_core1.GrpcService{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.ServiceSpecifier = &RateLimitServiceConfig_GrpcService{v}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UseDataPlaneProto", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRls
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.UseDataPlaneProto = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipRls(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRls
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipRls(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRls
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRls
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRls
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthRls
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowRls
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipRls(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthRls = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRls   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("envoy/config/ratelimit/v2/rls.proto", fileDescriptorRls) }

var fileDescriptorRls = []byte{
	// 301 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8e, 0xcd, 0x4a, 0xf3, 0x40,
	0x14, 0x86, 0xbf, 0x93, 0x7e, 0x8a, 0xa6, 0x5d, 0xd8, 0x20, 0x1a, 0xbb, 0x08, 0x41, 0x5d, 0x74,
	0x35, 0x91, 0xf4, 0x0e, 0x52, 0x41, 0x17, 0x22, 0xa5, 0xee, 0xdc, 0x84, 0xe3, 0xf4, 0xb4, 0x0c,
	0xa4, 0xc9, 0x30, 0x33, 0x1d, 0xf0, 0xce, 0xc4, 0x55, 0x97, 0x2e, 0xbd, 0x84, 0xd2, 0x5d, 0xef,
	0x42, 0x26, 0x53, 0x7f, 0x76, 0xe7, 0xf0, 0xfe, 0x3c, 0x6f, 0x78, 0x45, 0xb5, 0x6d, 0x5e, 0x33,
	0xde, 0xd4, 0x73, 0xb1, 0xc8, 0x14, 0x1a, 0xaa, 0xc4, 0x52, 0x98, 0xcc, 0xe6, 0x99, 0xaa, 0x34,
	0x93, 0xaa, 0x31, 0x4d, 0x74, 0xd1, 0x9a, 0x98, 0x37, 0xb1, 0x1f, 0x13, 0xb3, 0xf9, 0xe0, 0xda,
	0xe7, 0x51, 0x0a, 0x17, 0xe1, 0x8d, 0xa2, 0x6c, 0xa1, 0x24, 0x2f, 0x35, 0x29, 0x2b, 0x38, 0xf9,
	0x82, 0xc1, 0xb9, 0xc5, 0x4a, 0xcc, 0xd0, 0x50, 0xf6, 0x7d, 0x78, 0xe1, 0x72, 0x03, 0xe1, 0xd9,
	0x14, 0x0d, 0x3d, 0xb8, 0xbe, 0x27, 0x9f, 0x19, 0xb7, 0x98, 0xe8, 0x26, 0xec, 0xf1, 0x6a, 0xa5,
	0x0d, 0xa9, 0xb2, 0xc6, 0x25, 0xc5, 0x90, 0xc2, 0xf0, 0xb8, 0xe8, 0xbe, 0xef, 0xd6, 0x9d, 0xff,
	0x2a, 0x48, 0x21, 0x86, 0xfb, 0x7f, 0xd3, 0xee, 0xde, 0xf2, 0x88, 0x4b, 0x8a, 0xc6, 0x61, 0xef,
	0x2f, 0x3b, 0x0e, 0x52, 0x18, 0x76, 0xf3, 0x84, 0xf9, 0xf5, 0x28, 0x05, 0xb3, 0x39, 0x73, 0x13,
	0xd9, 0x9d, 0x92, 0x7c, 0x4f, 0x73, 0x25, 0x8b, 0xdf, 0x37, 0x1a, 0x85, 0xa7, 0x2b, 0x4d, 0xe5,
	0x0c, 0x0d, 0x96, 0xb2, 0xc2, 0x9a, 0xca, 0x76, 0x69, 0xdc, 0x49, 0x61, 0x78, 0x54, 0x04, 0x31,
	0x4c, 0xfb, 0x2b, 0x4d, 0xb7, 0x68, 0x70, 0xe2, 0xd4, 0x89, 0x13, 0x8b, 0x41, 0xd8, 0xdf, 0x43,
	0x4b, 0x2d, 0x89, 0x8b, 0xb9, 0x20, 0x15, 0x1d, 0xbc, 0xed, 0xd6, 0x1d, 0x28, 0x4e, 0x3e, 0xb6,
	0x09, 0x7c, 0x6e, 0x13, 0xd8, 0x6c, 0x13, 0x78, 0x0e, 0x6c, 0xfe, 0x72, 0xd8, 0x36, 0x8e, 0xbe,
	0x02, 0x00, 0x00, 0xff, 0xff, 0x67, 0x4c, 0xfc, 0xa8, 0x7c, 0x01, 0x00, 0x00,
}
