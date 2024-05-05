package extsvc

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	pbEmail "github.com/indikay/email-service/api/email"
	"github.com/indikay/notification-service/internal/conf"
)

type EmailService struct {
	emailClient pbEmail.EmailServiceClient
}

func NewEmailService(serverConfig *conf.Data) *EmailService {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(serverConfig.Email.Addr),
		grpc.WithTimeout(serverConfig.Email.Timeout.AsDuration()),
		grpc.WithMiddleware(
			recovery.Recovery(),
			validate.Validator(),
		),
	)

	if err != nil {
		panic(err)
	}

	client := pbEmail.NewEmailServiceClient(conn)
	return &EmailService{
		emailClient: client,
	}
}

type SendEmail struct {
	To     string `json:"to"`
	Data   string `json:"data"`
	Action int32  `json:"action"`
}

type ValidationData struct {
	UserName      string `json:"username"`
	TelegramName  string `json:"telegramName"`
	ValidationURL string `json:"validationURL"`
	Code          string `json:"code"`
}

func (s *EmailService) SendValidationEmail(ctx context.Context, email string, locale string, data *ValidationData) (*pbEmail.BaseResponse, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	emailData := &pbEmail.SendEmailData{
		To:     email,
		Locale: locale,
		Action: pbEmail.EnumEmailAction_EMAIL_ACTION_TELEGRAM_REGISTER,
		Data:   string(jsonData),
	}
	resp, err := s.emailClient.SendEmail(ctx, &pbEmail.SendEmailRequest{
		Data: []*pbEmail.SendEmailData{emailData},
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
