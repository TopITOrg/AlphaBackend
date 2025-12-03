package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sport_platform/application/controllers"
	"sport_platform/application/models/claims"
	"sport_platform/internal/configuration"
	"sport_platform/internal/env_loader"
	"sport_platform/internal/jwt"
	"sport_platform/internal/middleware"
	"sport_platform/internal/password"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Application struct {
	engine        *gin.Engine
	wrapper       *service_wrapper.Wrapper
	configuration configuration.IConfiguration
}

func (appl *Application) GetEnv() error {
	appEnvLoader := env_loader.CreateLoaderFromEnv()

	var dbConfig db.Config
	if err := appEnvLoader.LoadDataIntoStruct(&dbConfig); err != nil {
		return err
	}

	var passwordConfig password.PasswordConfig
	if err := appEnvLoader.LoadDataIntoStruct(&passwordConfig); err != nil {
		return err
	}

	var jwtConfig jwt.JwtConfig
	if err := appEnvLoader.LoadDataIntoStruct(&jwtConfig); err != nil {
		return err
	}

	appl.configuration.
		AddConfiguration(&dbConfig).
		AddConfiguration(&passwordConfig).
		AddConfiguration(&jwtConfig)

	return nil
}

func (appl *Application) ConstructClients() error {
	dbConfig, dbConfigGetterError := appl.configuration.Get(&db.Config{})
	if dbConfigGetterError != nil {
		return dbConfigGetterError
	}

	dbClient, dbConnectionError := db.CreateConnection(dbConfig.(*db.Config), context.Background())
	if dbConnectionError != nil {
		return dbConnectionError
	}
	appl.wrapper.Db = dbClient

	passwordConfig, passwordConfigGetterError := appl.configuration.Get(&password.PasswordConfig{})
	if passwordConfigGetterError != nil {
		return passwordConfigGetterError
	}
	appl.wrapper.PasswordHandler = password.CreateHandler(passwordConfig.(*password.PasswordConfig))

	jwtConfig, jwtConfigGetterError := appl.configuration.Get(&jwt.JwtConfig{})
	if jwtConfigGetterError != nil {
		return jwtConfigGetterError
	}
	appl.wrapper.JwtHandler = jwt.CreateHandler[claims.UserClaims](jwtConfig.(*jwt.JwtConfig))

	return nil
}

func (appl *Application) Configure(engine *gin.Engine) error {
	envGetterError := appl.GetEnv()
	if envGetterError != nil {
		return envGetterError
	}

	clientsConstructionError := appl.ConstructClients()
	if clientsConstructionError != nil {
		return clientsConstructionError
	}

	engine.Use(middleware.AuthMiddleware(appl.wrapper))

	controllers.UserController(engine, appl.wrapper)
	controllers.ClubController(engine, appl.wrapper)
	controllers.WorkoutController(engine, appl.wrapper)

	controllers.ClubController(engine, appl.wrapper)

	controllers.WorkoutController(engine, appl.wrapper)

	controllers.ClubJoinRequestController(engine, appl.wrapper)
	return nil
}

func (appl *Application) Run() {
	engine := gin.New()
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: engine,
	}

	if err := appl.Configure(engine); err != nil {
		panic(err)
	}

	go func() {
		fmt.Println("Server started on http://locahost:8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Received a shutdown call. Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := appl.wrapper.Close(); err != nil {
		panic(err)
	}

	fmt.Println("Closed all connections. Shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}

func CreateApplication() *Application {
	engine := gin.New()

	return &Application{
		engine:        engine,
		wrapper:       &service_wrapper.Wrapper{},
		configuration: configuration.CreateConfiguration(),
	}
}
