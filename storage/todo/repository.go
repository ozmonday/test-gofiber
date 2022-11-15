package todo

import (
	"database/sql"
	"testfiber/storage/entities"
)

type Repository interface {
	Create(*entities.Todo) error
	Read(map[string]string) (*entities.Todo, error)
	Reads(map[string]string) (*[]entities.Todo, error)
	Update(activity *entities.Todo) error
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

func (r *repository) Create(activity *entities.Todo) error {
	return nil
}

func (r *repository) Read(where map[string]string) (*entities.Todo, error) {
	result := entities.Todo{}
	return &result, nil
}

func (r *repository) Reads(where map[string]string) (*[]entities.Todo, error) {
	result := []entities.Todo{}
	return &result, nil
}

func (r *repository) Update(activity *entities.Todo) error {
	return nil
}

func (r *repository) Delete(where map[string]string) error {
	return nil
}
