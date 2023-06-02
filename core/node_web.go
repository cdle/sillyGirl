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
var listenCounter2 int64

func AddHttpListen(api, method string, vm *goja.Runtime, uuid string, resolve func(result interface{}), reject func(reason interface{})) {
	key := atomic.AddInt64(&listenCounter2, 1)
	hl := &HttpListen{
		Path:   api,
		Method: method,
		Chan:   make(chan *RR, 1),
	}
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
}
