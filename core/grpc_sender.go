package core

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/proto3/srpc"
	"github.com/cdle/sillyGirl/utils"
)

var senders sync.Map

func init() {
	//垃圾回收
	go func() {
		for {
			time.Sleep(time.Minute)
			senders.Range(func(key, value any) bool {
				s := value.(common.Sender)
				if s.GetTime().Add(time.Minute * 20).Before(time.Now()) {
					senders.Delete(s.GetID())
				}
				return true
			})
		}
	}()
}

func GetSender(uuid string) (common.Sender, error) {
	v, ok := senders.Load(uuid)
	if !ok {
		return nil, errors.New("not found sender")
	}
	return v.(common.Sender), nil
}

func (sg *SillyGirlService) SenderGetUserId(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := GetSender(req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetChatID()}, nil
}

func (sg *SillyGirlService) SenderGetUserName(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := GetSender(req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetUserName()}, nil
}

func (sg *SillyGirlService) SenderGetChatId(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := GetSender(req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetChatID()}, nil
}

func (sg *SillyGirlService) SenderGetChatName(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := GetSender(req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetChatName()}, nil
}

func (sg *SillyGirlService) SenderGetMessageId(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := GetSender(req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetMessageID()}, nil
}

func (sg *SillyGirlService) SenderGetPlatform(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := GetSender(req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetImType()}, nil
}

func (sg *SillyGirlService) SenderGetBotId(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := GetSender(req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetBotID()}, nil
}

func (sg *SillyGirlService) SenderGetContent(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := GetSender(req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: s.GetContent()}, nil
}

func (sg *SillyGirlService) SenderSetContent(ctx context.Context, req *srpc.SenderContentRequest) (*srpc.Empty, error) {
	s, err := GetSender(req.Uuid)
	if err != nil {
		return nil, err
	}
	s.SetContent(req.Content)
	return &srpc.Empty{}, nil
}

func (sg *SillyGirlService) SenderContinue(ctx context.Context, req *srpc.SenderRequest) (*srpc.Empty, error) {
	s, err := GetSender(req.Uuid)
	if err != nil {
		return nil, err
	}
	s.Continue()
	return &srpc.Empty{}, nil
}

// todo
func (sg *SillyGirlService) SenderListen(ctx context.Context, req *srpc.SenderListenRequest) (*srpc.Default, error) {
	s, err := GetSender(req.Uuid)
	if err != nil {
		return nil, err
	}

	options := []interface{}{}
	var carry = &Carry{
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
	options = append(options, carry)
	id := ""
	s.Await(s, func(s common.Sender) interface{} {
		id = s.GetID()
		return nil
	}, options...)
	return &srpc.Default{Value: id}, nil
}

func (sg *SillyGirlService) SenderEvent(ctx context.Context, req *srpc.SenderRequest) (*srpc.Default, error) {
	s, err := GetSender(req.Uuid)
	if err != nil {
		return nil, err
	}
	return &srpc.Default{Value: string(utils.JsonMarshal(s.Event()))}, nil
}

func (sg *SillyGirlService) SenderReply(ctx context.Context, req *srpc.ReplyRequest) (*srpc.Default, error) {
	s, err := GetSender(req.Uuid)
	if err != nil {
		return nil, err
	}
	message_id, err := s.Reply(req.Content)
	return &srpc.Default{Value: message_id}, err
}

func (sg *SillyGirlService) SenderAction(ctx context.Context, req *srpc.ReplyRequest) (*srpc.Default, error) {
	s, err := GetSender(req.Uuid)
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
