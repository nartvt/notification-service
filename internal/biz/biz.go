package biz

import (
	"context"
	"time"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewNotificationUseCase)

const (
	TELEGRAM = "TELEGRAM"
)

type UserSetting struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    string
	Type      string
	Nid       string
	Enabled   bool
}

type UserSettingRepo interface {
	Create(ctx context.Context, userSetting *UserSetting) (*UserSetting, error)
	Activate(ctx context.Context, id int) (*UserSetting, error)
	GetByUserId(ctx context.Context, userId string) ([]*UserSetting, error)
	GetByUserIdAndType(ctx context.Context, userId, notiType string) (*UserSetting, error)
}

type ValidateCodeRepo interface {
	Save(ctx context.Context, code, data string, timeoutInSecs int32) error
	SaveToken(ctx context.Context, token, data string, timeoutInSecs int32) error
	Get(ctx context.Context, code string) (string, error)
	GetToken(ctx context.Context, token string) (string, error)
}
