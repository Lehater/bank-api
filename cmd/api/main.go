package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"bank-api/config"
	"bank-api/handlers"
	"bank-api/middleware"
	"bank-api/repositories"
	"bank-api/services"
	"bank-api/scheduler"

	"github.com/gorilla/mux"
)


func main() {

    // Загружаем переменные окружения из .env файла
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

	// Подключаемся к базе данных.
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Создаем репозитории.
	userRepo := repositories.NewUserRepository(db)
	accountRepo := repositories.NewAccountRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	creditRepo := repositories.NewCreditRepository(db)
	paymentScheduleRepo := repositories.NewPaymentScheduleRepository(db)
	cardRepo := repositories.NewCardRepository(db) // должен быть реализован
	// Создаем сервисы.
	jwtSecret := os.Getenv("JWT_SECRET")
	userService := services.NewUserService(userRepo, jwtSecret)
    accountService := services.NewAccountService(accountRepo, db)
	creditService := services.NewCreditService(creditRepo, paymentScheduleRepo)
	cardService := services.NewCardService(cardRepo)
    analyticsService := services.NewAnalyticsService(
        transactionRepo,
        accountRepo,
        creditRepo,
        paymentScheduleRepo,
    )	

	// Создаем обработчики.
	userHandler := handlers.NewUserHandler(userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	creditHandler := handlers.NewCreditHandler(creditService)
	cardHandler := handlers.NewCardHandler(cardService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	// Настраиваем маршруты.
	r := mux.NewRouter()
	// Публичные маршруты.
	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	// Защищенные маршруты.
	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.RecoveryMiddleware(nil)) // можно передать логгер
	authRouter.Use(middleware.LoggingMiddleware(nil))
	authRouter.Use(middleware.AuthMiddleware(jwtSecret))
	authRouter.HandleFunc("/credits", creditHandler.ApplyForCredit).Methods("POST")
	authRouter.HandleFunc("/cards", cardHandler.CreateCard).Methods("POST")
	authRouter.HandleFunc("/cards/{id}", cardHandler.GetCard).Methods("GET")
	
    // endpoint для переводов
	authRouter.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	authRouter.HandleFunc("/transfer", accountHandler.Transfer).Methods("POST")
	// маршруты аналитики.
	authRouter.HandleFunc("/analytics", analyticsHandler.GetAnalytics).Methods("GET")
	authRouter.HandleFunc("/accounts/{accountId}/predict", analyticsHandler.PredictBalance).Methods("GET")
	 // endpoint для графика платежей по кредиту
	 authRouter.HandleFunc("/credits/{creditId}/schedule", creditHandler.GetSchedule).Methods("GET")
	// Запуск шедулера (если используется).
	paymentScheduler := scheduler.NewPaymentScheduler(creditService, accountService)
	paymentScheduler.Start()

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
