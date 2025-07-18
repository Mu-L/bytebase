// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: v1/annotation.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AuthMethod int32

const (
	AuthMethod_AUTH_METHOD_UNSPECIFIED AuthMethod = 0
	// IAM uses the standard IAM authorization check on the organizational resources.
	AuthMethod_IAM AuthMethod = 1
	// Custom authorization method.
	AuthMethod_CUSTOM AuthMethod = 2
)

// Enum value maps for AuthMethod.
var (
	AuthMethod_name = map[int32]string{
		0: "AUTH_METHOD_UNSPECIFIED",
		1: "IAM",
		2: "CUSTOM",
	}
	AuthMethod_value = map[string]int32{
		"AUTH_METHOD_UNSPECIFIED": 0,
		"IAM":                     1,
		"CUSTOM":                  2,
	}
)

func (x AuthMethod) Enum() *AuthMethod {
	p := new(AuthMethod)
	*p = x
	return p
}

func (x AuthMethod) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AuthMethod) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_annotation_proto_enumTypes[0].Descriptor()
}

func (AuthMethod) Type() protoreflect.EnumType {
	return &file_v1_annotation_proto_enumTypes[0]
}

func (x AuthMethod) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AuthMethod.Descriptor instead.
func (AuthMethod) EnumDescriptor() ([]byte, []int) {
	return file_v1_annotation_proto_rawDescGZIP(), []int{0}
}

var file_v1_annotation_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         100000,
		Name:          "bytebase.v1.allow_without_credential",
		Tag:           "varint,100000,opt,name=allow_without_credential",
		Filename:      "v1/annotation.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         100001,
		Name:          "bytebase.v1.permission",
		Tag:           "bytes,100001,opt,name=permission",
		Filename:      "v1/annotation.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*AuthMethod)(nil),
		Field:         100002,
		Name:          "bytebase.v1.auth_method",
		Tag:           "varint,100002,opt,name=auth_method,enum=bytebase.v1.AuthMethod",
		Filename:      "v1/annotation.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         100003,
		Name:          "bytebase.v1.audit",
		Tag:           "varint,100003,opt,name=audit",
		Filename:      "v1/annotation.proto",
	},
}

// Extension fields to descriptorpb.MethodOptions.
var (
	// optional bool allow_without_credential = 100000;
	E_AllowWithoutCredential = &file_v1_annotation_proto_extTypes[0]
	// optional string permission = 100001;
	E_Permission = &file_v1_annotation_proto_extTypes[1]
	// optional bytebase.v1.AuthMethod auth_method = 100002;
	E_AuthMethod = &file_v1_annotation_proto_extTypes[2]
	// optional bool audit = 100003;
	E_Audit = &file_v1_annotation_proto_extTypes[3]
)

var File_v1_annotation_proto protoreflect.FileDescriptor

const file_v1_annotation_proto_rawDesc = "" +
	"\n" +
	"\x13v1/annotation.proto\x12\vbytebase.v1\x1a google/protobuf/descriptor.proto*>\n" +
	"\n" +
	"AuthMethod\x12\x1b\n" +
	"\x17AUTH_METHOD_UNSPECIFIED\x10\x00\x12\a\n" +
	"\x03IAM\x10\x01\x12\n" +
	"\n" +
	"\x06CUSTOM\x10\x02:Z\n" +
	"\x18allow_without_credential\x12\x1e.google.protobuf.MethodOptions\x18\xa0\x8d\x06 \x01(\bR\x16allowWithoutCredential:@\n" +
	"\n" +
	"permission\x12\x1e.google.protobuf.MethodOptions\x18\xa1\x8d\x06 \x01(\tR\n" +
	"permission:Z\n" +
	"\vauth_method\x12\x1e.google.protobuf.MethodOptions\x18\xa2\x8d\x06 \x01(\x0e2\x17.bytebase.v1.AuthMethodR\n" +
	"authMethod:6\n" +
	"\x05audit\x12\x1e.google.protobuf.MethodOptions\x18\xa3\x8d\x06 \x01(\bR\x05auditB6Z4github.com/bytebase/bytebase/backend/generated-go/v1b\x06proto3"

var (
	file_v1_annotation_proto_rawDescOnce sync.Once
	file_v1_annotation_proto_rawDescData []byte
)

func file_v1_annotation_proto_rawDescGZIP() []byte {
	file_v1_annotation_proto_rawDescOnce.Do(func() {
		file_v1_annotation_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_annotation_proto_rawDesc), len(file_v1_annotation_proto_rawDesc)))
	})
	return file_v1_annotation_proto_rawDescData
}

var file_v1_annotation_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_v1_annotation_proto_goTypes = []any{
	(AuthMethod)(0),                    // 0: bytebase.v1.AuthMethod
	(*descriptorpb.MethodOptions)(nil), // 1: google.protobuf.MethodOptions
}
var file_v1_annotation_proto_depIdxs = []int32{
	1, // 0: bytebase.v1.allow_without_credential:extendee -> google.protobuf.MethodOptions
	1, // 1: bytebase.v1.permission:extendee -> google.protobuf.MethodOptions
	1, // 2: bytebase.v1.auth_method:extendee -> google.protobuf.MethodOptions
	1, // 3: bytebase.v1.audit:extendee -> google.protobuf.MethodOptions
	0, // 4: bytebase.v1.auth_method:type_name -> bytebase.v1.AuthMethod
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	4, // [4:5] is the sub-list for extension type_name
	0, // [0:4] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_v1_annotation_proto_init() }
func file_v1_annotation_proto_init() {
	if File_v1_annotation_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_annotation_proto_rawDesc), len(file_v1_annotation_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 4,
			NumServices:   0,
		},
		GoTypes:           file_v1_annotation_proto_goTypes,
		DependencyIndexes: file_v1_annotation_proto_depIdxs,
		EnumInfos:         file_v1_annotation_proto_enumTypes,
		ExtensionInfos:    file_v1_annotation_proto_extTypes,
	}.Build()
	File_v1_annotation_proto = out.File
	file_v1_annotation_proto_goTypes = nil
	file_v1_annotation_proto_depIdxs = nil
}
