package admin

// LoginRequest login params
type LoginRequest struct {
	Username    string `json:"username" binding:""`    // 账号
	Password    string `json:"password" binding:""`    // 密码
	CaptchaKey  string `json:"captchaKey" binding:""`  // 验证码key
	CaptchaCode string `json:"captchaCode" binding:""` // 验证码code
}

// LoginReply login reply
type LoginReply struct {
	Code int       `json:"code"` // return code
	Msg  string    `json:"msg"`  // return information description
	Data LoginItem `json:"data"` // return data
}

// CaptchaItem captcha item
type CaptchaItem struct {
	CaptchaKey    string `json:"captchaKey"`
	CaptchaBase64 string `json:"captchaBase64"`
}

// LoginItem login response item
type LoginItem struct {
	AccessToken string `json:"accessToken"` // access token
	Expires     int    `json:"expires"`     // expire time
	TokenType   string `json:"tokenType"`   // token type
}

// CaptchaReply captcha reply
type CaptchaReply struct {
	Code int         `json:"code"` // return code
	Msg  string      `json:"msg"`  // return information description
	Data CaptchaItem `json:"data"` // return data
}
