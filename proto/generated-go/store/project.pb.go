// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: store/project.proto

package store

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type Label struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Value         string                 `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Color         string                 `protobuf:"bytes,2,opt,name=color,proto3" json:"color,omitempty"`
	Group         string                 `protobuf:"bytes,3,opt,name=group,proto3" json:"group,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Label) Reset() {
	*x = Label{}
	mi := &file_store_project_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Label) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Label) ProtoMessage() {}

func (x *Label) ProtoReflect() protoreflect.Message {
	mi := &file_store_project_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Label.ProtoReflect.Descriptor instead.
func (*Label) Descriptor() ([]byte, []int) {
	return file_store_project_proto_rawDescGZIP(), []int{0}
}

func (x *Label) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Label) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

func (x *Label) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

type Project struct {
	state       protoimpl.MessageState `protogen:"open.v1"`
	IssueLabels []*Label               `protobuf:"bytes,2,rep,name=issue_labels,json=issueLabels,proto3" json:"issue_labels,omitempty"`
	// Force issue labels to be used when creating an issue.
	ForceIssueLabels bool `protobuf:"varint,3,opt,name=force_issue_labels,json=forceIssueLabels,proto3" json:"force_issue_labels,omitempty"`
	// Allow modifying statement after issue is created.
	AllowModifyStatement bool `protobuf:"varint,4,opt,name=allow_modify_statement,json=allowModifyStatement,proto3" json:"allow_modify_statement,omitempty"`
	// Enable auto resolve issue.
	AutoResolveIssue bool `protobuf:"varint,5,opt,name=auto_resolve_issue,json=autoResolveIssue,proto3" json:"auto_resolve_issue,omitempty"`
	// Enforce issue title created by user instead of generated by Bytebase.
	EnforceIssueTitle bool `protobuf:"varint,6,opt,name=enforce_issue_title,json=enforceIssueTitle,proto3" json:"enforce_issue_title,omitempty"`
	// Whether to automatically enable backup.
	AutoEnableBackup bool `protobuf:"varint,7,opt,name=auto_enable_backup,json=autoEnableBackup,proto3" json:"auto_enable_backup,omitempty"`
	// Whether to skip backup errors and continue the data migration.
	SkipBackupErrors bool `protobuf:"varint,8,opt,name=skip_backup_errors,json=skipBackupErrors,proto3" json:"skip_backup_errors,omitempty"`
	// Whether to enable the database tenant mode for PostgreSQL.
	// If enabled, the issue will be created with the prepend "set role <db_owner>" statement.
	PostgresDatabaseTenantMode bool `protobuf:"varint,9,opt,name=postgres_database_tenant_mode,json=postgresDatabaseTenantMode,proto3" json:"postgres_database_tenant_mode,omitempty"`
	// Whether to allow the issue creator to self-approve the issue.
	AllowSelfApproval bool `protobuf:"varint,10,opt,name=allow_self_approval,json=allowSelfApproval,proto3" json:"allow_self_approval,omitempty"`
	// Execution retry policy for the task run.
	ExecutionRetryPolicy *Project_ExecutionRetryPolicy `protobuf:"bytes,11,opt,name=execution_retry_policy,json=executionRetryPolicy,proto3" json:"execution_retry_policy,omitempty"`
	// The maximum number of databases to sample during CI data validation.
	// Without specification, sampling is disabled, resulting in a full validation.
	CiSamplingSize int32 `protobuf:"varint,12,opt,name=ci_sampling_size,json=ciSamplingSize,proto3" json:"ci_sampling_size,omitempty"`
	// The maximum number of parallel tasks to run during the rollout.
	ParallelTasksPerRollout int32 `protobuf:"varint,13,opt,name=parallel_tasks_per_rollout,json=parallelTasksPerRollout,proto3" json:"parallel_tasks_per_rollout,omitempty"`
	unknownFields           protoimpl.UnknownFields
	sizeCache               protoimpl.SizeCache
}

func (x *Project) Reset() {
	*x = Project{}
	mi := &file_store_project_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Project) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Project) ProtoMessage() {}

