package utility

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBContext struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (db *DBContext) Connect() (*sql.DB, error) {
	data := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.User, db.Password, db.Host, db.Port, db.DBName)
	conn, err := sql.Open("mysql", data)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func GetTime() string {
	t := time.Now()
	return fmt.Sprintf("%d-%d-%d %d:%d:%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func Migration(db *sql.DB) error {
	tabel_activity := `CREATE TABLE IF NOT EXISTS activities (
		id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
		email varchar(255) NOT NULL,
		title varchar(255) NOT NULL,
		created_at datetime NOT NULL,
		updated_at datetime NOT NULL,
		deleted_at datetime DEFAULT NULL
	) ENGINE=InnoDB;`
	_, err := db.Exec(tabel_activity)
	if err != nil {
		return err
	}

	tabel_todos := `CREATE TABLE IF NOT EXISTS todos (
		id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
		title varchar(255) NOT NULL,
		activity_group_id int NOT NULL,
		is_active bool NOT NULL,
		priority varchar(55) NOT NULL,
		created_at datetime NOT NULL,
		updated_at datetime NOT NULL,
		deleted_at datetime DEFAULT NULL
	) ENGINE=InnoDB;`

	_, err = db.Exec(tabel_todos)
	if err != nil {
		return err
	}

	fmt.Println("Minration Success")
	return nil
}