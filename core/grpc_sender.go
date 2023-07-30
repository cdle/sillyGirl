package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/proto3/srpc"
	"github.com/cdle/sillyGirl/utils"
)

func (sg *SillyGirlService) SenderGetUserId(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetUserID()}, nil
}

func (sg *SillyGirlService) SenderGetUserName(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetUserName()}, nil
}

func (sg *SillyGirlService) SenderGetChatId(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetChatID()}, nil
}

func (sg *SillyGirlService) SenderGetChatName(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetChatName()}, nil
}

func (sg *SillyGirlService) SenderGetMessageId(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetMessageID()}, nil
}

func (sg *SillyGirlService) SenderGetPlatform(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetImType()}, nil
}

func (sg *SillyGirlService) SenderGetBotId(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetBotID()}, nil
}

func (sg *SillyGirlService) SenderGetContent(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetContent()}, nil
}

func (sg *SillyGirlService) SenderSetContent(ctx context.Context, req *srpc.SenderContentRequest) (*srpc.Empty, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	s.SetContent(req.Content)
	return &srpc.Empty{}, nil
}

func (sg *SillyGirlService) SenderContinue(ctx context.Context, req *srpc.SenderRequest) (*srpc.Empty, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	s.Continue()
	return &srpc.Empty{}, nil
}

// todo
func (sg *SillyGirlService) SenderListen(stream srpc.SillyGirlService_SenderListenServer) error {
	runtime_id, register, err := getSenderRegisterByCtx(stream.Context())
	if err != nil {
		return err
	}
	// fmt.Println(runtime_id)
	var carry *Carry
	var echos sync.Map
	var persistent bool
	// defer fmt.Println("已关闭，", "===")
	for {
		req, err := stream.Recv()
		// fmt.Println("carry", carry, err)
		if err == io.EOF {
			break // 如果流已经关闭，则退出循环
		}
		if err != nil {
			return err
		}

		if carry != nil {
			// fmt.Println("已监听，", req)
			// fmt.Println("req.Uuid", req.Uuid)
			echo := req.GetUuid()
			value := req.GetValue()
			// fmt.Println("echo", echo, "value", value)
			v, ok := echos.Load(echo)
			if ok {
				select {
				case v.(chan string) <- value:
				case <-time.After(time.Millisecond):
				}
			}
			// if !persistent {
			// 	return nil
			// }
			continue
		}
		s, err := getRegisterSender(runtime_id, req.Uuid)

		if err != nil {
			return err
		}
		options := []interface{}{}
		carry = &Carry{
			ListenPrivate:     req.ListenPrivate,
			ListenGroup:       req.ListenGroup,
			RequireAdmin:      req.RequireAdmin,
			AllowPlatforms:    req.AllowPlatforms,
			ProhibitPlatforms: req.ProhibitPlatforms,
			AllowGroups:       req.AllowGroups,
			ProhibitGroups:    req.ProhibitGroups,
			AllowUsers:        req.AllowUsers,
			ProhibitUsers:     req.ProhibitUsers,
			UserID:            s.GetUserID(),
			ChatID:            s.GetChatID(),
			UUID:              req.PluginId,
		}
		if req.Timeout != 0 {
			options = append(options, time.Duration(req.Timeout)*time.Millisecond)
		}
		if len(req.Rules) != 0 {
			for _, rule := range req.Rules {
				_rs := formatRule(rule)
				if len(_rs) != 0 {
					carry.Function.Rules = append(carry.Function.Rules, _rs...)
				} else {
					carry.Function.Rules = append(carry.Function.Rules, rule)
				}
			}
		}
		if carry.AllowPlatforms == nil {
			plt := s.GetImType()
			if plt != "*" {
				carry.AllowPlatforms = []string{plt}
			}
		}
		options = append(options, carry)
		if req.Persistent {
			persistent = req.Persistent
			options = append(options, "persistent")
		} else {
			options = append(options, func(err error) {
				stream.Send(&srpc.SenderListenResponse{Echo: ""})
			})
		}
		level := s.GetLevel()

		go s.Await(s, func(s common.Sender) interface{} {
			s.SetLevel(level + 1)
			echo := utils.GenUUID()
			ch := make(chan string)
			echos.Store(echo, ch)
			defer echos.Delete(echo)
			stream.Send(&srpc.SenderListenResponse{Echo: echo, Uuid: register(s)})
			value := <-ch
			if !persistent {
				if strings.HasPrefix(value, "go_again_") {
					value = strings.Replace(value, "go_again_", "", 1)
					return GoAgain(value)
				} else {
					stream.Send(&srpc.SenderListenResponse{Echo: "END"})
				}
			} else {
				value = strings.Replace(value, "go_again_", "", 1)
			}
			return value
		}, options...)
	}
	return nil
}

func (sg *SillyGirlService) SenderEvent(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: string(utils.JsonMarshal(s.Event()))}, nil
}

func (sg *SillyGirlService) SenderReply(ctx context.Context, req *srpc.ReplyRequest) (*srpc.Default, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	message_id, err := s.Reply(req.Content)
	return &srpc.Default{Value: message_id}, err
}

func (sg *SillyGirlService) SenderParam(ctx context.Context, req *srpc.ReplyRequest) (*srpc.Default, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	value := ""
	i := utils.Int(req.Content)
	if fmt.Sprint(i) == req.Content {
		value = s.Get(i - 1)
	} else {
		value = s.Get(req.Content)
	}
	return &srpc.Default{Value: value}, err
}

func (sg *SillyGirlService) SenderAction(ctx context.Context, req *srpc.ReplyRequest) (*srpc.Default, error) {
	s, err := getRegisterSenderByCtx(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	var params = map[string]interface{}{}
	err = json.Unmarshal([]byte(req.Content), &params)
	if err != nil {
		return nil, err
	}
	result, err := s.Action(params)
	return &srpc.Default{Value: string(utils.JsonMarshal(result))}, err
}

func (sg *SillyGirlService) SenderDestroy(ctx context.Context, req *srpc.ReplyRequest) (*srpc.Empty, error) {
	return &srpc.Empty{}, nil
}
