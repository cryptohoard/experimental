// Code generated by protoc-gen-go. DO NOT EDIT.
// source: betservice.proto

/*
Package betservice is a generated protocol buffer package.

It is generated from these files:
	betservice.proto

It has these top-level messages:
	BetServiceRequest
	BetCancelRequest
	Bet
	BetFilter
	BetServiceResponse
*/
package betservice

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

type Product int32

const (
	Product_PRODUCT_NONE Product = 0
	Product_BTCUSD       Product = 1
	Product_BCHUSD       Product = 2
	Product_ETHUSD       Product = 3
	Product_LTCUSD       Product = 4
)

var Product_name = map[int32]string{
	0: "PRODUCT_NONE",
	1: "BTCUSD",
	2: "BCHUSD",
	3: "ETHUSD",
	4: "LTCUSD",
}
var Product_value = map[string]int32{
	"PRODUCT_NONE": 0,
	"BTCUSD":       1,
	"BCHUSD":       2,
	"ETHUSD":       3,
	"LTCUSD":       4,
}

func (x Product) String() string {
	return proto.EnumName(Product_name, int32(x))
}
func (Product) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type State int32

const (
	State_STATE_NONE State = 0
	State_PROCESSING State = 1
	State_PLACED     State = 2
	State_CASHOUT    State = 3
	State_PAYOUT     State = 4
)

var State_name = map[int32]string{
	0: "STATE_NONE",
	1: "PROCESSING",
	2: "PLACED",
	3: "CASHOUT",
	4: "PAYOUT",
}
var State_value = map[string]int32{
	"STATE_NONE": 0,
	"PROCESSING": 1,
	"PLACED":     2,
	"CASHOUT":    3,
	"PAYOUT":     4,
}

func (x State) String() string {
	return proto.EnumName(State_name, int32(x))
}
func (State) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type BetServiceRequest struct {
	BetId         string  `protobuf:"bytes,1,opt,name=bet_id,json=betId" json:"bet_id,omitempty"`
	CustomerId    string  `protobuf:"bytes,2,opt,name=customer_id,json=customerId" json:"customer_id,omitempty"`
	Exchange      string  `protobuf:"bytes,3,opt,name=exchange" json:"exchange,omitempty"`
	Product       Product `protobuf:"varint,4,opt,name=product,enum=betservice.Product" json:"product,omitempty"`
	Amount        int32   `protobuf:"varint,5,opt,name=amount" json:"amount,omitempty"`
	ProfitPercent int32   `protobuf:"varint,6,opt,name=profit_percent,json=profitPercent" json:"profit_percent,omitempty"`
	LossPercent   int32   `protobuf:"varint,7,opt,name=loss_percent,json=lossPercent" json:"loss_percent,omitempty"`
}

func (m *BetServiceRequest) Reset()                    { *m = BetServiceRequest{} }
func (m *BetServiceRequest) String() string            { return proto.CompactTextString(m) }
func (*BetServiceRequest) ProtoMessage()               {}
func (*BetServiceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *BetServiceRequest) GetBetId() string {
	if m != nil {
		return m.BetId
	}
	return ""
}

func (m *BetServiceRequest) GetCustomerId() string {
	if m != nil {
		return m.CustomerId
	}
	return ""
}

func (m *BetServiceRequest) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *BetServiceRequest) GetProduct() Product {
	if m != nil {
		return m.Product
	}
	return Product_PRODUCT_NONE
}

func (m *BetServiceRequest) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *BetServiceRequest) GetProfitPercent() int32 {
	if m != nil {
		return m.ProfitPercent
	}
	return 0
}

func (m *BetServiceRequest) GetLossPercent() int32 {
	if m != nil {
		return m.LossPercent
	}
	return 0
}

type BetCancelRequest struct {
	BetId string `protobuf:"bytes,1,opt,name=bet_id,json=betId" json:"bet_id,omitempty"`
}

func (m *BetCancelRequest) Reset()                    { *m = BetCancelRequest{} }
func (m *BetCancelRequest) String() string            { return proto.CompactTextString(m) }
func (*BetCancelRequest) ProtoMessage()               {}
func (*BetCancelRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *BetCancelRequest) GetBetId() string {
	if m != nil {
		return m.BetId
	}
	return ""
}

