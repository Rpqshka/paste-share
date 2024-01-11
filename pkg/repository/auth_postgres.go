package repository

import (
	pasteShare "PasteShare"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user pasteShare.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (nickname, email, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.NickName, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) CheckNickNameAndEmail(nickname, email string) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE nickname = $1 OR email = $2", usersTable)
	err := r.db.Get(&id, query, nickname, email)
	logrus.Print(id, err)
	return id, err
}

func (r *AuthPostgres) GetUser(nickname, password string) (pasteShare.User, error) {
	var user pasteShare.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE nickname = $1 AND password_hash = $2", usersTable)
	err := r.db.Get(&user, query, nickname, password)
	return user, err
}

func (r *AuthPostgres) GetPasswordHash(nickname string) (string, error) {
	var hash string
	query := fmt.Sprintf("SELECT password_hash FROM %s WHERE nickname = $1", usersTable)
	err := r.db.Get(&hash, query, nickname)
	return hash, err
}

func (r *AuthPostgres) CheckEmail(email string) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE email = $1", usersTable)
	err := r.db.Get(&id, query, email)
	return id, err
}

func (r *AuthPostgres) SendRecoveryMail(toEmail, recoveryCode, expiredAt string) error {
	query := fmt.Sprintf("UPDATE %s SET recovery_code = $1, expired_at = $2 WHERE email = $3", usersTable)
	_, err := r.db.Exec(query, recoveryCode, expiredAt, toEmail)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) CheckCode(code string) (string, string, error) {
	var id string
	var expiredAt string
	query := fmt.Sprintf("SELECT id, expired_at FROM %s WHERE recovery_code = $1", usersTable)
	err := r.db.QueryRow(query, code).Scan(&id, &expiredAt)
	if err != nil {
		return "", "", err
	}
	return id, expiredAt, err
}

func (r *AuthPostgres) SetNewPassword(id, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password_hash = $1 WHERE id = $2", usersTable)
	_, err := r.db.Exec(query, password, id)
	if err != nil {
		return err
	}
	return nil
}
