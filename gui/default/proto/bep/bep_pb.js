// @generated by protoc-gen-es v2.2.5
// @generated from file bep/bep.proto (package bep, syntax proto3)
/* eslint-disable */

import { enumDesc, fileDesc, messageDesc, tsEnum } from "@bufbuild/protobuf/codegenv1";

/**
 * Describes the file bep/bep.proto.
 */
export const file_bep_bep = /*@__PURE__*/
  fileDesc("Cg1iZXAvYmVwLnByb3RvEgNiZXAidQoFSGVsbG8SEwoLZGV2aWNlX25hbWUYASABKAkSEwoLY2xpZW50X25hbWUYAiABKAkSFgoOY2xpZW50X3ZlcnNpb24YAyABKAkSFwoPbnVtX2Nvbm5lY3Rpb25zGAQgASgFEhEKCXRpbWVzdGFtcBgFIAEoAyJWCgZIZWFkZXISHgoEdHlwZRgBIAEoDjIQLmJlcC5NZXNzYWdlVHlwZRIsCgtjb21wcmVzc2lvbhgCIAEoDjIXLmJlcC5NZXNzYWdlQ29tcHJlc3Npb24iQAoNQ2x1c3RlckNvbmZpZxIcCgdmb2xkZXJzGAEgAygLMgsuYmVwLkZvbGRlchIRCglzZWNvbmRhcnkYAiABKAgitQEKBkZvbGRlchIKCgJpZBgBIAEoCRINCgVsYWJlbBgCIAEoCRIRCglyZWFkX29ubHkYAyABKAgSGgoSaWdub3JlX3Blcm1pc3Npb25zGAQgASgIEhUKDWlnbm9yZV9kZWxldGUYBSABKAgSHAoUZGlzYWJsZV90ZW1wX2luZGV4ZXMYBiABKAgSDgoGcGF1c2VkGAcgASgIEhwKB2RldmljZXMYECADKAsyCy5iZXAuRGV2aWNlIvIBCgZEZXZpY2USCgoCaWQYASABKAwSDAoEbmFtZRgCIAEoCRIRCglhZGRyZXNzZXMYAyADKAkSJQoLY29tcHJlc3Npb24YBCABKA4yEC5iZXAuQ29tcHJlc3Npb24SEQoJY2VydF9uYW1lGAUgASgJEhQKDG1heF9zZXF1ZW5jZRgGIAEoAxISCgppbnRyb2R1Y2VyGAcgASgIEhAKCGluZGV4X2lkGAggASgEEiIKGnNraXBfaW50cm9kdWN0aW9uX3JlbW92YWxzGAkgASgIEiEKGWVuY3J5cHRpb25fcGFzc3dvcmRfdG9rZW4YCiABKAwiTAoFSW5kZXgSDgoGZm9sZGVyGAEgASgJEhwKBWZpbGVzGAIgAygLMg0uYmVwLkZpbGVJbmZvEhUKDWxhc3Rfc2VxdWVuY2UYAyABKAMiaQoLSW5kZXhVcGRhdGUSDgoGZm9sZGVyGAEgASgJEhwKBWZpbGVzGAIgAygLMg0uYmVwLkZpbGVJbmZvEhUKDWxhc3Rfc2VxdWVuY2UYAyABKAMSFQoNcHJldl9zZXF1ZW5jZRgEIAEoAyKGBAoIRmlsZUluZm8SDAoEbmFtZRgBIAEoCRIMCgRzaXplGAMgASgDEhIKCm1vZGlmaWVkX3MYBSABKAMSEwoLbW9kaWZpZWRfYnkYDCABKAQSHAoHdmVyc2lvbhgJIAEoCzILLmJlcC5WZWN0b3ISEAoIc2VxdWVuY2UYCiABKAMSHgoGYmxvY2tzGBAgAygLMg4uYmVwLkJsb2NrSW5mbxIWCg5zeW1saW5rX3RhcmdldBgRIAEoDBITCgtibG9ja3NfaGFzaBgSIAEoDBIRCgllbmNyeXB0ZWQYEyABKAwSHwoEdHlwZRgCIAEoDjIRLmJlcC5GaWxlSW5mb1R5cGUSEwoLcGVybWlzc2lvbnMYBCABKA0SEwoLbW9kaWZpZWRfbnMYCyABKAUSEgoKYmxvY2tfc2l6ZRgNIAEoBRIjCghwbGF0Zm9ybRgOIAEoCzIRLmJlcC5QbGF0Zm9ybURhdGESFAoLbG9jYWxfZmxhZ3MY6AcgASgNEhUKDHZlcnNpb25faGFzaBjpByABKAwSGAoPaW5vZGVfY2hhbmdlX25zGOoHIAEoAxIgChdlbmNyeXB0aW9uX3RyYWlsZXJfc2l6ZRjrByABKAUSDwoHZGVsZXRlZBgGIAEoCBIPCgdpbnZhbGlkGAcgASgIEhYKDm5vX3Blcm1pc3Npb25zGAggASgIIj0KCUJsb2NrSW5mbxIMCgRoYXNoGAMgASgMEg4KBm9mZnNldBgBIAEoAxIMCgRzaXplGAIgASgFSgQIBBAFIigKBlZlY3RvchIeCghjb3VudGVycxgBIAMoCzIMLmJlcC5Db3VudGVyIiQKB0NvdW50ZXISCgoCaWQYASABKAQSDQoFdmFsdWUYAiABKAQizgEKDFBsYXRmb3JtRGF0YRIbCgR1bml4GAEgASgLMg0uYmVwLlVuaXhEYXRhEiEKB3dpbmRvd3MYAiABKAsyEC5iZXAuV2luZG93c0RhdGESHQoFbGludXgYAyABKAsyDi5iZXAuWGF0dHJEYXRhEh4KBmRhcndpbhgEIAEoCzIOLmJlcC5YYXR0ckRhdGESHwoHZnJlZWJzZBgFIAEoCzIOLmJlcC5YYXR0ckRhdGESHgoGbmV0YnNkGAYgASgLMg4uYmVwLlhhdHRyRGF0YSJMCghVbml4RGF0YRISCgpvd25lcl9uYW1lGAEgASgJEhIKCmdyb3VwX25hbWUYAiABKAkSCwoDdWlkGAMgASgFEgsKA2dpZBgEIAEoBSI5CgtXaW5kb3dzRGF0YRISCgpvd25lcl9uYW1lGAEgASgJEhYKDm93bmVyX2lzX2dyb3VwGAIgASgIIicKCVhhdHRyRGF0YRIaCgZ4YXR0cnMYASADKAsyCi5iZXAuWGF0dHIiJAoFWGF0dHISDAoEbmFtZRgBIAEoCRINCgV2YWx1ZRgCIAEoDCKPAQoHUmVxdWVzdBIKCgJpZBgBIAEoBRIOCgZmb2xkZXIYAiABKAkSDAoEbmFtZRgDIAEoCRIOCgZvZmZzZXQYBCABKAMSDAoEc2l6ZRgFIAEoBRIMCgRoYXNoGAYgASgMEhYKDmZyb21fdGVtcG9yYXJ5GAcgASgIEhAKCGJsb2NrX25vGAkgASgFSgQICBAJIkIKCFJlc3BvbnNlEgoKAmlkGAEgASgFEgwKBGRhdGEYAiABKAwSHAoEY29kZRgDIAEoDjIOLmJlcC5FcnJvckNvZGUiVAoQRG93bmxvYWRQcm9ncmVzcxIOCgZmb2xkZXIYASABKAkSMAoHdXBkYXRlcxgCIAMoCzIfLmJlcC5GaWxlRG93bmxvYWRQcm9ncmVzc1VwZGF0ZSKxAQoaRmlsZURvd25sb2FkUHJvZ3Jlc3NVcGRhdGUSOAoLdXBkYXRlX3R5cGUYASABKA4yIy5iZXAuRmlsZURvd25sb2FkUHJvZ3Jlc3NVcGRhdGVUeXBlEgwKBG5hbWUYAiABKAkSHAoHdmVyc2lvbhgDIAEoCzILLmJlcC5WZWN0b3ISGQoNYmxvY2tfaW5kZXhlcxgEIAMoBUICEAASEgoKYmxvY2tfc2l6ZRgFIAEoBSIGCgRQaW5nIhcKBUNsb3NlEg4KBnJlYXNvbhgBIAEoCSrtAQoLTWVzc2FnZVR5cGUSHwobTUVTU0FHRV9UWVBFX0NMVVNURVJfQ09ORklHEAASFgoSTUVTU0FHRV9UWVBFX0lOREVYEAESHQoZTUVTU0FHRV9UWVBFX0lOREVYX1VQREFURRACEhgKFE1FU1NBR0VfVFlQRV9SRVFVRVNUEAMSGQoVTUVTU0FHRV9UWVBFX1JFU1BPTlNFEAQSIgoeTUVTU0FHRV9UWVBFX0RPV05MT0FEX1BST0dSRVNTEAUSFQoRTUVTU0FHRV9UWVBFX1BJTkcQBhIWChJNRVNTQUdFX1RZUEVfQ0xPU0UQBypPChJNZXNzYWdlQ29tcHJlc3Npb24SHAoYTUVTU0FHRV9DT01QUkVTU0lPTl9OT05FEAASGwoXTUVTU0FHRV9DT01QUkVTU0lPTl9MWjQQASpWCgtDb21wcmVzc2lvbhIYChRDT01QUkVTU0lPTl9NRVRBREFUQRAAEhUKEUNPTVBSRVNTSU9OX05FVkVSEAESFgoSQ09NUFJFU1NJT05fQUxXQVlTEAIqsAEKDEZpbGVJbmZvVHlwZRIXChNGSUxFX0lORk9fVFlQRV9GSUxFEAASHAoYRklMRV9JTkZPX1RZUEVfRElSRUNUT1JZEAESIwobRklMRV9JTkZPX1RZUEVfU1lNTElOS19GSUxFEAIaAggBEigKIEZJTEVfSU5GT19UWVBFX1NZTUxJTktfRElSRUNUT1JZEAMaAggBEhoKFkZJTEVfSU5GT19UWVBFX1NZTUxJTksQBCp2CglFcnJvckNvZGUSFwoTRVJST1JfQ09ERV9OT19FUlJPUhAAEhYKEkVSUk9SX0NPREVfR0VORVJJQxABEhsKF0VSUk9SX0NPREVfTk9fU1VDSF9GSUxFEAISGwoXRVJST1JfQ09ERV9JTlZBTElEX0ZJTEUQAyp+Ch5GaWxlRG93bmxvYWRQcm9ncmVzc1VwZGF0ZVR5cGUSLQopRklMRV9ET1dOTE9BRF9QUk9HUkVTU19VUERBVEVfVFlQRV9BUFBFTkQQABItCilGSUxFX0RPV05MT0FEX1BST0dSRVNTX1VQREFURV9UWVBFX0ZPUkdFVBABQnAKB2NvbS5iZXBCCEJlcFByb3RvUAFaL2dpdGh1Yi5jb20vc3luY3RoaW5nL3N5bmN0aGluZy9pbnRlcm5hbC9nZW4vYmVwogIDQlhYqgIDQmVwygIDQmVw4gIPQmVwXEdQQk1ldGFkYXRh6gIDQmVwYgZwcm90bzM");

