package router

import (
	"github.com/Furkanturan8/motorbike-rental-backend-v2/config"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/handler"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/middleware"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/repository"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/service"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/email"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/monitoring"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/uptrace/bun"
	"time"
)

type Router struct {
	app *fiber.App
	db  *bun.DB
	cfg *config.Config
}

var prometheusEndpoint string
var prometheusEnabled bool

func NewRouter(db *bun.DB, cfg *config.Config) *Router {
	prometheusEnabled = cfg.MonitoringConfig.Prometheus.Enabled
	prometheusEndpoint = cfg.MonitoringConfig.Prometheus.Endpoint

	return &Router{
		app: fiber.New(),
		db:  db,
		cfg: cfg,
	}
}

func (r *Router) SetupRoutes() {
	// Prometheus'un topladığı metrikleri görüntülemek için /metrics endpoint'i
	if prometheusEnabled {
		r.app.Get(prometheusEndpoint, monitoring.MetricsHandler())
	}

	// Middleware'leri ekle
	r.app.Use(logger.New())
	r.app.Use(recover.New())
	r.app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:63342,http://localhost:3005,http://localhost:5173",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))

	// Rate limiting middleware'i ekle (30 sn de 10 istek olsun)
	r.app.Use(limiter.New(limiter.Config{
		Max:        10,               // Maksimum istek sayısı
		Expiration: 30 * time.Second, // Zaman aralığı
		KeyGenerator: func(c *fiber.Ctx) string {
			// /metrics endpoint'i için rate limiting'i devre dışı bırak
			if c.Path() == prometheusEndpoint {
				return "metrics_no_limit"
			}
			// Her route'u ayrı ayrı sınırla (örneğin: "/users", "/users/:id", "/auth/login")
			return c.IP() + ":" + c.Path()
		},
	}))

	// Prometheus Middleware ekleyelim
	r.app.Use(monitoring.PrometheusMiddleware())

	// API versiyonu
	api := r.app.Group("/api")
	v1 := api.Group("/v1")

	// Dış paketler emailPkg
	emailPkg := email.NewEmail(
		r.cfg.MailConfig.FromEmail,
		r.cfg.MailConfig.SMTPPassword,
		r.cfg.MailConfig.SMTPHost,
		r.cfg.MailConfig.SMTPPort,
	)

	// Repository'ler
	userRepo := repository.NewUserRepository(r.db)
	authRepo := repository.NewAuthRepository(r.db)
	rideRepo := repository.NewRideRepository(r.db)
	motorbikeRepo := repository.NewMotorbikeRepository(r.db)
	bluetoothRepo := repository.NewBluetoothConnectionRepository(r.db)

	// Service'ler
	authService := service.NewAuthService(authRepo, userRepo)
	userService := service.NewUserService(userRepo)
	rideService := service.NewRideService(rideRepo, motorbikeRepo)
	motorbikeService := service.NewMotorbikeService(motorbikeRepo)
	bluetoothService := service.NewBluetoothConnectionService(bluetoothRepo)

	// Handler'lar
	authHandler := handler.NewAuthHandler(authService, emailPkg)
	userHandler := handler.NewUserHandler(userService)
	rideHandler := handler.NewRideHandler(rideService)
	motorbikeHandler := handler.NewMotorbikeHandler(motorbikeService)
	bluetoothHandler := handler.NewBluetoothConnectionHandler(bluetoothService)

	// Auth routes
	auth := v1.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/forgot-password", authHandler.ForgotPassword)
	auth.Post("/reset-password", authHandler.ResetPassword)
	auth.Post("/logout", middleware.AuthMiddleware(), authHandler.Logout)

	// User routes - Base group
	users := v1.Group("/users")

	// Normal user routes (profil yönetimi)
	userProfile := users.Group("/me")
	userProfile.Use(middleware.AuthMiddleware()) // Sadece authentication gerekli
	userProfile.Get("/", userHandler.GetProfile)
	userProfile.Put("/", userHandler.UpdateProfile)

	// Admin only routes
	adminUsers := users.Group("/")
	adminUsers.Use(middleware.AuthMiddleware(), middleware.AdminOnly()) // Admin yetkisi gerekli
	adminUsers.Post("/", userHandler.Create)
	adminUsers.Get("/", userHandler.List)
	adminUsers.Get("/:id", userHandler.GetByID)
	adminUsers.Put("/:id", userHandler.Update)
	adminUsers.Delete("/:id", userHandler.Delete)

	// Ride routes
	rides := v1.Group("/rides")
	adminRides := rides.Group("/")

	adminRides.Use(middleware.AuthMiddleware(), middleware.AdminOnly()) // Admin yetkisi gerekli
	adminRides.Get("/", rideHandler.List)
	adminRides.Get("/user/:userID", rideHandler.ListRideByUserID)
	adminRides.Get("/bike/:motorbikeID", rideHandler.ListRideByMotorbikeID)
	adminRides.Get("/:id", rideHandler.GetByID)
	adminRides.Put("/:id", rideHandler.Update)
	adminRides.Delete("/:id", rideHandler.Delete)

	rides.Use(middleware.AuthMiddleware()) // Sadece authentication gerekli (normal kullanıcılar için)
	rides.Post("/", rideHandler.Create)
	rides.Get("/me", rideHandler.ListMyRides)
	rides.Put("/finish/:id", rideHandler.FinishRide)
	rides.Post("/photo/:id", rideHandler.AddRidePhoto)

	// Motorbike routes
	motorbike := v1.Group("/motorbike")
	adminMotorbike := motorbike.Group("/")

	adminMotorbike.Use(middleware.AuthMiddleware(), middleware.AdminOnly()) // Admin yetkisi gerekli
	adminMotorbike.Post("/", motorbikeHandler.Create)
	adminMotorbike.Put("/:id", motorbikeHandler.Update)
	adminMotorbike.Delete("/:id", motorbikeHandler.Delete)
	adminMotorbike.Get("/maintenance", motorbikeHandler.GetMaintenanceMotors)
	adminMotorbike.Get("/rented-motorbikes", motorbikeHandler.GetRentedMotors)
	adminMotorbike.Get("/motorbike-photos/:id", motorbikeHandler.GetPhotosByID)

	motorbike.Use(middleware.AuthMiddleware()) // Sadece authentication gerekli (normal kullanıcılar için)
	motorbike.Get("/", motorbikeHandler.List)
	motorbike.Get("/available", motorbikeHandler.GetAvailableMotors)
	motorbike.Get("/:id", motorbikeHandler.GetByID)

	// Bluetooth routes
	bluetooth := v1.Group("/bluetooth")
	adminBluetooth := bluetooth.Group("/")
	adminBluetooth.Use(middleware.AuthMiddleware(), middleware.AdminOnly()) // Admin yetkisi gerekli
	adminBluetooth.Post("/", bluetoothHandler.Create)
	adminBluetooth.Put("/:id", bluetoothHandler.Update)
	adminBluetooth.Delete("/:id", bluetoothHandler.Delete)
	adminBluetooth.Get("/", bluetoothHandler.List)
	adminBluetooth.Get("/:id", bluetoothHandler.GetByID)
	bluetooth.Use(middleware.AuthMiddleware()) // Sadece authentication gerekli (normal kullanıcılar için)

}

func (r *Router) GetApp() *fiber.App {
	return r.app
}
