// Code generated by ICEBERG protoc-gen-go. DO NOT EDIT EXCEPET SERVER VERSION.
// source: cron.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	cron.proto
	cron_request.proto
	cron_response.proto

It has these top-level messages:
	RegisterReq
	DelTaskReq
	PauseTaskReq
	RestoreTaskReq
	GetTaskReq
	MotifyTaskReq
	RegisterResp
	DelTaskResp
	PauseTaskResp
	RestoreTaskResp
	GetTaskResp
	MotifyTaskResp
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	"context"
	"iceberg/frame"
	"iceberg/frame/config"
	"iceberg/frame/protocol"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.

// Client API for CronService service
// iceberg server version,relation to server uri.
var cronServiceVersion = frame.SrvVersionName[frame.SV1]

// 定时任务注册
func RegisterTask(ctx frame.Context, in *RegisterReq, opts ...frame.CallOption) (*RegisterResp, error) {
	task, err := frame.ReadyTask(ctx, "registertask", "cronservice", cronServiceVersion, in, opts...)
	if err != nil {
		return nil, err
	}
	back, err := frame.DeliverTo(task)
	if err != nil {
		return nil, err
	}

	var out RegisterResp
	if err := protocol.Unpack(back.GetFormat(), back.GetBody(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// 删除指定定时任务
func DelTask(ctx frame.Context, in *DelTaskReq, opts ...frame.CallOption) (*DelTaskResp, error) {
	task, err := frame.ReadyTask(ctx, "deltask", "cronservice", cronServiceVersion, in, opts...)
	if err != nil {
		return nil, err
	}
	back, err := frame.DeliverTo(task)
	if err != nil {
		return nil, err
	}

	var out DelTaskResp
	if err := protocol.Unpack(back.GetFormat(), back.GetBody(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// 定时任务暂停
func PauseTask(ctx frame.Context, in *PauseTaskReq, opts ...frame.CallOption) (*PauseTaskResp, error) {
	task, err := frame.ReadyTask(ctx, "pausetask", "cronservice", cronServiceVersion, in, opts...)
	if err != nil {
		return nil, err
	}
	back, err := frame.DeliverTo(task)
	if err != nil {
		return nil, err
	}

	var out PauseTaskResp
	if err := protocol.Unpack(back.GetFormat(), back.GetBody(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// 定时任务恢复
func RestoreTask(ctx frame.Context, in *RestoreTaskReq, opts ...frame.CallOption) (*RestoreTaskResp, error) {
	task, err := frame.ReadyTask(ctx, "restoretask", "cronservice", cronServiceVersion, in, opts...)
	if err != nil {
		return nil, err
	}
	back, err := frame.DeliverTo(task)
	if err != nil {
		return nil, err
	}

	var out RestoreTaskResp
	if err := protocol.Unpack(back.GetFormat(), back.GetBody(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// 查询指定定时任务 或者 所有定时任务
// 带task_name，返回指定定时任务查询结果
// 否则返回所有定时任务结果
func GetTask(ctx frame.Context, in *GetTaskReq, opts ...frame.CallOption) (*GetTaskResp, error) {
	task, err := frame.ReadyTask(ctx, "gettask", "cronservice", cronServiceVersion, in, opts...)
	if err != nil {
		return nil, err
	}
	back, err := frame.DeliverTo(task)
	if err != nil {
		return nil, err
	}

	var out GetTaskResp
	if err := protocol.Unpack(back.GetFormat(), back.GetBody(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// 修改定时任务
func MotifyTask(ctx frame.Context, in *MotifyTaskReq, opts ...frame.CallOption) (*MotifyTaskResp, error) {
	task, err := frame.ReadyTask(ctx, "motifytask", "cronservice", cronServiceVersion, in, opts...)
	if err != nil {
		return nil, err
	}
	back, err := frame.DeliverTo(task)
	if err != nil {
		return nil, err
	}

	var out MotifyTaskResp
	if err := protocol.Unpack(back.GetFormat(), back.GetBody(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CronServiceServer Server API for Hello service
type CronServiceServer interface {
	RegisterTask(c frame.Context) error

	DelTask(c frame.Context) error

	PauseTask(c frame.Context) error

	RestoreTask(c frame.Context) error

	GetTask(c frame.Context) error

	MotifyTask(c frame.Context) error
}

// RegisterCronServiceServer register CronServiceServer with etcd info
func RegisterCronServiceServer(srv CronServiceServer, cfg *config.BaseCfg) {
	frame.RegisterAndServe(&cronServiceServerDesc, srv, cfg)
}

// cronService server RegisterTask handler
func cronServiceRegisterTaskHandler(srv interface{}, ctx frame.Context) error {
	return srv.(CronServiceServer).RegisterTask(ctx)
}

// cronService server DelTask handler
func cronServiceDelTaskHandler(srv interface{}, ctx frame.Context) error {
	return srv.(CronServiceServer).DelTask(ctx)
}

// cronService server PauseTask handler
func cronServicePauseTaskHandler(srv interface{}, ctx frame.Context) error {
	return srv.(CronServiceServer).PauseTask(ctx)
}

// cronService server RestoreTask handler
func cronServiceRestoreTaskHandler(srv interface{}, ctx frame.Context) error {
	return srv.(CronServiceServer).RestoreTask(ctx)
}

// cronService server GetTask handler
func cronServiceGetTaskHandler(srv interface{}, ctx frame.Context) error {
	return srv.(CronServiceServer).GetTask(ctx)
}

// cronService server MotifyTask handler
func cronServiceMotifyTaskHandler(srv interface{}, ctx frame.Context) error {
	return srv.(CronServiceServer).MotifyTask(ctx)
}

// cronService server describe
var cronServiceServerDesc = frame.ServiceDesc{
	Version:     cronServiceVersion,
	ServiceName: "CronService",
	HandlerType: (*CronServiceServer)(nil),
	Methods: []frame.MethodDesc{
		{
			A:          frame.Internal,
			MethodName: "registertask",
			Handler:    cronServiceRegisterTaskHandler,
		},
		{
			A:          frame.Internal,
			MethodName: "deltask",
			Handler:    cronServiceDelTaskHandler,
		},
		{
			A:          frame.Internal,
			MethodName: "pausetask",
			Handler:    cronServicePauseTaskHandler,
		},
		{
			A:          frame.Internal,
			MethodName: "restoretask",
			Handler:    cronServiceRestoreTaskHandler,
		},
		{
			A:          frame.Internal,
			MethodName: "gettask",
			Handler:    cronServiceGetTaskHandler,
		},
		{
			A:          frame.Internal,
			MethodName: "motifytask",
			Handler:    cronServiceMotifyTaskHandler,
		},
	},
	ServiceURI: []string{
		"/services/" + cronServiceVersion + "/cronservice",
	},
	Metadata: "pb.CronService",
}

func init() { proto.RegisterFile("cron.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 219 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0xcd, 0x6a, 0xc2, 0x40,
	0x14, 0x46, 0xdb, 0x2c, 0x5a, 0x7a, 0x53, 0x9a, 0x76, 0xb2, 0xcb, 0xb2, 0x0f, 0x30, 0x94, 0xb6,
	0x22, 0xb8, 0x55, 0x70, 0x25, 0x48, 0x74, 0x2f, 0x49, 0xb8, 0x4a, 0x50, 0x32, 0x93, 0xb9, 0x13,
	0xc1, 0xa7, 0xf4, 0x95, 0x64, 0x7e, 0x9c, 0x90, 0xb8, 0x3c, 0x87, 0x7b, 0x86, 0x8f, 0x01, 0xa8,
	0x94, 0x68, 0xb8, 0x54, 0x42, 0x0b, 0x16, 0xc9, 0x32, 0x63, 0x86, 0x77, 0x0a, 0xdb, 0x0e, 0x49,
	0x3b, 0x9f, 0xa5, 0xde, 0x91, 0x14, 0x0d, 0xa1, 0x93, 0xbf, 0xd7, 0x08, 0xe2, 0xb9, 0x12, 0xcd,
	0x06, 0xd5, 0xb9, 0xae, 0x90, 0x4d, 0xe0, 0x3d, 0xc7, 0x43, 0x4d, 0x1a, 0xd5, 0xb6, 0xa0, 0x23,
	0x4b, 0xb8, 0x2c, 0xf9, 0xdd, 0xe4, 0xd8, 0x66, 0x9f, 0x43, 0x41, 0xf2, 0xfb, 0xe9, 0xe7, 0x99,
	0x71, 0x78, 0x5d, 0xe0, 0xc9, 0x16, 0x1f, 0xe6, 0xc0, 0x83, 0x09, 0x92, 0x01, 0xfb, 0xfb, 0x7f,
	0x78, 0x5b, 0x17, 0x1d, 0xa1, 0x2d, 0xec, 0x93, 0x01, 0x4d, 0xf3, 0x35, 0x32, 0xbe, 0x9a, 0x41,
	0x9c, 0x23, 0x69, 0xa1, 0x5c, 0xc7, 0xdc, 0x94, 0x20, 0x4c, 0x99, 0x3e, 0xb8, 0x7e, 0xe1, 0x12,
	0x75, 0xbf, 0xd0, 0x43, 0x58, 0x18, 0xd8, 0xdf, 0x4f, 0x01, 0x56, 0x42, 0xd7, 0xfb, 0x8b, 0x4d,
	0xec, 0xa0, 0x9e, 0x4d, 0xc5, 0xc6, 0xca, 0x85, 0xe5, 0x8b, 0xfd, 0xd8, 0xbf, 0x5b, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x2c, 0x33, 0xf3, 0x78, 0x93, 0x01, 0x00, 0x00,
}
