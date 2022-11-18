package main

import (
	"log"
	"os"
	"testfiber/api/routes"
	"testfiber/storage/activity"
	"testfiber/storage/todo"
	"testfiber/utility"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	mysql := utility.DBContext{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		DBName:   os.Getenv("MYSQL_DBNAME"),
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	conn, err := mysql.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	if err := utility.Migration(conn); err != nil {
		log.Fatalln(err)
	}

	activityRepo := activity.NewRepository(conn)
	activitySess := activity.NewSession(rdb)
	activityService := activity.NewService(activityRepo, activitySess)

	todoRepo := todo.NewRepository(conn)
	todoSess := todo.NewSession(rdb)
	todoService := todo.NewService(todoRepo, todoSess)

	routes.ActivityRouter(app, activityService)
	routes.TodoRouter(app, todoService)

	app.Listen(os.Getenv("PORT"))
}
