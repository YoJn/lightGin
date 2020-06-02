package common

import (
	"reflect"
	"runtime"
)

// NameOfFunction 获得最后一个函数
func NameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
