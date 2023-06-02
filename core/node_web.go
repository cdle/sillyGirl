package core

import (
	"sync"
	"sync/atomic"

	"github.com/dop251/goja"
)

type RR struct {
	Req *Request
	Res *Response
	End func()
}

type HttpListen struct {
	Path   string
	Method string
	Chan   chan *RR
	UUID   string
	Closed bool
	// Handle func(*Request, *Response)
}

var httpListens sync.Map
var httpListensAny sync.Map
var listenCounter2 int64

func AddHttpListen(api, method string, vm *goja.Runtime, uuid string, resolve func(result interface{}), reject func(reason interface{})) {
	if method == "ANY" {
		var hl *HttpListen
		v, ok := httpListensAny.Load(uuid)
		if ok {
			// logs.Debug("load", uuid)
			hl = v.(*HttpListen)
		} else {
			hl = &HttpListen{
				Path:   api,
				Method: method,
				Chan:   make(chan *RR, 1000),
			}
			// logs.Debug("crate", uuid)
			httpListensAny.Store(uuid, hl)
		}
		f, ok := <-hl.Chan
		if ok {
			// logs.Debug("http resolved")
			resolve(f)
		} else {
			// logs.Debug("script is stopped")
			reject(Error(vm, "script is stopped"))
		}
		return
	}
	hl := &HttpListen{
		Path:   api,
		Method: method,
		UUID:   uuid,
	}
	key := atomic.AddInt64(&listenCounter2, 1)
	hl.Chan = make(chan *RR, 1)
	httpListens.Store(key, hl)
	f, ok := <-hl.Chan
	if ok {
		resolve(f)
	} else {
		reject(Error(vm, "script is stopped"))
	}
}

func CancelHttpListen(uuid string) {
	httpListens.Range(func(key, value any) bool {
		hl := value.(*HttpListen)
		if hl.UUID == uuid {
			httpListens.Delete(key)
			if !hl.Closed {
				hl.Closed = true
				close(hl.Chan)
			}
		}
		return true
	})
	httpListensAny.Range(func(key, value any) bool {
		hl := value.(*HttpListen)
		if hl.UUID == uuid {
			httpListens.Delete(key)
			if !hl.Closed {
				hl.Closed = true
				close(hl.Chan)
			}
		}
		return true
	})
}
