package emoji

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
)

func init() {
	for i, emoji := range emojiMap {
		code_points := []string{}
		for _, v := range strings.Split(emojiMap[i].CodePoint, " ") {
			switch len(v) {
			case 4:
				code_points = append(code_points, v)
			case 5:
				code_points = append(code_points, transfer(v)...)
			default:
			}
		}
		if len(code_points) == 0 {
			for _, v := range strings.Split(emojiMap[i].CodePoint, "-") {
				switch len(v) {
				case 4:
					code_points = append(code_points, v)
				case 5:
					code_points = append(code_points, transfer(v)...)
				default:
				}
			}
		}
		if len(code_points) == 0 {
			delete(emojiMap, i)
		} else {
			emoji.CodePoint2 = code_points
			emojiMap[i] = emoji
		}
	}
}

func ConvertRegex(input string) string {
	// 定义要转义的字符集合
	escapeSet := []string{"\\", "[", "]", "(", ")", "{", "}", ".", "*", "+", "?", "^", "$", "|"}

	// 定义正则表达式模式和替换字符串
	var pattern = `%s`
	var replaceStr = `([0-9A-Za-z]{4})`

	// 对可能引起正则表达式冲突的字符进行转义
	for _, esc := range escapeSet {
		input = strings.ReplaceAll(input, esc, "\\"+esc)
	}

	// 将子匹配内容替换为指定格式的正则表达式
	replaced := strings.ReplaceAll(input, pattern, replaceStr)

	// 返回转换结果和是否大写的标记
	return replaced
}

// \[emoji=[0-9A-Z]{4}\]
func ReplaceToEmojis(str string, pattern string) string {
	format := pattern
	pattern = ConvertRegex(pattern)
	str = regexp.MustCompile("[0-9]"+pattern+pattern).ReplaceAllStringFunc(str, func(s string) string {
		if !strings.EqualFold(strings.ToUpper(s[1:]), strings.ToUpper(fmt.Sprintf(format+format, "FE0F", "20E3"))) { //"[emoji=FE0F][emoji=20E3]"
			return s
		}
		return fmt.Sprintf(format, fmt.Sprintf("%04X", s[0])) + s[1:]
	})
	res := str
	h := func(ssr *[]string) string {
		chs := []string{}
		ch := ""
		l := 100
		lft := []string{}
	again:
		for j := range emojiMap {
			ok, left := matchPrefix(emojiMap[j].CodePoint2, *ssr)
			if ok {
				if ln := len(left); ln < l {
					ch = emojiMap[j].Character
					l = ln
					lft = left
				}
			}
		}
		if ch != "" {
			chs = append(chs, ch)
		}
		if len(lft) != 0 && ch != "" {
			ch = ""
			l = 100
			*ssr = lft
			lft = []string{}
			goto again
		}
		return strings.Join(chs, "")
	}
	sss := regexp.MustCompile(pattern).FindAllStringSubmatchIndex(str, -1)
	ssr := []string{}
	ssrs := []string{}
	e := -1
	for i := range sss {
		ss := sss[i]
		a, b, c, d := ss[0], ss[1], ss[2], ss[3]

		//结算
		if e != a && len(ssr) != 0 {
			rs := h(&ssr)
			res = strings.Replace(res, strings.Join(ssrs, ""), rs, 1)
			ssr = []string{}
			ssrs = []string{}
			e = -1
		}

		ssrs = append(ssrs, str[a:b])
		// if upper {
		// 	ssr = append(ssr, str[c:d])
		// } else {
		ssr = append(ssr, strings.ToUpper(str[c:d]))
		// }
		e = b
	}
	if len(ssr) != 0 {
		rs := h(&ssr)
		res = strings.Replace(res, strings.Join(ssrs, ""), rs, 1)
		e = -1
	}
	return res
}

func transfer(str string) []string {
	// 将输入字符串解析为十六进制数
	unicodeValue, err := strconv.ParseUint(str, 16, 32)
	if err != nil {
		panic(err)
	}
	// 将 Unicode 编码值转换为 UTF-16 编码单元
	utf16Units := utf16.Encode([]rune{rune(unicodeValue)})

	// 将 UTF-16 编码单元格式化为字符串，并返回字符串切片
	return []string{fmt.Sprintf("%04X", utf16Units[0]), fmt.Sprintf("%04X", utf16Units[1])}
}

func matchPrefix(s1 []string, s2 []string) (bool, []string) {
	if len(s1) > len(s2) {
		return false, s2
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false, nil
		}
	}
	return true, s2[len(s1):]
}

func RuneToHx(a rune) string {
	return fmt.Sprintf("%04X", a)
}

func HexToRune(hexStr string) (rune, error) {
	// 解析十六进制字符串，获取对应的无符号整数
	u, err := strconv.ParseUint(hexStr, 16, 32)
	if err != nil {
		return 0, err
	}

	// 将无符号整数转换为 rune 类型的 Unicode 编码值
	r := rune(u)

	// 返回转换结果
	return r, nil
}
