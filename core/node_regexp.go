package core

import (
	"regexp"

	"github.com/dop251/goja"
)

func Regexp(call goja.ConstructorCall) *goja.Object {
	pattern := call.Arguments[0].String()
	call.This.Set("find", regexp.MustCompile(pattern).FindString)
	call.This.Set("findSubmatch", regexp.MustCompile(pattern).FindStringSubmatch)
	call.This.Set("findAll", regexp.MustCompile(pattern).FindAllString)
	call.This.Set("findAllSubmatch", regexp.MustCompile(pattern).FindAllStringSubmatch)
	call.This.Set("replaceAll", regexp.MustCompile(pattern).ReplaceAllString)
	return nil
}
