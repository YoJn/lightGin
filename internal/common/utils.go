package common

import (
	"reflect"
	"runtime"
	"strings"
)

// NameOfFunction 获得最后一个函数
func NameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func FilterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func ChooseData(custom, wildcard interface{}) interface{} {
	if custom != nil {
		return custom
	}
	if wildcard != nil {
		return wildcard
	}
	panic("negotiation config is invalid")
}

func Assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func ParseAccept(acceptHeader string) []string {
	parts := strings.Split(acceptHeader, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(strings.Split(part, ";")[0]); part != "" {
			out = append(out, part)
		}
	}
	return out
}
