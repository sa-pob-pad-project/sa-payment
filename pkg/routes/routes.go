package routes

import (
	// "user-service/pkg/context"
	// "user-service/pkg/dto"
	_ "payment-service/docs"
	"payment-service/pkg/handlers"
	"payment-service/pkg/jwt"
	"payment-service/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, paymentHandler *handlers.PaymentHandler, jwtSvc *jwt.JwtService) {

	api := app.Group("/api")

	// Payment Routes
	payment := api.Group("/payment")
	payment.Get("/swagger/*", swagger.HandlerDefault)

	paymentV1 := payment.Group("/v1")
	paymentV1.Use(middleware.JwtMiddleware(jwtSvc))
	// payment
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
}