/**
 * Describes the message bep.Hello.
 * Use `create(HelloSchema)` to create a new message.
 */
export const HelloSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 0);

/**
 * Describes the message bep.Header.
 * Use `create(HeaderSchema)` to create a new message.
 */
export const HeaderSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 1);

/**
 * Describes the message bep.ClusterConfig.
 * Use `create(ClusterConfigSchema)` to create a new message.
 */
export const ClusterConfigSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 2);

/**
 * Describes the message bep.Folder.
 * Use `create(FolderSchema)` to create a new message.
 */
export const FolderSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 3);

/**
 * Describes the message bep.Device.
 * Use `create(DeviceSchema)` to create a new message.
 */
export const DeviceSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 4);

/**
 * Describes the message bep.Index.
 * Use `create(IndexSchema)` to create a new message.
 */
export const IndexSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 5);

/**
 * Describes the message bep.IndexUpdate.
 * Use `create(IndexUpdateSchema)` to create a new message.
 */
export const IndexUpdateSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 6);

/**
 * Describes the message bep.FileInfo.
 * Use `create(FileInfoSchema)` to create a new message.
 */
export const FileInfoSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 7);

/**
 * Describes the message bep.BlockInfo.
 * Use `create(BlockInfoSchema)` to create a new message.
 */
