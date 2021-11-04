package core

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

//随机生成uuid
func GetUUID() string {
	u2, _ := uuid.NewUUID()
	return u2.String()
}

func Float64(str interface{}) float64 {
	f, _ := strconv.ParseFloat(fmt.Sprint(str), 64)
	return f
}

func RegistIm(i interface{}) Bucket {
	return NewBucket(regexp.MustCompile("[^/]+$").FindString(reflect.TypeOf(i).PkgPath()))
}
