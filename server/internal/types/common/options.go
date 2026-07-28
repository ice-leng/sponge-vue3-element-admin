package common

// Options 下拉选项
type Options struct {
	Label    string      `json:"label"` // 标签
	Value    interface{} `json:"value"` // 值
	Other    interface{} `json:"other"`
	Children []Options   `json:"children"`
}

// OptionsReply only for api docs
type OptionsReply struct {
	Code int       `json:"code"` // return code
	Msg  string    `json:"msg"`  // return information description
	Data []Options `json:"data"` // return data
}
