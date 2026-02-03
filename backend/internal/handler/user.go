package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/request"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/response"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/constant"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	nullvalue "github.com/o-ga09/zenn-hackthon-2026/pkg/null_value"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/ptr"
	"gorm.io/gorm"
)

type IUserServer interface {
	List(c echo.Context) error
	GetByID(c echo.Context) error
	GetByUID(c echo.Context) error
	GetByName(c echo.Context) error
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type UserServer struct {
	repo    domain.IUserRepository
	storage domain.IUserStorage
}

func NewUserServer(repo domain.IUserRepository, storage domain.IUserStorage) IUserServer {
	return &UserServer{
		repo:    repo,
		storage: storage,
	}
}

// List ユーザー一覧取得
func (s *UserServer) List(c echo.Context) error {
	ctx := c.Request().Context()

	var query request.ListQuery
	if err := c.Bind(&query); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := c.Validate(&query); err != nil {
		return errors.Wrap(ctx, err)
	}

	users, err := s.repo.FindAll(ctx, &domain.FindOptions{
		Limit:  ptr.PtrToInt(query.Limit),
		Offset: ptr.PtrToInt(query.Offset),
	})
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	responses := make([]*response.UserResponse, len(users))
	for i, user := range users {
		responses[i] = response.ToResponse(user)
	}

	return c.JSON(http.StatusOK, responses)
}

// GetByID IDでユーザー取得
func (s *UserServer) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	var param request.GetByIDParam
	if err := c.Bind(&param); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := c.Validate(&param); err != nil {
		return errors.Wrap(ctx, err)
	}

	user, err := s.repo.FindByID(ctx, &domain.User{BaseModel: domain.BaseModel{ID: param.ID}})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.MakeNotFoundError(ctx, "User not found")
		}
		return errors.Wrap(ctx, err)
	}

	if user.ProfileImage.Valid && !strings.HasPrefix(user.ProfileImage.String, "https://") {
		user.ProfileImage.String, err = s.storage.Get(ctx, user.ProfileImage.String)
		if err != nil {
			return errors.Wrap(ctx, err)
		}
	}

	return c.JSON(http.StatusOK, response.ToResponse(user))
}

func (s *UserServer) GetByName(c echo.Context) error {
	ctx := c.Request().Context()

	var param request.GetByNameParam
	if err := c.Bind(&param); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := c.Validate(&param); err != nil {
		return errors.Wrap(ctx, err)
	}

	user, err := s.repo.FindByName(ctx, &domain.User{Name: param.Name})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.MakeNotFoundError(ctx, "User not found")
		}
		return errors.Wrap(ctx, err)
	}

	if user.ProfileImage.Valid && !strings.HasPrefix(user.ProfileImage.String, "https://") {
		user.ProfileImage.String, err = s.storage.Get(ctx, user.ProfileImage.String)
		if err != nil {
			return errors.Wrap(ctx, err)
		}
	}

	return c.JSON(http.StatusOK, response.ToResponse(user))
}

// GetByUID Firebase UIDでユーザー取得
func (s *UserServer) GetByUID(c echo.Context) error {
	ctx := c.Request().Context()

	var query request.GetByUIDQuery
	if err := c.Bind(&query); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := c.Validate(&query); err != nil {
		return errors.Wrap(ctx, err)
	}

	user, err := s.repo.FindByUID(ctx, &domain.User{UID: query.UID})
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	if user.ProfileImage.Valid && !strings.HasPrefix(user.ProfileImage.String, "https://") {
		user.ProfileImage.String, err = s.storage.Get(ctx, user.ProfileImage.String)
		if err != nil {
			return errors.Wrap(ctx, err)
		}
	}

	return c.JSON(http.StatusOK, response.ToResponse(user))
}

// Create ユーザー作成
func (s *UserServer) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var req request.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := c.Validate(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	// 既存ユーザーの確認
	existingUser, _ := s.repo.FindByUID(ctx, &domain.User{UID: req.UID})
	if existingUser != nil {
		return errors.MakeConflictError(ctx, "User with the same UID already exists")
	}

	// ユーザー作成
	user := req.ToUser()

	// プランに応じた初期トークン残高設定
	if req.Plan == constant.UserPlanFree {
		user.TokenBalance = nullvalue.ToNullInt64(0)
	} else if req.Plan == constant.UserPlanPremium {
		user.TokenBalance = nullvalue.ToNullInt64(10000)
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.JSON(http.StatusCreated, response.ToResponse(user))
}

// Update ユーザー更新
func (s *UserServer) Update(c echo.Context) error {
	ctx := c.Request().Context()

	var req request.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := c.Validate(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	// 既存ユーザーの取得
	_, err := s.repo.FindByID(ctx, &domain.User{BaseModel: domain.BaseModel{ID: req.ID}})
	if err != nil {
		return errors.MakeNotFoundError(ctx, "指定されたユーザーは存在しません")
	}

	updateUser := req.ToUser()
	if req.ProfileImage != nil && !strings.HasPrefix(*req.ProfileImage, "https://") {
		key, err := s.storage.Upload(ctx, fmt.Sprintf("profile_images/%s", updateUser.Name), *req.ProfileImage)
		if err != nil {
			return errors.Wrap(ctx, err)
		}
		updateUser.ProfileImage = nullvalue.ToNullString(key)
	}
	if err := s.repo.Update(ctx, updateUser); err != nil {
		return errors.Wrap(ctx, err)
	}

	updatedUser, err := s.repo.FindByID(ctx, &domain.User{BaseModel: domain.BaseModel{ID: updateUser.ID}})
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	if updateUser.ProfileImage.Valid && !strings.HasPrefix(updateUser.ProfileImage.String, "https://") {
		updateUser.ProfileImage.String, err = s.storage.Get(ctx, updateUser.ProfileImage.String)
		if err != nil {
			return errors.Wrap(ctx, err)
		}
	}

	return c.JSON(http.StatusOK, response.ToResponse(updatedUser))
}

// Delete ユーザー削除
func (s *UserServer) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	// パスパラメータのバインドとバリデーション
	var param request.DeleteUserParam
	if err := c.Bind(&param); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := c.Validate(&param); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := s.repo.Delete(ctx, &domain.User{BaseModel: domain.BaseModel{ID: param.ID}}); err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.NoContent(http.StatusNoContent)
}
