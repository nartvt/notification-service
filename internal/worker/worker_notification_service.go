package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/indikay/notification-service/internal/biz"
	"github.com/indikay/notification-service/internal/botsvc"
	"github.com/indikay/notification-service/internal/conf"
	"github.com/indikay/notification-service/internal/constants"
	"github.com/indikay/notification-service/internal/service"
	"github.com/nats-io/nats.go"

	"github.com/go-kratos/kratos/v2/log"
)

type NotificationWorkerService struct {
	log                 *log.Helper
	nc                  *nats.Conn
	notificationService *service.NotificationService
	telegramBot         botsvc.BotService
}

type EventRealtime struct {
	EventName string      `json:"event_name"`
	EventType []string    `json:"event_type,omitempty"`
	UserID    string      `json:"user_id"`
	EventData interface{} `json:"data"`
}

func NewWorkerNotification(cfg *conf.Data, notificationService *service.NotificationService, telegramBot botsvc.BotService) *NotificationWorkerService {
	nc, err := nats.Connect(cfg.Nats.NatsHost)
	if err != nil {
		panic(err)
	}

	return &NotificationWorkerService{
		log:                 log.NewHelper(log.DefaultLogger),
		nc:                  nc,
		notificationService: notificationService,
		telegramBot:         telegramBot,
	}
}

func (w *NotificationWorkerService) Run() {
	// Initialize NATS connection
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Handle OS signals to gracefully shutdown
	handleSignals(cancel, &wg)
	// Subscribe to Nats channel
	_, err := w.nc.QueueSubscribe(constants.NATS_EVENT_NAME_NOTIFICATION, "", w.consumerMsg)
	if err != nil {
		w.log.Error("Subscribe failed", "err", err)
		panic(err)
	}
}

// Handle OS signals for graceful shutdown
func handleSignals(cancel context.CancelFunc, wg *sync.WaitGroup) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.DefaultLogger.Log(log.LevelInfo, "Received signal: %s", sig.String())
		cancel()
		wg.Wait() // Wait for goroutines to finish before exiting
		os.Exit(0)
	}()
}

// Consumer message
func (w *NotificationWorkerService) consumerMsg(msg *nats.Msg) {
	fmt.Println("MESSAGE: ", string(msg.Data))
	if msg == nil {
		return
	}

	if len(msg.Data) == 0 {
		return
	}

	event := &EventRealtime{}
	err := json.Unmarshal(msg.Data, event)
	if err != nil {
		log.Errorf("[CONSUME MESSAGE] %s", err.Error())
		return
	}
	// handle
	switch event.EventName {
	case constants.NATS_EVENT_NAME_NOTIFICATION,
		constants.NATS_EVENT_NAME_REWARDBACK_COMMISSION,
		constants.NATS_EVENT_NAME_CASHBACK_COMMISSION,
		constants.NATS_EVENT_NAME_REFERRAL_COMMISSION:
		// Log err handled inside function
		err := w.processUserNotification(msg)
		if err != nil {
			log.Error("FAIL_CONSUME_MESSAGE: ", err.Error())
		}
	case constants.NATS_EVENT_NAME_INDI_ALERT_REALTIME:
		err := w.processIndicatorAlert(event, msg.Data)
		if err != nil {
			log.Error("FAIL_CONSUME_ALERT: ", err.Error())
		}
	}
}

func (w *NotificationWorkerService) processUserNotification(msg *nats.Msg) error {
	log.Infof("MESSAGE RECEIVE: %s", string(msg.Data))
	var notificationData biz.Notification
	if err := json.Unmarshal(msg.Data, &notificationData); err != nil {
		log.Error("UNMARSHAL_EVENT_DATA_ERROR: ", err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := w.notificationService.SaveNotification(ctx, &notificationData)
	if err != nil {
		log.Error("SAVE NOTI ERROR: ", err.Error())
		return err
	}

	err = w.PublishRealtime(&notificationData)
	if err != nil {
		return err
	}
	return nil
}

func (w *NotificationWorkerService) PublishRealtime(data *biz.Notification) error {
	eventRealtime := EventRealtime{
		UserID:    data.UserID,
		EventName: constants.NATS_EVENT_NAME_NOTIFICATION,
		EventData: data,
	}

	publishData, err := json.Marshal(eventRealtime)
	fmt.Println("PUBLISH REALTIME: ", string(publishData))
	if err != nil {
		return err
	}

	err = w.nc.Publish(constants.NATS_SUBJECT_NAME_REALTIME, publishData)
	if err != nil {
		return err
	}
	return nil
}

func (w *NotificationWorkerService) processIndicatorAlert(event *EventRealtime, rawData []byte) error {
	if slices.Contains(event.EventType, constants.EVENT_TYPE_SOCKET) {
		w.nc.Publish(constants.NATS_SUBJECT_NAME_REALTIME, rawData)
	}

	if slices.Contains(event.EventType, constants.EVENT_TYPE_TELEGRAM) {
		dataMsg, err := json.Marshal(event.EventData)
		if err != nil {
			return err
		}
		alert := &botsvc.IndicatorAlert{}
		json.Unmarshal(dataMsg, alert)
		w.telegramBot.SendIndicatorAlertMessage(context.Background(), event.UserID, alert)
	}

	return nil
}
