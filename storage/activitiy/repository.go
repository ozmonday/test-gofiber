package activitiy

import (
	"database/sql"
	"fmt"
	"strconv"
	"testfiber/storage/entities"
	"testfiber/utility"
)

type Repository interface {
	Create(*entities.Activity) error
	Read(map[string]string) (*entities.Activity, error)
	Reads(map[string]string) (*[]entities.Activity, error)
	Update(activity *entities.Activity) error
	Delete(where map[string]string) error
}

type repository struct {
	conn *sql.DB
}

func NewRepository(connection *sql.DB) Repository {
	return &repository{
		conn: connection,
	}
}

func (r *repository) Create(activity *entities.Activity) error {
	t := utility.GetTime()
	activity.CreateAt = t
	activity.UpdateAt = t
	query := fmt.Sprintf("INSERT INTO activities (email, title, created_at, updated_at) VALUES ('%s', '%s', '%s', '%s');", activity.Email, activity.Title, activity.CreateAt, activity.UpdateAt)
	_, err := r.conn.Exec(query)
	if err != nil {
		return err
	}

	result := "SELECT LAST_INSERT_ID();"
	row := r.conn.QueryRow(result)
	if err := row.Scan(&activity.ID); err != nil {
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

func (r *repository) Update(activity *entities.Activity) error {
	data := fmt.Sprintf("updated_at='%s'", utility.GetTime())

	if activity.Email != "" {
		data = fmt.Sprintf("%s, email='%s'", data, activity.Email)
	}

	if activity.Title != "" {
		data = fmt.Sprintf("%s, title='%s'", data, activity.Title)
	}

	query := fmt.Sprintf("UPDATE activities SET %s WHERE id=%s;", data, fmt.Sprint(activity.ID))
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
