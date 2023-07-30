package core

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/utils"
	"google.golang.org/grpc/metadata"
)

var senderRegisters sync.Map

func getRegisterSenderByCtx(ctx context.Context, uuid string) (common.Sender, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	runtimeID := md.Get("RUNTIME_ID")
	if runtimeID == nil {
		return nil, fmt.Errorf("missing runtime_id")
	}
	runtime_id := runtimeID[0]
	return getRegisterSender(runtime_id, uuid)
}

func getRegisterSender(runtime_id, uuid string) (common.Sender, error) {
	if uuid == "" {
		return &CustomSender{
			F: &Factory{
				botid:  "default",
				botplt: "*",
			},
		}, nil
	}
	var sm *sync.Map
	v, ok := senderRegisters.Load(runtime_id)
	if ok {
		sm = v.(*sync.Map)
		v, ok := sm.Load(uuid)
		if ok {
			return v.(common.Sender), nil
		}
	}
	return nil, errors.New("NOT FOUND SENDER")
}

func createSenderRegister(runtime_id string) func(common.Sender) string {
	var sm = new(sync.Map)
	senderRegisters.Store(runtime_id, sm)
	return func(s common.Sender) string {
		uuid := utils.GenUUID()
		sm.Store(uuid, s)
		return uuid
	}
}

func getSenderRegister(runtime_id string) (func(common.Sender) string, error) {
	var sm *sync.Map
	v, ok := senderRegisters.Load(runtime_id)
	if !ok {
		return nil, errors.New("INVALID RUNTIME")

	}
	sm = v.(*sync.Map)
	return func(s common.Sender) string {
		uuid := utils.GenUUID()
		sm.Store(uuid, s)
		return uuid
	}, nil
}

func deleteSenderRegister(runtime_id string) {
	v, ok := senderRegisters.Load(runtime_id)
	if ok {
		sm := v.(*sync.Map)
		sm.Range(func(key, value any) bool {
			fmt.Println(key)
			sm.Delete(key)
			return true
		})
	}

	senderRegisters.Delete(runtime_id)
}

func getSenderRegisterByCtx(ctx context.Context) (string, func(common.Sender) string, error) {
	var senderRegister func(common.Sender) string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", senderRegister, fmt.Errorf("missing metadata")
	}
	runtimeID := md.Get("RUNTIME_ID")
	if runtimeID == nil {
		return "", senderRegister, fmt.Errorf("missing runtime_id")
	}
	runtime_id := runtimeID[0]
	f, err := getSenderRegister(runtime_id)
	return runtime_id, f, err
}
