package initialize

import (
	"github.com/ntquang/ecommerce/global"
	"github.com/ntquang/ecommerce/internal/database"
	"github.com/ntquang/ecommerce/internal/services"
	"github.com/ntquang/ecommerce/internal/services/imple"
)

func InitServiceInterface() {
	queries := database.New(global.Pdbc)
	// Init Interface here

	// event service interface
	services.InitEvent(imple.NewEventImpl(queries))
	// menu function service interface
	services.InitMenuFunction(imple.NewMenuFunctionImpl(queries))
	// contact message service interface
	services.InitContactMessage(imple.NewContactMessage(queries))
}