type Bet struct {
	BetId          string                     `protobuf:"bytes,1,opt,name=bet_id,json=betId" json:"bet_id,omitempty"`
	CustomerId     string                     `protobuf:"bytes,2,opt,name=customer_id,json=customerId" json:"customer_id,omitempty"`
	Exchange       string                     `protobuf:"bytes,3,opt,name=exchange" json:"exchange,omitempty"`
	Product        Product                    `protobuf:"varint,4,opt,name=product,enum=betservice.Product" json:"product,omitempty"`
	InitialAmount  int32                      `protobuf:"varint,5,opt,name=initial_amount,json=initialAmount" json:"initial_amount,omitempty"`
	CurrentAmount  int32                      `protobuf:"varint,6,opt,name=current_amount,json=currentAmount" json:"current_amount,omitempty"`
	CryptoCurrency float32                    `protobuf:"fixed32,7,opt,name=crypto_currency,json=cryptoCurrency" json:"crypto_currency,omitempty"`
	ProfitPercent  int32                      `protobuf:"varint,8,opt,name=profit_percent,json=profitPercent" json:"profit_percent,omitempty"`
	LossPercent    int32                      `protobuf:"varint,9,opt,name=loss_percent,json=lossPercent" json:"loss_percent,omitempty"`
	State          State                      `protobuf:"varint,10,opt,name=state,enum=betservice.State" json:"state,omitempty"`
	CreationTime   *google_protobuf.Timestamp `protobuf:"bytes,11,opt,name=creation_time,json=creationTime" json:"creation_time,omitempty"`
	ClosedTime     *google_protobuf.Timestamp `protobuf:"bytes,12,opt,name=closed_time,json=closedTime" json:"closed_time,omitempty"`
}

func (m *Bet) Reset()                    { *m = Bet{} }
func (m *Bet) String() string            { return proto.CompactTextString(m) }
func (*Bet) ProtoMessage()               {}
func (*Bet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Bet) GetBetId() string {
	if m != nil {
		return m.BetId
	}
	return ""
}

func (m *Bet) GetCustomerId() string {
	if m != nil {
		return m.CustomerId
	}
	return ""
}

func (m *Bet) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *Bet) GetProduct() Product {
	if m != nil {
		return m.Product
	}
	return Product_PRODUCT_NONE
}

func (m *Bet) GetInitialAmount() int32 {
	if m != nil {
		return m.InitialAmount
	}
	return 0
}

func (m *Bet) GetCurrentAmount() int32 {
	if m != nil {
		return m.CurrentAmount
	}
	return 0
}

func (m *Bet) GetCryptoCurrency() float32 {
	if m != nil {
		return m.CryptoCurrency
	}
	return 0
}

func (m *Bet) GetProfitPercent() int32 {
	if m != nil {
		return m.ProfitPercent
	}
	return 0
}

func (m *Bet) GetLossPercent() int32 {
	if m != nil {
		return m.LossPercent
	}
	return 0
}

func (m *Bet) GetState() State {
	if m != nil {
		return m.State
	}
	return State_STATE_NONE
}

func (m *Bet) GetCreationTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.CreationTime
	}
	return nil
}

func (m *Bet) GetClosedTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.ClosedTime
	}
	return nil
}

type BetFilter struct {
	Product Product `protobuf:"varint,1,opt,name=product,enum=betservice.Product" json:"product,omitempty"`
	State   State   `protobuf:"varint,2,opt,name=state,enum=betservice.State" json:"state,omitempty"`
}

func (m *BetFilter) Reset()                    { *m = BetFilter{} }
func (m *BetFilter) String() string            { return proto.CompactTextString(m) }
func (*BetFilter) ProtoMessage()               {}
func (*BetFilter) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *BetFilter) GetProduct() Product {
	if m != nil {
		return m.Product
	}
	return Product_PRODUCT_NONE
}

func (m *BetFilter) GetState() State {
	if m != nil {
		return m.State
	}
	return State_STATE_NONE
}

type BetServiceResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
	Bet     *Bet `protobuf:"bytes,2,opt,name=bet" json:"bet,omitempty"`
}

func (m *BetServiceResponse) Reset()                    { *m = BetServiceResponse{} }
func (m *BetServiceResponse) String() string            { return proto.CompactTextString(m) }
func (*BetServiceResponse) ProtoMessage()               {}
func (*BetServiceResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *BetServiceResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *BetServiceResponse) GetBet() *Bet {
	if m != nil {
		return m.Bet
	}
	return nil
}

