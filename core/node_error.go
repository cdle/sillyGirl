package core

import (
	"errors"

	"github.com/dop251/goja"
)

func Error(vm *goja.Runtime, errs ...interface{}) *goja.Object {
	var msg error
	for _, err := range errs {
		switch err := err.(type) {
		case string:
			msg = errors.New(err)
		case error:
			msg = err
		}
	}
	return vm.NewGoError(msg)
}
