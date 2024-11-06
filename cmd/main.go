package main

import (
	"fmt"
	"log"
	"os"

	"github.com/LordMoMA/Hexagonal-Architecture/internal/adapters/grpccon"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/adapters/handler"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/adapters/repository"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/domain"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/services"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/logger"
	"github.com/casbin/casbin"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/middleware"
	userPb "github.com/LordMoMA/Hexagonal-Architecture/internal/core/proto/user"
)

var (
	userRepo    repository.UserRepositoryService
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

	jwtSecret := os.Getenv("JWT_SECRET")

	userRepo = repository.NewUserRepository(db)
	userService = services.NewUserUsecase(jwtSecret, userRepo)

	db.AutoMigrate(&domain.User{}, &domain.ForgetPassword{})
	countUser, _ := userRepo.CountUserByEmail("admin@gmail.com")
	if countUser == 0 {
		if db.HasTable(&domain.User{}) {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password-admin"), bcrypt.DefaultCost)

			db.Create(&domain.User{
				ID:       "32c2b3ed-d9d5-11ee-8f2e-523b287e1657",
				Email:    "admin@gmail.com",
				Password: string(hashedPassword),
				Role:     "admin",
			})
		}
	}

	InitRoutes()
}

func InitRoutes() {
	router := gin.Default()

	pprof.Register(router)

	enforcer := casbin.NewEnforcer("../model.conf", "../policy.csv")
	casbinMiddleware := middleware.NewCasbinMiddleware(enforcer, userRepo)

	userHttpHandler := handler.NewUserHttpHandler(userService)
	userGrpcHandler := handler.NewUserGrpcHandler(userService)

	go func() {
		host := "localhost:4444"
		grpcServer, lis := grpccon.NewGrpcServer(host)
		userPb.RegisterUserGrpcServiceServer(grpcServer, userGrpcHandler)

		log.Printf("User gRPC server listening on %s", host)
		grpcServer.Serve(lis)
	}()

	v1 := router.Group("/v1")

	v1.GET("/users", casbinMiddleware.CasbinMiddleware(), userHttpHandler.ReadUser)
	v1.POST("/login", userHttpHandler.LoginUser)
	v1.POST("/forgetpassword", userHttpHandler.ForgetPassword)
	v1.PUT("/reset-password", userHttpHandler.ResetPassword)
	v1.POST("/users", userHttpHandler.CreateUser)

	err := router.Run(":4242")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
