package core

import (
	"errors"

	"github.com/dop251/goja"
)

func Buffer(vm *goja.Runtime, call goja.ConstructorCall) *goja.Object {
	var buffer []byte
	switch arg := call.Arguments[0].Export().(type) {
	case string:
		buffer = []byte(arg)
	case []byte:
		buffer = arg
	case int:
		buffer = make([]byte, arg)
	case int64:
		buffer = make([]byte, arg)
	default:
		panic(vm.NewTypeError("invalid argument type"))
	}
	obj := call.This.ToObject(vm)
	obj.Set("length", func() int {
		return len(buffer)
	})

	obj.Set("write", func(value interface{}, offset int, length int) {
		var buf []byte
		switch value := value.(type) {
		case string:
			buf = []byte(value)
		case []byte:
			buf = value
		default:
			panic(vm.NewTypeError("invalid argument type"))
		}

		if offset+length > len(buffer) {
			panic(vm.NewGoError(errors.New("out of range")))
		}

		copy(buffer[offset:], buf[:length])
	})

	obj.Set("toString", func() string {
		return string(buffer)
	})

	obj.Set("slice", func(start int, end int) []byte {
		if start < 0 {
			start += len(buffer)
		}
		if end < 0 {
			end += len(buffer)
		}
		if end > len(buffer) {
			end = len(buffer)
		}
		if start >= end || start >= len(buffer) {
			return nil
		}
		return buffer[start:end]
	})

	obj.Set("copy", func(target []byte, sourceStart int, targetStart int, sourceEnd int) int {
		if sourceStart < 0 {
			sourceStart += len(buffer)
		}
		if sourceEnd < 0 {
			sourceEnd += len(buffer)
		}
		if targetStart < 0 {
			targetStart += len(target)
		}
		if sourceEnd > len(buffer) {
			sourceEnd = len(buffer)
		}
		if sourceStart >= sourceEnd || sourceStart >= len(buffer) {
			return 0
		}
		return copy(target[targetStart:], buffer[sourceStart:sourceEnd])
	})

	obj.Set("join", func(sep []byte, buffers ...[]byte) []byte {
		var totalLen int
		for _, buf := range buffers {
			totalLen += len(buf)
		}
		if len(sep) > 0 && len(buffers) > 1 {
			totalLen += len(sep) * (len(buffers) - 1)
		}

		result := make([]byte, totalLen)
		offset := 0
		for i, buf := range buffers {
			copy(result[offset:], buf)
			offset += len(buf)
			if i < len(buffers)-1 && len(sep) > 0 {
				copy(result[offset:], sep)
				offset += len(sep)
			}
		}

		return result
	})

	return call.This
}
