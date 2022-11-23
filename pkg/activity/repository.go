package activity

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testfiber/pkg/entities"
)

type Repository interface {
	Create(entities.Activity) error
	Read(map[string]string) (*entities.Activity, error)
	Reads(map[string]string) (*[]entities.Activity, error)
	Update(entities.Activity) error
	Delete(map[string]string) error
}

type repository struct {
	conn *sql.DB
}

func NewRepository(connection *sql.DB) Repository {
	return &repository{
		conn: connection,
	}
}

func (r *repository) Create(activity entities.Activity) error {
	query := fmt.Sprintf("INSERT INTO activities (id, email, title, created_at, updated_at) VALUES (%d, '%s', '%s', '%s', '%s');", activity.ID, activity.Email, activity.Title, activity.CreateAt, activity.UpdateAt)
	_, err := r.conn.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Read(where map[string]string) (*entities.Activity, error) {
	id, err := strconv.Atoi(where["id"])
	if err != nil {
		return nil, fmt.Errorf("Activity with ID %d Not Found", id)
	}

	query := fmt.Sprintf("SELECT id, email, title, created_at, updated_at, deleted_at FROM activities WHERE id=%d;", id)
	row := r.conn.QueryRow(query)
	res := entities.Activity{}
	del := sql.NullString{}
	err = row.Scan(&res.ID, &res.Email, &res.Title, &res.CreateAt, &res.UpdateAt, &del)
	if err != nil {
		return nil, fmt.Errorf("Activity with ID %d Not Found", id)
	}

	res.DeleteAt = del.String

	return &res, nil
}

func (r *repository) Reads(where map[string]string) (*[]entities.Activity, error) {
	result := []entities.Activity{}
	query := "SELECT id, email, title, created_at, updated_at, deleted_at FROM activities;"
	rows, err := r.conn.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		res := entities.Activity{}
		del := sql.NullString{}
		err := rows.Scan(&res.ID, &res.Email, &res.Title, &res.CreateAt, &res.UpdateAt, &del)
		if err != nil {
			continue
		}
		res.DeleteAt = del.String
		result = append(result, res)
	}

	return &result, nil
}

func (r *repository) Update(activity entities.Activity) error {
	data := []string{}
	v := reflect.ValueOf(activity)
	t := reflect.TypeOf(activity)

	for i := 1; i < t.NumField(); i++ {
		var d string
		if fmt.Sprint(v.Field(i)) == "" {
			continue
		}

		if v.Field(i).Kind() == reflect.String {
			d = fmt.Sprintf("%s='%v'", t.Field(i).Tag.Get("json"), v.Field(i))
		} else {
			d = fmt.Sprintf("%s=%v", t.Field(i).Tag.Get("json"), v.Field(i))
		}
		data = append(data, d)
	}

	query := fmt.Sprintf("UPDATE activities SET %s WHERE id=%s;", strings.Join(data, ", "), fmt.Sprint(activity.ID))
	_, err := r.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("Activity with ID %d Not Found", activity.ID)
	}
	return nil
}

func (r *repository) Delete(where map[string]string) error {
	id, err := strconv.Atoi(where["id"])
	if err != nil {
		return fmt.Errorf("Activity with ID %d Not Found", id)
	}
	query := fmt.Sprintf("UPDATE activities SET deleted_at=CURRENT_TIMESTAMP WHERE id=%d", id)
	row, err := r.conn.Exec(query)
	i, _ := row.RowsAffected()
	if i == 0 {
		return fmt.Errorf("Activity with ID %d Not Found", id)
	}
	if err != nil {
		return err
	}
	return nil
}
