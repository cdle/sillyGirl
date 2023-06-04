package core

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Strings struct {
}

func (sender *Strings) Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func (sender *Strings) Replace(s string, old string, new string, n int) string {
	if n == 0 {
		n = -1
	}
	return strings.Replace(s, old, new, n)
}

func (sender *Strings) ReplaceAll(s string, old string, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func (sender *Strings) Split(s string, sep string, n int) []string {
	return strings.SplitN(s, sep, n)
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

func (sender *Strings) DecodeQueryString(querystring string) map[string]interface{} {
	values, err := url.ParseQuery(querystring)
	if err != nil {
		panic(err)
	}
	params := make(map[string]interface{})
	for key, values := range values {
		if len(values) > 0 {
			value := values[0]
			if intValue, err := strconv.Atoi(value); err == nil {
				params[key] = intValue
			} else if boolValue, err := strconv.ParseBool(value); err == nil {
				params[key] = boolValue
			} else {
				params[key] = value
			}
		}
	}
	return params
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
