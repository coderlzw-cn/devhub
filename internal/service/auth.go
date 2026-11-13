package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"gin-enterprise/internal/dto/request"
	"gin-enterprise/internal/dto/response"
	httpresponse "gin-enterprise/internal/handler/http/response"
	pkgresponse "gin-enterprise/pkg/response"
)

// AuthConfig 是签发和校验访问令牌所需的配置。
type AuthConfig struct {
	ClientID     string
	ClientSecret string
	JWTSecret    string
	Issuer       string
	TTL          time.Duration
}

// AccessToken 是 OAuth2 客户端凭证模式的令牌响应。
type AccessToken struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int
}

// AuthService 签发并校验 JWT 访问令牌。
type AuthService struct {
	cfg   AuthConfig
	users *UserService
}

// NewAuthService 创建认证服务。
func NewAuthService(cfg AuthConfig, users *UserService) *AuthService {
	return &AuthService{cfg: cfg, users: users}
}

// Login 校验用户名和密码后签发访问令牌。
func (s *AuthService) Login(ctx context.Context, input *request.LoginInput) (*response.TokenOutput, error) {
	userID, err := s.users.Authenticate(ctx, input.Body.Username, input.Body.Password)
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	token, err := s.IssueFor(strconv.FormatUint(userID, 10))
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	return tokenOutput(token), nil
}

// Issue 校验客户端凭证后签发访问令牌。
func (s *AuthService) Issue(_ context.Context, input *request.TokenInput) (*response.TokenOutput, error) {
	body := input.Body
	if body.GrantType != "client_credentials" {
		return nil, httpresponse.NewError(nil, http.StatusBadRequest, "Unsupported grant type")
	}
	if !hmac.Equal([]byte(body.ClientID), []byte(s.cfg.ClientID)) ||
		!hmac.Equal([]byte(body.ClientSecret), []byte(s.cfg.ClientSecret)) {
		return nil, httpresponse.NewError(nil, http.StatusUnauthorized, "Invalid client")
	}

	token, err := s.IssueFor(body.ClientID)
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	return tokenOutput(token), nil
}

// IssueFor 为指定主体签发访问令牌。
func (s *AuthService) IssueFor(subject string) (AccessToken, error) {
	now := time.Now()
	jti, err := newTokenID()
	if err != nil {
		return AccessToken{}, fmt.Errorf("create token id: %w", err)
	}

	claims := jwt.RegisteredClaims{
		Issuer:    s.cfg.Issuer,
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.TTL)),
		ID:        jti,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return AccessToken{}, fmt.Errorf("sign access token: %w", err)
	}

	return AccessToken{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.cfg.TTL.Seconds()),
	}, nil
}

// Subject 校验访问令牌并返回 sub。
func (s *AuthService) Subject(raw string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("Unauthorized")
		}
		return []byte(s.cfg.JWTSecret), nil
	},
		jwt.WithIssuer(s.cfg.Issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !token.Valid || claims.Subject == "" {
		return "", errors.New("Unauthorized")
	}
	return claims.Subject, nil
}

func tokenOutput(token AccessToken) *response.TokenOutput {
	return &response.TokenOutput{
		Body: pkgresponse.SuccessData(response.Token{
			AccessToken: token.AccessToken,
			TokenType:   token.TokenType,
			ExpiresIn:   token.ExpiresIn,
		}, "").APIBody(),
	}
}

func newTokenID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
