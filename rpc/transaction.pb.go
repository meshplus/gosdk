// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: transaction.proto

package rpc

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
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
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type TransactionValue_Opcode int32

const (
	TransactionValue_NORMAL   TransactionValue_Opcode = 0
	TransactionValue_UPDATE   TransactionValue_Opcode = 1
	TransactionValue_FREEZE   TransactionValue_Opcode = 2
	TransactionValue_UNFREEZE TransactionValue_Opcode = 3
	TransactionValue_SKIPVM   TransactionValue_Opcode = 4
	TransactionValue_ARCHIVE  TransactionValue_Opcode = 100
)

var TransactionValue_Opcode_name = map[int32]string{
	0:   "NORMAL",
	1:   "UPDATE",
	2:   "FREEZE",
	3:   "UNFREEZE",
	4:   "SKIPVM",
	100: "ARCHIVE",
}

var TransactionValue_Opcode_value = map[string]int32{
	"NORMAL":   0,
	"UPDATE":   1,
	"FREEZE":   2,
	"UNFREEZE": 3,
	"SKIPVM":   4,
	"ARCHIVE":  100,
}

func (x TransactionValue_Opcode) String() string {
	return proto.EnumName(TransactionValue_Opcode_name, int32(x))
}

func (TransactionValue_Opcode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{0, 0}
}

type TransactionValue_VmType int32

const (
	TransactionValue_EVM      TransactionValue_VmType = 0
	TransactionValue_JVM      TransactionValue_VmType = 1
	TransactionValue_HVM      TransactionValue_VmType = 2
	TransactionValue_BVM      TransactionValue_VmType = 3
	TransactionValue_TRANSFER TransactionValue_VmType = 4
	TransactionValue_KVSQL    TransactionValue_VmType = 5
)

var TransactionValue_VmType_name = map[int32]string{
	0: "EVM",
	1: "JVM",
	2: "HVM",
	3: "BVM",
	4: "TRANSFER",
}

var TransactionValue_VmType_value = map[string]int32{
	"EVM":      0,
	"JVM":      1,
	"HVM":      2,
	"BVM":      3,
	"TRANSFER": 4,
}

func (x TransactionValue_VmType) String() string {
	return proto.EnumName(TransactionValue_VmType_name, int32(x))
}

func (TransactionValue_VmType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{0, 1}
}

type TransactionValue struct {
	Price              int64                   `protobuf:"varint,1,opt,name=price,proto3" json:"price,omitempty"`
	GasLimit           int64                   `protobuf:"varint,2,opt,name=gasLimit,proto3" json:"gasLimit,omitempty"`
	Amount             int64                   `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
	Payload            []byte                  `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
	EncryptedAmount    []byte                  `protobuf:"bytes,5,opt,name=encryptedAmount,proto3" json:"encryptedAmount,omitempty"`
	HomomorphicAmount  []byte                  `protobuf:"bytes,6,opt,name=homomorphicAmount,proto3" json:"homomorphicAmount,omitempty"`
	HomomorphicBalance []byte                  `protobuf:"bytes,7,opt,name=homomorphicBalance,proto3" json:"homomorphicBalance,omitempty"`
	Op                 TransactionValue_Opcode `protobuf:"varint,8,opt,name=op,proto3,enum=rpc.TransactionValue_Opcode" json:"op,omitempty"`
	Extra              []byte                  `protobuf:"bytes,9,opt,name=extra,proto3" json:"extra,omitempty"`
	ExtraId            []byte                  `protobuf:"bytes,10,opt,name=extraId,proto3" json:"extraId,omitempty"`
	VmType             TransactionValue_VmType `protobuf:"varint,11,opt,name=vmType,proto3,enum=rpc.TransactionValue_VmType" json:"vmType,omitempty"`
}

func (m *TransactionValue) Reset()         { *m = TransactionValue{} }
func (m *TransactionValue) String() string { return proto.CompactTextString(m) }
func (*TransactionValue) ProtoMessage()    {}
func (*TransactionValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{0}
}
func (m *TransactionValue) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TransactionValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TransactionValue.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TransactionValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionValue.Merge(m, src)
}
func (m *TransactionValue) XXX_Size() int {
	return m.Size()
}
func (m *TransactionValue) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionValue.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionValue proto.InternalMessageInfo

