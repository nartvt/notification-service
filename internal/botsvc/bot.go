package botsvc

import (
	"context"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewTelegramBot)

type IndicatorAlert struct {
	Id           int     `json:"id"`
	SignalId     int     `json:"signalId"`
	Indicator    string  `json:"indicatorName"`
	Source       string  `json:"source"`
	SignalStatus string  `json:"signalStatus"`
	IsUpdate     bool    `json:"isUpdate"`
	Symbol       string  `json:"symbol"`
	TimeFrame    string  `json:"timeFrame"`
	Signal       string  `json:"signal"`
	CloseTime    int64   `json:"closeTime"`
	TimeStamp    int64   `json:"timestamp"`
	ClosePrice   float32 `json:"closePrice"`
	Type         string  `json:"type"` // close / realtime
	SourceDown   string  `json:"sourceDown"`
	IsTest       bool    `json:"isTest"`
}

type BotService interface {
	Initial(ctx context.Context, data interface{}) error
	ReceiveMessage(ctx context.Context, msg interface{}) error
	SendIndicatorAlertMessage(ctx context.Context, userId string, alert *IndicatorAlert) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
