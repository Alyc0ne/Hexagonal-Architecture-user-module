package main

import (
	"fmt"
	"log"
	"os"

	"github.com/LordMoMA/Hexagonal-Architecture/internal/adapters/handler"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/adapters/repository"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/domain"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/services"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/logger"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var (
	userService services.UserUsecaseService
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	con := "%s:%s@/%s?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", fmt.Sprintf(con, user, password, dbname))
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	logger.SetupLogger()

	db.AutoMigrate(&domain.User{}, &domain.ForgetPassword{})

	userRepo := repository.NewUserRepository(db)
	userService = services.NewUserUsecase(userRepo)

	InitRoutes()
}

func InitRoutes() {
	router := gin.Default()

	pprof.Register(router)

	v1 := router.Group("/v1")

	userHandler := handler.NewUserHttpHandler(userService)

	v1.POST("/login", userHandler.LoginUser)
	v1.POST("/forgetpassword", userHandler.ForgetPassword)
	v1.PUT("/reset-password", userHandler.ResetPassword)
	v1.POST("/users", userHandler.CreateUser)

	err := router.Run(":4242")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