func (m *TransactionValue) GetPrice() int64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *TransactionValue) GetGasLimit() int64 {
	if m != nil {
		return m.GasLimit
	}
	return 0
}

func (m *TransactionValue) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *TransactionValue) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *TransactionValue) GetEncryptedAmount() []byte {
	if m != nil {
		return m.EncryptedAmount
	}
	return nil
}

func (m *TransactionValue) GetHomomorphicAmount() []byte {
	if m != nil {
		return m.HomomorphicAmount
	}
	return nil
}

func (m *TransactionValue) GetHomomorphicBalance() []byte {
	if m != nil {
		return m.HomomorphicBalance
	}
	return nil
}

func (m *TransactionValue) GetOp() TransactionValue_Opcode {
	if m != nil {
		return m.Op
	}
	return TransactionValue_NORMAL
}

func (m *TransactionValue) GetExtra() []byte {
	if m != nil {
		return m.Extra
	}
	return nil
}

func (m *TransactionValue) GetExtraId() []byte {
	if m != nil {
		return m.ExtraId
	}
	return nil
}

func (m *TransactionValue) GetVmType() TransactionValue_VmType {
	if m != nil {
		return m.VmType
	}
	return TransactionValue_EVM
}

func init() {
	proto.RegisterEnum("rpc.TransactionValue_Opcode", TransactionValue_Opcode_name, TransactionValue_Opcode_value)
	proto.RegisterEnum("rpc.TransactionValue_VmType", TransactionValue_VmType_name, TransactionValue_VmType_value)
	proto.RegisterType((*TransactionValue)(nil), "rpc.TransactionValue")
}

func init() { proto.RegisterFile("transaction.proto", fileDescriptor_2cc4e03d2c28c490) }

