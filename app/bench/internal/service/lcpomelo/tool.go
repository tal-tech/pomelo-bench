package lcpomelo

import "github.com/zeromicro/go-zero/core/jsonx"

func structToJsonStr(s interface{}) string {
	v, _ := jsonx.MarshalToString(s)

	return v
}
