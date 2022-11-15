package entities

import "github.com/gofiber/fiber/v2"

type Todo struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	ActivityID int64  `json:"activity_group_id"`
	IsActive   bool   `json:"is_active"`
	Priority   string `json:"priority"`
	CreateAt   string `json:"created_at"`
	UpdateAt   string `json:"updated_at"`
	DeleteAt   string `json:"deleted_at"`
}

func (t *Todo) Map() *fiber.Map {
	result := fiber.Map{
		"id":                t.ID,
		"title":             t.Title,
		"activity_group_id": t.ActivityID,
		"is_active":         t.IsActive,
		"priority":          t.Priority,
		"created_at":        t.CreateAt,
		"updated_at":        t.UpdateAt,
		"deleted_at":        t.DeleteAt,
	}

	if t.DeleteAt == "" {
		result["deleted_at"] = nil
	}

	return &result
}
