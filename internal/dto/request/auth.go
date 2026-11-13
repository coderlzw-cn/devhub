package request

// LoginBody 是用户名密码登录的请求体。
type LoginBody struct {
	Username string `json:"username" minLength:"1" maxLength:"32" example:"admin" doc:"用户名"`
	Password string `json:"password" minLength:"1" maxLength:"72" example:"111111" doc:"登录密码"`
}

// LoginInput 是用户名密码登录请求。
type LoginInput struct {
	Body LoginBody
}

// TokenBody 是 OAuth2 客户端凭证模式的请求体。
type TokenBody struct {
	GrantType    string `json:"grant_type" enum:"client_credentials" example:"client_credentials" doc:"授权类型，固定为 client_credentials"`
	ClientID     string `json:"client_id" minLength:"1" example:"gin-enterprise" doc:"客户端 ID"`
	ClientSecret string `json:"client_secret" minLength:"1" example:"client-secret" doc:"客户端密钥"`
}

// TokenInput 是 OAuth2 客户端凭证模式的令牌请求。
type TokenInput struct {
	Body TokenBody
}
