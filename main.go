package main

import (
	"log"
	"os"
	"testfiber/api/routes"
	"testfiber/storage/activitiy"
	"testfiber/storage/todo"
	"testfiber/utility"

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

	conn, err := mysql.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	if err := utility.Migration(conn); err != nil {
		log.Fatalln(err)
	}

	activityRepo := activitiy.NewRepository(conn)
	activityService := activitiy.NewService(activityRepo)

	todoRepo := todo.NewRepository(conn)
	todoService := todo.NewService(todoRepo)

	routes.ActivityRouter(app, activityService)
	routes.TodoRouter(app, todoService)

	app.Listen(os.Getenv("PORT"))
}
