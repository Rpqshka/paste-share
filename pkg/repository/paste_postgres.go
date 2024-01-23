package repository

import (
	pasteShare "PasteShare"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PastePostgres struct {
	db *sqlx.DB
}

func NewPastePostgres(db *sqlx.DB) *PastePostgres {
	return &PastePostgres{db: db}
}

func (r *PastePostgres) CreatePaste(userId int, paste pasteShare.Paste) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createPasteQuery := fmt.Sprintf("INSERT INTO %s (title, description, data, paste_date, likes) VALUES ($1, $2, $3, $4, $5) RETURNING id", pastesTable)
	row := r.db.QueryRow(createPasteQuery, paste.Title, paste.Description, paste.Data, paste.Date, paste.Likes)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersPastesQuery := fmt.Sprintf("INSERT INTO %s (user_id, paste_id) VALUES ($1, $2)", usersPastesTable)
	_, err = tx.Exec(createUsersPastesQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *PastePostgres) GetAll(userId int) ([]pasteShare.Paste, error) {
	var pastes []pasteShare.Paste
	query := fmt.Sprintf(`SELECT p.id, p.title, p.description, p.data, p.paste_date, p.likes FROM %s p
		INNER JOIN %s u ON p.id = u.paste_id WHERE u.user_id = $1`, pastesTable, usersPastesTable)
	err := r.db.Select(&pastes, query, userId)

	return pastes, err
}