func init() {
	proto.RegisterType((*BetServiceRequest)(nil), "betservice.BetServiceRequest")
	proto.RegisterType((*BetCancelRequest)(nil), "betservice.BetCancelRequest")
	proto.RegisterType((*Bet)(nil), "betservice.Bet")
	proto.RegisterType((*BetFilter)(nil), "betservice.BetFilter")
	proto.RegisterType((*BetServiceResponse)(nil), "betservice.BetServiceResponse")
	proto.RegisterEnum("betservice.Product", Product_name, Product_value)
	proto.RegisterEnum("betservice.State", State_name, State_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for BetService service

type BetServiceClient interface {
	ListBets(ctx context.Context, in *BetFilter, opts ...grpc.CallOption) (BetService_ListBetsClient, error)
	CreateBet(ctx context.Context, in *BetServiceRequest, opts ...grpc.CallOption) (*BetServiceResponse, error)
	CancelBet(ctx context.Context, in *BetCancelRequest, opts ...grpc.CallOption) (*BetServiceResponse, error)
}

type betServiceClient struct {
	cc *grpc.ClientConn
}

func NewBetServiceClient(cc *grpc.ClientConn) BetServiceClient {
	return &betServiceClient{cc}
}

func (c *betServiceClient) ListBets(ctx context.Context, in *BetFilter, opts ...grpc.CallOption) (BetService_ListBetsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_BetService_serviceDesc.Streams[0], c.cc, "/betservice.BetService/ListBets", opts...)
	if err != nil {
		return nil, err
	}
	x := &betServiceListBetsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BetService_ListBetsClient interface {
	Recv() (*Bet, error)
	grpc.ClientStream
}

type betServiceListBetsClient struct {
	grpc.ClientStream
}

func (x *betServiceListBetsClient) Recv() (*Bet, error) {
	m := new(Bet)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *betServiceClient) CreateBet(ctx context.Context, in *BetServiceRequest, opts ...grpc.CallOption) (*BetServiceResponse, error) {
	out := new(BetServiceResponse)
	err := grpc.Invoke(ctx, "/betservice.BetService/CreateBet", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *betServiceClient) CancelBet(ctx context.Context, in *BetCancelRequest, opts ...grpc.CallOption) (*BetServiceResponse, error) {
	out := new(BetServiceResponse)
	err := grpc.Invoke(ctx, "/betservice.BetService/CancelBet", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for BetService service

type BetServiceServer interface {
	ListBets(*BetFilter, BetService_ListBetsServer) error
	CreateBet(context.Context, *BetServiceRequest) (*BetServiceResponse, error)
	CancelBet(context.Context, *BetCancelRequest) (*BetServiceResponse, error)
}

func RegisterBetServiceServer(s *grpc.Server, srv BetServiceServer) {
	s.RegisterService(&_BetService_serviceDesc, srv)
}

func _BetService_ListBets_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(BetFilter)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BetServiceServer).ListBets(m, &betServiceListBetsServer{stream})
}

type BetService_ListBetsServer interface {
	Send(*Bet) error
	grpc.ServerStream
}

type betServiceListBetsServer struct {
	grpc.ServerStream
}

func (x *betServiceListBetsServer) Send(m *Bet) error {
	return x.ServerStream.SendMsg(m)
}

func _BetService_CreateBet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BetServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BetServiceServer).CreateBet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/betservice.BetService/CreateBet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BetServiceServer).CreateBet(ctx, req.(*BetServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BetService_CancelBet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BetCancelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BetServiceServer).CancelBet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/betservice.BetService/CancelBet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BetServiceServer).CancelBet(ctx, req.(*BetCancelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _BetService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "betservice.BetService",
	HandlerType: (*BetServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateBet",
			Handler:    _BetService_CreateBet_Handler,
		},
		{
			MethodName: "CancelBet",
			Handler:    _BetService_CancelBet_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListBets",
			Handler:       _BetService_ListBets_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "betservice.proto",
}

func init() { proto.RegisterFile("betservice.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 626 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x54, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xad, 0x93, 0xe6, 0x6b, 0x9c, 0xa6, 0xe9, 0xa2, 0x22, 0x2b, 0x02, 0x9a, 0x5a, 0xaa, 0x1a,
	0x2a, 0x91, 0xa2, 0x20, 0x71, 0xe1, 0x80, 0x62, 0x37, 0xd0, 0xaa, 0x51, 0x13, 0x6c, 0xf7, 0xc0,
	0x29, 0x72, 0x36, 0xd3, 0x62, 0x29, 0xb1, 0x8d, 0x77, 0x8d, 0xe8, 0xaf, 0xe5, 0x2f, 0x70, 0xe4,
	0x88, 0x76, 0xd7, 0x4e, 0xd3, 0x00, 0x6d, 0x8f, 0xdc, 0x66, 0xdf, 0x7b, 0x3b, 0xde, 0xf7, 0x76,
	0xbc, 0xd0, 0x9c, 0x22, 0x67, 0x98, 0x7c, 0x0b, 0x28, 0x76, 0xe3, 0x24, 0xe2, 0x11, 0x81, 0x5b,
	0xa4, 0xb5, 0x77, 0x1d, 0x45, 0xd7, 0x73, 0x3c, 0x96, 0xcc, 0x34, 0xbd, 0x3a, 0xe6, 0xc1, 0x02,
	0x19, 0xf7, 0x17, 0xb1, 0x12, 0x9b, 0xbf, 0x34, 0xd8, 0xb1, 0x90, 0xbb, 0x4a, 0xef, 0xe0, 0xd7,
	0x14, 0x19, 0x27, 0xbb, 0x50, 0x9e, 0x22, 0x9f, 0x04, 0x33, 0x43, 0x6b, 0x6b, 0x9d, 0x9a, 0x53,
	0x9a, 0x22, 0x3f, 0x9b, 0x91, 0x3d, 0xd0, 0x69, 0xca, 0x78, 0xb4, 0xc0, 0x44, 0x70, 0x05, 0xc9,
	0x41, 0x0e, 0x9d, 0xcd, 0x48, 0x0b, 0xaa, 0xf8, 0x9d, 0x7e, 0xf1, 0xc3, 0x6b, 0x34, 0x8a, 0x92,
	0x5d, 0xae, 0xc9, 0x2b, 0xa8, 0xc4, 0x49, 0x34, 0x4b, 0x29, 0x37, 0x36, 0xdb, 0x5a, 0xa7, 0xd1,
	0x7b, 0xd2, 0x5d, 0x39, 0xfa, 0x58, 0x51, 0x4e, 0xae, 0x21, 0x4f, 0xa1, 0xec, 0x2f, 0xa2, 0x34,
	0xe4, 0x46, 0xa9, 0xad, 0x75, 0x4a, 0x4e, 0xb6, 0x22, 0x07, 0xd0, 0x88, 0x93, 0xe8, 0x2a, 0xe0,
	0x93, 0x18, 0x13, 0x8a, 0x21, 0x37, 0xca, 0x92, 0xdf, 0x52, 0xe8, 0x58, 0x81, 0x64, 0x1f, 0xea,
	0xf3, 0x88, 0xb1, 0xa5, 0xa8, 0x22, 0x45, 0xba, 0xc0, 0x32, 0x89, 0xf9, 0x12, 0x9a, 0x16, 0x72,
	0xdb, 0x0f, 0x29, 0xce, 0xef, 0x37, 0x6e, 0xfe, 0x2c, 0x42, 0xd1, 0xc2, 0xff, 0x22, 0x97, 0x03,
	0x68, 0x04, 0x61, 0xc0, 0x03, 0x7f, 0x3e, 0xb9, 0x93, 0xcf, 0x56, 0x86, 0xf6, 0x97, 0x31, 0xd1,
	0x34, 0x49, 0x30, 0xe4, 0xb9, 0x2c, 0x8b, 0x29, 0x43, 0x33, 0xd9, 0x21, 0x6c, 0xd3, 0xe4, 0x26,
	0xe6, 0xd1, 0x44, 0xe1, 0xf4, 0x46, 0x26, 0x55, 0x70, 0x1a, 0x0a, 0xb6, 0x33, 0xf4, 0x2f, 0xb1,
	0x57, 0x1f, 0x13, 0x7b, 0xed, 0x8f, 0xd8, 0xc9, 0x21, 0x94, 0x18, 0xf7, 0x39, 0x1a, 0x20, 0xdd,
	0xee, 0xac, 0xba, 0x75, 0x05, 0xe1, 0x28, 0x9e, 0xbc, 0x87, 0x2d, 0x9a, 0xa0, 0xcf, 0x83, 0x28,
	0x9c, 0x88, 0xb1, 0x35, 0xf4, 0xb6, 0xd6, 0xd1, 0x7b, 0xad, 0xae, 0x9a, 0xe9, 0x6e, 0x3e, 0xd3,
	0x5d, 0x2f, 0x9f, 0x69, 0xa7, 0x9e, 0x6f, 0x10, 0x10, 0x79, 0x07, 0x3a, 0x9d, 0x47, 0x0c, 0x67,
	0x6a, 0x7b, 0xfd, 0xc1, 0xed, 0xa0, 0xe4, 0x02, 0x30, 0x29, 0xd4, 0x2c, 0xe4, 0x1f, 0x82, 0x39,
	0xc7, 0x64, 0xf5, 0x8e, 0xb4, 0x47, 0xdc, 0xd1, 0xd2, 0x62, 0xe1, 0x7e, 0x8b, 0xe6, 0x27, 0x20,
	0xab, 0x3f, 0x1f, 0x8b, 0xa3, 0x90, 0x21, 0x31, 0xa0, 0xc2, 0x52, 0x4a, 0x91, 0x31, 0xf9, 0xb5,
	0xaa, 0x93, 0x2f, 0xc9, 0x3e, 0x14, 0xa7, 0xc8, 0x65, 0x5b, 0xbd, 0xb7, 0xbd, 0xda, 0xd6, 0x42,
	0xee, 0x08, 0xee, 0xe8, 0x1c, 0x2a, 0xd9, 0x79, 0x48, 0x13, 0xea, 0x63, 0x67, 0x74, 0x72, 0x69,
	0x7b, 0x93, 0x8b, 0xd1, 0xc5, 0xa0, 0xb9, 0x41, 0x00, 0xca, 0x96, 0x67, 0x5f, 0xba, 0x27, 0x4d,
	0x4d, 0xd6, 0xf6, 0xa9, 0xa8, 0x0b, 0xa2, 0x1e, 0x78, 0xb2, 0x2e, 0x8a, 0x7a, 0xa8, 0x34, 0x9b,
	0x47, 0x43, 0x28, 0xc9, 0xf3, 0x92, 0x06, 0x80, 0xeb, 0xf5, 0xbd, 0x41, 0xde, 0xa8, 0x01, 0x30,
	0x76, 0x46, 0xf6, 0xc0, 0x75, 0xcf, 0x2e, 0x3e, 0xaa, 0x66, 0xe3, 0x61, 0xdf, 0x1e, 0x88, 0x66,
	0x3a, 0x54, 0xec, 0xbe, 0x7b, 0x3a, 0xba, 0xf4, 0x54, 0xb7, 0x71, 0xff, 0xb3, 0xa8, 0x37, 0x7b,
	0x3f, 0x34, 0x80, 0x5b, 0xbb, 0xe4, 0x2d, 0x54, 0x87, 0x01, 0xe3, 0x16, 0x72, 0x46, 0x76, 0xd7,
	0xbc, 0xa8, 0xdc, 0x5b, 0xeb, 0x16, 0xcd, 0x8d, 0xd7, 0x1a, 0x19, 0x42, 0xcd, 0x16, 0xd7, 0x8c,
	0xe2, 0x8f, 0x7c, 0xbe, 0xa6, 0xb8, 0xfb, 0x90, 0xb5, 0x5e, 0xfc, 0x8b, 0x56, 0x51, 0x9b, 0x1b,
	0xe4, 0x1c, 0x6a, 0xea, 0x09, 0x10, 0xdd, 0x9e, 0xad, 0xc9, 0xef, 0x3c, 0x0e, 0x0f, 0x37, 0x9b,
	0x96, 0xe5, 0x50, 0xbd, 0xf9, 0x1d, 0x00, 0x00, 0xff, 0xff, 0xcc, 0x7b, 0xe3, 0xa3, 0x95, 0x05,
	0x00, 0x00,
}
