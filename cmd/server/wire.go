//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	coreConf "github.com/indikay/go-core/conf"
	"github.com/indikay/go-core/server"
	coreService "github.com/indikay/go-core/service"
	"github.com/indikay/notification-service/internal/biz"
	"github.com/indikay/notification-service/internal/botsvc"
	"github.com/indikay/notification-service/internal/conf"
	data "github.com/indikay/notification-service/internal/data"
	"github.com/indikay/notification-service/internal/extsvc"
	"github.com/indikay/notification-service/internal/service"
	"github.com/indikay/notification-service/internal/worker"
)

// initApp init kratos application.
func initApp(*coreConf.Server, *conf.Data, log.Logger) (coreService.Service, func(), error) {
	panic(wire.Build(data.ProviderSet, extsvc.ProviderSet, biz.ProviderSet, service.ProviderSet, server.ProviderSet, botsvc.ProviderSet, worker.ProviderSet, initService))
}
