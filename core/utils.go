package core

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

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

func TrimHiddenCharacter(originStr string) string {
	srcRunes := []rune(originStr)
	dstRunes := make([]rune, 0, len(srcRunes))
	for _, c := range srcRunes {
		if c >= 0 && c <= 31 && c != 10 {
			continue
		}
		if c == 127 {
			continue
		}

		dstRunes = append(dstRunes, c)
	}
	return strings.ReplaceAll(string(dstRunes), "￼", "")
}

func ForCQ(content string, callback func(key string, values map[string]string)) {

}
