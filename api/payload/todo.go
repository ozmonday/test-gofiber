package payload

import "database/sql"

type Todo struct {
	ID         sql.NullInt64  `json:"id"`
	Title      sql.NullString `json:"title"`
	ActivityID sql.NullInt64  `json:"activity_group_id"`
	IsActive   sql.NullBool   `json:"is_active"`
	Priority   sql.NullString `json:"priority"`
	Create     sql.NullString `json:"created_at"`
	Update     sql.NullString `json:"updated_at"`
	Delete     sql.NullString `json:"deleted_at"`
}
