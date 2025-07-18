// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: v1/role_service.proto

package v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
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

type Role_Type int32

const (
	Role_TYPE_UNSPECIFIED Role_Type = 0
	Role_BUILT_IN         Role_Type = 1
	Role_CUSTOM           Role_Type = 2
)

// Enum value maps for Role_Type.
var (
	Role_Type_name = map[int32]string{
		0: "TYPE_UNSPECIFIED",
		1: "BUILT_IN",
		2: "CUSTOM",
	}
	Role_Type_value = map[string]int32{
		"TYPE_UNSPECIFIED": 0,
		"BUILT_IN":         1,
		"CUSTOM":           2,
	}
)

func (x Role_Type) Enum() *Role_Type {
	p := new(Role_Type)
	*p = x
	return p
}

func (x Role_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Role_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_role_service_proto_enumTypes[0].Descriptor()
}

func (Role_Type) Type() protoreflect.EnumType {
	return &file_v1_role_service_proto_enumTypes[0]
}

func (x Role_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Role_Type.Descriptor instead.
func (Role_Type) EnumDescriptor() ([]byte, []int) {
	return file_v1_role_service_proto_rawDescGZIP(), []int{6, 0}
}

type ListRolesRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Not used.
	// The maximum number of roles to return. The service may return fewer than
	// this value.
	// If unspecified, at most 10 reviews will be returned.
	// The maximum value is 1000; values above 1000 will be coerced to 1000.
	PageSize int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// Not used.
	// A page token, received from a previous `ListRoles` call.
	// Provide this to retrieve the subsequent page.
	//
	// When paginating, all other parameters provided to `ListRoles` must match
	// the call that provided the page token.
	PageToken     string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListRolesRequest) Reset() {
	*x = ListRolesRequest{}
	mi := &file_v1_role_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRolesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRolesRequest) ProtoMessage() {}

