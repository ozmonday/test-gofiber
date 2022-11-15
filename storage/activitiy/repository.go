package activitiy

import (
	"database/sql"
	"testfiber/storage/entities"
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
	return nil
}

func (r *repository) Read(where map[string]string) (*entities.Activity, error) {
	result := entities.Activity{}
	return &result, nil
}

func (r *repository) Reads(where map[string]string) (*[]entities.Activity, error) {
	result := []entities.Activity{}
	return &result, nil
}

func (r *repository) Update(activity *entities.Activity) error {
	return nil
}

func (r *repository) Delete(where map[string]string) error {
	return nil
}
