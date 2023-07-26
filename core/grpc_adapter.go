package core

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/proto3/srpc"
	"github.com/cdle/sillyGirl/utils"
)

func (sg *SillyGirlService) AdapterRegist(stream srpc.SillyGirlService_AdapterRegistServer) error {
	var adapter *Factory
	var echos sync.Map
	// defer func() {
	// 	fmt.Println("test defer")
	// }()
	for {
		// fmt.Println("test for")
		req, err := stream.Recv()
		// fmt.Println("test fored", err)
		if err != nil {
			return err
		}
		if adapter == nil {
			bot_id := req.GetBotId()
			platform := req.GetPlatform()
			adapter = &Factory{}
			defer adapter.Destroy()
			adapter.Init(platform, bot_id, nil)
			adapter.SetReplyHandler(func(m map[string]interface{}) string {
				m["__type__"] = "reply"
				echo := utils.GenUUID()
				ch := make(chan string)
				echos.Store(echo, ch)
				defer echos.Delete(echo)
				m["echo"] = echo
				stream.Send(&srpc.Default{
					Value: string(utils.JsonMarshal(m)),
				})
				select {
				case v := <-ch:
					return v
				case <-time.After(time.Second * 5):
				}
				return ""
			})
			adapter.action = func(m map[string]interface{}) string {
				m["__type__"] = "action"
				echo := utils.GenUUID()
				ch := make(chan string)
				echos.Store(echo, ch)
				defer echos.Delete(echo)
				m["echo"] = echo
				stream.Send(&srpc.Default{
					Value: string(utils.JsonMarshal(m)),
				})
				select {
				case v := <-ch:
					return v
				case <-time.After(time.Second * 5):
				}
				return ""
			}
			// fmt.Println("test start")
		} else {
			echo := req.GetBotId()
			message_id := req.GetPlatform()
			v, ok := echos.Load(echo)
			if ok {
				select {
				case v.(chan string) <- message_id:
				case <-time.After(time.Millisecond):
				}
			}
		}
	}
}

func (sg *SillyGirlService) AdapterReceive(ctx context.Context, req *srpc.AdapterRequest) (*srpc.Default, error) {
	msgs := map[string]interface{}{}
	bot_id := req.GetBotId()
	platform := req.GetPlatform()
	// fmt.Println("a ...any", bot_id, "=", platform, string(utils.JsonMarshal(msgs)))
	json.Unmarshal([]byte(req.Value), &msgs)
	adapter, err := GetAdapter(platform, bot_id)
	if err == nil {
		s := adapter.Receive(msgs)
		return &srpc.Default{Value: s.SetID()}, nil
	}
	return &srpc.Default{Value: ""}, err
}

func (sg *SillyGirlService) AdapterPush(ctx context.Context, req *srpc.AdapterRequest) (*srpc.Default, error) {
	msgs := map[string]string{}
	bot_id := req.GetBotId()
	platform := req.GetPlatform()
	json.Unmarshal([]byte(req.Value), &msgs)
	adapter, err := GetAdapter(platform, bot_id)
	if err == nil {
		result := adapter.Push(msgs)
		message_id := result["message_id"]
		errst := result["error"]
		if errst != "" {
			return &srpc.Default{Value: ""}, errors.New(errst)
		}
		return &srpc.Default{Value: message_id}, nil
	}
	return &srpc.Default{Value: ""}, nil
}

func (sg *SillyGirlService) AdapterSender(ctx context.Context, req *srpc.AdapterRequest) (*srpc.Default, error) {
	msgs := map[string]string{}
	bot_id := req.GetBotId()
	platform := req.GetPlatform()
	json.Unmarshal([]byte(req.Value), &msgs)
	adapter, err := GetAdapter(platform, bot_id)
	if err == nil {
		s := adapter.Sender2(msgs)
		return &srpc.Default{Value: s.SetID()}, nil
	}
	return &srpc.Default{Value: ""}, err
}

func (sg *SillyGirlService) AdapterDestroy(ctx context.Context, req *srpc.AdapterRequest) (*srpc.Empty, error) {
	bot_id := req.GetBotId()
	platform := req.GetPlatform()
	adapter, err := GetAdapter(platform, bot_id)
	if err == nil {
		adapter.Destroy()
	}
	return &srpc.Empty{}, err
}
