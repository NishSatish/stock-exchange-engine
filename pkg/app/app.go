package app

import (
	dbPackage "stock-exchange-simulator/pkg/db"
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

// Creating a NestJS/Angular type dependency injection
// constructor based injection, mainly i want to eliminate circular dependency

func NewAppServices(db *pgxpool.Pool) *AppServices {
	dbService := dbPackage.NewRepositoryFactory(db)

	tickBusService := tickbus.NewTickBusService(db)
	tradeMasterService := trademaster.NewTradeMasterService(db)
	nexusService := service.NewNexusService(tradeMasterService, dbService)

	return &AppServices{
		Nexus:       nexusService,
		TradeMaster: tradeMasterService,
		TickBus:     tickBusService,
	}
}