func (x *ListRolesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_role_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRolesRequest.ProtoReflect.Descriptor instead.
func (*ListRolesRequest) Descriptor() ([]byte, []int) {
	return file_v1_role_service_proto_rawDescGZIP(), []int{0}
}

func (x *ListRolesRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListRolesRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

type ListRolesResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Roles []*Role                `protobuf:"bytes,1,rep,name=roles,proto3" json:"roles,omitempty"`
	// A token, which can be sent as `page_token` to retrieve the next page.
	// If this field is omitted, there are no subsequent pages.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListRolesResponse) Reset() {
	*x = ListRolesResponse{}
	mi := &file_v1_role_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRolesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRolesResponse) ProtoMessage() {}

func (x *ListRolesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_role_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRolesResponse.ProtoReflect.Descriptor instead.
func (*ListRolesResponse) Descriptor() ([]byte, []int) {
	return file_v1_role_service_proto_rawDescGZIP(), []int{1}
}

func (x *ListRolesResponse) GetRoles() []*Role {
	if x != nil {
		return x.Roles
	}
	return nil
}

func (x *ListRolesResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type CreateRoleRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Role  *Role                  `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
	// The ID to use for the role, which will become the final component
	// of the role's resource name.
	//
	// This value should be 4-63 characters, and valid characters
	// are /[a-z][A-Z][0-9]/.
	RoleId        string `protobuf:"bytes,2,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateRoleRequest) Reset() {
	*x = CreateRoleRequest{}
	mi := &file_v1_role_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRoleRequest) ProtoMessage() {}

func (x *CreateRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_role_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRoleRequest.ProtoReflect.Descriptor instead.
func (*CreateRoleRequest) Descriptor() ([]byte, []int) {
	return file_v1_role_service_proto_rawDescGZIP(), []int{2}
}

func (x *CreateRoleRequest) GetRole() *Role {
	if x != nil {
		return x.Role
	}
	return nil
}

func (x *CreateRoleRequest) GetRoleId() string {
	if x != nil {
		return x.RoleId
	}
	return ""
}

type GetRoleRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The name of the role to retrieve.
	// Format: roles/{role}
	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRoleRequest) Reset() {
	*x = GetRoleRequest{}
	mi := &file_v1_role_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRoleRequest) ProtoMessage() {}

func (x *GetRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_role_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRoleRequest.ProtoReflect.Descriptor instead.
func (*GetRoleRequest) Descriptor() ([]byte, []int) {
	return file_v1_role_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetRoleRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type UpdateRoleRequest struct {
	state      protoimpl.MessageState `protogen:"open.v1"`
	Role       *Role                  `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
	UpdateMask *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	// If set to true, and the role is not found, a new role will be created.
	AllowMissing  bool `protobuf:"varint,3,opt,name=allow_missing,json=allowMissing,proto3" json:"allow_missing,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateRoleRequest) Reset() {
	*x = UpdateRoleRequest{}
	mi := &file_v1_role_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRoleRequest) ProtoMessage() {}

func (x *UpdateRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_role_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRoleRequest.ProtoReflect.Descriptor instead.
func (*UpdateRoleRequest) Descriptor() ([]byte, []int) {
	return file_v1_role_service_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateRoleRequest) GetRole() *Role {
	if x != nil {
		return x.Role
	}
	return nil
}

func (x *UpdateRoleRequest) GetUpdateMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.UpdateMask
	}
	return nil
}

func (x *UpdateRoleRequest) GetAllowMissing() bool {
	if x != nil {
		return x.AllowMissing
	}
	return false
}

type DeleteRoleRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Format: roles/{role}
	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteRoleRequest) Reset() {
	*x = DeleteRoleRequest{}
	mi := &file_v1_role_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRoleRequest) ProtoMessage() {}

func (x *DeleteRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_role_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRoleRequest.ProtoReflect.Descriptor instead.
func (*DeleteRoleRequest) Descriptor() ([]byte, []int) {
	return file_v1_role_service_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteRoleRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Role struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Format: roles/{role}
	Name          string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Title         string    `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description   string    `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Permissions   []string  `protobuf:"bytes,4,rep,name=permissions,proto3" json:"permissions,omitempty"`
	Type          Role_Type `protobuf:"varint,5,opt,name=type,proto3,enum=bytebase.v1.Role_Type" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Role) Reset() {
	*x = Role{}
	mi := &file_v1_role_service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Role) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Role) ProtoMessage() {}

func (x *Role) ProtoReflect() protoreflect.Message {
	mi := &file_v1_role_service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Role.ProtoReflect.Descriptor instead.
func (*Role) Descriptor() ([]byte, []int) {
	return file_v1_role_service_proto_rawDescGZIP(), []int{6}
}

func (x *Role) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Role) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Role) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Role) GetPermissions() []string {
	if x != nil {
		return x.Permissions
	}
	return nil
}

func (x *Role) GetType() Role_Type {
	if x != nil {
		return x.Type
	}
	return Role_TYPE_UNSPECIFIED
}

var File_v1_role_service_proto protoreflect.FileDescriptor

const file_v1_role_service_proto_rawDesc = "" +
	"\n" +
	"\x15v1/role_service.proto\x12\vbytebase.v1\x1a\x1cgoogle/api/annotations.proto\x1a\x17google/api/client.proto\x1a\x1fgoogle/api/field_behavior.proto\x1a\x19google/api/resource.proto\x1a\x1bgoogle/protobuf/empty.proto\x1a google/protobuf/field_mask.proto\x1a\x13v1/annotation.proto\"N\n" +
	"\x10ListRolesRequest\x12\x1b\n" +
	"\tpage_size\x18\x01 \x01(\x05R\bpageSize\x12\x1d\n" +
	"\n" +
	"page_token\x18\x02 \x01(\tR\tpageToken\"d\n" +
	"\x11ListRolesResponse\x12'\n" +
	"\x05roles\x18\x01 \x03(\v2\x11.bytebase.v1.RoleR\x05roles\x12&\n" +
	"\x0fnext_page_token\x18\x02 \x01(\tR\rnextPageToken\"]\n" +
	"\x11CreateRoleRequest\x12*\n" +
	"\x04role\x18\x01 \x01(\v2\x11.bytebase.v1.RoleB\x03\xe0A\x02R\x04role\x12\x1c\n" +
	"\arole_id\x18\x02 \x01(\tB\x03\xe0A\x02R\x06roleId\"?\n" +
	"\x0eGetRoleRequest\x12-\n" +
	"\x04name\x18\x01 \x01(\tB\x19\xe0A\x02\xfaA\x13\n" +
	"\x11bytebase.com/RoleR\x04name\"\x9c\x01\n" +
	"\x11UpdateRoleRequest\x12%\n" +
	"\x04role\x18\x01 \x01(\v2\x11.bytebase.v1.RoleR\x04role\x12;\n" +
	"\vupdate_mask\x18\x02 \x01(\v2\x1a.google.protobuf.FieldMaskR\n" +
	"updateMask\x12#\n" +
	"\rallow_missing\x18\x03 \x01(\bR\fallowMissing\"B\n" +
	"\x11DeleteRoleRequest\x12-\n" +
	"\x04name\x18\x01 \x01(\tB\x19\xe0A\x02\xfaA\x13\n" +
	"\x11bytebase.com/RoleR\x04name\"\xfe\x01\n" +
	"\x04Role\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x14\n" +
	"\x05title\x18\x02 \x01(\tR\x05title\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x12 \n" +
	"\vpermissions\x18\x04 \x03(\tR\vpermissions\x12*\n" +
	"\x04type\x18\x05 \x01(\x0e2\x16.bytebase.v1.Role.TypeR\x04type\"6\n" +
	"\x04Type\x12\x14\n" +
	"\x10TYPE_UNSPECIFIED\x10\x00\x12\f\n" +
	"\bBUILT_IN\x10\x01\x12\n" +
	"\n" +
	"\x06CUSTOM\x10\x02:$\xeaA!\n" +
	"\x11bytebase.com/Role\x12\froles/{role}2\x84\x05\n" +
	"\vRoleService\x12r\n" +
	"\tListRoles\x12\x1d.bytebase.v1.ListRolesRequest\x1a\x1e.bytebase.v1.ListRolesResponse\"&\x8a\xea0\rbb.roles.list\x90\xea0\x01\x82\xd3\xe4\x93\x02\v\x12\t/v1/roles\x12p\n" +
	"\aGetRole\x12\x1b.bytebase.v1.GetRoleRequest\x1a\x11.bytebase.v1.Role\"5\xdaA\x04name\x8a\xea0\fbb.roles.get\x90\xea0\x01\x82\xd3\xe4\x93\x02\x14\x12\x12/v1/{name=roles/*}\x12s\n" +
	"\n" +
	"CreateRole\x12\x1e.bytebase.v1.CreateRoleRequest\x1a\x11.bytebase.v1.Role\"2\x8a\xea0\x0fbb.roles.create\x90\xea0\x01\x98\xea0\x01\x82\xd3\xe4\x93\x02\x11:\x04role\"\t/v1/roles\x12\x94\x01\n" +
	"\n" +
	"UpdateRole\x12\x1e.bytebase.v1.UpdateRoleRequest\x1a\x11.bytebase.v1.Role\"S\xdaA\x10role,update_mask\x8a\xea0\x0fbb.roles.update\x90\xea0\x01\x98\xea0\x01\x82\xd3\xe4\x93\x02\x1f:\x04role2\x17/v1/{role.name=roles/*}\x12\x82\x01\n" +
	"\n" +
	"DeleteRole\x12\x1e.bytebase.v1.DeleteRoleRequest\x1a\x16.google.protobuf.Empty\"<\xdaA\x04name\x8a\xea0\x0fbb.roles.delete\x90\xea0\x01\x98\xea0\x01\x82\xd3\xe4\x93\x02\x14*\x12/v1/{name=roles/*}B6Z4github.com/bytebase/bytebase/backend/generated-go/v1b\x06proto3"

var (
	file_v1_role_service_proto_rawDescOnce sync.Once
	file_v1_role_service_proto_rawDescData []byte
)

func file_v1_role_service_proto_rawDescGZIP() []byte {
	file_v1_role_service_proto_rawDescOnce.Do(func() {
		file_v1_role_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_role_service_proto_rawDesc), len(file_v1_role_service_proto_rawDesc)))
	})
	return file_v1_role_service_proto_rawDescData
}

var file_v1_role_service_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_v1_role_service_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_v1_role_service_proto_goTypes = []any{
	(Role_Type)(0),                // 0: bytebase.v1.Role.Type
	(*ListRolesRequest)(nil),      // 1: bytebase.v1.ListRolesRequest
	(*ListRolesResponse)(nil),     // 2: bytebase.v1.ListRolesResponse
	(*CreateRoleRequest)(nil),     // 3: bytebase.v1.CreateRoleRequest
	(*GetRoleRequest)(nil),        // 4: bytebase.v1.GetRoleRequest
	(*UpdateRoleRequest)(nil),     // 5: bytebase.v1.UpdateRoleRequest
	(*DeleteRoleRequest)(nil),     // 6: bytebase.v1.DeleteRoleRequest
	(*Role)(nil),                  // 7: bytebase.v1.Role
	(*fieldmaskpb.FieldMask)(nil), // 8: google.protobuf.FieldMask
	(*emptypb.Empty)(nil),         // 9: google.protobuf.Empty
}
var file_v1_role_service_proto_depIdxs = []int32{
	7,  // 0: bytebase.v1.ListRolesResponse.roles:type_name -> bytebase.v1.Role
	7,  // 1: bytebase.v1.CreateRoleRequest.role:type_name -> bytebase.v1.Role
	7,  // 2: bytebase.v1.UpdateRoleRequest.role:type_name -> bytebase.v1.Role
	8,  // 3: bytebase.v1.UpdateRoleRequest.update_mask:type_name -> google.protobuf.FieldMask
	0,  // 4: bytebase.v1.Role.type:type_name -> bytebase.v1.Role.Type
	1,  // 5: bytebase.v1.RoleService.ListRoles:input_type -> bytebase.v1.ListRolesRequest
	4,  // 6: bytebase.v1.RoleService.GetRole:input_type -> bytebase.v1.GetRoleRequest
	3,  // 7: bytebase.v1.RoleService.CreateRole:input_type -> bytebase.v1.CreateRoleRequest
	5,  // 8: bytebase.v1.RoleService.UpdateRole:input_type -> bytebase.v1.UpdateRoleRequest
	6,  // 9: bytebase.v1.RoleService.DeleteRole:input_type -> bytebase.v1.DeleteRoleRequest
	2,  // 10: bytebase.v1.RoleService.ListRoles:output_type -> bytebase.v1.ListRolesResponse
	7,  // 11: bytebase.v1.RoleService.GetRole:output_type -> bytebase.v1.Role
	7,  // 12: bytebase.v1.RoleService.CreateRole:output_type -> bytebase.v1.Role
	7,  // 13: bytebase.v1.RoleService.UpdateRole:output_type -> bytebase.v1.Role
	9,  // 14: bytebase.v1.RoleService.DeleteRole:output_type -> google.protobuf.Empty
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_v1_role_service_proto_init() }
func file_v1_role_service_proto_init() {
	if File_v1_role_service_proto != nil {
		return
	}
	file_v1_annotation_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_role_service_proto_rawDesc), len(file_v1_role_service_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_role_service_proto_goTypes,
		DependencyIndexes: file_v1_role_service_proto_depIdxs,
		EnumInfos:         file_v1_role_service_proto_enumTypes,
		MessageInfos:      file_v1_role_service_proto_msgTypes,
	}.Build()
	File_v1_role_service_proto = out.File
	file_v1_role_service_proto_goTypes = nil
	file_v1_role_service_proto_depIdxs = nil
}
