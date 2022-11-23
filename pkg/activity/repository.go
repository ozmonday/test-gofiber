package activity

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testfiber/pkg/entities"
)

type Repository interface {
	Create(context.Context, entities.Activity) (int64, error)
	Read(context.Context, map[string]string) (*entities.Activity, error)
	Reads(context.Context, map[string]string) (*[]entities.Activity, error)
	Update(context.Context, entities.Activity) error
	Delete(context.Context, map[string]string) error
}

type repository struct {
	conn *sql.DB
}

func NewRepository(connection *sql.DB) Repository {
	return &repository{
		conn: connection,
	}
}

func (r *repository) Create(ctx context.Context, activity entities.Activity) (int64, error) {
	id := sql.NullInt64{}

	query := fmt.Sprintf("INSERT INTO activities (email, title, created_at, updated_at) VALUES ('%s', '%s', '%s', '%s');", activity.Email, activity.Title, activity.CreateAt, activity.UpdateAt)
	_, err := r.conn.ExecContext(ctx, query)
	if err != nil {
		return -1, err
	}

	// result := "SELECT LAST_INSERT_ID();"
	// row := r.conn.QueryRowContext(ctx, result)
	// if err := row.Scan(&id); err != nil {
	// 	return -1, err
	// }

	return id.Int64, nil
}

func (r *repository) Read(ctx context.Context, where map[string]string) (*entities.Activity, error) {
	id, err := strconv.Atoi(where["id"])
	if err != nil {
		return nil, fmt.Errorf("Activity with ID %d Not Found", id)
	}

	query := fmt.Sprintf("SELECT id, email, title, created_at, updated_at, deleted_at FROM activities WHERE id=%d;", id)
	row := r.conn.QueryRowContext(ctx, query)
	res := entities.Activity{}
	del := sql.NullString{}
	err = row.Scan(&res.ID, &res.Email, &res.Title, &res.CreateAt, &res.UpdateAt, &del)
	if err != nil {
		return nil, fmt.Errorf("Activity with ID %d Not Found", id)
	}

	res.DeleteAt = del.String

	return &res, nil
}

func (r *repository) Reads(ctx context.Context, where map[string]string) (*[]entities.Activity, error) {
	result := []entities.Activity{}
	query := "SELECT id, email, title, created_at, updated_at, deleted_at FROM activities;"
	rows, err := r.conn.QueryContext(ctx, query)
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

func (r *repository) Update(ctx context.Context, activity entities.Activity) error {
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
	_, err := r.conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("Activity with ID %d Not Found", activity.ID)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, where map[string]string) error {
	id, err := strconv.Atoi(where["id"])
	if err != nil {
		return fmt.Errorf("Activity with ID %d Not Found", id)
	}
	query := fmt.Sprintf("UPDATE activities SET deleted_at=CURRENT_TIMESTAMP WHERE id=%d", id)
	row, err := r.conn.ExecContext(ctx, query)
	i, _ := row.RowsAffected()
	if i == 0 {
		return fmt.Errorf("Activity with ID %d Not Found", id)
	}
	if err != nil {
		return err
	}
	return nil
}
