package data

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/indikay/notification-service/ent"
	"github.com/indikay/notification-service/ent/notification"
	"github.com/indikay/notification-service/ent/schema"
	"github.com/indikay/notification-service/internal/biz"
)

type notificationRepo struct {
	data *Data
	log  *log.Helper
}

func NewNotificationRepo(data *Data) biz.NotificationRepo {
	return &notificationRepo{
		data: data,
		log:  log.NewHelper(log.DefaultLogger),
	}
}

func (repo *notificationRepo) SaveNotification(ctx context.Context, notification *biz.Notification) error {
	var (
		userID   = uuid.MustParse(notification.UserID)
		dataJSON = schema.NotificationData{
			Amount:   notification.Data.Amount,
			Symbol:   notification.Data.Symbol,
			Name:     notification.Data.PackageName,
			Code:     notification.Data.Code,
			Tx:       notification.Data.Tx,
			Message:  notification.Data.Message,
			Referral: notification.Data.Referral,
		}
	)
	_, err := repo.data.db.Notification.Create().
		SetUserID(userID).
		SetTitleKey(notification.TitleKey).
		SetData(dataJSON).Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *notificationRepo) GetListUserNotification(ctx context.Context, userID uuid.UUID) ([]*ent.Notification, error) {
	data, err := repo.data.db.Notification.Query().Where(notification.UserID(userID)).Order(notification.ByCreatedAt(sql.OrderDesc())).All(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *notificationRepo) UpdateReadAll(ctx context.Context, userID uuid.UUID) error {
	_, err := repo.data.db.Notification.Update().Where(notification.UserID(userID)).SetRead(true).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *notificationRepo) UpdateRead(ctx context.Context, notiID uuid.UUID) error {
	_, err := repo.data.db.Notification.Update().Where(notification.ID(notiID)).SetRead(true).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}
