package types

type UploadLocalReply struct {
	Code int        `json:"code"` // return code
	Msg  string     `json:"msg"`  // return information description
	Data UploadItem `json:"data"` // return data
}

type UploadItem struct {
	Name string `json:"name"` // 文件名称
	Url  string `json:"url"`  // url
	Path string `json:"path"` // path
}
