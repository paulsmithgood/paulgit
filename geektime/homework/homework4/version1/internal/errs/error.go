package errs

import "errors"

var (
	ErrParseType   = errors.New("不符合解析类型")
	ErrInvalidTag  = errors.New("非法的标签")
	ErrUnknowField = errors.New("未知Field")
)
