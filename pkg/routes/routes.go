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
	paymentV1.Post("/", paymentHandler.CreatePayment)
	paymentV1.Get("/", paymentHandler.GetAllPayments)
	paymentV1.Get("/:id", paymentHandler.GetPaymentByID)
	// payment info routes
	paymentV1.Post("/info", paymentHandler.CreatePaymentInfo)
	paymentV1.Put("/info", paymentHandler.UpdatePaymentInfo)
	paymentV1.Delete("/info", paymentHandler.DeletePaymentInfo)
	paymentV1.Get("/info", paymentHandler.GetAllPaymentInfos)
	paymentV1.Get("/info/:id", paymentHandler.GetPaymentInfo)
	// payment attempt routes
	paymentV1.Post("/attempt", paymentHandler.CreatePaymentAttempt)
	paymentV1.Get("/attempt/:id", paymentHandler.GetPaymentAttempt)
	paymentV1.Patch("/attempt", paymentHandler.UpdatePaymentAttempt)
	// payment
	paymentV1.Post("/", paymentHandler.CreatePayment)
	paymentV1.Get("/", paymentHandler.GetAllPayments)
	paymentV1.Get("/:id", paymentHandler.GetPaymentByID)
}
