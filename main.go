package main

import (
	"Users/diggi/Documents/Go_tutorials/handlers"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file!")
	}
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USERNAME") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database!")
	}
	// db.Debug().AutoMigrate(&models.User{}, &models.EmailVerify{})
	handlers.DB = db
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "temp-redis.redis.cache.windows.net:6379",
		Password: "XcSYdtzZpmrBVDEfW7hwpissobhvXrvArAzCaPXmxUM=",
		DB:       0,
	})
	er := rdb.Ping(ctx)
	if er != nil {
		fmt.Println(er)
	}
	handlers.Rds = rdb
	handlers.ConnectSmtp()

	// register routes
	RegisterRoutes(app)

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
