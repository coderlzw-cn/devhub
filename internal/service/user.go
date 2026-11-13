// Package service 实现业务逻辑。
package service

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"gin-enterprise/internal/dto/request"
	"gin-enterprise/internal/dto/response"
	httpresponse "gin-enterprise/internal/handler/http/response"
	"gin-enterprise/internal/model"
	"gin-enterprise/internal/repository"
	pkgresponse "gin-enterprise/pkg/response"
)

// ClientNotifier 向在线连接发送消息。
type ClientNotifier interface {
	Send(ctx context.Context, userID string, data []byte) int
	Broadcast(ctx context.Context, data []byte) int
}

// UserService 处理用户业务。
type UserService struct {
	repo   *repository.User
	notify ClientNotifier
}

// NewUserService 创建用户服务。notify 为空时不推送消息。
func NewUserService(repo *repository.User, notify ClientNotifier) *UserService {
	return &UserService{repo: repo, notify: notify}
}

// Create 创建用户。
func (s *UserService) Create(ctx context.Context, input *request.CreateUserInput) (*response.CreateUserOutput, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	user, err := s.repo.Create(ctx, model.User{
		Username: strings.TrimSpace(input.Body.Username),
		Email:    strings.TrimSpace(input.Body.Email),
		Age:      input.Body.Age,
		Password: string(hash),
	})
	if repository.IsUniqueViolation(err) {
		return nil, httpresponse.NewError(nil, http.StatusConflict, "Username or email already exists")
	}
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	return &response.CreateUserOutput{
		Status: http.StatusCreated,
		Body:   pkgresponse.SuccessData(user, "创建成功").APIBody(),
	}, nil
}

// Authenticate 校验用户名和密码，成功时返回用户 ID。
func (s *UserService) Authenticate(ctx context.Context, username, password string) (uint64, error) {
	// 1. 去除用户名的首尾空格，并根据处理后的用户名从数据库（Repository）中查询用户信息
	user, err := s.repo.GetByUsername(ctx, strings.TrimSpace(username))

	// 2. 如果查询报错，直接返回原始错误，以便上层调用者处理系统异常
	if err != nil {
		return 0, err
	}
	// 3. 如果查询结果为 nil，说明用户不存在，返回“凭证无效”错误码
	if user == nil {
		return 0, httpresponse.NewError(nil, http.StatusUnauthorized, "用户名或密码错误")
	}

	// 4. 用户存在，使用 bcrypt 算法比对用户输入的密码（明文）与数据库中存储的密码哈希值
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return 0, httpresponse.NewError(nil, http.StatusUnauthorized, "用户名或密码错误")
	}

	// 5. 用户名存在且密码比对成功，返回该用户的唯一 ID，错误返回 nil 表示认证成功
	return user.ID, nil
}

// List 按分页参数返回用户。
func (s *UserService) List(ctx context.Context, input *request.ListUsersInput) (*response.ListUsersOutput, error) {
	items, total, err := s.repo.List(ctx, input.Skip(), input.Limit())
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	page, pageSize := input.PageNumber(), int(input.Limit())
	return &response.ListUsersOutput{
		Body: pkgresponse.SuccessData(response.NewPage(items, total, page, pageSize), "").APIBody(),
	}, nil
}

// Get 按 ID 查询用户。
func (s *UserService) Get(ctx context.Context, input *request.GetUserInput) (*response.UserOutput, error) {
	user, err := s.repo.Get(ctx, input.ID)
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	if user == nil {
		return nil, httpresponse.NewError(nil, http.StatusNotFound, "User not found")
	}
	return &response.UserOutput{Body: pkgresponse.SuccessData(*user, "").APIBody()}, nil
}

// Update 按补丁更新用户。
func (s *UserService) Update(ctx context.Context, input *request.UpdateUserInput) (*response.UserOutput, error) {
	body := input.Body
	if body.Username == nil && body.Email == nil && body.Age == nil {
		return nil, httpresponse.NewError(nil, http.StatusUnprocessableEntity, "Invalid user input")
	}

	if body.Username != nil {
		value := strings.TrimSpace(*body.Username)
		body.Username = &value
	}
	if body.Email != nil {
		value := strings.TrimSpace(*body.Email)
		body.Email = &value
	}

	user, err := s.repo.Update(ctx, input.ID, repository.UserPatch{
		Username: body.Username,
		Email:    body.Email,
		Age:      body.Age,
	})
	if repository.IsUniqueViolation(err) {
		return nil, httpresponse.NewError(nil, http.StatusConflict, "Username or email already exists")
	}
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	if user == nil {
		return nil, httpresponse.NewError(nil, http.StatusNotFound, "User not found")
	}
	return &response.UserOutput{Body: pkgresponse.SuccessData(*user, "更新成功").APIBody()}, nil
}

// Notify 向指定用户当前在线的 WebSocket 连接发送一条文本消息。
func (s *UserService) Notify(ctx context.Context, input *request.NotifyUserInput) (*response.NotifyUserOutput, error) {
	user, err := s.repo.Get(ctx, input.ID)
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	if user == nil {
		return nil, httpresponse.NewError(nil, http.StatusNotFound, "User not found")
	}

	delivered := 0
	if s.notify != nil {
		payload, err := json.Marshal(map[string]string{"message": input.Body.Message})
		if err != nil {
			return nil, httpresponse.Error(err)
		}
		delivered = s.notify.Send(ctx, strconv.FormatUint(user.ID, 10), payload)
	}
	return &response.NotifyUserOutput{
		Body: pkgresponse.SuccessData(response.NotifyResult{Delivered: delivered}, "已发送").APIBody(),
	}, nil
}

// Broadcast 向当前全部在线的 WebSocket 连接发送一条文本消息。
func (s *UserService) Broadcast(ctx context.Context, input *request.BroadcastInput) (*response.NotifyUserOutput, error) {
	delivered := 0
	if s.notify != nil {
		payload, err := json.Marshal(map[string]string{"message": input.Body.Message})
		if err != nil {
			return nil, httpresponse.Error(err)
		}
		delivered = s.notify.Broadcast(ctx, payload)
	}
	return &response.NotifyUserOutput{
		Body: pkgresponse.SuccessData(response.NotifyResult{Delivered: delivered}, "已广播").APIBody(),
	}, nil
}

// Delete 按 ID 删除用户。
func (s *UserService) Delete(ctx context.Context, input *request.GetUserInput) (*response.DeleteUserOutput, error) {
	deleted, err := s.repo.Delete(ctx, input.ID)
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	if deleted == nil {
		return nil, httpresponse.NewError(nil, http.StatusNotFound, "User not found")
	}
	return &response.DeleteUserOutput{
		Body: pkgresponse.SuccessData(*deleted, "删除成功").APIBody(),
	}, nil
}
