package todo

import (
	"database/sql"
	"fmt"
	"strconv"
	"testfiber/pkg/entities"
	"testfiber/pkg/utility"
)

type Repository interface {
	Create(*entities.Todo) error
	Read(map[string]string) (*entities.Todo, error)
	Reads(map[string]string) (*[]entities.Todo, error)
	Update(*entities.Todo, string) error
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

func (r *repository) Create(todo *entities.Todo) error {
	t := utility.GetTime()
	todo.UpdateAt = t
	todo.CreateAt = t

	if todo.Priority == "" {
		todo.Priority = "very-high"
	}

	if !todo.IsActive {
		todo.IsActive = true
	}
	query := fmt.Sprintf("INSERT INTO todos (title, activity_group_id, priority, is_active, created_at, updated_at) VALUES ('%s', %d, '%s', %v, '%s', '%s');", todo.Title, todo.ActivityID, todo.Priority, todo.IsActive, todo.CreateAt, todo.UpdateAt)
	_, err := r.conn.Exec(query)
	if err != nil {
		return err
	}

	result := `SELECT LAST_INSERT_ID();`
	row := r.conn.QueryRow(result)
	err = row.Scan(&todo.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Read(where map[string]string) (*entities.Todo, error) {
	id, err := strconv.Atoi(where["id"])
	if err != nil {
		return nil, fmt.Errorf("Todo with ID %d Not Found", id)
	}

	query := fmt.Sprintf("SELECT id, title, is_active, activity_group_id, priority, created_at, updated_at, deleted_at FROM todos WHERE id=%d;", id)
	row := r.conn.QueryRow(query)
	res := entities.Todo{}
	del := sql.NullString{}
	err = row.Scan(&res.ID, &res.Title, &res.IsActive, &res.ActivityID, &res.Priority, &res.CreateAt, &res.UpdateAt, &del)
	if err != nil {
		return nil, fmt.Errorf("Todo with ID %d Not Found", id)
	}
	return &res, nil
}

func (r *repository) Reads(where map[string]string) (*[]entities.Todo, error) {
	result := []entities.Todo{}
	query := fmt.Sprintf("SELECT id, title, is_active, activity_group_id, priority, created_at, updated_at, deleted_at FROM todos WHERE activity_group_id=%s;", where["activity_group_id"])
	rows, err := r.conn.Query(query)
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

func (r *repository) Update(todo *entities.Todo, time string) error {
	data := fmt.Sprintf("updated_at='%s'", time)
	if todo.Title != "" {
		data = fmt.Sprintf("%s, title='%s'", data, todo.Title)
	}

	if todo.IsActive {
		data = fmt.Sprintf("%s, is_active=%v", data, todo.IsActive)
	}

	if todo.Priority != "" {
		data = fmt.Sprintf("%s, priority='%s'", data, todo.Priority)
	}

	query := fmt.Sprintf("UPDATE todos SET %s WHERE id=%d;", data, todo.ID)
	_, err := r.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("Todo with ID %d Not Found", todo.ID)
	}

	return nil
}

func (r *repository) Delete(where map[string]string) error {
	id, err := strconv.Atoi(where["id"])
	if err != nil {
		return fmt.Errorf("Todo with ID %d Not Found", id)
	}
	query := fmt.Sprintf("UPDATE todos SET deleted_at=CURRENT_TIMESTAMP WHERE id=%d", id)
	row, err := r.conn.Exec(query)
	i, _ := row.RowsAffected()
	if i == 0 {
		return fmt.Errorf("Todo with ID %d Not Found", id)
	}
	if err != nil {
		return err
	}
	return nil
}
