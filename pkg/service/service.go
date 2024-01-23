package service

import (
	pasteShare "PasteShare"
	"PasteShare/pkg/repository"
)

type Authorization interface {
	CreateUser(user pasteShare.User) (int, error)
	CheckNickNameAndEmail(nickname, email string) (int, error)
	GetPasswordHash(nickname string) (string, error)
	GenerateToken(nickname, passwordHash string) (string, error)
	ParseToken(accessToken string) (int, error)
	CheckEmail(email string) (int, error)
	SendRecoveryMail(toEmail, code, expiredAt string) error
	CheckCode(code string) (string, string, error)
	SetNewPassword(id, password string) error
}

type Paste interface {
	CreatePaste(userId int, paste pasteShare.Paste) (int, error)
	GetAll(userId int) ([]pasteShare.Paste, error)
}

type TodoItem interface {
}

type Service struct {
	Authorization
	Paste
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Paste:         NewPasteService(repos.Paste),
	}
}