export const BlockInfoSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 8);

/**
 * Describes the message bep.Vector.
 * Use `create(VectorSchema)` to create a new message.
 */
export const VectorSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 9);

/**
 * Describes the message bep.Counter.
 * Use `create(CounterSchema)` to create a new message.
 */
export const CounterSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 10);

/**
 * Describes the message bep.PlatformData.
 * Use `create(PlatformDataSchema)` to create a new message.
 */
export const PlatformDataSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 11);

/**
 * Describes the message bep.UnixData.
 * Use `create(UnixDataSchema)` to create a new message.
 */
export const UnixDataSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 12);

/**
 * Describes the message bep.WindowsData.
 * Use `create(WindowsDataSchema)` to create a new message.
 */
export const WindowsDataSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 13);

/**
 * Describes the message bep.XattrData.
 * Use `create(XattrDataSchema)` to create a new message.
 */
export const XattrDataSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 14);

/**
 * Describes the message bep.Xattr.
 * Use `create(XattrSchema)` to create a new message.
 */
export const XattrSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 15);

/**
 * Describes the message bep.Request.
 * Use `create(RequestSchema)` to create a new message.
 */
export const RequestSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 16);

/**
 * Describes the message bep.Response.
 * Use `create(ResponseSchema)` to create a new message.
 */
