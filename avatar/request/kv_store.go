package request

import "avatar/utils"

// ctx.BindJSON	绑定 JSON Body，调用 json.Unmarshal() 进行反序列化，需要 Body 为 application/json 格式
type UpsertKVStore struct {
	KeyName string `json:"key_name"`
	Value   string `json:"value"`
}

type ListKVStore struct {
	KeyName string `json:"key_name"`
	utils.PageReq
}
