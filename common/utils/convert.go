package utils

import (
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

func Int642Str(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Str2Int64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		logx.Error(err)
	}
	return i
}
