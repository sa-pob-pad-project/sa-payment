package routes

import (
	// "user-service/pkg/context"
	// "user-service/pkg/dto"
	_ "order-service/docs"
	"order-service/pkg/handlers"
	"order-service/pkg/jwt"
	"order-service/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, paymentHandler *handlers.PaymentHandler, jwtSvc *jwt.JwtService) {

	api := app.Group("/api")
	api.Get("/swagger/*", swagger.HandlerDefault)

	// Payment Routes
	payment := api.Group("/payment")
	paymentV1 := payment.Group("/v1")
	paymentV1.Use(middleware.JwtMiddleware(jwtSvc))

}
