package domain

type Order struct {
	ID     string `db:"id"`
	UserID string `db:"user_id"`
	Status string `db:"status"`
}
