package logger

import (
	"coding-guidelines/common/pkg/enum/logleveltype"
)

type LogLevel = logleveltype.Variant

func ParseLogLevel(s string) LogLevel {
	res := logleveltype.Parse(s)
	if res.IsSuccess() {
		return res.Data()
	}

	return LevelUnknown
}
