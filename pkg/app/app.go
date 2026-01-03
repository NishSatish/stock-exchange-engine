package app

import (
	"context"
	"github.com/hibiken/asynq"
	dbPackage "stock-exchange-simulator/pkg/db"
	"stock-exchange-simulator/pkg/libs"
	"stock-exchange-simulator/pkg/libs/taskqueue/dto"
	"stock-exchange-simulator/pkg/nexus/service"
	tickbus "stock-exchange-simulator/pkg/tickbus/service"
	trademaster "stock-exchange-simulator/pkg/trademaster/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppServices struct {
	Nexus       service.INexusServiceInterface
	TradeMaster trademaster.ITradeMasterServiceInterface
	TickBus     tickbus.ITickBusServiceInterface
}

type IAppServicesInterface interface {
	AsynqEventHandlerMap(mux *asynq.ServeMux) *asynq.ServeMux
}

// Creating a NestJS/Angular type dependency injection
// constructor based injection, mainly i want to eliminate circular dependency

func NewAppServices(ctx context.Context, db *pgxpool.Pool, libsService *libs.LibsFactory) *AppServices {
	dbService := dbPackage.NewRepositoryFactory(db)

	tickBusService := tickbus.NewTickBusService(ctx, libsService)
	tickBusService.Start()
	tradeMasterService := trademaster.NewTradeMasterService(dbService, libsService)
	nexusService := service.NewNexusService(tradeMasterService, dbService, libsService)

	return &AppServices{
		Nexus:       nexusService,
		TradeMaster: tradeMasterService,
		TickBus:     tickBusService,
	}
}

func (this *AppServices) AsynqEventHandlerMap(mux *asynq.ServeMux) *asynq.ServeMux {
	// Just take the mux object, assign handlers and send it back
	// assigns each "topic" or message handler to a processor

	/*
	 * ORDER HANDLERS
	 */
	mux.HandleFunc(dto.EnqueueOrderPlacedTopic, this.TradeMaster.OrderProcessor)

	return mux
}