func (x *Project) ProtoReflect() protoreflect.Message {
	mi := &file_store_project_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Project.ProtoReflect.Descriptor instead.
func (*Project) Descriptor() ([]byte, []int) {
	return file_store_project_proto_rawDescGZIP(), []int{1}
}

func (x *Project) GetIssueLabels() []*Label {
	if x != nil {
		return x.IssueLabels
	}
	return nil
}

func (x *Project) GetForceIssueLabels() bool {
	if x != nil {
		return x.ForceIssueLabels
	}
	return false
}

func (x *Project) GetAllowModifyStatement() bool {
	if x != nil {
		return x.AllowModifyStatement
	}
	return false
}

func (x *Project) GetAutoResolveIssue() bool {
	if x != nil {
		return x.AutoResolveIssue
	}
	return false
}

func (x *Project) GetEnforceIssueTitle() bool {
	if x != nil {
		return x.EnforceIssueTitle
	}
	return false
}

func (x *Project) GetAutoEnableBackup() bool {
	if x != nil {
		return x.AutoEnableBackup
	}
	return false
}

func (x *Project) GetSkipBackupErrors() bool {
	if x != nil {
		return x.SkipBackupErrors
	}
	return false
}

func (x *Project) GetPostgresDatabaseTenantMode() bool {
	if x != nil {
		return x.PostgresDatabaseTenantMode
	}
	return false
}

func (x *Project) GetAllowSelfApproval() bool {
	if x != nil {
		return x.AllowSelfApproval
	}
	return false
}

func (x *Project) GetExecutionRetryPolicy() *Project_ExecutionRetryPolicy {
	if x != nil {
		return x.ExecutionRetryPolicy
	}
	return nil
}

func (x *Project) GetCiSamplingSize() int32 {
	if x != nil {
		return x.CiSamplingSize
	}
	return 0
}

func (x *Project) GetParallelTasksPerRollout() int32 {
	if x != nil {
		return x.ParallelTasksPerRollout
	}
	return 0
}

type Project_ExecutionRetryPolicy struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The maximum number of retries for the lock timeout issue.
	MaximumRetries int32 `protobuf:"varint,1,opt,name=maximum_retries,json=maximumRetries,proto3" json:"maximum_retries,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *Project_ExecutionRetryPolicy) Reset() {
	*x = Project_ExecutionRetryPolicy{}
	mi := &file_store_project_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Project_ExecutionRetryPolicy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Project_ExecutionRetryPolicy) ProtoMessage() {}

func (x *Project_ExecutionRetryPolicy) ProtoReflect() protoreflect.Message {
	mi := &file_store_project_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Project_ExecutionRetryPolicy.ProtoReflect.Descriptor instead.
func (*Project_ExecutionRetryPolicy) Descriptor() ([]byte, []int) {
	return file_store_project_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Project_ExecutionRetryPolicy) GetMaximumRetries() int32 {
	if x != nil {
		return x.MaximumRetries
	}
	return 0
}

var File_store_project_proto protoreflect.FileDescriptor

const file_store_project_proto_rawDesc = "" +
	"\n" +
	"\x13store/project.proto\x12\x0ebytebase.store\"I\n" +
	"\x05Label\x12\x14\n" +
	"\x05value\x18\x01 \x01(\tR\x05value\x12\x14\n" +
	"\x05color\x18\x02 \x01(\tR\x05color\x12\x14\n" +
	"\x05group\x18\x03 \x01(\tR\x05group\"\xe6\x05\n" +
	"\aProject\x128\n" +
	"\fissue_labels\x18\x02 \x03(\v2\x15.bytebase.store.LabelR\vissueLabels\x12,\n" +
	"\x12force_issue_labels\x18\x03 \x01(\bR\x10forceIssueLabels\x124\n" +
	"\x16allow_modify_statement\x18\x04 \x01(\bR\x14allowModifyStatement\x12,\n" +
	"\x12auto_resolve_issue\x18\x05 \x01(\bR\x10autoResolveIssue\x12.\n" +
	"\x13enforce_issue_title\x18\x06 \x01(\bR\x11enforceIssueTitle\x12,\n" +
	"\x12auto_enable_backup\x18\a \x01(\bR\x10autoEnableBackup\x12,\n" +
	"\x12skip_backup_errors\x18\b \x01(\bR\x10skipBackupErrors\x12A\n" +
	"\x1dpostgres_database_tenant_mode\x18\t \x01(\bR\x1apostgresDatabaseTenantMode\x12.\n" +
	"\x13allow_self_approval\x18\n" +
	" \x01(\bR\x11allowSelfApproval\x12b\n" +
	"\x16execution_retry_policy\x18\v \x01(\v2,.bytebase.store.Project.ExecutionRetryPolicyR\x14executionRetryPolicy\x12(\n" +
	"\x10ci_sampling_size\x18\f \x01(\x05R\x0eciSamplingSize\x12;\n" +
	"\x1aparallel_tasks_per_rollout\x18\r \x01(\x05R\x17parallelTasksPerRollout\x1a?\n" +
	"\x14ExecutionRetryPolicy\x12'\n" +
	"\x0fmaximum_retries\x18\x01 \x01(\x05R\x0emaximumRetriesJ\x04\b\x01\x10\x02B\x14Z\x12generated-go/storeb\x06proto3"

var (
	file_store_project_proto_rawDescOnce sync.Once
	file_store_project_proto_rawDescData []byte
)

func file_store_project_proto_rawDescGZIP() []byte {
	file_store_project_proto_rawDescOnce.Do(func() {
		file_store_project_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_store_project_proto_rawDesc), len(file_store_project_proto_rawDesc)))
	})
	return file_store_project_proto_rawDescData
}

var file_store_project_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_store_project_proto_goTypes = []any{
	(*Label)(nil),                        // 0: bytebase.store.Label
	(*Project)(nil),                      // 1: bytebase.store.Project
	(*Project_ExecutionRetryPolicy)(nil), // 2: bytebase.store.Project.ExecutionRetryPolicy
}
var file_store_project_proto_depIdxs = []int32{
	0, // 0: bytebase.store.Project.issue_labels:type_name -> bytebase.store.Label
	2, // 1: bytebase.store.Project.execution_retry_policy:type_name -> bytebase.store.Project.ExecutionRetryPolicy
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_store_project_proto_init() }
func file_store_project_proto_init() {
	if File_store_project_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_store_project_proto_rawDesc), len(file_store_project_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_store_project_proto_goTypes,
		DependencyIndexes: file_store_project_proto_depIdxs,
		MessageInfos:      file_store_project_proto_msgTypes,
	}.Build()
	File_store_project_proto = out.File
	file_store_project_proto_goTypes = nil
	file_store_project_proto_depIdxs = nil
}
