package response

import pkgresponse "gin-enterprise/pkg/response"

// Token 是访问令牌数据。
type Token struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" doc:"JWT 访问令牌"`
	TokenType   string `json:"token_type" example:"Bearer" doc:"令牌类型，固定为 Bearer"`
	ExpiresIn   int    `json:"expires_in" example:"3600" doc:"有效期，单位秒"`
}

// TokenOutput 是访问令牌响应。
type TokenOutput struct {
	Body pkgresponse.Body[Token]
}
