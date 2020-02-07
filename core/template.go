package core

import (
	"regexp"
	"strings"
)

const (
	startFlag = '{'
	endFlag = '}'
)

// NewTemplate 解析模板，若模板内容中定义的参数多于实际提供的参数，则未提供的参数不会解析
//
// 示例1：
//	args := make(map[string]string)
//	args["url"] = "http://www.baidu.com"
//	tmpl := "{url }/aa/bb"
//	result := core.NewTemplate().Parse(tmpl, args)
// 结果1：http://www.baidu.com/aa/bb
//
// 示例2：
//	args := make(map[string]string)
//	args["url"] = "http://www.sina.com"
//	tmpl := "{url }/11/{test}"
//	result := core.NewTemplate().Parse(tmpl, args)
// 结果2：http://www.sina.com/11/{test}
func NewTemplate() *Template{
	return &Template{}
}

type Template struct {

}

func (t *Template) findArgVal(args map[string]string, argName string) string{
	for k,v := range args {
		reg := regexp.MustCompile(`{\s*` + k + `\s*}`)
		if reg.MatchString(argName) {
			return v
		}
	}
	return ""
}

// Parse 解析模板，若模板内容中定义的参数多于实际提供的参数，则为提供的参数不会解析
//    tmpl: 模板内容，参数以 {a} 形式
//    args: 模板参数
func (t *Template) Parse(tmpl string, args map[string]string) string{
	canPush := false
	argNameStack := MakeStringStack()

	argMap := make(map[string]string)

	for _, s := range tmpl {
		if s==startFlag {
			canPush = true
		}
		if canPush{
			_, _ = argNameStack.PushByte(byte(s))
		}
		if s== endFlag {
			canPush = false
			argName, _ := argNameStack.Pop()

			argMap[argName] = t.findArgVal(args, argName)
		}
	}
	for k,v := range argMap {
		if v!="" {
			tmpl = strings.ReplaceAll(tmpl, k, v)
		}
	}
	return tmpl
}