export const ResponseSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 17);

/**
 * Describes the message bep.DownloadProgress.
 * Use `create(DownloadProgressSchema)` to create a new message.
 */
export const DownloadProgressSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 18);

/**
 * Describes the message bep.FileDownloadProgressUpdate.
 * Use `create(FileDownloadProgressUpdateSchema)` to create a new message.
 */
export const FileDownloadProgressUpdateSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 19);

/**
 * Describes the message bep.Ping.
 * Use `create(PingSchema)` to create a new message.
 */
export const PingSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 20);

/**
 * Describes the message bep.Close.
 * Use `create(CloseSchema)` to create a new message.
 */
export const CloseSchema = /*@__PURE__*/
  messageDesc(file_bep_bep, 21);

/**
 * Describes the enum bep.MessageType.
 */
export const MessageTypeSchema = /*@__PURE__*/
  enumDesc(file_bep_bep, 0);

/**
 * @generated from enum bep.MessageType
 */
export const MessageType = /*@__PURE__*/
  tsEnum(MessageTypeSchema);

/**
 * Describes the enum bep.MessageCompression.
 */
export const MessageCompressionSchema = /*@__PURE__*/
  enumDesc(file_bep_bep, 1);

/**
 * @generated from enum bep.MessageCompression
 */
export const MessageCompression = /*@__PURE__*/
  tsEnum(MessageCompressionSchema);

/**
 * Describes the enum bep.Compression.
 */
export const CompressionSchema = /*@__PURE__*/
  enumDesc(file_bep_bep, 2);

/**
 * @generated from enum bep.Compression
 */
export const Compression = /*@__PURE__*/
  tsEnum(CompressionSchema);

/**
 * Describes the enum bep.FileInfoType.
 */
export const FileInfoTypeSchema = /*@__PURE__*/
  enumDesc(file_bep_bep, 3);

/**
 * @generated from enum bep.FileInfoType
 */
export const FileInfoType = /*@__PURE__*/
  tsEnum(FileInfoTypeSchema);

/**
 * Describes the enum bep.ErrorCode.
 */
export const ErrorCodeSchema = /*@__PURE__*/
  enumDesc(file_bep_bep, 4);

/**
 * @generated from enum bep.ErrorCode
 */
export const ErrorCode = /*@__PURE__*/
  tsEnum(ErrorCodeSchema);

/**
 * Describes the enum bep.FileDownloadProgressUpdateType.
 */
export const FileDownloadProgressUpdateTypeSchema = /*@__PURE__*/
  enumDesc(file_bep_bep, 5);

/**
 * @generated from enum bep.FileDownloadProgressUpdateType
 */
export const FileDownloadProgressUpdateType = /*@__PURE__*/
  tsEnum(FileDownloadProgressUpdateTypeSchema);

