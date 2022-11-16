package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testfiber/api/routes"
	"testfiber/storage/activitiy"
	"testfiber/storage/todo"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func connect(host string, port string, user string, password string, dbname string) *sql.DB {
	data := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)
	conn, err := sql.Open("mysql", data)
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(cors.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	conn := connect(os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DBNAME"))

	activityRepo := activitiy.NewRepository(conn)
	activityService := activitiy.NewService(activityRepo)

	todoRepo := todo.NewRepository(conn)
	todoService := todo.NewService(todoRepo)

	routes.ActivityRouter(app, activityService)
	routes.TodoRouter(app, todoService)

	app.Listen(":3000")
}
