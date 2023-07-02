package core

import (
	"errors"

	"github.com/dop251/goja"
)

type BufferObj struct {
	vm *goja.Runtime
}

func (b *BufferObj) From(v interface{}, EF string) *Buffer {
	var buffer []byte
	switch v := v.(type) {
	case string:
		if EF != "" {
			buffer = Convert(b.vm, v, EF, "bytes").([]byte)
		} else {
			buffer = []byte(v)
		}
	case []byte:
		if EF != "" {
			buffer = Convert(b.vm, v, EF, "bytes").([]byte)
		} else {
			buffer = v
		}
	case int:
		buffer = make([]byte, v)
	case int64:
		buffer = make([]byte, v)
	default:
		panic(b.vm.NewTypeError("invalid argument type"))
	}
	return &Buffer{
		vm:    b.vm,
		value: buffer,
	}
}

func (b *BufferObj) Alloc(v int) *Buffer {
	return &Buffer{
		vm:    b.vm,
		value: make([]byte, v),
	}
}

type Buffer struct {
	value []byte
	vm    *goja.Runtime
}

func (b *Buffer) Length() int {
	return len(b.value)
}

func (b *Buffer) Write(value interface{}, offset int, length int) int {
	var buf []byte
	switch value := value.(type) {
	case string:
		buf = []byte(value)
	case []byte:
		buf = value
	default:
		panic(b.vm.NewTypeError("invalid argument type"))
	}

	if offset+length > len(b.value) {
		panic(b.vm.NewGoError(errors.New("out of range")))
	}
	copy(b.value[offset:], buf[:length])
	return len(b.value)
}

func (b *Buffer) ToString(EF string) interface{} {
	if EF != "" {
		return Convert(b.vm, b.value, "", EF)
	}
	return string(b.value)
}

func (b *Buffer) Slice(start int, end int) *Buffer {
	if start < 0 {
		start += len(b.value)
	}
	if end <= 0 {
		end += len(b.value)
	}
	if end > len(b.value) {
		end = len(b.value)
	}
	if start >= end || start >= len(b.value) {
		return nil
	}
	return &Buffer{
		value: b.value[start:end],
	}
}

func (b *Buffer) Copy(target []byte, sourceStart int, targetStart int, sourceEnd int) int {
	if sourceStart < 0 {
		sourceStart += len(b.value)
	}
	if sourceEnd < 0 {
		sourceEnd += len(b.value)
	}
	if targetStart < 0 {
		targetStart += len(target)
	}
	if sourceEnd > len(b.value) {
		sourceEnd = len(b.value)
	}
	if sourceStart >= sourceEnd || sourceStart >= len(b.value) {
		return 0
	}
	return copy(target[targetStart:], b.value[sourceStart:sourceEnd])
}

func (b *Buffer) Join(sep []byte, buffers ...[]byte) []byte {
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
}
