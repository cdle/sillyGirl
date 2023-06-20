package core

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cdle/sillyGirl/emoji"
	"github.com/cdle/sillyGirl/utils"
)

type Strings struct {
}

func (sender *Strings) Diff(a, b []interface{}) []interface{} {
	m := make(map[interface{}]bool)
	c := make([]interface{}, 0)
	for _, v := range b {
		m[v] = true
	}
	for _, v := range a {
		if !m[v] {
			c = append(c, v)
		}
	}
	return c
}

// 寻找所有共同拥有的连续最长字符串
func (sender *Strings) Intersects(sa, sb string) []string {
	var common []string
	for i := 0; i < len(sa); i++ {
		for j := i + 1; j <= len(sa); j++ {
			substr := sa[i:j]
			if len(substr) > 0 && strings.Contains(sb, substr) && !Contains(common, substr) {
				common = append(common, substr)
			}
		}
	}
	return common
}

func (sender *Strings) Intersect(a, b interface{}) interface{} {
	// 判断是否为字符串类型
	if sa, ok := a.(string); ok {
		if sb, ok := b.(string); ok {
			// 如果是字符串类型，则寻找共同拥有的最长字符串
			var longest string
			for i := 0; i < len(sa); i++ {
				for j := i + 1; j <= len(sa); j++ {
					substr := sa[i:j]
					if strings.Contains(sb, substr) && len(substr) > len(longest) {
						longest = substr
					}
				}
			}
			return longest
		}
	}

	// 否则，使用普通的交集算法
	m := make(map[interface{}]bool)
	c := make([]interface{}, 0)
	for _, v := range a.([]interface{}) {
		m[v] = true
	}
	for _, v := range b.([]interface{}) {
		if m[v] {
			c = append(c, v)
		}
	}
	return c
}

func (sender *Strings) Union(a, b []interface{}) []interface{} {
	m := make(map[interface{}]bool)
	c := make([]interface{}, 0)
	for _, v := range a {
		m[v] = true
		c = append(c, v)
	}
	for _, v := range b {
		if !m[v] {
			c = append(c, v)
		}
	}
	return c
}

