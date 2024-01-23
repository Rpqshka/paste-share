package service

import (
	pasteShare "PasteShare"
	"PasteShare/pkg/repository"
)

type PasteService struct {
	repo repository.Paste
}

func NewPasteService(repo repository.Paste) *PasteService {
	return &PasteService{repo: repo}
}

func (s *PasteService) CreatePaste(userId int, paste pasteShare.Paste) (int, error) {
	return s.repo.CreatePaste(userId, paste)
}

func (s *PasteService) GetAll(userId int) ([]pasteShare.Paste, error) {
	return s.repo.GetAll(userId)
}
