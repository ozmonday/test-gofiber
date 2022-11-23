package todo

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
	Create(context.Context, entities.Todo) (int64, error)
	Read(context.Context, map[string]string) (*entities.Todo, error)
	Reads(context.Context, map[string]string) (*[]entities.Todo, error)
	Update(context.Context, entities.Todo) error
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

func (r *repository) Create(ctx context.Context, todo entities.Todo) (int64, error) {

	query := fmt.Sprintf("INSERT INTO todos (title, activity_group_id, priority, is_active, created_at, updated_at) VALUES ('%s', %d, '%s', %v, '%s', '%s');", todo.Title, todo.ActivityID, todo.Priority, todo.IsActive, todo.CreateAt, todo.UpdateAt)
	id := sql.NullInt64{}
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

func (r *repository) Read(ctx context.Context, where map[string]string) (*entities.Todo, error) {
	id, err := strconv.Atoi(where["id"])
	if err != nil {
		return nil, fmt.Errorf("Todo with ID %d Not Found", id)
	}

	query := fmt.Sprintf("SELECT id, title, is_active, activity_group_id, priority, created_at, updated_at, deleted_at FROM todos WHERE id=%d;", id)
	row := r.conn.QueryRowContext(ctx, query)
	res := entities.Todo{}
	del := sql.NullString{}
	err = row.Scan(&res.ID, &res.Title, &res.IsActive, &res.ActivityID, &res.Priority, &res.CreateAt, &res.UpdateAt, &del)
	if err != nil {
		return nil, fmt.Errorf("Todo with ID %d Not Found", id)
	}
	return &res, nil
}

func (r *repository) Reads(ctx context.Context, where map[string]string) (*[]entities.Todo, error) {
	result := []entities.Todo{}
	query := fmt.Sprintf("SELECT id, title, is_active, activity_group_id, priority, created_at, updated_at, deleted_at FROM todos WHERE activity_group_id=%s;", where["activity_group_id"])
	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		r := entities.Todo{}
		del := sql.NullString{}
		err := rows.Scan(&r.ID, &r.Title, &r.IsActive, &r.ActivityID, &r.Priority, &r.CreateAt, &r.UpdateAt, &del)
		if err != nil {
			continue
		}
		r.DeleteAt = del.String
		result = append(result, r)
	}

	return &result, nil
}

func (r *repository) Update(ctx context.Context, todo entities.Todo) error {
	data := []string{}
	v := reflect.ValueOf(todo)
	t := reflect.TypeOf(todo)

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

	query := fmt.Sprintf("UPDATE todos SET %s WHERE id=%d;", strings.Join(data, ", "), todo.ID)
	_, err := r.conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("Todo with ID %d Not Found", todo.ID)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, where map[string]string) error {
	id, err := strconv.Atoi(where["id"])
	if err != nil {
		return fmt.Errorf("Todo with ID %d Not Found", id)
	}
	query := fmt.Sprintf("UPDATE todos SET deleted_at=CURRENT_TIMESTAMP WHERE id=%d", id)
	row, err := r.conn.ExecContext(ctx, query)
	i, _ := row.RowsAffected()
	if i == 0 {
		return fmt.Errorf("Todo with ID %d Not Found", id)
	}
	if err != nil {
		return err
	}
	return nil
}