func (sender *Strings) Random(length int, substr string) string {
	ws := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if substr != "" {
		ws = substr
	}
	rand.Seed(time.Now().UnixNano())
	letters := []rune(ws)
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (sender *Strings) JoinFilepath(elem ...string) string {
	return filepath.Join(elem...)
}

func ForeachObject(m map[string]interface{}, callback func(key, value interface{}) bool) bool {
	for k, v := range m {
		switch v := v.(type) {
		case map[string]interface{}:
			if ForeachObject(v, callback) {
				return true
			}
		case []interface{}:
			for _, u := range v {
				if um, ok := u.(map[string]interface{}); ok {
					if ForeachObject(um, callback) {
						return true
					}
				}
			}
		case string:
			if callback(k, v) {
				return true
			}
		}
	}
	return false
}

func (sender *Strings) Trim(s, cutset string) string {
	return strings.Trim(s, cutset)
}
func (sender *Strings) TrimLeft(s, cutset string) string {
	return strings.TrimLeft(s, cutset)
}

func (sender *Strings) TrimRight(s, cutset string) string {
	return strings.TrimRight(s, cutset)
}

func (sender *Strings) Filename(path string) string {
	re := regexp.MustCompile(`[\\/]+`)
	parts := re.Split(path, -1)
	filename := parts[len(parts)-1]
	return filename
}

func (sender *Strings) Dir(path string) string {
	re := regexp.MustCompile(`[\\/]+`)
	parts := re.Split(path, -1)
	filename := parts[len(parts)-1]
	dir := path[:len(path)-len(filename)]
	dir = filepath.Clean(dir)
	return dir
}

func (sender *Strings) Contains(s string, substr interface{}) bool {
	switch substr := substr.(type) {
	case string:
		return strings.Contains(s, substr)
	case []string:
		for _, sub := range substr {
			if strings.Contains(s, sub) {
				return true
			}
		}
		return false
	case []interface{}:
		for _, sub := range substr {
			if strings.Contains(s, sub.(string)) {
				return true
			}
		}
		return false
	}
	return false
}

func (sender *Strings) ToLower(s string) string {
	return strings.ToLower(s)
}

func (sender *Strings) ToUpper(s string) string {
	return strings.ToUpper(s)
}

func (sender *Strings) Remove(ss []string, s string) []string {
	return utils.Remove(ss, s)
}

func (sender *Strings) Append(ss []string, s string) []string {
	return append(ss, s)
}
func (sender *Strings) HasPrefix(s, substr string) bool {
	return strings.HasPrefix(s, substr)
}
func (sender *Strings) HasSuffix(s, substr string) bool {
	return strings.HasSuffix(s, substr)
}

func (sender *Strings) Replace(s_ interface{}, old string, new string, n int) interface{} {
	switch s := s_.(type) {
	case string:
		if n == 0 {
			n = -1
		}
		return strings.Replace(s, old, new, n)
	case []string:
		for i := range s {
			if s[i] == old {
				s[i] = new
			}
		}
		return s
	case []interface{}:
		for i := range s {
			if s[i] == old {
				s[i] = new
			}
		}
		return s
	}
	return ""
}

func (sender *Strings) ReplaceAll(s string, old string, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func (sender *Strings) Split(s string, sep string, n int) []string {
	return strings.SplitN(s, sep, n)
}

func (sender *Strings) BuildQuery(params map[string]interface{}) string {
	return sender.EncodeQueryString(params)
}

func (sender *Strings) EncodeQueryString(params map[string]interface{}) string {
	var buf bytes.Buffer
	for key, value := range params {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(key))
		buf.WriteByte('=')
		switch v := value.(type) {
		case string:
			buf.WriteString(url.QueryEscape(v))
		case int:
			buf.WriteString(strconv.Itoa(v))
		case int64:
			buf.WriteString(strconv.Itoa(int(v)))
		case int32:
			buf.WriteString(strconv.Itoa(int(v)))
		case bool:
			buf.WriteString(strconv.FormatBool(v))
		default:
			buf.WriteString(url.QueryEscape(fmt.Sprintf("%v", v)))
		}
	}
	return buf.String()
}

func (sender *Strings) ParseQuery(querystring string) map[string]interface{} {
	return sender.DecodeQueryString(querystring)
}

func (sender *Strings) DecodeQueryString(querystring string) map[string]interface{} {
	u, err := url.Parse(querystring)

	if err != nil {
		panic(err)
	}
	params := make(map[string]interface{})
	for key, values := range u.Query() {
		if len(values) > 0 {
			value := values[0]
			// if intValue, err := strconv.Atoi(value); err == nil {
			// 	params[key] = intValue
			// } else if boolValue, err := strconv.ParseBool(value); err == nil {
			// 	params[key] = boolValue
			// } else {
			// 	params[key] = value
			// }
			params[key] = value
		}
	}
	return params
}

func (sender *Strings) HideCQEmoji(text string) map[string]interface{} {
	i := 0
	var ms = map[string]string{}
	text = regexp.MustCompile(`\[CQ:(\w+)(.*?)\]`).ReplaceAllStringFunc(text, func(s string) string {
		v := fmt.Sprintf("#%d#", i)
		i++
		ms[v] = s
		return v
	})
	text = emoji.ReplaceEmojisWithFunc(text, func(e emoji.Emoji) string {
		v := fmt.Sprintf("#%d#", i)
		i++
		ms[v] = e.Character
		return v
	})
	return map[string]interface{}{
		"text": text,
		"recover": func(text string) string {
			return regexp.MustCompile(`#\d{1,4}#`).ReplaceAllStringFunc(text, func(s string) string {
				return ms[s]
			})
		},
	}
}

// 构建CQ码
func (sender *Strings) BuildCQCode(cqType string, params map[string]interface{}) string {
	var sb strings.Builder
	sb.WriteString("[CQ:" + cqType)
	for k, v := range params {
		sb.WriteString(", ")
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(fmt.Sprintf("%v", v))
	}
	sb.WriteString("]")
	return sb.String()
}

// 将含有 CQ码 的文本解析成文本和 CQ 对象数组
func (sender *Strings) ParseCQText(text string) []interface{} {
	cqRegex := regexp.MustCompile(`\[CQ:(\w+)(.*?)\]`)
	cqMatches := cqRegex.FindAllStringSubmatch(text, -1)
	result := make([]interface{}, 0, len(cqMatches)*2+1)
	// 依次解析 CQ 码和文本
	lastIndex := 0
	for _, match := range cqMatches {
		// 添加 CQ 码前的文本
		if matchIndex := strings.Index(text[lastIndex:], match[0]); matchIndex > 0 {
			result = append(result, text[lastIndex:lastIndex+matchIndex])
			lastIndex += matchIndex
		}

		// 解析 CQ 码
		params := make(map[string]string)
		paramRegex := regexp.MustCompile(`(\w+)=([^,]+)`)
		paramMatches := paramRegex.FindAllStringSubmatch(match[2], -1)
		for _, paramMatch := range paramMatches {
			params[paramMatch[1]] = strings.TrimSpace(paramMatch[2])
		}
		result = append(result, CQ{
			Type:   match[1],
			Params: params,
		})

		lastIndex += len(match[0])
	}

	// 添加最后一个 CQ 码后的文本
	if lastIndex < len(text) {
		result = append(result, text[lastIndex:])
	}

	return result
}

// CQ 对象
type CQ struct {
	Type   string
	Params map[string]string
}

// 将 CQ 对象数组转换回原始文本
func (sender *Strings) StringifyCQText(cqList []interface{}) string {
	var sb strings.Builder
	for _, item := range cqList {
		switch item := item.(type) {
		case string:
			sb.WriteString(item)
		case CQ:
			sb.WriteString(fmt.Sprintf("[CQ:%s", item.Type))
			for k, v := range item.Params {
				sb.WriteString(fmt.Sprintf(",%s=%s", k, v))
			}
			sb.WriteString("]")
		case map[string]interface{}:
			cq := CQ{
				Type:   item["Type"].(string),
				Params: convertParams(item["Params"].(map[string]interface{})),
			}
			sb.WriteString(fmt.Sprintf("[CQ:%s", cq.Type))
			for k, v := range cq.Params {
				sb.WriteString(fmt.Sprintf(",%s=%s", k, v))
			}
			sb.WriteString("]")
		}
	}
	return sb.String()
}

// 将 map[string]interface{} 类型的 params 转换为 map[string]string 类型
func convertParams(params map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range params {
		if s, ok := v.(string); ok {
			result[k] = s
		}
	}
	return result
}

func (sender *Strings) ReplaceEmojis(str string, f func([]string) string) string {
	return emoji.ReplaceEmojisWithFunc(str, func(e emoji.Emoji) string {
		return f(e.CodePoint2)
	})
}

// `\[emoji=([0-9A-Z]{4})\]`
func (sender *Strings) ReplaceToEmojis(str string, pattern string) string {
	return emoji.ReplaceToEmojis(str, pattern)
}

func (sender *Strings) ExtractAddress(input string) string {
	return regexp.MustCompile(`http[s]?://[\w.]+:?\d*`).FindString(input)
}

func (sender *Strings) Unique(str ...interface{}) []string {
	return utils.Unique(str...)
}

func (sender *Strings) Longest(args ...interface{}) string {
	var longest string
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			if len(v) > len(longest) {
				longest = v
			}
		case []string:
			for _, s := range v {
				if len(s) > len(longest) {
					longest = s
				}
			}
		case []interface{}:
			for _, s := range v {
				if len(s.(string)) > len(longest) {
					longest = s.(string)
				}
			}
		case [][]string:
			for _, s := range v {
				longest = sender.Longest(s)
			}
		case [][]interface{}:
			for _, s := range v {
				longest = sender.Longest(s)
			}

		}
	}
	return longest
}

func (sender *Strings) Shortest(args ...interface{}) string {
	var longest string
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			if len(v) < len(longest) {
				longest = v
			}
		case []string:
			for _, s := range v {
				if len(s) < len(longest) {
					longest = s
				}
			}
		case []interface{}:
			for _, s := range v {
				if len(s.(string)) < len(longest) {
					longest = s.(string)
				}
			}
		case [][]string:
			for _, s := range v {
				longest = sender.Shortest(s)
			}
		case [][]interface{}:
			for _, s := range v {
				longest = sender.Shortest(s)
			}

		}
	}
	return longest
}
