package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/indikay/notification-service/ent"
	"github.com/indikay/notification-service/internal/extsvc"
	"github.com/indikay/notification-service/internal/utils"
)

type NotificationRepo interface {
	SaveNotification(ctx context.Context, notification *Notification) error
	GetListUserNotification(ctx context.Context, userID uuid.UUID) ([]*ent.Notification, error)
	UpdateReadAll(ctx context.Context, userID uuid.UUID) error
	UpdateRead(ctx context.Context, notiID uuid.UUID) error
}

type NotificationUseCase struct {
	notificationRepo NotificationRepo
	validateRepo     ValidateCodeRepo
	userSettingRepo  UserSettingRepo
	profileSvc       *extsvc.ProfileService
	log              *log.Helper
}

func NewNotificationUseCase(notificationRepo NotificationRepo, userSettingRepo UserSettingRepo, validateRepo ValidateCodeRepo, profileSvc *extsvc.ProfileService) *NotificationUseCase {
	return &NotificationUseCase{
		notificationRepo: notificationRepo,
		validateRepo:     validateRepo,
		userSettingRepo:  userSettingRepo,
		profileSvc:       profileSvc,
		log:              log.NewHelper(log.DefaultLogger)}
}

type Notification struct {
	TitleKey string           `json:"title_key"`
	UserID   string           `json:"user_id"`
	Data     NotificationData `json:"data"`
}

type NotificationData struct {
	Amount      string `json:"amount"`
	Symbol      string `json:"symbol"`
	Code        string `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	Tx          string `json:"tx,omitempty"`
	PackageName string `json:"package_name,omitempty"`
	Referral    string `json:"referral,omitempty"`
}

func (uc *NotificationUseCase) SaveNotification(ctx context.Context, notification *Notification) error {
	if len(notification.UserID) == 0 {
		return errors.New("MISSING_USER_ID")
	}

	if len(notification.TitleKey) == 0 {
		return errors.New("MISSING_TITLE_KEY")
	}

	return uc.notificationRepo.SaveNotification(ctx, notification)
}

func (uc *NotificationUseCase) GetListUserNotification(ctx context.Context, userID string) ([]*ent.Notification, error) {
	if len(userID) == 0 {
		return nil, errors.New("MISSING_USER_ID")
	}

	userUUID := uuid.MustParse(userID)
	notifications, err := uc.notificationRepo.GetListUserNotification(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
func (uc *NotificationUseCase) UpdateReadAllNotification(ctx context.Context, userID string) error {
	if len(userID) == 0 {
		return errors.New("MISSING_USER_ID")
	}

	userUUID := uuid.MustParse(userID)
	return uc.notificationRepo.UpdateReadAll(ctx, userUUID)
}

func (uc *NotificationUseCase) UpdateReadNotification(ctx context.Context, notiID string) error {
	if len(notiID) == 0 {
		return errors.New("MISSING_NOTIFICATION_ID")
	}
	notiUUID := uuid.MustParse(notiID)
	return uc.notificationRepo.UpdateRead(ctx, notiUUID)
}

func (uc *NotificationUseCase) ActiveTelegramBot(ctx context.Context, userId, token string) error {
	data, err := uc.validateRepo.GetToken(ctx, token)
	if err != nil {
		return err
	}
	activationData := utils.GetActivationData(data)
	profile, err := uc.profileSvc.GetProfileByEmail(ctx, activationData[0])
	if err != nil {
		return err
	}

	// if profile.Data.Id != userId {
	// 	return errors.New("INVALID_USER")
	// }

	nid := activationData[1]
	_, err = uc.userSettingRepo.Create(ctx, &UserSetting{UserID: profile.Data.Id, Type: TELEGRAM, Nid: nid, Enabled: true})
	if err != nil {
		return err
	}

	return nil
}

func (uc *NotificationUseCase) GetSettingsByUserId(ctx context.Context, userId string) ([]*UserSetting, error) {
	return uc.userSettingRepo.GetByUserId(ctx, userId)
}
