package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/indikay/notification-service/ent"
	"github.com/indikay/notification-service/ent/usersetting"
	"github.com/indikay/notification-service/internal/biz"
)

type userSettingRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserSettingRepo(data *Data) biz.UserSettingRepo {
	return &userSettingRepo{
		data: data,
		log:  log.NewHelper(log.DefaultLogger),
	}
}

func (u *userSettingRepo) Create(ctx context.Context, userSetting *biz.UserSetting) (*biz.UserSetting, error) {
	item, err := u.data.db.UserSetting.Create().SetUserID(userSetting.UserID).SetNid(userSetting.Nid).SetType(userSetting.Type).SetEnabled(userSetting.Enabled).Save(ctx)
	if err != nil {
		if err.Error() != "pq: duplicate key value" {
			return nil, nil
		}
		return nil, err
	}

	return &biz.UserSetting{ID: item.ID, UserID: item.UserID, Type: item.Type, Nid: item.Nid, Enabled: item.Enabled, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}, nil
}

func (u *userSettingRepo) Activate(ctx context.Context, id int) (*biz.UserSetting, error) {
	item, err := u.data.db.UserSetting.UpdateOne(&ent.UserSetting{ID: id}).SetEnabled(true).Save(ctx)
	if err != nil {
		return nil, err
	}

	return &biz.UserSetting{ID: item.ID, UserID: item.UserID, Type: item.Type, Nid: item.Nid, Enabled: item.Enabled, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}, nil
}

func (u *userSettingRepo) GetByUserId(ctx context.Context, userId string) ([]*biz.UserSetting, error) {
	settings, err := u.data.db.UserSetting.Query().Where(usersetting.UserID(userId)).All(ctx)
	if err != nil {
		return nil, err
	}

	var rs []*biz.UserSetting
	for _, item := range settings {
		rs = append(rs, &biz.UserSetting{ID: item.ID, UserID: item.UserID, Type: item.Type, Nid: item.Nid, Enabled: item.Enabled, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt})
	}

	return rs, nil
}

func (u *userSettingRepo) GetByUserIdAndType(ctx context.Context, userId, notiType string) (*biz.UserSetting, error) {
	item, err := u.data.db.UserSetting.Query().Where(usersetting.UserID(userId), usersetting.TypeEQ(notiType)).First(ctx)
	if err != nil {
		return nil, err
	}

	return &biz.UserSetting{ID: item.ID, UserID: item.UserID, Type: item.Type, Nid: item.Nid, Enabled: item.Enabled, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}, nil
}
