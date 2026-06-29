package main

import (
	"log"
	"os"
	"spotsync-api/config"
	"spotsync-api/handler"
	"spotsync-api/models"
	"spotsync-api/repository"
	"spotsync-api/service"
	"spotsync-api/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize Connection
	db := config.ConnectDB()

	// Run GORM Automigrations
	err := db.AutoMigrate(&models.User{}, &models.ParkingZone{}, &models.Reservation{})
	if err != nil {
		log.Fatalf("AutoMigration Failed: %v", err)
	}

	e := echo.New()

	// Centralized Global Error Handler Registration
	e.HTTPErrorHandler = utils.CustomHTTPErrorHandler

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Clean Architecture Manual Dependency Injection Setup
	userRepo := repository.NewUserRepository(db)
	zoneRepo := repository.NewZoneRepository(db)
	resRepo := repository.NewReservationRepository(db)

	authService := service.NewAuthService(userRepo)
	zoneService := service.NewZoneService(zoneRepo)
	resService := service.NewReservationService(resRepo)

	authHandler := handler.NewAuthHandler(authService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	resHandler := handler.NewReservationHandler(resService)
	homeHandler := handler.NewHomeHandler()

	// Home Route
	e.GET("/", homeHandler.Home)

	// API Routing Setup (Matching Blueprint Spec Exactly)
	v1 := e.Group("/api/v1")

	// Authentication Module
	v1.POST("/auth/register", authHandler.Register)
	v1.POST("/auth/login", authHandler.Login)

	// Parking Zones Module
	v1.POST("/zones", zoneHandler.Create, handler.JWTMiddleware, handler.RoleBlock("admin"))
	v1.GET("/zones", zoneHandler.GetAll)
	v1.GET("/zones/:id", zoneHandler.GetOne)

	// Reservations Module
	v1.POST("/reservations", resHandler.Create, handler.JWTMiddleware)
	v1.GET("/reservations/my-reservations", resHandler.GetMy, handler.JWTMiddleware)
	v1.DELETE("/reservations/:id", resHandler.Cancel, handler.JWTMiddleware)
	v1.GET("/reservations", resHandler.GetAll, handler.JWTMiddleware, handler.RoleBlock("admin"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
