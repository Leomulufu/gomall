// Code generated by Kitex v0.8.0. DO NOT EDIT.

package emailservice

import (
	"context"
	email "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/email"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	proto "google.golang.org/protobuf/proto"
)

func serviceInfo() *kitex.ServiceInfo {
	return emailServiceServiceInfo
}

var emailServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "EmailService"
	handlerType := (*email.EmailService)(nil)
	methods := map[string]kitex.MethodInfo{
		"Send": kitex.NewMethodInfo(sendHandler, newSendArgs, newSendResult, false),
	}
	extra := map[string]interface{}{
		"PackageName":     "email",
		"ServiceFilePath": ``,
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Protobuf,
		KiteXGenVersion: "v0.8.0",
		Extra:           extra,
	}
	return svcInfo
}

func sendHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(email.EmailReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(email.EmailService).Send(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *SendArgs:
		success, err := handler.(email.EmailService).Send(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*SendResult)
		realResult.Success = success
	}
	return nil
}
func newSendArgs() interface{} {
	return &SendArgs{}
}

func newSendResult() interface{} {
	return &SendResult{}
}

type SendArgs struct {
	Req *email.EmailReq
}

func (p *SendArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(email.EmailReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *SendArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *SendArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *SendArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *SendArgs) Unmarshal(in []byte) error {
	msg := new(email.EmailReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var SendArgs_Req_DEFAULT *email.EmailReq

func (p *SendArgs) GetReq() *email.EmailReq {
	if !p.IsSetReq() {
		return SendArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *SendArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *SendArgs) GetFirstArgument() interface{} {
	return p.Req
}

type SendResult struct {
	Success *email.EmailResp
}

var SendResult_Success_DEFAULT *email.EmailResp

func (p *SendResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(email.EmailResp)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *SendResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *SendResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *SendResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *SendResult) Unmarshal(in []byte) error {
	msg := new(email.EmailResp)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *SendResult) GetSuccess() *email.EmailResp {
	if !p.IsSetSuccess() {
		return SendResult_Success_DEFAULT
	}
	return p.Success
}

func (p *SendResult) SetSuccess(x interface{}) {
	p.Success = x.(*email.EmailResp)
}

func (p *SendResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *SendResult) GetResult() interface{} {
	return p.Success
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) Send(ctx context.Context, Req *email.EmailReq) (r *email.EmailResp, err error) {
	var _args SendArgs
	_args.Req = Req
	var _result SendResult
	if err = p.c.Call(ctx, "Send", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