var fileDescriptor_2cc4e03d2c28c490 = []byte{
	// 403 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0x4f, 0x6f, 0xd3, 0x30,
	0x18, 0xc6, 0xe3, 0xa4, 0x4b, 0xcb, 0xbb, 0x09, 0x3c, 0x0b, 0x21, 0x0b, 0xa1, 0xa8, 0xea, 0xa9,
	0x87, 0x29, 0x07, 0xe0, 0xc4, 0x2d, 0x85, 0x4c, 0x2b, 0x2c, 0xdd, 0xe4, 0x76, 0x3e, 0x70, 0x33,
	0x8e, 0xc5, 0x22, 0x35, 0xb1, 0xe5, 0x65, 0x88, 0x7e, 0x0b, 0xbe, 0x03, 0x5f, 0x86, 0xe3, 0x8e,
	0x1c, 0x51, 0xfb, 0x45, 0x90, 0x9d, 0xf0, 0x47, 0x05, 0xed, 0xf6, 0xfc, 0xde, 0xe7, 0xf1, 0xfb,
	0xea, 0x7d, 0x65, 0x38, 0x6e, 0xad, 0x68, 0x6e, 0x84, 0x6c, 0x2b, 0xdd, 0xa4, 0xc6, 0xea, 0x56,
	0x93, 0xc8, 0x1a, 0x39, 0xf9, 0x3a, 0x00, 0xbc, 0xfa, 0x63, 0x71, 0xb1, 0xbe, 0x55, 0xe4, 0x31,
	0x1c, 0x18, 0x5b, 0x49, 0x45, 0xd1, 0x18, 0x4d, 0x23, 0xd6, 0x01, 0x79, 0x0a, 0xa3, 0x8f, 0xe2,
	0xe6, 0xbc, 0xaa, 0xab, 0x96, 0x86, 0xde, 0xf8, 0xcd, 0xe4, 0x09, 0xc4, 0xa2, 0xd6, 0xb7, 0x4d,
	0x4b, 0x23, 0xef, 0xf4, 0x44, 0x28, 0x0c, 0x8d, 0xd8, 0xac, 0xb5, 0x28, 0xe9, 0x60, 0x8c, 0xa6,
	0x47, 0xec, 0x17, 0x92, 0x29, 0x3c, 0x52, 0x8d, 0xb4, 0x1b, 0xd3, 0xaa, 0x32, 0xeb, 0x9e, 0x1e,
	0xf8, 0xc4, 0x7e, 0x99, 0x9c, 0xc0, 0xf1, 0xb5, 0xae, 0x75, 0xad, 0xad, 0xb9, 0xae, 0x64, 0x9f,
	0x8d, 0x7d, 0xf6, 0x5f, 0x83, 0xa4, 0x40, 0xfe, 0x2a, 0xce, 0xc4, 0x5a, 0x34, 0x52, 0xd1, 0xa1,
	0x8f, 0xff, 0xc7, 0x21, 0x27, 0x10, 0x6a, 0x43, 0x47, 0x63, 0x34, 0x7d, 0xf8, 0xfc, 0x59, 0x6a,
	0x8d, 0x4c, 0xf7, 0xcf, 0x91, 0x5e, 0x18, 0xa9, 0x4b, 0xc5, 0x42, 0x6d, 0xdc, 0x65, 0xd4, 0xe7,
	0xd6, 0x0a, 0xfa, 0xc0, 0x37, 0xec, 0xc0, 0x6d, 0xe9, 0xc5, 0xbc, 0xa4, 0xd0, 0x6d, 0xd9, 0x23,
	0x79, 0x09, 0xf1, 0xa7, 0x7a, 0xb5, 0x31, 0x8a, 0x1e, 0xde, 0x37, 0x81, 0xfb, 0x0c, 0xeb, 0xb3,
	0x93, 0x25, 0xc4, 0xdd, 0x4c, 0x02, 0x10, 0x2f, 0x2e, 0x58, 0x91, 0x9d, 0xe3, 0xc0, 0xe9, 0xab,
	0xcb, 0x37, 0xd9, 0x2a, 0xc7, 0xc8, 0xe9, 0x53, 0x96, 0xe7, 0xef, 0x73, 0x1c, 0x92, 0x23, 0x18,
	0x5d, 0x2d, 0x7a, 0x8a, 0x9c, 0xb3, 0x7c, 0x37, 0xbf, 0xe4, 0x05, 0x1e, 0x90, 0x43, 0x18, 0x66,
	0xec, 0xf5, 0xd9, 0x9c, 0xe7, 0xb8, 0x9c, 0xbc, 0x82, 0xb8, 0x1b, 0x43, 0x86, 0x10, 0xe5, 0xbc,
	0xc0, 0x81, 0x13, 0x6f, 0x79, 0x81, 0x91, 0x13, 0x67, 0xbc, 0xc0, 0xa1, 0x13, 0x33, 0x5e, 0xe0,
	0xc8, 0x35, 0x5d, 0xb1, 0x6c, 0xb1, 0x3c, 0xcd, 0x19, 0x1e, 0xcc, 0xe8, 0xb7, 0x6d, 0x82, 0xee,
	0xb6, 0x09, 0xfa, 0xb1, 0x4d, 0xd0, 0x97, 0x5d, 0x12, 0xdc, 0xed, 0x92, 0xe0, 0xfb, 0x2e, 0x09,
	0x3e, 0xc4, 0xfe, 0x2f, 0xbd, 0xf8, 0x19, 0x00, 0x00, 0xff, 0xff, 0xbc, 0x2a, 0x6e, 0x6d, 0x60,
	0x02, 0x00, 0x00,
}

func (m *TransactionValue) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TransactionValue) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Price != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintTransaction(dAtA, i, uint64(m.Price))
	}
	if m.GasLimit != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintTransaction(dAtA, i, uint64(m.GasLimit))
	}
	if m.Amount != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintTransaction(dAtA, i, uint64(m.Amount))
	}
	if len(m.Payload) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintTransaction(dAtA, i, uint64(len(m.Payload)))
		i += copy(dAtA[i:], m.Payload)
	}
	if len(m.EncryptedAmount) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintTransaction(dAtA, i, uint64(len(m.EncryptedAmount)))
		i += copy(dAtA[i:], m.EncryptedAmount)
	}
	if len(m.HomomorphicAmount) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintTransaction(dAtA, i, uint64(len(m.HomomorphicAmount)))
		i += copy(dAtA[i:], m.HomomorphicAmount)
	}
	if len(m.HomomorphicBalance) > 0 {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintTransaction(dAtA, i, uint64(len(m.HomomorphicBalance)))
		i += copy(dAtA[i:], m.HomomorphicBalance)
	}
	if m.Op != 0 {
		dAtA[i] = 0x40
		i++
		i = encodeVarintTransaction(dAtA, i, uint64(m.Op))
	}
	if len(m.Extra) > 0 {
		dAtA[i] = 0x4a
		i++
		i = encodeVarintTransaction(dAtA, i, uint64(len(m.Extra)))
		i += copy(dAtA[i:], m.Extra)
	}
	if len(m.ExtraId) > 0 {
		dAtA[i] = 0x52
		i++
		i = encodeVarintTransaction(dAtA, i, uint64(len(m.ExtraId)))
		i += copy(dAtA[i:], m.ExtraId)
	}
	if m.VmType != 0 {
		dAtA[i] = 0x58
		i++
		i = encodeVarintTransaction(dAtA, i, uint64(m.VmType))
	}
	return i, nil
}

