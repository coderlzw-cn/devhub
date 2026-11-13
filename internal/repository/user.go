package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgconn"

	"gin-enterprise/internal/db/sqlc"
	"gin-enterprise/internal/model"
)

// UserPatch 是用户的局部更新。nil 字段保持不变。
type UserPatch struct {
	Username *string
	Email    *string
	Age      *int
}

// User 是用户表的数据访问对象。
type User struct {
	queries *sqlc.Queries
}

// NewUser 使用给定连接创建用户仓储。
func NewUser(db *sql.DB) *User {
	return &User{queries: sqlc.New(db)}
}

// Create 写入用户。
func (r *User) Create(ctx context.Context, user model.User) (model.User, error) {
	age, err := toInt32(user.Age)
	if err != nil {
		return model.User{}, fmt.Errorf("user age: %w", err)
	}
	var created sqlc.User
	if user.ID == 0 {
		created, err = r.queries.CreateUser(ctx, sqlc.CreateUserParams{
			Username: user.Username,
			Email:    user.Email,
			Age:      age,
			Password: user.Password,
		})
	} else {
		id, err := toInt64(user.ID)
		if err != nil {
			return model.User{}, fmt.Errorf("user id: %w", err)
		}
		created, err = r.queries.CreateUserWithID(ctx, sqlc.CreateUserWithIDParams{
			ID:       id,
			Username: user.Username,
			Email:    user.Email,
			Age:      age,
			Password: user.Password,
		})
	}
	if err != nil {
		return model.User{}, fmt.Errorf("insert user: %w", err)
	}
	result, err := userModel(created)
	if err != nil {
		return model.User{}, err
	}
	return result, nil
}

// List 按 ID 升序返回一页用户和总数。
func (r *User) List(ctx context.Context, skip, limit int64) ([]model.User, int, error) {
	total, err := r.queries.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	users, err := r.queries.ListUsers(ctx, sqlc.ListUsersParams{SkipCount: skip, PageLimit: limit})
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	result := make([]model.User, 0, len(users))
	for _, user := range users {
		item, err := userModel(user)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, item)
	}
	return result, int(total), nil
}

// GetByUsername 按用户名查询用户。不存在时返回 nil。
func (r *User) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	result, err := userModel(user)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Get 按 ID 查询用户。不存在时返回 nil。
func (r *User) Get(ctx context.Context, id uint64) (*model.User, error) {
	dbID, err := toInt64(id)
	if err != nil {
		return nil, fmt.Errorf("user id: %w", err)
	}
	user, err := r.queries.GetUser(ctx, dbID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	result, err := userModel(user)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update 按补丁更新用户。不存在时返回 nil。
func (r *User) Update(ctx context.Context, id uint64, patch UserPatch) (*model.User, error) {
	dbID, err := toInt64(id)
	if err != nil {
		return nil, fmt.Errorf("user id: %w", err)
	}
	params := sqlc.UpdateUserParams{ID: dbID}
	if patch.Username != nil {
		params.SetUsername, params.Username = true, *patch.Username
	}
	if patch.Email != nil {
		params.SetEmail, params.Email = true, *patch.Email
	}
	if patch.Age != nil {
		age, err := toInt32(*patch.Age)
		if err != nil {
			return nil, fmt.Errorf("user age: %w", err)
		}
		params.SetAge, params.Age = true, age
	}
	if !params.SetUsername && !params.SetEmail && !params.SetAge {
		return r.Get(ctx, id)
	}

	if _, err := r.queries.UpdateUser(ctx, params); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return r.Get(ctx, id)
}

// Delete 按 ID 删除用户，并返回被删除的用户。不存在时返回 nil。
func (r *User) Delete(ctx context.Context, id uint64) (*model.User, error) {
	dbID, err := toInt64(id)
	if err != nil {
		return nil, fmt.Errorf("user id: %w", err)
	}
	user, err := r.queries.DeleteUser(ctx, dbID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("delete user: %w", err)
	}
	result, err := userModel(user)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func userModel(user sqlc.User) (model.User, error) {
	if user.ID < 0 {
		return model.User{}, fmt.Errorf("invalid negative user id %d", user.ID)
	}
	return model.User{
		ID:       uint64(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Age:      int(user.Age),
		Password: user.Password,
	}, nil
}

// IsUniqueViolation reports whether a PostgreSQL constraint rejected a duplicate value.
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func toInt64(value uint64) (int64, error) {
	if value > math.MaxInt64 {
		return 0, fmt.Errorf("%d exceeds PostgreSQL BIGINT range", value)
	}
	return int64(value), nil
}
