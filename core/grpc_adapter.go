package core

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/proto3/srpc"
	"github.com/cdle/sillyGirl/utils"
)

func (sg *SillyGirlService) AdapterRegist(stream srpc.SillyGirlService_AdapterRegistServer) error {
	var adapter *Factory
	var echos sync.Map
	for {
		req, err := stream.Recv()
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
	json.Unmarshal([]byte(req.Value), &msgs)
	adapter, err := GetAdapter(bot_id, platform)
	if err == nil {
		s := adapter.Receive(msgs)
		return &srpc.Default{Value: s.GetID()}, nil
	}
	return &srpc.Default{Value: ""}, err
}

func (sg *SillyGirlService) AdapterPush(ctx context.Context, req *srpc.AdapterRequest) (*srpc.Default, error) {
	msgs := map[string]string{}
	bot_id := req.GetBotId()
	platform := req.GetPlatform()
	json.Unmarshal([]byte(req.Value), &msgs)
	adapter, err := GetAdapter(bot_id, platform)
	if err == nil {
		result := adapter.Push(msgs)
		message_id := result["message_id"]
		error := result["error"]
		if error != "" {
			return &srpc.Default{Value: ""}, err
		}
		return &srpc.Default{Value: message_id}, err
	}
	return &srpc.Default{Value: ""}, err
}

func (sg *SillyGirlService) AdapterSender(ctx context.Context, req *srpc.AdapterRequest) (*srpc.Default, error) {
	msgs := map[string]string{}
	bot_id := req.GetBotId()
	platform := req.GetPlatform()
	json.Unmarshal([]byte(req.Value), &msgs)
	adapter, err := GetAdapter(bot_id, platform)
	if err == nil {
		s := adapter.Sender2(msgs)
		return &srpc.Default{Value: s.GetID()}, nil
	}
	return &srpc.Default{Value: ""}, err
}

func (sg *SillyGirlService) AdapterDestroy(ctx context.Context, req *srpc.AdapterRequest) (*srpc.Empty, error) {
	bot_id := req.GetBotId()
	platform := req.GetPlatform()
	adapter, err := GetAdapter(bot_id, platform)
	if err == nil {
		adapter.Destroy()
	}
	return &srpc.Empty{}, err
}