func encodeVarintTransaction(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *TransactionValue) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Price != 0 {
		n += 1 + sovTransaction(uint64(m.Price))
	}
	if m.GasLimit != 0 {
		n += 1 + sovTransaction(uint64(m.GasLimit))
	}
	if m.Amount != 0 {
		n += 1 + sovTransaction(uint64(m.Amount))
	}
	l = len(m.Payload)
	if l > 0 {
		n += 1 + l + sovTransaction(uint64(l))
	}
	l = len(m.EncryptedAmount)
	if l > 0 {
		n += 1 + l + sovTransaction(uint64(l))
	}
	l = len(m.HomomorphicAmount)
	if l > 0 {
		n += 1 + l + sovTransaction(uint64(l))
	}
	l = len(m.HomomorphicBalance)
	if l > 0 {
		n += 1 + l + sovTransaction(uint64(l))
	}
	if m.Op != 0 {
		n += 1 + sovTransaction(uint64(m.Op))
	}
	l = len(m.Extra)
	if l > 0 {
		n += 1 + l + sovTransaction(uint64(l))
	}
	l = len(m.ExtraId)
	if l > 0 {
		n += 1 + l + sovTransaction(uint64(l))
	}
	if m.VmType != 0 {
		n += 1 + sovTransaction(uint64(m.VmType))
	}
	return n
}

func sovTransaction(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozTransaction(x uint64) (n int) {
	return sovTransaction(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TransactionValue) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransaction
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TransactionValue: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TransactionValue: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			m.Price = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransaction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Price |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasLimit", wireType)
			}
			m.GasLimit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransaction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasLimit |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransaction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransaction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTransaction
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTransaction
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
			if m.Payload == nil {
				m.Payload = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EncryptedAmount", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransaction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTransaction
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTransaction
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EncryptedAmount = append(m.EncryptedAmount[:0], dAtA[iNdEx:postIndex]...)
			if m.EncryptedAmount == nil {
				m.EncryptedAmount = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HomomorphicAmount", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransaction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTransaction
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTransaction
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HomomorphicAmount = append(m.HomomorphicAmount[:0], dAtA[iNdEx:postIndex]...)
			if m.HomomorphicAmount == nil {
				m.HomomorphicAmount = []byte{}
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HomomorphicBalance", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransaction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTransaction
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTransaction
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HomomorphicBalance = append(m.HomomorphicBalance[:0], dAtA[iNdEx:postIndex]...)
			if m.HomomorphicBalance == nil {
				m.HomomorphicBalance = []byte{}
			}
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Op", wireType)
			}
			m.Op = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransaction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Op |= TransactionValue_Opcode(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Extra", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransaction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTransaction
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTransaction
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Extra = append(m.Extra[:0], dAtA[iNdEx:postIndex]...)
			if m.Extra == nil {
				m.Extra = []byte{}
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtraId", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransaction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTransaction
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTransaction
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExtraId = append(m.ExtraId[:0], dAtA[iNdEx:postIndex]...)
			if m.ExtraId == nil {
				m.ExtraId = []byte{}
			}
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VMType", wireType)
			}
			m.VmType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransaction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VmType |= TransactionValue_VmType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTransaction(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTransaction
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTransaction
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
func skipTransaction(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTransaction
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
					return 0, ErrIntOverflowTransaction
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
					return 0, ErrIntOverflowTransaction
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
			if length < 0 {
				return 0, ErrInvalidLengthTransaction
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthTransaction
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowTransaction
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
				next, err := skipTransaction(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthTransaction
				}
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
	ErrInvalidLengthTransaction = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTransaction   = fmt.Errorf("proto: integer overflow")
)
