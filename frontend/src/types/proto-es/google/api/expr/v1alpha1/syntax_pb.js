// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// @generated by protoc-gen-es v2.5.2
// @generated from file google/api/expr/v1alpha1/syntax.proto (package google.api.expr.v1alpha1, syntax proto3)
/* eslint-disable */

import { fileDesc, messageDesc } from "@bufbuild/protobuf/codegenv2";
import { file_google_protobuf_duration, file_google_protobuf_struct, file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";

/**
 * Describes the file google/api/expr/v1alpha1/syntax.proto.
 */
export const file_google_api_expr_v1alpha1_syntax = /*@__PURE__*/
  fileDesc("CiVnb29nbGUvYXBpL2V4cHIvdjFhbHBoYTEvc3ludGF4LnByb3RvEhhnb29nbGUuYXBpLmV4cHIudjFhbHBoYTEiigsKBEV4cHISCgoCaWQYAiABKAMSOAoKY29uc3RfZXhwchgDIAEoCzIiLmdvb2dsZS5hcGkuZXhwci52MWFscGhhMS5Db25zdGFudEgAEjoKCmlkZW50X2V4cHIYBCABKAsyJC5nb29nbGUuYXBpLmV4cHIudjFhbHBoYTEuRXhwci5JZGVudEgAEjwKC3NlbGVjdF9leHByGAUgASgLMiUuZ29vZ2xlLmFwaS5leHByLnYxYWxwaGExLkV4cHIuU2VsZWN0SAASOAoJY2FsbF9leHByGAYgASgLMiMuZ29vZ2xlLmFwaS5leHByLnYxYWxwaGExLkV4cHIuQ2FsbEgAEj4KCWxpc3RfZXhwchgHIAEoCzIpLmdvb2dsZS5hcGkuZXhwci52MWFscGhhMS5FeHByLkNyZWF0ZUxpc3RIABJCCgtzdHJ1Y3RfZXhwchgIIAEoCzIrLmdvb2dsZS5hcGkuZXhwci52MWFscGhhMS5FeHByLkNyZWF0ZVN0cnVjdEgAEkoKEmNvbXByZWhlbnNpb25fZXhwchgJIAEoCzIsLmdvb2dsZS5hcGkuZXhwci52MWFscGhhMS5FeHByLkNvbXByZWhlbnNpb25IABoVCgVJZGVudBIMCgRuYW1lGAEgASgJGlsKBlNlbGVjdBIvCgdvcGVyYW5kGAEgASgLMh4uZ29vZ2xlLmFwaS5leHByLnYxYWxwaGExLkV4cHISDQoFZmllbGQYAiABKAkSEQoJdGVzdF9vbmx5GAMgASgIGnYKBENhbGwSLgoGdGFyZ2V0GAEgASgLMh4uZ29vZ2xlLmFwaS5leHByLnYxYWxwaGExLkV4cHISEAoIZnVuY3Rpb24YAiABKAkSLAoEYXJncxgDIAMoCzIeLmdvb2dsZS5hcGkuZXhwci52MWFscGhhMS5FeHByGlgKCkNyZWF0ZUxpc3QSMAoIZWxlbWVudHMYASADKAsyHi5nb29nbGUuYXBpLmV4cHIudjFhbHBoYTEuRXhwchIYChBvcHRpb25hbF9pbmRpY2VzGAIgAygFGpkCCgxDcmVhdGVTdHJ1Y3QSFAoMbWVzc2FnZV9uYW1lGAEgASgJEkIKB2VudHJpZXMYAiADKAsyMS5nb29nbGUuYXBpLmV4cHIudjFhbHBoYTEuRXhwci5DcmVhdGVTdHJ1Y3QuRW50cnkargEKBUVudHJ5EgoKAmlkGAEgASgDEhMKCWZpZWxkX2tleRgCIAEoCUgAEjEKB21hcF9rZXkYAyABKAsyHi5nb29nbGUuYXBpLmV4cHIudjFhbHBoYTEuRXhwckgAEi0KBXZhbHVlGAQgASgLMh4uZ29vZ2xlLmFwaS5leHByLnYxYWxwaGExLkV4cHISFgoOb3B0aW9uYWxfZW50cnkYBSABKAhCCgoIa2V5X2tpbmQayAIKDUNvbXByZWhlbnNpb24SEAoIaXRlcl92YXIYASABKAkSEQoJaXRlcl92YXIyGAggASgJEjIKCml0ZXJfcmFuZ2UYAiABKAsyHi5nb29nbGUuYXBpLmV4cHIudjFhbHBoYTEuRXhwchIQCghhY2N1X3ZhchgDIAEoCRIxCglhY2N1X2luaXQYBCABKAsyHi5nb29nbGUuYXBpLmV4cHIudjFhbHBoYTEuRXhwchI2Cg5sb29wX2NvbmRpdGlvbhgFIAEoCzIeLmdvb2dsZS5hcGkuZXhwci52MWFscGhhMS5FeHByEjEKCWxvb3Bfc3RlcBgGIAEoCzIeLmdvb2dsZS5hcGkuZXhwci52MWFscGhhMS5FeHByEi4KBnJlc3VsdBgHIAEoCzIeLmdvb2dsZS5hcGkuZXhwci52MWFscGhhMS5FeHByQgsKCWV4cHJfa2luZCLNAgoIQ29uc3RhbnQSMAoKbnVsbF92YWx1ZRgBIAEoDjIaLmdvb2dsZS5wcm90b2J1Zi5OdWxsVmFsdWVIABIUCgpib29sX3ZhbHVlGAIgASgISAASFQoLaW50NjRfdmFsdWUYAyABKANIABIWCgx1aW50NjRfdmFsdWUYBCABKARIABIWCgxkb3VibGVfdmFsdWUYBSABKAFIABIWCgxzdHJpbmdfdmFsdWUYBiABKAlIABIVCgtieXRlc192YWx1ZRgHIAEoDEgAEjcKDmR1cmF0aW9uX3ZhbHVlGAggASgLMhkuZ29vZ2xlLnByb3RvYnVmLkR1cmF0aW9uQgIYAUgAEjkKD3RpbWVzdGFtcF92YWx1ZRgJIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBCAhgBSABCDwoNY29uc3RhbnRfa2luZEJuChxjb20uZ29vZ2xlLmFwaS5leHByLnYxYWxwaGExQgtTeW50YXhQcm90b1ABWjxnb29nbGUuZ29sYW5nLm9yZy9nZW5wcm90by9nb29nbGVhcGlzL2FwaS9leHByL3YxYWxwaGExO2V4cHL4AQFiBnByb3RvMw", [file_google_protobuf_duration, file_google_protobuf_struct, file_google_protobuf_timestamp]);

/**
 * Describes the message google.api.expr.v1alpha1.Expr.
 * Use `create(ExprSchema)` to create a new message.
 */
export const ExprSchema = /*@__PURE__*/
  messageDesc(file_google_api_expr_v1alpha1_syntax, 0);

/**
 * Describes the message google.api.expr.v1alpha1.Expr.Ident.
 * Use `create(Expr_IdentSchema)` to create a new message.
 */
export const Expr_IdentSchema = /*@__PURE__*/
  messageDesc(file_google_api_expr_v1alpha1_syntax, 0, 0);

/**
 * Describes the message google.api.expr.v1alpha1.Expr.Select.
 * Use `create(Expr_SelectSchema)` to create a new message.
 */
export const Expr_SelectSchema = /*@__PURE__*/
  messageDesc(file_google_api_expr_v1alpha1_syntax, 0, 1);

/**
 * Describes the message google.api.expr.v1alpha1.Expr.Call.
 * Use `create(Expr_CallSchema)` to create a new message.
 */
export const Expr_CallSchema = /*@__PURE__*/
  messageDesc(file_google_api_expr_v1alpha1_syntax, 0, 2);

/**
 * Describes the message google.api.expr.v1alpha1.Expr.CreateList.
 * Use `create(Expr_CreateListSchema)` to create a new message.
 */
export const Expr_CreateListSchema = /*@__PURE__*/
  messageDesc(file_google_api_expr_v1alpha1_syntax, 0, 3);

/**
 * Describes the message google.api.expr.v1alpha1.Expr.CreateStruct.
 * Use `create(Expr_CreateStructSchema)` to create a new message.
 */
export const Expr_CreateStructSchema = /*@__PURE__*/
  messageDesc(file_google_api_expr_v1alpha1_syntax, 0, 4);

/**
 * Describes the message google.api.expr.v1alpha1.Expr.CreateStruct.Entry.
 * Use `create(Expr_CreateStruct_EntrySchema)` to create a new message.
 */
export const Expr_CreateStruct_EntrySchema = /*@__PURE__*/
  messageDesc(file_google_api_expr_v1alpha1_syntax, 0, 4, 0);

/**
 * Describes the message google.api.expr.v1alpha1.Expr.Comprehension.
 * Use `create(Expr_ComprehensionSchema)` to create a new message.
 */
export const Expr_ComprehensionSchema = /*@__PURE__*/
  messageDesc(file_google_api_expr_v1alpha1_syntax, 0, 5);

/**
 * Describes the message google.api.expr.v1alpha1.Constant.
 * Use `create(ConstantSchema)` to create a new message.
 */
export const ConstantSchema = /*@__PURE__*/
  messageDesc(file_google_api_expr_v1alpha1_syntax, 1);

