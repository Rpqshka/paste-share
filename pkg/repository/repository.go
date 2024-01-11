package repository

import (
	pasteShare "PasteShare"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user pasteShare.User) (int, error)
	CheckNickNameAndEmail(nickname, email string) (int, error)
	GetPasswordHash(nickname string) (string, error)
	GetUser(nickname, password string) (pasteShare.User, error)
	CheckEmail(email string) (int, error)
	SendRecoveryMail(toEmail, recoveryCode, expiredAt string) error
	CheckCode(code string) (string, string, error)
	SetNewPassword(id, password string) error
}

type TodoList interface {
}

type TodoItem interface {
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
