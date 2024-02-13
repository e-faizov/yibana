// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.8
// source: proto/metrics.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Metric_Type int32

const (
	Metric_GAUGE   Metric_Type = 0
	Metric_COUNTER Metric_Type = 1
)

// Enum value maps for Metric_Type.
var (
	Metric_Type_name = map[int32]string{
		0: "GAUGE",
		1: "COUNTER",
	}
	Metric_Type_value = map[string]int32{
		"GAUGE":   0,
		"COUNTER": 1,
	}
)

func (x Metric_Type) Enum() *Metric_Type {
	p := new(Metric_Type)
	*p = x
	return p
}

func (x Metric_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Metric_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_metrics_proto_enumTypes[0].Descriptor()
}

func (Metric_Type) Type() protoreflect.EnumType {
	return &file_proto_metrics_proto_enumTypes[0]
}

func (x Metric_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Metric_Type.Descriptor instead.
func (Metric_Type) EnumDescriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{0, 0}
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	MType Metric_Type `protobuf:"varint,2,opt,name=mType,proto3,enum=demo.Metric_Type" json:"mType,omitempty"`
	Delta *int64      `protobuf:"varint,3,opt,name=delta,proto3,oneof" json:"delta,omitempty"`
	Value *float64    `protobuf:"fixed64,4,opt,name=value,proto3,oneof" json:"value,omitempty"`
	Hash  string      `protobuf:"bytes,5,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metric) GetMType() Metric_Type {
	if x != nil {
		return x.MType
	}
	return Metric_GAUGE
}

func (x *Metric) GetDelta() int64 {
	if x != nil && x.Delta != nil {
		return *x.Delta
	}
	return 0
}

func (x *Metric) GetValue() float64 {
	if x != nil && x.Value != nil {
		return *x.Value
	}
	return 0
}

func (x *Metric) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

type Metrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*Metric `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *Metrics) Reset() {
	*x = Metrics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metrics) ProtoMessage() {}

func (x *Metrics) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metrics.ProtoReflect.Descriptor instead.
func (*Metrics) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{1}
}

func (x *Metrics) GetData() []*Metric {
	if x != nil {
		return x.Data
	}
	return nil
}

type AddMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *AddMetricResponse) Reset() {
	*x = AddMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddMetricResponse) ProtoMessage() {}

func (x *AddMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddMetricResponse.ProtoReflect.Descriptor instead.
func (*AddMetricResponse) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{2}
}

func (x *AddMetricResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type GetMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data  *Metric `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Error string  `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *GetMetricResponse) Reset() {
	*x = GetMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricResponse) ProtoMessage() {}

func (x *GetMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetricResponse.ProtoReflect.Descriptor instead.
func (*GetMetricResponse) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{3}
}

func (x *GetMetricResponse) GetData() *Metric {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *GetMetricResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_proto_metrics_proto protoreflect.FileDescriptor

var file_proto_metrics_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x64, 0x65, 0x6d, 0x6f, 0x22, 0xbf, 0x01, 0x0a, 0x06,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x27, 0x0a, 0x05, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x19, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00,
	0x52, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x48, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x88, 0x01, 0x01, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x22, 0x1e, 0x0a, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x09, 0x0a, 0x05, 0x47, 0x41, 0x55, 0x47, 0x45, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07,
	0x43, 0x4f, 0x55, 0x4e, 0x54, 0x45, 0x52, 0x10, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x64, 0x65,
	0x6c, 0x74, 0x61, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x2b, 0x0a,
	0x07, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x20, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x29, 0x0a, 0x11, 0x41, 0x64,
	0x64, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x4b, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x32, 0xae, 0x01, 0x0a, 0x0e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x09, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x12, 0x0c, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x1a, 0x17, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x0a, 0x41, 0x64, 0x64,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x0d, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x1a, 0x17, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x41, 0x64,
	0x64, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x32, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x0c, 0x2e, 0x64,
	0x65, 0x6d, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x1a, 0x17, 0x2e, 0x64, 0x65, 0x6d,
	0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x10, 0x5a, 0x0e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_metrics_proto_rawDescOnce sync.Once
	file_proto_metrics_proto_rawDescData = file_proto_metrics_proto_rawDesc
)

func file_proto_metrics_proto_rawDescGZIP() []byte {
	file_proto_metrics_proto_rawDescOnce.Do(func() {
		file_proto_metrics_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_metrics_proto_rawDescData)
	})
	return file_proto_metrics_proto_rawDescData
}

var file_proto_metrics_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_metrics_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_metrics_proto_goTypes = []interface{}{
	(Metric_Type)(0),          // 0: demo.Metric.Type
	(*Metric)(nil),            // 1: demo.Metric
	(*Metrics)(nil),           // 2: demo.Metrics
	(*AddMetricResponse)(nil), // 3: demo.AddMetricResponse
	(*GetMetricResponse)(nil), // 4: demo.GetMetricResponse
}
var file_proto_metrics_proto_depIdxs = []int32{
	0, // 0: demo.Metric.mType:type_name -> demo.Metric.Type
	1, // 1: demo.Metrics.data:type_name -> demo.Metric
	1, // 2: demo.GetMetricResponse.data:type_name -> demo.Metric
	1, // 3: demo.MetricsService.AddMetric:input_type -> demo.Metric
	2, // 4: demo.MetricsService.AddMetrics:input_type -> demo.Metrics
	1, // 5: demo.MetricsService.GetMetric:input_type -> demo.Metric
	3, // 6: demo.MetricsService.AddMetric:output_type -> demo.AddMetricResponse
	3, // 7: demo.MetricsService.AddMetrics:output_type -> demo.AddMetricResponse
	4, // 8: demo.MetricsService.GetMetric:output_type -> demo.GetMetricResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_metrics_proto_init() }
func file_proto_metrics_proto_init() {
	if File_proto_metrics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_metrics_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_metrics_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metrics); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_metrics_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddMetricResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_metrics_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMetricResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_proto_metrics_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_metrics_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_metrics_proto_goTypes,
		DependencyIndexes: file_proto_metrics_proto_depIdxs,
		EnumInfos:         file_proto_metrics_proto_enumTypes,
		MessageInfos:      file_proto_metrics_proto_msgTypes,
	}.Build()
	File_proto_metrics_proto = out.File
	file_proto_metrics_proto_rawDesc = nil
	file_proto_metrics_proto_goTypes = nil
	file_proto_metrics_proto_depIdxs = nil
}