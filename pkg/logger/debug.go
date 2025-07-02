package logger

import (
	"github.com/astaxie/beego/logs"
)

var Debug bool = true

func DebugLog(format string, args ...interface{}) {
	if Debug {
		logs.Debug("[DEBUG] "+format+"\n", args...)
	}
}

func ErrorLog(format string, args ...interface{}) {
	if Debug {
		logs.Error("[ERROR] "+format+"\n", args...)
	}
}
