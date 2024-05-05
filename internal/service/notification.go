package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/indikay/go-core/middleware/jwt"
	"github.com/indikay/notification-service/api/notifications"
	"github.com/indikay/notification-service/ent"
	"github.com/indikay/notification-service/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NotificationService struct {
	notifications.UnimplementedNotificationServer
	notificationUC *biz.NotificationUseCase
	log            *log.Helper
}

func NewNotificationService(notificationUC *biz.NotificationUseCase) *NotificationService {
	return &NotificationService{
		log:            log.NewHelper(log.DefaultLogger),
		notificationUC: notificationUC,
	}
}

func (s *NotificationService) SaveNotification(ctx context.Context, notification *biz.Notification) error {
	err := s.notificationUC.SaveNotification(ctx, notification)
	if err != nil {
		return err
	}

	return nil
}

func (s *NotificationService) GetListUserNotification(ctx context.Context, _ *notifications.GetListUserNotificationRequest) (*notifications.GetListUserNotificationResponse, error) {
	userID, err := jwt.GetUserId(ctx)
	if err != nil {
		return &notifications.GetListUserNotificationResponse{
			Code:    1,
			Message: "UNAUTHORIZED",
			Key:     "FAIL",
		}, nil
	}
	fmt.Println("USER_ID_REQUEST: ", userID)
	data, err := s.notificationUC.GetListUserNotification(ctx, userID)
	if err != nil && !ent.IsNotFound(err) {
		return &notifications.GetListUserNotificationResponse{
			Code:    1,
			Message: err.Error(),
			Key:     "FAIL",
		}, nil
	}

	var notiResponse []*notifications.NotificationData
	for _, noti := range data {
		notificationData := s.toNotificationResponse(noti)
		notiResponse = append(notiResponse, notificationData)
	}

	return &notifications.GetListUserNotificationResponse{
		Code:    0,
		Message: "SUCCESS",
		Key:     "SUCCESS",
		Data:    notiResponse,
		Total:   int64(len(notiResponse)),
	}, nil
}

func (s *NotificationService) toNotificationResponse(notification *ent.Notification) *notifications.NotificationData {
	return &notifications.NotificationData{
		NotificationId: notification.ID.String(),
		UserId:         notification.UserID.String(),
		TitleKey:       notification.TitleKey,
		Read:           notification.Read,
		Data: &notifications.NotificationData_Data{
			Amount:   notification.Data.Amount,
			Symbol:   notification.Data.Symbol,
			Code:     notification.Data.Code,
			Message:  notification.Data.Message,
			Referral: notification.Data.Referral,
			Tx:       notification.Data.Tx,
		},
		CreatedTime: notification.CreatedAt.Format(time.RFC3339),
	}
}

func (s *NotificationService) ReadAllNotification(ctx context.Context, _ *notifications.ReadNotificationAllRequest) (*notifications.ReadNotificationAllResponse, error) {
	userID, err := jwt.GetUserId(ctx)
	if err != nil {
		return &notifications.ReadNotificationAllResponse{
			Code:    1,
			Message: "UNAUTHORIZED",
			Key:     "FAIL",
		}, nil
	}

	err = s.notificationUC.UpdateReadAllNotification(ctx, userID)
	if err != nil {
		return &notifications.ReadNotificationAllResponse{
			Code:    1,
			Message: err.Error(),
			Key:     "FAIL",
		}, nil
	}

	return &notifications.ReadNotificationAllResponse{
		Code:    0,
		Message: "SUCCESS",
		Key:     "SUCCESS",
	}, nil
}

func (s *NotificationService) ReadNotification(ctx context.Context, request *notifications.ReadNotificationRequest) (*notifications.ReadNotificationResponse, error) {
	err := s.notificationUC.UpdateReadNotification(ctx, request.NotificationId)
	if err != nil {
		return &notifications.ReadNotificationResponse{
			Code:    1,
			Message: err.Error(),
			Key:     "FAIL",
		}, nil
	}

	return &notifications.ReadNotificationResponse{
		Code:    0,
		Message: "SUCCESS",
		Key:     "SUCCESS",
	}, nil
}

func (s *NotificationService) TelegramActivation(ctx context.Context, req *notifications.TelegramActivationRequest) (*notifications.TelegramActivationResponse, error) {
	userID, err := jwt.GetUserId(ctx)
	if err != nil {
		return &notifications.TelegramActivationResponse{
			Code:   1,
			Msg:    "UNAUTHORIZED",
			MsgKey: "UNAUTHORIZED",
		}, nil
	}
	err = s.notificationUC.ActiveTelegramBot(ctx, userID, req.Token)
	if err != nil {
		return &notifications.TelegramActivationResponse{
			Code:   1,
			Msg:    err.Error(),
			MsgKey: "ERROR",
		}, nil
	}
	return &notifications.TelegramActivationResponse{}, nil
}

func (s *NotificationService) GetNotificationSettings(ctx context.Context, req *emptypb.Empty) (*notifications.GetNotificationSettingsResponse, error) {
	userID, err := jwt.GetUserId(ctx)
	if err != nil {
		return &notifications.GetNotificationSettingsResponse{
			Code:   1,
			Msg:    "UNAUTHORIZED",
			MsgKey: "UNAUTHORIZED",
		}, nil
	}
	settings, err := s.notificationUC.GetSettingsByUserId(ctx, userID)
	if err != nil {
		return &notifications.GetNotificationSettingsResponse{
			Code:   1,
			Msg:    err.Error(),
			MsgKey: "ERROR",
		}, nil
	}

	resp := make([]*notifications.GetNotificationSettingsResponse_NotificationSetting, len(settings))
	for _, v := range settings {
		resp = append(resp, &notifications.GetNotificationSettingsResponse_NotificationSetting{Type: v.Type, Active: v.Enabled})
	}

	return &notifications.GetNotificationSettingsResponse{
		Data: resp,
	}, nil
}
