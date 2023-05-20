// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: api/inhouse/v1/entgen.proto

package v1

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

// EnumField is a field of an enum.
type EnumField struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// type is the type of the enum field.
	// e.g. v1alpha.MachineType
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// name is the name of the enum field.
	// e.g. MachineType
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *EnumField) Reset() {
	*x = EnumField{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_inhouse_v1_entgen_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnumField) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnumField) ProtoMessage() {}

func (x *EnumField) ProtoReflect() protoreflect.Message {
	mi := &file_api_inhouse_v1_entgen_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnumField.ProtoReflect.Descriptor instead.
func (*EnumField) Descriptor() ([]byte, []int) {
	return file_api_inhouse_v1_entgen_proto_rawDescGZIP(), []int{0}
}

func (x *EnumField) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *EnumField) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Ent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// name is the name of the ent resource.
	// e.g. Cluster
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// import_path is the import path of the proto resource.
	// e.g. github.com/nokamoto/2pf23/pkg/api/ke/v1alpha
	ImportPath *ImportPath `protobuf:"bytes,2,opt,name=import_path,json=importPath,proto3" json:"import_path,omitempty"`
	// fields is the list of fields of the resource.
	// e.g. DisplayName, NumNodes
	Fields []string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
	// enum_fields is the list of enum fields of the resource.
	EnumFields []*EnumField `protobuf:"bytes,4,rep,name=enum_fields,json=enumFields,proto3" json:"enum_fields,omitempty"`
}

func (x *Ent) Reset() {
	*x = Ent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_inhouse_v1_entgen_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ent) ProtoMessage() {}

func (x *Ent) ProtoReflect() protoreflect.Message {
	mi := &file_api_inhouse_v1_entgen_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ent.ProtoReflect.Descriptor instead.
func (*Ent) Descriptor() ([]byte, []int) {
	return file_api_inhouse_v1_entgen_proto_rawDescGZIP(), []int{1}
}

func (x *Ent) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Ent) GetImportPath() *ImportPath {
	if x != nil {
		return x.ImportPath
	}
	return nil
}

func (x *Ent) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *Ent) GetEnumFields() []*EnumField {
	if x != nil {
		return x.EnumFields
	}
	return nil
}

var File_api_inhouse_v1_entgen_proto protoreflect.FileDescriptor

var file_api_inhouse_v1_entgen_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x69, 0x6e, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2f, 0x76, 0x31,
	0x2f, 0x65, 0x6e, 0x74, 0x67, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x61,
	0x70, 0x69, 0x2e, 0x69, 0x6e, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x18, 0x61,
	0x70, 0x69, 0x2f, 0x69, 0x6e, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x65,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x33, 0x0a, 0x09, 0x45, 0x6e, 0x75, 0x6d, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0xaa, 0x01, 0x0a,
	0x03, 0x45, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x69, 0x6d, 0x70, 0x6f,
	0x72, 0x74, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x69, 0x6e, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x49,
	0x6d, 0x70, 0x6f, 0x72, 0x74, 0x50, 0x61, 0x74, 0x68, 0x52, 0x0a, 0x69, 0x6d, 0x70, 0x6f, 0x72,
	0x74, 0x50, 0x61, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x3a, 0x0a,
	0x0b, 0x65, 0x6e, 0x75, 0x6d, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x69, 0x6e, 0x68, 0x6f, 0x75, 0x73, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x0a, 0x65,
	0x6e, 0x75, 0x6d, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x6f, 0x6b, 0x61, 0x6d, 0x6f, 0x74, 0x6f,
	0x2f, 0x32, 0x70, 0x66, 0x32, 0x33, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x69,
	0x6e, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_api_inhouse_v1_entgen_proto_rawDescOnce sync.Once
	file_api_inhouse_v1_entgen_proto_rawDescData = file_api_inhouse_v1_entgen_proto_rawDesc
)

func file_api_inhouse_v1_entgen_proto_rawDescGZIP() []byte {
	file_api_inhouse_v1_entgen_proto_rawDescOnce.Do(func() {
		file_api_inhouse_v1_entgen_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_inhouse_v1_entgen_proto_rawDescData)
	})
	return file_api_inhouse_v1_entgen_proto_rawDescData
}

var file_api_inhouse_v1_entgen_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_inhouse_v1_entgen_proto_goTypes = []interface{}{
	(*EnumField)(nil),  // 0: api.inhouse.v1.EnumField
	(*Ent)(nil),        // 1: api.inhouse.v1.Ent
	(*ImportPath)(nil), // 2: api.inhouse.v1.ImportPath
}
var file_api_inhouse_v1_entgen_proto_depIdxs = []int32{
	2, // 0: api.inhouse.v1.Ent.import_path:type_name -> api.inhouse.v1.ImportPath
	0, // 1: api.inhouse.v1.Ent.enum_fields:type_name -> api.inhouse.v1.EnumField
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_inhouse_v1_entgen_proto_init() }
func file_api_inhouse_v1_entgen_proto_init() {
	if File_api_inhouse_v1_entgen_proto != nil {
		return
	}
	file_api_inhouse_v1_gen_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_inhouse_v1_entgen_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnumField); i {
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
		file_api_inhouse_v1_entgen_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ent); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_inhouse_v1_entgen_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_inhouse_v1_entgen_proto_goTypes,
		DependencyIndexes: file_api_inhouse_v1_entgen_proto_depIdxs,
		MessageInfos:      file_api_inhouse_v1_entgen_proto_msgTypes,
	}.Build()
	File_api_inhouse_v1_entgen_proto = out.File
	file_api_inhouse_v1_entgen_proto_rawDesc = nil
	file_api_inhouse_v1_entgen_proto_goTypes = nil
	file_api_inhouse_v1_entgen_proto_depIdxs = nil
}
