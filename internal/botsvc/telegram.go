package botsvc

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/indikay/notification-service/internal/biz"
	"github.com/indikay/notification-service/internal/conf"
	"github.com/indikay/notification-service/internal/constants"
	"github.com/indikay/notification-service/internal/extsvc"
	"github.com/indikay/notification-service/internal/utils"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

const (
	REGISTER = "/register"
	CODE     = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type TelegramBot struct {
	botSvc           *bot.Bot
	alertSettingRepo biz.UserSettingRepo
	validateRepo     biz.ValidateCodeRepo
	emailSvc         *extsvc.EmailService
	profileSvc       *extsvc.ProfileService
	logger           *log.Helper
	botConf          *conf.Data_TelegramBot
}

func NewTelegramBot(botConf *conf.Data, alertSettingRepo biz.UserSettingRepo, validateRepo biz.ValidateCodeRepo, emailSvc *extsvc.EmailService, profileSvc *extsvc.ProfileService) BotService {

	inst := &TelegramBot{logger: log.NewHelper(log.DefaultLogger), alertSettingRepo: alertSettingRepo, validateRepo: validateRepo, emailSvc: emailSvc, profileSvc: profileSvc, botConf: botConf.Telegram}
	opts := []bot.Option{
		bot.WithDefaultHandler(telegramDefaultHandler),
	}
	inst.botSvc, _ = bot.New(inst.botConf.Token, opts...)
	inst.botSvc.RegisterHandler(bot.HandlerTypeMessageText, "/register", bot.MatchTypePrefix, inst.registerHandler)
	return inst
}

func (t *TelegramBot) Start(ctx context.Context) error {
	go t.botSvc.Start(ctx)
	// t.SendIndicatorAlertMessage(ctx, "1eb1a218-62cf-43dc-9d15-f1f0318c7813", &IndicatorAlert{Indicator: "4H", ClosePrice: 12345, Symbol: "BTC", Signal: "UP", Type: "close"})
	return nil
}

func (t *TelegramBot) Stop(ctx context.Context) error {
	_, err := t.botSvc.Close(ctx)
	return err
}

func (t *TelegramBot) SendIndicatorAlertMessage(ctx context.Context, userId string, alert *IndicatorAlert) error {
	candleType := "Candle Close"
	if alert.Type == "real_time" {
		candleType = "Realtime"
	}

	p := message.NewPrinter(language.English)

	msg := p.Sprintf(`
<b>INDIKAY ALERT - $%s</b>
				💰 Price: <b>$%v</b>
				🕘 Time Frame: <b>%s</b>
				📈 Signal: <b>%s</b>
				🕹 Candle conditions: <b>%s</b>
<b><a href="https://indikay.com/crypto/%s">Visit On INDIKAY</a></b>
	`, alert.Symbol, number.Decimal(alert.ClosePrice), strings.ToUpper(alert.TimeFrame), strings.ToUpper(alert.Signal), candleType, alert.Symbol)
	return t.SendMessage(ctx, userId, msg)
}

func (t *TelegramBot) SendMessage(ctx context.Context, userId string, msg string) error {
	setting, err := t.alertSettingRepo.GetByUserIdAndType(ctx, userId, constants.EVENT_TYPE_TELEGRAM)
	if err != nil {
		t.logger.Error("SendMessage ", err)
	}
	_, err = t.botSvc.SendMessage(ctx, &bot.SendMessageParams{ChatID: setting.Nid, Text: msg, ParseMode: models.ParseModeHTML})
	if err != nil {
		t.logger.Error("SendMessage ", err)
	}
	return err
}

// Initial implements BotService.
func (t *TelegramBot) Initial(ctx context.Context, data interface{}) error {
	msg := data.(*models.Update)
	if msg.Message == nil {
		return nil
	}

	t.botSvc.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: msg.Message.Chat.ID,
		Text: `
To register signal alert from INDIKAY you need to type:
	/register {email} with {email} is the email which registered on indikay.com`,
	})

	return nil
}

// ReceiveMessage implements BotService.
func (t *TelegramBot) ReceiveMessage(ctx context.Context, msg interface{}) error {
	data := msg.(*models.Update)
	if data.Message == nil {
		return nil
	}

	return nil
}

func (t *TelegramBot) registerHandler(ctx context.Context, b *bot.Bot, data *models.Update) {
	msg := data.Message.Text
	email := strings.ReplaceAll(msg, REGISTER, "")
	email = strings.ReplaceAll(email, " ", "")
	code, err := gonanoid.Generate(CODE, 6)
	if err != nil {
		t.logger.Error("registerHandler ", err)
	}
	activationData := utils.BuildActivationData(email, fmt.Sprintf("%d", data.Message.Chat.ID))
	err = t.validateRepo.Save(ctx, code, activationData, int32(t.botConf.Timeout.Seconds))
	if err != nil {
		t.logger.Error("registerHandler ", err)
	}

	token := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s:%s", code, email))))
	err = t.validateRepo.SaveToken(ctx, token, activationData, int32(t.botConf.Timeout.Seconds))
	if err != nil {
		t.logger.Error("registerHandler ", err)
	}

	activationUrl := fmt.Sprintf("%s/%s", t.botConf.UrlActivation, token)

	profile, err := t.profileSvc.GetProfileByEmail(ctx, email)
	if err != nil {
		t.logger.Error("registerHandler ", err)
	}

	name := profile.Data.FullName
	if len(name) == 0 {
		name = "INDIKAY member"
		if profile.Data.Language == "vi" {
			name = "Thành viên INDIKAY"
		}
	}

	_, err = t.emailSvc.SendValidationEmail(ctx, email, profile.Data.Language, &extsvc.ValidationData{ValidationURL: activationUrl, Code: code, TelegramName: fmt.Sprintf("%s %s", data.Message.Chat.FirstName, data.Message.Chat.LastName), UserName: name})
	if err != nil {
		t.logger.Error("registerHandler ", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: data.Message.Chat.ID, Text: "You are registered successful.\nWe sent a code to your email,  please check your email"})
	if err != nil {
		t.logger.Error("registerHandler", err)
	}
}

func telegramDefaultHandler(ctx context.Context, b *bot.Bot, data *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: data.Message.Chat.ID,
		Text: `
To register signal alert from INDIKAY you need to type:
/register {email} with {email} is the email which register on indikay.com`,
	})
}
