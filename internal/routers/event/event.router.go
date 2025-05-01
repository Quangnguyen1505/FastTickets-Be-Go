package event

import (
	"github.com/gin-gonic/gin"
	"github.com/ntquang/ecommerce/global"
	event "github.com/ntquang/ecommerce/internal/controller/event"
	"github.com/ntquang/ecommerce/internal/middlewares"
	"github.com/ntquang/ecommerce/internal/middlewares/grpc"
	pb "github.com/ntquang/ecommerce/proto"
)

type EventRouter struct{}

func (mktC *EventRouter) InitEventRouter(Router *gin.RouterGroup) {
	publicEventRouter := Router.Group("/events")
	{
		publicEventRouter.GET("", event.Event.GetAllEvents)
		publicEventRouter.GET("/:id", event.Event.GetEventById)
	}

	uploadClient := pb.NewUploadServiceClient(global.Grpc)
	privateEventRouter := Router.Group("/events")
	privateEventRouter.Use(middlewares.Authentication())
	privateEventRouter.Use(middlewares.CheckPermission())
	{
		privateEventRouter.PUT(
			"/:id",
			grpc.UploadImageMiddleware(uploadClient),
			event.Event.EditEvent,
		)
		privateEventRouter.POST(
			"",
			grpc.UploadImageMiddleware(uploadClient),
			event.Event.NewEvent,
		)
		privateEventRouter.DELETE("/:id", event.Event.DeleteEvent)
	}
}
